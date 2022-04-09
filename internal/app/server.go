package app

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4"
	"github.com/tim3-p/go-link-shortener/internal/configs"
	"github.com/tim3-p/go-link-shortener/internal/models"
	"github.com/tim3-p/go-link-shortener/internal/pkg"
	"github.com/tim3-p/go-link-shortener/internal/storage"
)

type AppHandler struct {
	storage storage.Repository
	tChan   chan *models.Task
}

func NewAppHandler(s storage.Repository, tChan chan *models.Task) *AppHandler {
	return &AppHandler{storage: s, tChan: make(chan *models.Task)}
}

func NewRouter(handler *AppHandler) chi.Router {
	r := chi.NewRouter()
	r.Use(GzipHandle, AuthHandle)
	r.Get("/{ID}", handler.GetHandler)
	r.Post("/", handler.PostHandler)
	r.Post("/api/shorten", handler.ShortenHandler)
	r.Get("/api/user/urls", handler.UserUrls)
	r.Get("/ping", handler.DBPing)
	r.Post("/api/shorten/batch", handler.ShortenBatchHandler)
	r.Delete("/api/user/urls", handler.DeleteBatchHandler)

	return r
}

func (h *AppHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "ID")

	v, err := h.storage.Get(urlID, userIDVar)
	if err != nil {
		/*
			if errors.Is(err, errors.New("URL is deleted")) {
				w.WriteHeader(http.StatusGone)
				return
			}
		*/
		http.Error(w, "ID not found", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Location", v)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *AppHandler) PostHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	urlHash := pkg.HashURL(b)

	err = h.storage.Add(urlHash, string(b), userIDVar)
	status, err := pkg.CheckDBError(err)

	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	w.WriteHeader(status)
	w.Write([]byte(configs.EnvConfig.BaseURL + "/" + urlHash))
}

func (h *AppHandler) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	var req models.ShortenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	urlHash := pkg.HashURL([]byte(req.URL))
	err := h.storage.Add(urlHash, string(req.URL), userIDVar)
	status, err := pkg.CheckDBError(err)

	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	res := models.ShortenResponse{Result: configs.EnvConfig.BaseURL + "/" + urlHash}

	jsonRes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Add("Accept", "application/json")
	w.WriteHeader(status)
	w.Write(jsonRes)
}

func (h *AppHandler) UserUrls(w http.ResponseWriter, r *http.Request) {

	mapRes, err := h.storage.GetUserURLs(userIDVar)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(mapRes) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	res := []models.UserURL{}

	for key, element := range mapRes {
		item := models.UserURL{ShortURL: configs.EnvConfig.BaseURL + "/" + key, OriginalURL: element}
		res = append(res, item)
	}

	jsonRes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Add("Accept", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonRes)
}

func (h *AppHandler) DBPing(w http.ResponseWriter, r *http.Request) {
	conn, err := pgx.Connect(context.Background(), configs.EnvConfig.DatabaseDSN)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Close(context.Background())
	w.WriteHeader(http.StatusOK)
}

func (h *AppHandler) ShortenBatchHandler(w http.ResponseWriter, r *http.Request) {
	var req []models.ShortenBatchRequest
	var res []models.ShortenBatchResponse

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, value := range req {
		urlHash := pkg.HashURL([]byte(value.OriginalURL))
		h.storage.Add(urlHash, string(value.OriginalURL), userIDVar)

		res = append(res, models.ShortenBatchResponse{CorrelationID: value.CorrelationID, ShortURL: configs.EnvConfig.BaseURL + "/" + urlHash})
	}

	jsonRes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Add("Accept", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonRes)
}

func (h *AppHandler) DeleteBatchHandler(w http.ResponseWriter, r *http.Request) {
	var req []string

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.tChan <- &models.Task{
		URLs:   req,
		UserID: userIDVar,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusAccepted)
}

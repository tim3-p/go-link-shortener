package app

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tim3-p/go-link-shortener/internal/configs"
	"github.com/tim3-p/go-link-shortener/internal/models"
	"github.com/tim3-p/go-link-shortener/internal/pkg"
	"github.com/tim3-p/go-link-shortener/internal/storage"
)

type AppHandler struct {
	storage storage.Repository
}

func NewAppHandler(s storage.Repository) *AppHandler {
	return &AppHandler{storage: s}
}

func NewRouter(handler *AppHandler) chi.Router {
	r := chi.NewRouter()
	r.Use(GzipHandle)
	r.Get("/{ID}", handler.GetHandler)
	r.Post("/", handler.PostHandler)
	r.Post("/api/shorten", handler.ShortenHandler)
	r.Get("/api/user/urls", handler.UserUrls)
	return r
}

func (h *AppHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "ID")

	v, err := h.storage.Get(urlID)
	if err != nil {
		http.Error(w, "ID not found", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, v, http.StatusTemporaryRedirect)
	w.Write(nil)
}

func (h *AppHandler) PostHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	urlHash := pkg.HashURL(b)
	h.storage.Add(urlHash, string(b))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(configs.EnvConfig.BaseURL + "/" + urlHash))
}

func (h *AppHandler) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	var req models.ShortenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	urlHash := pkg.HashURL([]byte(req.URL))
	h.storage.Add(urlHash, string(req.URL))

	res := models.ShortenResponse{Result: configs.EnvConfig.BaseURL + "/" + urlHash}

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

func (h *AppHandler) UserUrls(w http.ResponseWriter, r *http.Request) {

	mapRes, err := h.storage.GetUserURLs()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(mapRes) == 0 {
		w.WriteHeader(http.StatusNoContent)
		w.Write(nil)
		return
	}

	res := models.UserUrlsResponse{}

	for key, element := range mapRes {
		item := models.UserUrl{ShortUrl: configs.EnvConfig.BaseURL + "/" + key, OriginalUrl: element}
		res.UserUrls = append(res.UserUrls, item)
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

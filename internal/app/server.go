package app

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4"
	"github.com/tim3-p/go-link-shortener/internal/configs"
	"github.com/tim3-p/go-link-shortener/internal/models"
	"github.com/tim3-p/go-link-shortener/internal/pkg"
	"github.com/tim3-p/go-link-shortener/internal/storage"
)

// Channel fot batch operations
var TChan = make(chan *models.Task)

// Application repo
type AppHandler struct {
	storage storage.Repository
}

// App Handler constructor
func NewAppHandler(s storage.Repository) *AppHandler {
	return &AppHandler{storage: s}
}

// Http router constructor
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
	r.Get("/debug/pprof/", pprof.Index)
	r.Get("/debug/pprof/cmdline", pprof.Cmdline)
	r.Get("/debug/pprof/profile", pprof.Profile)
	r.Get("/debug/pprof/symbol", pprof.Symbol)
	r.Get("/debug/pprof/trace", pprof.Trace)
	r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	r.Get("/api/internal/stats", handler.StatsHandler)

	return r
}

// GetHandler - Returns full link.
// Output:
// # Request
// GET /{ID}
//
// # Response
// HTTP/1.1 307 OK
// Content-Type: text/plain; charset=utf-8
// Redirect to URL

func (h *AppHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "ID")

	v, err := h.storage.Get(urlID, userIDVar)
	if err != nil {

		if errors.Is(err, storage.ErrURLDeleted) {
			http.Error(w, "URL deleted", http.StatusGone)
			return
		}

		http.Error(w, "ID not found", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Location", v)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// PostHandler - Crate and store short link.
// Output:
// # Request
// POST /
//
// # Response
// HTTP/1.1 201 OK
// Content-Type: text/plain; charset=utf-8
// Response short URL
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

// ShortenHandler - Crate and store short link in json data.
// Output:
// # Request
// POST /api/shorten
//
// {
//   url: “<full URL>”
// }
//
// # Response
// HTTP/1.1 201 OK
// Content-Type: application/json; charset=UTF-8
// {
//   “result”: “<short URL>”
// }
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

// UserUrls - Returns list of all user URLs.
// Output:
// # Request
// GET /api/user/urls
//
//
// # Response
// HTTP/1.1 200 OK
// Content-Type: application/json; charset=UTF-8
// [
//   {
//     “short_url“: “<short URL>”,
//	   “original_url“: “<original URL>”
//   },
//    ...
// ]
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

// DBPing - test connection with DB.
// Output:
// # Request
// GET /
//
// # Response
// HTTP/1.1 200 OK
// Content-Type: text/plain; charset=utf-8
func (h *AppHandler) DBPing(w http.ResponseWriter, r *http.Request) {
	conn, err := pgx.Connect(context.Background(), configs.EnvConfig.DatabaseDSN)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Close(context.Background())
	w.WriteHeader(http.StatusOK)
}

// ShortenBatchHandler - batch add new user URLs
// Output:
// # Request
// POST /api/shorten/batch
//[
// {
//   “correlation_id”: “<string ID>”,
//   “original_url”: “<URL for shorting>”,
// },
// ...
// ]
//
// # Response
// HTTP/1.1 201 OK
// [
//   {
//      “correlation_id”: “<string ID>”,
//      “short_url”: “<short URL>”
//   },
//    ...
// ]
// Content-Type: application/json; charset=UTF-8
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

// DeleteBatchHandler - batch delete list of URLs
// Output:
// # Request
// DELETE /api/user/urls
//
// # Response
// HTTP/1.1 202 OK
// Content-Type: application/json; charset=UTF-8
func (h *AppHandler) DeleteBatchHandler(w http.ResponseWriter, r *http.Request) {
	var req []string

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	TChan <- &models.Task{
		URLs:   req,
		UserID: userIDVar,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusAccepted)
}

// StatsHandler - Returns links and users stats.
// Output:
// # Request
// GET /api/internal/stats
//
// # Response
//   {
//     “urls“: “<int>”,
//	   “users“: “<int>”
//   }
// HTTP/1.1 200 OK
// Content-Type: application/json; charset=UTF-8

func (h *AppHandler) StatsHandler(w http.ResponseWriter, r *http.Request) {

	if configs.EnvConfig.TrustedSubnet == "" {
		http.Error(w, "You don't have access to this handler", http.StatusForbidden)
		return
	}

	_, IPNet, err := net.ParseCIDR(configs.EnvConfig.TrustedSubnet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	if !IPNet.Contains(net.ParseIP(r.Header.Get("X-Real-IP"))) {
		http.Error(w, "You subnet don't have access to this handler", http.StatusForbidden)
		return
	}

	urls, users, err := h.storage.GetStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := models.StatsResponse{URLs: urls, Users: users}

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

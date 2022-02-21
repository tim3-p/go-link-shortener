package app

import (
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tim3-p/go-link-shortener/configs"
	"github.com/tim3-p/go-link-shortener/internal/pkg"
	"github.com/tim3-p/go-link-shortener/internal/storage"
)

func NewRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/{ID}", GetHandler)
	r.Post("/", PostHandler)
	return r
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "ID")

	v, err := storage.Get(urlID)
	if err != nil {
		http.Error(w, "ID not found", http.StatusBadRequest)
		return
	}
	//w.Header().Set("Location", configs.DefaultAddress+v)
	//w.WriteHeader(http.StatusTemporaryRedirect)
	http.Redirect(w, r, configs.DefaultAddress+v, http.StatusTemporaryRedirect)
	w.Write(nil)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	urlHash := pkg.HashURL(b)
	storage.Add(urlHash, string(b))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(configs.DefaultAddress + urlHash))
}

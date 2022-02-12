package app

import (
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tim3-p/go-link-shortener/configs"
	"github.com/tim3-p/go-link-shortener/internal/pkg"
)

var (
	urlBase = make(map[string]string)
)

func NewRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/{ID}", GetHandler)
	r.Post("/", PostHandler)
	return r
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "ID")

	if v, found := urlBase[urlID]; found {
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Header().Set("Location", configs.DefaultAddress+v)
	} else {
		http.Error(w, "ID not found", http.StatusBadRequest)
	}
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	urlHash := pkg.HashURL(b)
	urlBase[urlHash] = string(b)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(urlHash))
}

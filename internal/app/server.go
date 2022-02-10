package app

import (
	"io"
	"net/http"
)

var (
	urlBase = map[string]string{
		"fgRth": "http://habr.com/aaa",
		"tYrg4": "http://codewars.com/aaa",
		"l5Fg3": "http://ru-tracker.com/aaa",
		"tXO2A": "http://gitlab.com/aaa",
		"DtSX":  "http://ovaop0.biz/aa",
	}

	urlBaseReverse = map[string]string{
		"http://habr.com":       "http://localhost:8080/fgRth",
		"http://codewars.com":   "http://localhost:8080/tYrg4",
		"http://ru-tracker.com": "http://localhost:8080/l5Fg3",
		"http://gitlab.com":     "http://localhost:8080/tXO2A",
		"http://ovaop0.biz":     "http://localhost:8080/DtSX",
	}
)

func CommonHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:

		urlID := r.URL.Path[1:]

		if v, found := urlBase[urlID]; found {
			w.WriteHeader(http.StatusTemporaryRedirect)
			w.Header().Set("Location", "http://localhost:8080/"+v)
		} else {
			http.Error(w, "ID not found", http.StatusBadRequest)
		}
		return

	case http.MethodPost:
		b, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if v, found := urlBaseReverse[string(b)]; found {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(v))
		} else {
			http.Error(w, "This logic will be implemented in future", http.StatusBadRequest)
		}

		return

	default:
		http.Error(w, "Unsupported method, use POST or GET method", http.StatusBadRequest)
		return
	}
}

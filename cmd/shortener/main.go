package main

import (
	"net/http"

	"github.com/tim3-p/go-link-shortener/internal/app"
)

func main() {
	http.HandleFunc("/", app.CommonHandler)
	http.ListenAndServe(":8080", nil)
}

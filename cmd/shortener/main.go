package main

import (
	"net/http"

	"github.com/tim3-p/go-link-shortener/internal/app"
)

func main() {
	r := app.NewRouter()
	http.ListenAndServe(":8080", r)
}

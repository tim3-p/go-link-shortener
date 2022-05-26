package main

import (
	"log"
	"net/http"

	"github.com/caarlos0/env"
	"github.com/tim3-p/go-link-shortener/internal/app"
	"github.com/tim3-p/go-link-shortener/internal/configs"
	"github.com/tim3-p/go-link-shortener/internal/storage"
)

func InitConfig() {
	err := env.Parse(&configs.EnvConfig)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	InitConfig()

	var repository storage.Repository
	if configs.EnvConfig.FileStoragePath == "" {
		repository = storage.NewMapRepository()
	} else {
		repository = storage.NewFileRepository(configs.EnvConfig.FileStoragePath)
	}

	handler := app.NewAppHandler(repository)

	r := app.NewRouter(handler)
	http.ListenAndServe(configs.EnvConfig.ServerAddress, r)
}

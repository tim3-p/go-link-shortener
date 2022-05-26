package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/caarlos0/env"
	"github.com/tim3-p/go-link-shortener/internal/app"
	"github.com/tim3-p/go-link-shortener/internal/configs"
	"github.com/tim3-p/go-link-shortener/internal/storage"
)

func SetCommandLineFlags() {
	flag.StringVar(&configs.EnvConfig.ServerAddress, "a", configs.EnvConfig.ServerAddress, "server http address")
	flag.StringVar(&configs.EnvConfig.BaseURL, "b", configs.EnvConfig.BaseURL, "base url of shortener")
	flag.StringVar(&configs.EnvConfig.FileStoragePath, "f", configs.EnvConfig.FileStoragePath, "file storage path")
	flag.Parse()
}

func InitConfig() error {
	err := env.Parse(&configs.EnvConfig)
	if err != nil {
		return err
	}
	SetCommandLineFlags()
	return nil
}

func main() {
	err := InitConfig()

	if err != nil {
		log.Fatal(err)
	}

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

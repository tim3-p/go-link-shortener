package main

import (
	"log"
	"net/http"

	"github.com/caarlos0/env"
	"github.com/tim3-p/go-link-shortener/internal/app"
	"github.com/tim3-p/go-link-shortener/internal/configs"
)

func InitConfig() {
	err := env.Parse(&configs.EnvConfig)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	InitConfig()
	r := app.NewRouter()
	http.ListenAndServe(configs.EnvConfig.ServerAddress, r)
}

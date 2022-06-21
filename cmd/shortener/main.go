package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/caarlos0/env"
	"github.com/jackc/pgx/v4"
	"github.com/tim3-p/go-link-shortener/internal/app"
	"github.com/tim3-p/go-link-shortener/internal/configs"
	"github.com/tim3-p/go-link-shortener/internal/storage"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func SetCommandLineFlags() {
	flag.StringVar(&configs.EnvConfig.ServerAddress, "a", configs.EnvConfig.ServerAddress, "server http address")
	flag.StringVar(&configs.EnvConfig.BaseURL, "b", configs.EnvConfig.BaseURL, "base url of shortener")
	flag.StringVar(&configs.EnvConfig.FileStoragePath, "f", configs.EnvConfig.FileStoragePath, "file storage path")
	flag.StringVar(&configs.EnvConfig.DatabaseDSN, "d", configs.EnvConfig.DatabaseDSN, "database connection string")
	flag.Parse()
}

func InitConfig() error {
	err := env.Parse(&configs.EnvConfig)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

	err := InitConfig()
	if err != nil {
		log.Fatal(err)
	}
	SetCommandLineFlags()
	var repository storage.Repository

	if configs.EnvConfig.DatabaseDSN != "" {
		conn, err := pgx.Connect(context.Background(), configs.EnvConfig.DatabaseDSN)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close(context.Background())

		repository, err = storage.NewDBRepository(conn)
		if err != nil {
			log.Fatal(err)
		}
	} else if configs.EnvConfig.FileStoragePath == "" {
		repository = storage.NewMapRepository()
	} else {
		repository = storage.NewFileRepository(configs.EnvConfig.FileStoragePath)
	}

	for k := 0; k <= 3; k++ {

		go func() {
			for task := range app.TChan {
				repository.Delete(task.URLs, task.UserID)
			}
		}()
	}

	handler := app.NewAppHandler(repository)

	r := app.NewRouter(handler)
	http.ListenAndServe(configs.EnvConfig.ServerAddress, r)
}

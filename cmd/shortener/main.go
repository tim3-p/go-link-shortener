package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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
	flag.BoolVar(&configs.EnvConfig.EnableHTTPS, "s", configs.EnvConfig.EnableHTTPS, "HTTPS server")
	flag.StringVar(&configs.EnvConfig.ConfigJson, "с", configs.EnvConfig.ConfigJson, "read config from json file")
	flag.Parse()
}

func InitConfig() error {
	err := env.Parse(&configs.EnvConfig)
	if err != nil {
		return err
	}
	return nil
}

func InitJsonConfig() error {
	jsonFile, err := os.OpenFile(configs.EnvConfig.ConfigJson, os.O_RDONLY, 0644)
	if err != nil {
		return err
	} else {
		jsonBody, err := io.ReadAll(jsonFile)
		if err != nil {
			return err
		} else if err = json.Unmarshal(jsonBody, &configs.EnvConfig); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := InitConfig()
	if err != nil {
		log.Fatal(err)
	}
	SetCommandLineFlags()

	if configs.EnvConfig.ConfigJson != "" {
		err = InitJsonConfig()
		if err != nil {
			log.Fatal(err)
		}
	}
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

	srv := &http.Server{
		Addr:    configs.EnvConfig.ServerAddress,
		Handler: r,
	}

	if configs.EnvConfig.EnableHTTPS {
		err = app.GenerateCert()
		if err != nil {
			log.Fatal(err)
		}
		go srv.ListenAndServeTLS(app.CertFile, app.KeyFile)
	} else {
		go srv.ListenAndServe()
	}

	sh := make(chan os.Signal, 1)
	signal.Notify(sh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	v := <-sh

	log.Printf("Recived signal: %v", v)

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown with error: %v", err)
	}
}

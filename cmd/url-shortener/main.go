package main

import (
	"github.com/rshelekhov/url-shortener/internal/config"
	"github.com/rshelekhov/url-shortener/internal/lib/logger/sl"
	"github.com/rshelekhov/url-shortener/internal/storage/postgres"
	"github.com/rshelekhov/url-shortener/pkg/logs"
	"log/slog"
	"os"
)

func main() {
	cfg := config.MustLoad()

	// TODO: save logs to file
	log := logs.SetupLogger(cfg.Env)

	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := postgres.NewStorage(cfg.DatabaseURL)
	if err != nil {
		log.Error("failed to init storage: ", sl.Err(err))
		os.Exit(1)
	}

	defer func(storage *postgres.Storage) {
		err := storage.Close()
		if err != nil {
			log.Error("failed to close storage: ", sl.Err(err))
			os.Exit(1)
		}
	}(storage)

	// TODO: init router: chi, "chi render"

	// TODO: run server
}

package main

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/rshelekhov/url-shortener/internal/config"
	"github.com/rshelekhov/url-shortener/internal/http-server/handlers/url/redirect"
	"github.com/rshelekhov/url-shortener/internal/http-server/handlers/url/remove"
	"github.com/rshelekhov/url-shortener/internal/http-server/handlers/url/save"
	mwLogger "github.com/rshelekhov/url-shortener/internal/http-server/middleware/logger"
	"github.com/rshelekhov/url-shortener/internal/lib/logger/sl"
	"github.com/rshelekhov/url-shortener/internal/storage/postgres"
	"github.com/rshelekhov/url-shortener/pkg/logs"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	cfg := config.MustLoad()

	// TODO: save logs to file
	log := logs.SetupLogger(cfg.Env)

	log.Info(
		"starting url-shortener",
		slog.String("env", cfg.Env),
		slog.String("version", "1"))
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

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))

		r.Post("/", save.New(log, storage))
		r.Delete("/{alias}", remove.Url(log, storage))
	})

	router.Get("/{alias}", redirect.New(log, storage))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

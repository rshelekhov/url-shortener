package main

import (
	"github.com/rshelekhov/url-shortener/internal/config"
	"github.com/rshelekhov/url-shortener/pkg/logs"
	"log/slog"
)

func main() {
	cfg := config.MustLoad()

	log := logs.SetupLogger(cfg.Env)

	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// TODO: init storage: sqlite

	// TODO: init router: chi, "chi render"

	// TODO: run server
}

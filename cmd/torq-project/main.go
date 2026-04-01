package main

import (
	"log/slog"
	"net/http"
	"os"

	"zivusomer/torq-project/internal/api/findcountry"
	"zivusomer/torq-project/internal/config"
	"zivusomer/torq-project/internal/httpserver"
	"zivusomer/torq-project/internal/logging"
	"zivusomer/torq-project/internal/ratelimit"
	"zivusomer/torq-project/internal/store/factory"
)

func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		logger := logging.New("info")
		logger.Error("failed to load configuration", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger := logging.New(cfg.LogLevel)

	datastore, err := factory.New(cfg.DatastoreType, cfg.DatastorePath)
	if err != nil {
		logger.Error("failed to initialize datastore", slog.String("error", err.Error()))
		os.Exit(1)
	}

	limiter := ratelimit.New(cfg.RequestsPerSecond)
	findCountryHandler := findcountry.NewHandler(datastore, limiter)
	server := httpserver.New(findCountryHandler)
	addr := ":" + cfg.Port

	logger.Info("application starting",
		slog.String("app", cfg.AppName),
		slog.String("env", cfg.Env),
		slog.String("addr", addr),
		slog.String("datastore_type", cfg.DatastoreType),
		slog.String("datastore_path", cfg.DatastorePath),
		slog.Int("requests_per_second", cfg.RequestsPerSecond),
	)

	if err := http.ListenAndServe(addr, server.Handler()); err != nil {
		logger.Error("server stopped", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

package main

import (
	"log/slog"
	"os"

	"github.com/zivusomer/torq-project/internal/config"
	"github.com/zivusomer/torq-project/internal/logging"
)

func main() {
	cfg := config.LoadFromEnv()

	logger := logging.New(cfg.LogLevel)
	logger.Info("application starting",
		slog.String("app", cfg.AppName),
		slog.String("env", cfg.Env),
	)

	// Keep stdout output deterministic for first-run verification.
	_, _ = os.Stdout.WriteString(cfg.AppName + " is running\n")
}

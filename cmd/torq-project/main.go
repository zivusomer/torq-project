package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"zivusomer/torq-project/internal/app"
	"zivusomer/torq-project/internal/config"
	"zivusomer/torq-project/internal/logging"
)

func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		logging.Logger.Error("failed to load configuration: " + err.Error())
		os.Exit(1)
	}

	application, err := app.New(cfg)
	if err != nil {
		logging.Logger.Error("failed to initialize application: " + err.Error())
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := application.Run(ctx); err != nil {
		logging.Logger.Error("server stopped: " + err.Error())
		os.Exit(1)
	}
}

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
	cfg := initConfig()
	ctx, stop := initSignalContext()
	defer stop()
	application := initApplication(cfg)
	runApplication(application, ctx)
}

func runApplication(application *app.App, ctx context.Context) {
	if err := application.Run(ctx); err != nil {
		logging.Logger.Error("server stopped: " + err.Error())
		os.Exit(1)
	}
}

func initConfig() config.Config {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		logging.Logger.Error("failed to load configuration: " + err.Error())
		os.Exit(1)
	}
	return cfg
}

func initApplication(cfg config.Config) *app.App {
	application, err := app.New(cfg)
	if err != nil {
		logging.Logger.Error("failed to initialize application: " + err.Error())
		os.Exit(1)
	}

	return application
}

func initSignalContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
}

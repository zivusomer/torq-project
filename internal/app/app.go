package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"zivusomer/torq-project/internal/api/findcountry"
	"zivusomer/torq-project/internal/config"
	"zivusomer/torq-project/internal/httpserver"
	"zivusomer/torq-project/internal/logging"
	"zivusomer/torq-project/internal/ratelimit"
	"zivusomer/torq-project/internal/store/factory"
)

type App struct {
	cfg    config.Config
	server *httpserver.Server
}

func New(cfg config.Config) (*App, error) {
	datastore, err := factory.New(cfg.DatastoreType, cfg.DatastorePath)
	if err != nil {
		return nil, fmt.Errorf("initialize datastore: %w", err)
	}

	if err := ratelimit.Init(cfg.RequestsPerSecond); err != nil {
		return nil, fmt.Errorf("initialize rate limiter: %w", err)
	}

	findCountryHandler, err := findcountry.NewHandler(datastore)
	if err != nil {
		return nil, fmt.Errorf("initialize API handlers: %w", err)
	}

	routes := []httpserver.Route{
		{Path: "/v1/find-country", Handler: findCountryHandler},
	}

	server, err := httpserver.New(routes)
	if err != nil {
		return nil, fmt.Errorf("initialize HTTP server: %w", err)
	}

	return &App{
		cfg:    cfg,
		server: server,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	addr := ":" + a.cfg.Port
	httpServer := &http.Server{
		Addr:    addr,
		Handler: a.server.Handler(),
	}

	logging.Logger.Info("application starting")

	errCh := make(chan error, 1)
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			errCh <- fmt.Errorf("shutdown server: %w", err)
			return
		}
		errCh <- nil
	}()

	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("run server: %w", err)
	}

	return <-errCh
}

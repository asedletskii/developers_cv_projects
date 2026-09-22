package app

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"go.uber.org/zap"

	httpserver "developers-cv/backend/internal/platform/http"
	"developers-cv/backend/internal/platform/telemetry"
)

type Application struct {
	server      *httpserver.Server
	telemetry   *telemetry.Service
	logger      *zap.Logger
	once        sync.Once
	shutdownErr error
}

func NewApplication(
	server *httpserver.Server,
	telemetryService *telemetry.Service,
	logger *zap.Logger,
) *Application {
	return &Application{
		server:    server,
		telemetry: telemetryService,
		logger:    logger,
	}
}

func (a *Application) Run(ctx context.Context) error {
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- a.server.Run()
	}()

	select {
	case <-ctx.Done():
		return a.Shutdown(context.Background())
	case err := <-serverErr:
		if err != nil {
			return errors.Join(err, a.Shutdown(context.Background()))
		}
		return nil
	}
}

func (a *Application) Shutdown(ctx context.Context) error {
	a.once.Do(func() {
		serverErr := a.server.Shutdown(ctx)
		telemetryErr := a.telemetry.Shutdown(ctx)
		if serverErr != nil || telemetryErr != nil {
			a.shutdownErr = fmt.Errorf("shutdown application: %w", errors.Join(serverErr, telemetryErr))
		}
	})
	return a.shutdownErr
}

func (a *Application) Logger() *zap.Logger { return a.logger }

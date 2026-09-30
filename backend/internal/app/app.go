package app

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"

	"go.uber.org/zap"

	"developers-cv/backend/internal/platform/database"
	httpserver "developers-cv/backend/internal/platform/http"
	"developers-cv/backend/internal/platform/telemetry"
)

type Application struct {
	server          *httpserver.Server
	telemetry       *telemetry.Service
	logger          *zap.Logger
	gormDB          *gorm.DB
	shutdownTimeout time.Duration
	worker          backgroundWorker
	workerCancel    context.CancelFunc
	workerDone      chan struct{}
	once            sync.Once
	shutdownErr     error
}

func NewApplication(
	server *httpserver.Server,
	telemetryService *telemetry.Service,
	logger *zap.Logger,
	gormDB *gorm.DB,
	shutdownTimeout time.Duration,
	worker backgroundWorker,
) *Application {
	return &Application{
		server:          server,
		telemetry:       telemetryService,
		logger:          logger,
		gormDB:          gormDB,
		shutdownTimeout: shutdownTimeout,
		worker:          worker,
	}
}

type backgroundWorker interface {
	Run(context.Context)
}

func (a *Application) Run(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	workerCtx, workerCancel := context.WithCancel(ctx)
	a.workerCancel = workerCancel
	if a.worker != nil {
		a.workerDone = make(chan struct{})
		go func() {
			defer close(a.workerDone)
			a.worker.Run(workerCtx)
		}()
	}

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- a.server.Run()
	}()

	select {
	case <-ctx.Done():
		return a.Shutdown(context.Background())
	case err := <-serverErr:
		return errors.Join(err, a.Shutdown(context.Background()))
	}
}

func (a *Application) Shutdown(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	shutdownTimeout := a.shutdownTimeout
	if shutdownTimeout <= 0 {
		shutdownTimeout = 10 * time.Second
	}
	shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	a.once.Do(func() {
		serverErr := a.server.Shutdown(shutdownCtx)
		if a.workerCancel != nil {
			a.workerCancel()
			if a.workerDone != nil {
				<-a.workerDone
			}
		}
		telemetryErr := a.telemetry.Shutdown(shutdownCtx)
		var dbErr error
		if a.gormDB != nil {
			dbErr = database.Close(a.gormDB)
		}
		if serverErr != nil || telemetryErr != nil || dbErr != nil {
			a.shutdownErr = fmt.Errorf("shutdown application: %w", errors.Join(serverErr, telemetryErr, dbErr))
		}
	})
	return a.shutdownErr
}

func (a *Application) Logger() *zap.Logger { return a.logger }

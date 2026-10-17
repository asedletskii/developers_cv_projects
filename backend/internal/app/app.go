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
)

type HTTPServerLifecycle interface {
	Run() error
	Shutdown(context.Context) error
}

type TelemetryLifecycle interface {
	Shutdown(context.Context) error
}

type Application struct {
	server          HTTPServerLifecycle
	telemetry       TelemetryLifecycle
	logger          *zap.Logger
	gormDB          *gorm.DB
	shutdownTimeout time.Duration
	workers         []backgroundWorker
	workerCancel    context.CancelFunc
	workerDone      chan struct{}
	once            sync.Once
	shutdownErr     error
}

func NewApplication(
	server HTTPServerLifecycle,
	telemetryService TelemetryLifecycle,
	logger *zap.Logger,
	gormDB *gorm.DB,
	shutdownTimeout time.Duration,
	workers []backgroundWorker,
) *Application {
	return &Application{
		server:          server,
		telemetry:       telemetryService,
		logger:          logger,
		gormDB:          gormDB,
		shutdownTimeout: shutdownTimeout,
		workers:         workers,
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
	activeWorkers := make([]backgroundWorker, 0, len(a.workers))
	for _, worker := range a.workers {
		if worker != nil {
			activeWorkers = append(activeWorkers, worker)
		}
	}
	if len(activeWorkers) > 0 {
		a.workerDone = make(chan struct{})
		var workerGroup sync.WaitGroup
		workerGroup.Add(len(activeWorkers))
		for _, worker := range activeWorkers {
			go func(w backgroundWorker) {
				defer workerGroup.Done()
				w.Run(workerCtx)
			}(worker)
		}
		go func() {
			workerGroup.Wait()
			close(a.workerDone)
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
		var workerErr error
		if a.workerCancel != nil {
			a.workerCancel()
			if err := waitForWorkers(shutdownCtx, a.workerDone); err != nil {
				workerErr = fmt.Errorf("wait for background workers: %w", err)
			}
		}
		var telemetryErr error
		if a.telemetry != nil {
			telemetryErr = a.telemetry.Shutdown(shutdownCtx)
		}
		var dbErr error
		if a.gormDB != nil {
			dbErr = database.Close(a.gormDB)
		}
		if serverErr != nil || workerErr != nil || telemetryErr != nil || dbErr != nil {
			a.shutdownErr = fmt.Errorf("shutdown application: %w", errors.Join(serverErr, workerErr, telemetryErr, dbErr))
		}
	})
	return a.shutdownErr
}

func (a *Application) Logger() *zap.Logger { return a.logger }

func waitForWorkers(ctx context.Context, done <-chan struct{}) error {
	if done == nil {
		return nil
	}
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

package app

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/gorm"

	"go.uber.org/zap"

	"developers-cv/backend/internal/platform/database"
	httpserver "developers-cv/backend/internal/platform/http"
	"developers-cv/backend/internal/platform/telemetry"
)

type Application struct {
	server       *httpserver.Server
	telemetry    *telemetry.Service
	logger       *zap.Logger
	postgresPool *pgxpool.Pool
	gormDB       *gorm.DB
	once         sync.Once
	shutdownErr  error
}

func NewApplication(
	server *httpserver.Server,
	telemetryService *telemetry.Service,
	logger *zap.Logger,
	postgresPool *pgxpool.Pool,
	gormDB *gorm.DB,
) *Application {
	return &Application{
		server:       server,
		telemetry:    telemetryService,
		logger:       logger,
		postgresPool: postgresPool,
		gormDB:       gormDB,
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
		var dbErr error
		if a.gormDB != nil {
			dbErr = database.Close(a.gormDB)
		}
		if a.postgresPool != nil {
			a.postgresPool.Close()
		}
		if serverErr != nil || telemetryErr != nil || dbErr != nil {
			a.shutdownErr = fmt.Errorf("shutdown application: %w", errors.Join(serverErr, telemetryErr, dbErr))
		}
	})
	return a.shutdownErr
}

func (a *Application) Logger() *zap.Logger { return a.logger }

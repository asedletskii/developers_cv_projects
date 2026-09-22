//go:build wireinject

package app

import (
	"context"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"

	"developers-cv/backend/internal/platform/config"
	"developers-cv/backend/internal/platform/health"
	httpserver "developers-cv/backend/internal/platform/http"
	"developers-cv/backend/internal/platform/logger"
	"developers-cv/backend/internal/platform/postgres"
	"developers-cv/backend/internal/platform/telemetry"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumepersistence "developers-cv/backend/internal/modules/resume/infrastructure/persistence/reader"
	resumehttp "developers-cv/backend/internal/modules/resume/transport/http"
)

func Initialize() (*Application, error) {
	panic(wire.Build(
		config.Load,
		func(cfg config.Config) (*pgxpool.Pool, error) {
			return postgres.NewPool(context.Background(), postgres.Config{DSN: cfg.Postgres.ConnectionString(), MaxConns: cfg.Postgres.MaxConns, MinConns: cfg.Postgres.MinConns, MaxConnLifetime: cfg.Postgres.MaxConnLifetime, MaxConnIdleTime: cfg.Postgres.MaxConnIdleTime})
		},
		logger.New,
		health.NewService,
		telemetry.New,
		telemetry.NewMetrics,
		resumepersistence.NewPersonReader,
		postgres.NewHealthChecker,
		resumequeries.NewGetPersonHandler,
		resumehttp.NewHandler,
		httpserver.NewRouter,
		httpserver.NewServer,
		NewApplication,
	))
}

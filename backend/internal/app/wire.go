//go:build wireinject

package app

import (
	"context"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumehttp "developers-cv/backend/internal/modules/resume/transport/http"

	"developers-cv/backend/internal/platform/config"
	httpserver "developers-cv/backend/internal/platform/http"
	"developers-cv/backend/internal/platform/logger"
	"developers-cv/backend/internal/platform/telemetry"
)

func Initialize() (*Application, error) {
	panic(wire.Build(
		config.Load,

		func(cfg config.Config) (*pgxpool.Pool, error) {
			return newPostgresPool(context.Background(), cfg)
		},

		logger.New,
		newHealthService,
		telemetry.New,
		telemetry.NewMetrics,

		newPersonReader,

		resumequeries.NewGetPersonHandler,
		resumehttp.NewHandler,
		httpserver.NewRouter,
		httpserver.NewServer,
		NewApplication,
	))
}

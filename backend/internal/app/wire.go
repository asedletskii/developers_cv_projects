//go:build wireinject

package app

import (
	"github.com/google/wire"

	"developers-cv/backend/internal/platform/config"
	"developers-cv/backend/internal/platform/health"
	httpserver "developers-cv/backend/internal/platform/http"
	"developers-cv/backend/internal/platform/logger"
	"developers-cv/backend/internal/platform/telemetry"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumemock "developers-cv/backend/internal/modules/resume/infrastructure/mock"
	resumehttp "developers-cv/backend/internal/modules/resume/transport/http"
)

func Initialize() (*Application, error) {
	panic(wire.Build(
		config.Load,
		logger.New,
		health.NewService,
		telemetry.New,
		telemetry.NewMetrics,
		resumemock.NewPersonReader,
		resumequeries.NewGetPersonHandler,
		resumehttp.NewHandler,
		httpserver.NewRouter,
		httpserver.NewServer,
		NewApplication,
	))
}

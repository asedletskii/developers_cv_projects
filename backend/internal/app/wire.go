//go:build wireinject

package app

import (
	"context"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"

	resumecommands "developers-cv/backend/internal/modules/resume/application/commands"
	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumerepository "developers-cv/backend/internal/modules/resume/infrastructure/persistence/repository"
	resumehttp "developers-cv/backend/internal/modules/resume/transport/http"
	userapplication "developers-cv/backend/internal/modules/user/application"
	userauth "developers-cv/backend/internal/modules/user/infrastructure/auth"
	userrepository "developers-cv/backend/internal/modules/user/infrastructure/persistence/repository"

	"developers-cv/backend/internal/platform/config"
	"developers-cv/backend/internal/platform/database"
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
		newGetResumeReader,
		newListResumesReader,
		database.NewGORMFromConfig,
		database.NewUnitOfWork,
		resumerepository.NewResumeRepository,
		wire.Bind(new(resumecommands.ResumeRepository), new(*resumerepository.ResumeRepository)),
		wire.Bind(new(resumecommands.UnitOfWork), new(*database.UnitOfWork)),
		userrepository.NewUserRepository,
		wire.Bind(new(userapplication.UserRepository), new(*userrepository.UserRepository)),
		userauth.NewMockAuthenticator,
		newAuthenticate,
		resumequeries.NewGetResumeHandler,
		resumequeries.NewListResumesHandler,
		resumecommands.NewCreateResumeHandler,
		resumecommands.NewUpdateResumeHandler,
		resumecommands.NewPatchResumeHandler,
		resumehttp.NewHandler,
		httpserver.NewRouter,
		httpserver.NewServer,
		NewApplication,
	))
}

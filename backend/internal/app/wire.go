//go:build wireinject

package app

import (
	"context"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"

	resumecommands "developers-cv/backend/internal/modules/resume/application/commands"
	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumerepository "developers-cv/backend/internal/modules/resume/infrastructure/persistence/repository/resume"
	resumehttp "developers-cv/backend/internal/modules/resume/transport/http"
	userqueries "developers-cv/backend/internal/modules/user/application/queries"
	userauth "developers-cv/backend/internal/modules/user/infrastructure/auth"
	userreaders "developers-cv/backend/internal/modules/user/infrastructure/persistence/readers"
	userhttp "developers-cv/backend/internal/modules/user/transport/http"

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
		newGetProfileReader,
		database.NewUnitOfWork,
		resumerepository.NewResumeRepository,
		wire.Bind(new(resumecommands.CreateResumeRepository), new(*resumerepository.ResumeRepository)),
		wire.Bind(new(resumecommands.UpdateResumeRepository), new(*resumerepository.ResumeRepository)),
		wire.Bind(new(resumecommands.GetResumeRepository), new(*resumerepository.ResumeRepository)),
		wire.Bind(new(resumecommands.UpdateResumeTitleRepository), new(*resumerepository.ResumeRepository)),
		wire.Bind(new(resumecommands.UnitOfWork), new(*database.UnitOfWork)),
		userreaders.NewUserReader,
		wire.Bind(new(userqueries.UserReader), new(*userreaders.UserReader)),
		userauth.NewMockAuthenticator,
		newAuthenticate,
		resumequeries.NewGetResumeHandler,
		resumequeries.NewListResumesHandler,
		userqueries.NewGetProfileHandler,
		userhttp.NewHandler,
		resumecommands.NewCreateResumeHandler,
		resumecommands.NewUpdateResumeHandler,
		resumecommands.NewPatchResumeHandler,
		resumehttp.NewHandler,
		httpserver.NewRouter,
		httpserver.NewServer,
		NewApplication,
	))
}

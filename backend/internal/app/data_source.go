package app

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumereaders "developers-cv/backend/internal/modules/resume/infrastructure/persistence/readers"
	userqueries "developers-cv/backend/internal/modules/user/application/queries"
	userauth "developers-cv/backend/internal/modules/user/infrastructure/auth"
	userreaders "developers-cv/backend/internal/modules/user/infrastructure/persistence/readers"
	"developers-cv/backend/internal/platform/config"
	"developers-cv/backend/internal/platform/health"
	httpmiddleware "developers-cv/backend/internal/platform/http/middleware"
	"developers-cv/backend/internal/platform/postgres"
)

func newPostgresPool(
	ctx context.Context,
	cfg config.Config,
) (*pgxpool.Pool, error) {
	return postgres.NewPool(ctx,
		postgres.Config{
			DSN:             cfg.Postgres.ConnectionString(),
			MaxConns:        cfg.Postgres.MaxConns,
			MinConns:        cfg.Postgres.MinConns,
			MaxConnLifetime: cfg.Postgres.MaxConnLifetime,
			MaxConnIdleTime: cfg.Postgres.MaxConnIdleTime,
		},
	)
}

func newPublicResumeReader(pool *pgxpool.Pool) resumequeries.PublicResumeReader {
	return resumereaders.NewPublicResumeReader(pool)
}

func newListResumesReader(pool *pgxpool.Pool) resumequeries.ListResumesReader {
	return resumereaders.NewListResumesReader(pool)
}

func newResumeMetadataReader(pool *pgxpool.Pool) resumequeries.ResumeMetadataReader {
	return resumereaders.NewResumeMetadataReader(pool)
}

func newResumeOverviewReader(pool *pgxpool.Pool) resumequeries.OverviewReader {
	return resumereaders.NewOverviewReader(pool)
}

func newResumeReader(
	metadata resumequeries.ResumeMetadataReader,
	overview resumequeries.OverviewReader,
	profile resumequeries.ResumeProfileReader,
	contacts resumequeries.ResumeContactsReader,
	skills resumequeries.SkillsReader,
	languages resumequeries.LanguagesReader,
	technologies resumequeries.TechnologiesReader,
	experiences resumequeries.ExperiencesReader,
	educations resumequeries.EducationsReader,
	courses resumequeries.CoursesReader,
	activities resumequeries.ActivitiesReader,
	projects resumequeries.ProjectsReader,
) resumequeries.ResumeReader {
	return resumereaders.NewResumeReader(metadata, overview, profile, contacts, skills, languages, technologies, experiences, educations, courses, activities, projects)
}

func newGetProfileReader(pool *pgxpool.Pool) userqueries.GetProfileReader {
	return userreaders.NewGetProfileReader(pool)
}

func newAuthenticate(authenticator *userauth.MockAuthenticator) httpmiddleware.AuthenticateFunc {
	return func(ctx context.Context, request *http.Request) (uuid.UUID, error) {
		userID, err := authenticator.Authenticate(ctx, request)
		if err != nil {
			if errors.Is(err, userauth.ErrUnauthenticated) {
				return uuid.Nil, httpmiddleware.ErrUnauthenticated
			}
			return uuid.Nil, err
		}
		return userID, nil
	}
}

func newHealthService(pool *pgxpool.Pool) *health.Service {
	return health.NewService(postgres.NewHealthChecker(pool))
}

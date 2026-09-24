package app

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/reader"
	userqueries "developers-cv/backend/internal/modules/user/application/queries"
	userauth "developers-cv/backend/internal/modules/user/infrastructure/auth"
	userreader "developers-cv/backend/internal/modules/user/infrastructure/persistence/reader"
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

func newGetResumeReader(pool *pgxpool.Pool) resumequeries.GetResumeReader {
	return reader.NewGetResumeReader(pool)
}

func newListResumesReader(pool *pgxpool.Pool) resumequeries.ListResumesReader {
	return reader.NewListResumesReader(pool)
}

func newGetProfileReader(pool *pgxpool.Pool) userqueries.GetProfileReader {
	return userreader.NewGetProfileReader(pool)
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

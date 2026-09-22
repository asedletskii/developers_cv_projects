package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumereader "developers-cv/backend/internal/modules/resume/infrastructure/persistence/reader"
	"developers-cv/backend/internal/platform/config"
	"developers-cv/backend/internal/platform/health"
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

func newResumeReader(pool *pgxpool.Pool) resumequeries.ResumeReader {
	return resumereader.NewResumeReader(pool)
}

func newHealthService(pool *pgxpool.Pool) *health.Service {
	return health.NewService(postgres.NewHealthChecker(pool))
}

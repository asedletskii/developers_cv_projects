package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumemock "developers-cv/backend/internal/modules/resume/infrastructure/mock"
	resumereader "developers-cv/backend/internal/modules/resume/infrastructure/persistence/reader"
	"developers-cv/backend/internal/platform/config"
	"developers-cv/backend/internal/platform/health"
	"developers-cv/backend/internal/platform/postgres"
)

type DataSource uint8

const (
	DataSourcePostgres DataSource = iota
	DataSourceMock
)

const dataSource = DataSourceMock

func newPostgresPool(
	ctx context.Context,
	cfg config.Config,
) (*pgxpool.Pool, error) {
	if dataSource == DataSourceMock {
		return nil, nil
	}

	return postgres.NewPool(ctx, postgres.Config{
		DSN:             cfg.Postgres.ConnectionString(),
		MaxConns:        cfg.Postgres.MaxConns,
		MinConns:        cfg.Postgres.MinConns,
		MaxConnLifetime: cfg.Postgres.MaxConnLifetime,
		MaxConnIdleTime: cfg.Postgres.MaxConnIdleTime,
	})
}

func newPersonReader(pool *pgxpool.Pool) resumequeries.PersonReader {
	switch dataSource {
	case DataSourceMock:
		return resumemock.NewPersonReader()
	case DataSourcePostgres:
		return resumereader.NewPersonReader(pool)
	default:
		panic("unsupported resume data source")
	}
}

func newHealthService(pool *pgxpool.Pool) *health.Service {
	if dataSource == DataSourceMock {
		return health.NewService()
	}

	return health.NewService(postgres.NewHealthChecker(pool))
}

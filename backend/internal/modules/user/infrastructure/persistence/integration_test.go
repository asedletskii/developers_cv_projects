//go:build integration

package persistence_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"developers-cv/backend/internal/modules/user/domain"
	"developers-cv/backend/internal/modules/user/domain/valueobjects"
	"developers-cv/backend/internal/modules/user/infrastructure/persistence/repository"
	"developers-cv/backend/internal/platform/database"
)

func TestPostgresUserRepositoryRoundTrip(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx := context.Background()
	container, err := testpostgres.Run(ctx, "postgres:17-alpine",
		testpostgres.WithDatabase("developers_cv"),
		testpostgres.WithUsername("developers_cv"),
		testpostgres.WithPassword("test-password"),
		testpostgres.BasicWaitStrategies(),
	)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, container.Terminate(ctx)) })

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	applyUsersMigration(t, ctx, dsn)

	gormDB, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(gormDB)) })

	repo := repository.NewUserRepository(gormDB)
	id, err := valueobjects.IDFromUUID(uuid.New())
	require.NoError(t, err)
	createdAt := time.Date(2026, 9, 23, 8, 0, 0, 0, time.UTC)
	user, err := domain.NewUser(domain.UserParams{ID: id, CreatedAt: createdAt, UpdatedAt: createdAt})
	require.NoError(t, err)

	require.NoError(t, repo.Save(ctx, user))

	loaded, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, user.ID().String(), loaded.ID().String())
	require.Equal(t, user.CreatedAt(), loaded.CreatedAt())
}

func applyUsersMigration(t *testing.T, ctx context.Context, dsn string) {
	t.Helper()

	config, err := pgx.ParseConfig(dsn)
	require.NoError(t, err)
	conn, err := pgx.ConnectConfig(ctx, config)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	content, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "..", "db", "migrations", "000008_users.up.sql"))
	require.NoError(t, err)
	_, err = conn.Exec(ctx, string(content))
	require.NoError(t, err)
}

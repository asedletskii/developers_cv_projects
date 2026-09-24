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

	"developers-cv/backend/internal/modules/user/domain/valueobjects"
	"developers-cv/backend/internal/modules/user/infrastructure/persistence/readers"
	"developers-cv/backend/internal/platform/database"
)

func TestPostgresUserReaderGetByID(t *testing.T) {
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

	config, err := pgx.ParseConfig(dsn)
	require.NoError(t, err)
	conn, err := pgx.ConnectConfig(ctx, config)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	userID := uuid.New()
	createdAt := time.Date(2026, 9, 23, 8, 0, 0, 0, time.UTC)
	_, err = conn.Exec(ctx,
		`INSERT INTO users (id, created_at, updated_at) VALUES ($1, $2, $2)`,
		userID, createdAt,
	)
	require.NoError(t, err)

	identity, err := valueobjects.IDFromUUID(userID)
	require.NoError(t, err)

	user, err := readers.NewUserReader(gormDB).GetByID(ctx, identity)
	require.NoError(t, err)
	require.Equal(t, identity, user.ID())
	require.Equal(t, createdAt, user.CreatedAt())
	require.Equal(t, createdAt, user.UpdatedAt())
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

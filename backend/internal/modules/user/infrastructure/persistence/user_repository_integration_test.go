//go:build integration

package persistence_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"developers-cv/backend/internal/modules/user/application/commands"
	"developers-cv/backend/internal/modules/user/infrastructure/persistence/repository"
	"developers-cv/backend/internal/platform/database"
)

func TestPostgresUserRepositoryCreatesUserInsideUnitOfWork(t *testing.T) {
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
	applyUserRepositoryMigration(t, ctx, dsn)

	db, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(db)) })

	repo := repository.NewUserRepository(db)
	handler := commands.NewCreateUserHandler(repo)
	uow := database.NewUnitOfWork(db)

	var userID string
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		id, err := handler.Handle(tx)
		if err != nil {
			return err
		}
		userID = id.UUID().String()
		return nil
	}))

	var count int64
	require.NoError(t, db.Table("users").Where("id = ?", userID).Count(&count).Error)
	require.Equal(t, int64(1), count)
}

func TestPostgresUserRepositoryRollsBackWithUnitOfWork(t *testing.T) {
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
	applyUserRepositoryMigration(t, ctx, dsn)

	db, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(db)) })

	repo := repository.NewUserRepository(db)
	handler := commands.NewCreateUserHandler(repo)
	uow := database.NewUnitOfWork(db)

	expectedErr := context.Canceled
	err = uow.Do(ctx, func(tx context.Context) error {
		_, err := handler.Handle(tx)
		require.NoError(t, err)
		return expectedErr
	})
	require.ErrorIs(t, err, expectedErr)

	var count int64
	require.NoError(t, db.Table("users").Count(&count).Error)
	require.Equal(t, int64(0), count)
}

func applyUserRepositoryMigration(t *testing.T, ctx context.Context, dsn string) {
	t.Helper()
	config, err := pgx.ParseConfig(dsn)
	require.NoError(t, err)
	conn, err := pgx.ConnectConfig(ctx, config)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })
	_, err = conn.Exec(ctx, `CREATE TABLE users (id uuid PRIMARY KEY, created_at timestamptz NOT NULL, updated_at timestamptz NOT NULL)`)
	require.NoError(t, err)
}

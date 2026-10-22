//go:build integration

package persistence_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/user/application/queries"
	"developers-cv/backend/internal/modules/user/infrastructure/persistence/readers"
	"developers-cv/backend/internal/platform/database"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
	"developers-cv/backend/internal/testutil/pgtest"
)

func TestPostgresUserReaderGetUser(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx, dsn := pgtest.NewDatabase(t)

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
		`INSERT INTO users.users (id, created_at, updated_at) VALUES ($1, $2, $2)`,
		userID, createdAt,
	)
	require.NoError(t, err)

	identity, err := sharedvalueobjects.IDFromUUID(userID)
	require.NoError(t, err)

	user, err := readers.NewUserReader(gormDB).GetUser(ctx, queries.GetUserQuery{ID: identity})
	require.NoError(t, err)
	require.Equal(t, identity, user.ID())
	require.True(t, createdAt.Equal(user.CreatedAt()))
	require.True(t, createdAt.Equal(user.UpdatedAt()))
}

func TestPostgresUserReaderReturnsNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx, dsn := pgtest.NewDatabase(t)
	gormDB, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(gormDB)) })

	identity := sharedvalueobjects.NewID()
	_, err = readers.NewUserReader(gormDB).GetUser(ctx, queries.GetUserQuery{ID: identity})
	require.ErrorIs(t, err, queries.ErrUserNotFound)
}

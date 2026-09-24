//go:build integration

package persistence_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"developers-cv/backend/internal/modules/user/application/queries"
	"developers-cv/backend/internal/modules/user/domain/valueobjects"
	"developers-cv/backend/internal/modules/user/infrastructure/persistence/reader"
)

func TestPostgresGetUserReaderReturnsUserProfile(t *testing.T) {
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
	applyUserProfileMigrations(t, ctx, dsn)

	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { pool.Close() })

	firstID := uuid.New()
	secondID := uuid.New()
	insertUser(t, ctx, dsn, firstID)
	insertUser(t, ctx, dsn, secondID)
	insertProfile(t, ctx, dsn, firstID, "Антон", "Седлецкий", "Frontend", "Новосибирск")
	insertProfile(t, ctx, dsn, secondID, "Second", "User", "Backend", "Москва")

	getUserReader := reader.NewGetUserReader(pool)
	firstIdentity, err := valueobjects.IDFromUUID(firstID)
	require.NoError(t, err)
	secondIdentity, err := valueobjects.IDFromUUID(secondID)
	require.NoError(t, err)

	firstUser, err := getUserReader.GetUser(ctx, queries.GetUserQuery{UserID: firstIdentity})
	require.NoError(t, err)
	require.Equal(t, "Антон", firstUser.Name)
	require.Equal(t, "Frontend", *firstUser.Position)

	secondUser, err := getUserReader.GetUser(ctx, queries.GetUserQuery{UserID: secondIdentity})
	require.NoError(t, err)
	require.Equal(t, "Second", secondUser.Name)
	require.Equal(t, "Backend", *secondUser.Position)
}

func TestPostgresGetUserReaderReturnsNotFoundForMissingProfile(t *testing.T) {
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
	applyUserProfileMigrations(t, ctx, dsn)

	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { pool.Close() })

	userID := uuid.New()
	insertUser(t, ctx, dsn, userID)
	identity, err := valueobjects.IDFromUUID(userID)
	require.NoError(t, err)

	_, err = reader.NewGetUserReader(pool).GetUser(ctx, queries.GetUserQuery{UserID: identity})
	require.ErrorIs(t, err, queries.ErrUserNotFound)
}

func applyUserProfileMigrations(t *testing.T, ctx context.Context, dsn string) {
	t.Helper()
	config, err := pgx.ParseConfig(dsn)
	require.NoError(t, err)
	conn, err := pgx.ConnectConfig(ctx, config)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	for _, migration := range []string{"000008_users.up.sql", "000011_user_profiles.up.sql"} {
		content, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "..", "db", "migrations", migration))
		require.NoError(t, err)
		_, err = conn.Exec(ctx, string(content))
		require.NoError(t, err)
	}
}

func insertUser(t *testing.T, ctx context.Context, dsn string, id uuid.UUID) {
	t.Helper()
	config, err := pgx.ParseConfig(dsn)
	require.NoError(t, err)
	conn, err := pgx.ConnectConfig(ctx, config)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })
	_, err = conn.Exec(ctx, `INSERT INTO users (id) VALUES ($1)`, id)
	require.NoError(t, err)
}

func insertProfile(t *testing.T, ctx context.Context, dsn string, userID uuid.UUID, name, surname, position, location string) {
	t.Helper()
	config, err := pgx.ParseConfig(dsn)
	require.NoError(t, err)
	conn, err := pgx.ConnectConfig(ctx, config)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })
	_, err = conn.Exec(ctx, `
		INSERT INTO user_profiles (user_id, name, surname, position, location)
		VALUES ($1, $2, $3, $4, $5)
	`, userID, name, surname, position, location)
	require.NoError(t, err)
}

//go:build integration

package persistence_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"developers-cv/backend/internal/modules/user/domain"
	"developers-cv/backend/internal/modules/user/domain/valueobjects"
	"developers-cv/backend/internal/modules/user/infrastructure/persistence/repository"
	"developers-cv/backend/internal/platform/database"
)

func TestPostgresProfileRepositoryReplacesLocalizedProfile(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx := context.Background()
	container, err := testpostgres.Run(ctx, "postgres:17-alpine", testpostgres.WithDatabase("developers_cv"), testpostgres.WithUsername("developers_cv"), testpostgres.WithPassword("test-password"), testpostgres.BasicWaitStrategies())
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, container.Terminate(ctx)) })
	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	applyProfileRepositoryMigrations(t, ctx, dsn)

	db, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(db)) })

	userID := uuid.New()
	insertProfileRepositoryUser(t, ctx, dsn, userID)
	identity, err := valueobjects.IDFromUUID(userID)
	require.NoError(t, err)
	name, enName, position := "Дмитрий", "Dmitry", "Backend Developer"
	profile, err := domain.NewProfile(identity, domain.LocalizedProfile{Name: &name}, domain.LocalizedProfile{Name: &enName, Position: &position}, nil)
	require.NoError(t, err)

	uow := database.NewUnitOfWork(db)
	repo := repository.NewProfileRepository(db)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error { return repo.ReplaceProfile(tx, profile) }))

	var count int64
	require.NoError(t, db.Table("user_profile_localizations").Where("user_id = ?", userID).Count(&count).Error)
	require.Equal(t, int64(2), count)

	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return repo.ReplaceProfile(tx, &domain.Profile{UserID: identity})
	}))
	require.NoError(t, db.Table("user_profile_localizations").Where("user_id = ?", userID).Count(&count).Error)
	require.Equal(t, int64(0), count)
}

func applyProfileRepositoryMigrations(t *testing.T, ctx context.Context, dsn string) {
	t.Helper()
	config, err := pgx.ParseConfig(dsn)
	require.NoError(t, err)
	conn, err := pgx.ConnectConfig(ctx, config)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })
	for _, migration := range []string{"000001_users.up.sql", "000002_user_profiles.up.sql", "000011_user_profile_localizations.up.sql"} {
		content, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "..", "db", "migrations", migration))
		require.NoError(t, err)
		_, err = conn.Exec(ctx, string(content))
		require.NoError(t, err)
	}
}

func insertProfileRepositoryUser(t *testing.T, ctx context.Context, dsn string, userID uuid.UUID) {
	t.Helper()
	config, err := pgx.ParseConfig(dsn)
	require.NoError(t, err)
	conn, err := pgx.ConnectConfig(ctx, config)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })
	_, err = conn.Exec(ctx, `INSERT INTO users (id) VALUES ($1)`, userID)
	require.NoError(t, err)
}

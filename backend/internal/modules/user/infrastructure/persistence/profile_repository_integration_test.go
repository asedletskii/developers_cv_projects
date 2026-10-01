//go:build integration

package persistence_test

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"path/filepath"

	"developers-cv/backend/internal/modules/user/domain"
	"developers-cv/backend/internal/modules/user/infrastructure/persistence/repository"
	"developers-cv/backend/internal/platform/database"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
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
	identity, err := sharedvalueobjects.IDFromUUID(userID)
	require.NoError(t, err)
	name, enName, position := "Дмитрий", "Dmitry", "Backend Developer"
	ru, err := domain.NewLocalizedProfile(&name, nil, nil, nil, nil)
	require.NoError(t, err)
	en, err := domain.NewLocalizedProfile(&enName, nil, nil, &position, nil)
	require.NoError(t, err)
	photoURL := "https://example.com/profiles/dmitry.svg"
	profile, err := domain.NewProfile(identity, ru, en, &photoURL)
	require.NoError(t, err)

	uow := database.NewUnitOfWork(db)
	repo := repository.NewProfileRepository(db)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error { return repo.ReplaceProfile(tx, profile) }))

	var count int64
	require.NoError(t, db.Table("user_profile_localizations").Where("user_id = ?", userID).Count(&count).Error)
	require.Equal(t, int64(2), count)

	var storedPhotoURL string
	require.NoError(t, db.Table("user_profiles").Select("photo_path").Where("user_id = ?", userID).Scan(&storedPhotoURL).Error)
	require.Equal(t, photoURL, storedPhotoURL)

	emptyRU, err := domain.NewLocalizedProfile(nil, nil, nil, nil, nil)
	require.NoError(t, err)
	emptyEN, err := domain.NewLocalizedProfile(nil, nil, nil, nil, nil)
	require.NoError(t, err)
	emptyProfile, err := domain.NewProfile(identity, emptyRU, emptyEN, nil)
	require.NoError(t, err)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return repo.ReplaceProfile(tx, emptyProfile)
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
	for _, migration := range []string{"000001_users.up.sql", "000002_user_profiles.up.sql"} {
		content, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "..", "db", "migrations", migration))
		require.NoError(t, err)
		_, err = conn.Exec(ctx, string(content))
		require.NoError(t, err)
	}

	_, err = conn.Exec(ctx, `
		ALTER TABLE user_profiles DROP CONSTRAINT IF EXISTS user_profiles_photo_path_check;
		UPDATE user_profiles
		SET photo_path = CASE photo_path
			WHEN '/profiles/anton.svg' THEN 'https://example.com/profiles/anton.svg'
			WHEN '/profiles/dmitry.svg' THEN 'https://example.com/profiles/dmitry.svg'
			ELSE photo_path
		END;
		ALTER TABLE user_profiles
		ADD CONSTRAINT user_profiles_photo_path_check
		CHECK (photo_path IS NULL OR photo_path ~ '^https?://[^/?#[:space:]]+([/?#].*)?$');
	`)
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "..", "db", "migrations", "000011_user_profile_localizations.up.sql"))
	require.NoError(t, err)
	_, err = conn.Exec(ctx, string(content))
	require.NoError(t, err)
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

func TestPostgresProfileRepositoryRollsBackReplacementAsOneUnit(t *testing.T) {
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
	identity, err := sharedvalueobjects.IDFromUUID(userID)
	require.NoError(t, err)
	name, position, location := "Dmitry", "Backend", "Paris"
	ru, err := domain.NewLocalizedProfile(&name, nil, nil, &position, &location)
	require.NoError(t, err)
	profile, err := domain.NewProfile(identity, ru, domain.LocalizedProfile{}, stringPtr("https://example.com/profile.jpg"))
	require.NoError(t, err)

	repo := repository.NewProfileRepository(db)
	uow := database.NewUnitOfWork(db)
	expectedErr := context.Canceled
	err = uow.Do(ctx, func(tx context.Context) error {
		require.NoError(t, repo.ReplaceProfile(tx, profile))
		return expectedErr
	})
	require.ErrorIs(t, err, expectedErr)

	var profileCount int64
	require.NoError(t, db.Table("user_profiles").Where("user_id = ?", userID).Count(&profileCount).Error)
	require.Equal(t, int64(0), profileCount)
	var localizationCount int64
	require.NoError(t, db.Table("user_profile_localizations").Where("user_id = ?", userID).Count(&localizationCount).Error)
	require.Equal(t, int64(0), localizationCount)
}

func stringPtr(value string) *string { return &value }

func TestPostgresProfileRepositoryRollbackPreservesExistingProfile(t *testing.T) {
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
	identity, err := sharedvalueobjects.IDFromUUID(userID)
	require.NoError(t, err)
	oldName, oldPhoto := "Old", "https://example.com/old.jpg"
	oldRU, err := domain.NewLocalizedProfile(&oldName, nil, nil, nil, nil)
	require.NoError(t, err)
	oldProfile, err := domain.NewProfile(identity, oldRU, domain.LocalizedProfile{}, &oldPhoto)
	require.NoError(t, err)
	newName, newPhoto := "New", "https://example.com/new.jpg"
	newRU, err := domain.NewLocalizedProfile(&newName, nil, nil, nil, nil)
	require.NoError(t, err)
	newProfile, err := domain.NewProfile(identity, newRU, domain.LocalizedProfile{}, &newPhoto)
	require.NoError(t, err)

	repo := repository.NewProfileRepository(db)
	uow := database.NewUnitOfWork(db)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error { return repo.ReplaceProfile(tx, oldProfile) }))

	expectedErr := context.Canceled
	err = uow.Do(ctx, func(tx context.Context) error {
		require.NoError(t, repo.ReplaceProfile(tx, newProfile))
		return expectedErr
	})
	require.ErrorIs(t, err, expectedErr)

	var storedPhoto string
	require.NoError(t, db.Table("user_profiles").Select("photo_path").Where("user_id = ?", userID).Scan(&storedPhoto).Error)
	require.Equal(t, oldPhoto, storedPhoto)
	var storedName string
	require.NoError(t, db.Table("user_profile_localizations").Select("name").Where("user_id = ? AND language = 'ru'", userID).Scan(&storedName).Error)
	require.Equal(t, oldName, storedName)
}

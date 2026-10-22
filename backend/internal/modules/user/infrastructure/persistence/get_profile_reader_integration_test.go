//go:build integration

package persistence_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/user/application/queries"
	"developers-cv/backend/internal/modules/user/infrastructure/persistence/readers"
	"developers-cv/backend/internal/platform/database"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
	"developers-cv/backend/internal/testutil/pgtest"
)

func TestPostgresGetProfileReaderReturnsUserProfile(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx, dsn := pgtest.NewDatabase(t)

	gormDB, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(gormDB)) })

	firstID := uuid.New()
	secondID := uuid.New()
	insertUser(t, ctx, dsn, firstID)
	insertUser(t, ctx, dsn, secondID)
	insertProfile(t, ctx, dsn, firstID, "Антон", "Седлецкий", "Frontend", "Новосибирск")
	insertProfile(t, ctx, dsn, secondID, "Second", "User", "Backend", "Москва")

	getProfileReader := readers.NewGetProfileReader(gormDB)
	firstIdentity, err := sharedvalueobjects.IDFromUUID(firstID)
	require.NoError(t, err)
	secondIdentity, err := sharedvalueobjects.IDFromUUID(secondID)
	require.NoError(t, err)

	firstProfile, err := getProfileReader.GetProfile(ctx, queries.GetProfileQuery{UserID: firstIdentity})
	require.NoError(t, err)
	require.Equal(t, "Антон", *firstProfile.RU.Name)
	require.Equal(t, "Frontend", *firstProfile.RU.Position)

	secondProfile, err := getProfileReader.GetProfile(ctx, queries.GetProfileQuery{UserID: secondIdentity})
	require.NoError(t, err)
	require.Equal(t, "Second", *secondProfile.RU.Name)
	require.Equal(t, "Backend", *secondProfile.RU.Position)
}

func TestPostgresGetProfileReaderReturnsNotFoundForMissingProfile(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx, dsn := pgtest.NewDatabase(t)

	gormDB, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(gormDB)) })

	userID := uuid.New()
	insertUser(t, ctx, dsn, userID)
	identity, err := sharedvalueobjects.IDFromUUID(userID)
	require.NoError(t, err)

	_, err = readers.NewGetProfileReader(gormDB).GetProfile(ctx, queries.GetProfileQuery{UserID: identity})
	require.ErrorIs(t, err, queries.ErrProfileNotFound)
}

func insertUser(t *testing.T, ctx context.Context, dsn string, id uuid.UUID) {
	t.Helper()
	config, err := pgx.ParseConfig(dsn)
	require.NoError(t, err)
	conn, err := pgx.ConnectConfig(ctx, config)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })
	_, err = conn.Exec(ctx, `INSERT INTO users.users (id) VALUES ($1)`, id)
	require.NoError(t, err)
}

func insertProfile(t *testing.T, ctx context.Context, dsn string, userID uuid.UUID, name, surname, position, location string) {
	t.Helper()
	config, err := pgx.ParseConfig(dsn)
	require.NoError(t, err)
	conn, err := pgx.ConnectConfig(ctx, config)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })
	_, err = conn.Exec(ctx, `INSERT INTO users.user_profiles (user_id) VALUES ($1)`, userID)
	require.NoError(t, err)
	_, err = conn.Exec(ctx, `INSERT INTO users.user_profile_localizations (user_id, language, name, surname, position, location) VALUES ($1, 'ru', $2, $3, $4, $5)`, userID, name, surname, position, location)
	require.NoError(t, err)
}

func TestPostgresGetProfileReaderDoesNotSeedDemoProfiles(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx, dsn := pgtest.NewDatabase(t)

	gormDB, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(gormDB)) })

	_, err = readers.NewGetProfileReader(gormDB).GetProfile(ctx, queries.GetProfileQuery{
		UserID: mustUserID(t, "10000000-0000-0000-0000-000000000001"),
	})
	require.ErrorIs(t, err, queries.ErrProfileNotFound)
}

func mustUserID(t *testing.T, value string) sharedvalueobjects.ID {
	t.Helper()
	id, err := uuid.Parse(value)
	require.NoError(t, err)
	identity, err := sharedvalueobjects.IDFromUUID(id)
	require.NoError(t, err)
	return identity
}

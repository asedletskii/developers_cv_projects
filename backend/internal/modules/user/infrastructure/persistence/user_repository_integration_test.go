//go:build integration

package persistence_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/user/application/commands"
	"developers-cv/backend/internal/modules/user/infrastructure/persistence/repository"
	"developers-cv/backend/internal/platform/database"
	"developers-cv/backend/internal/testutil/pgtest"
)

func TestPostgresUserRepositoryCreatesUserInsideUnitOfWork(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx, dsn := pgtest.NewDatabase(t)

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
	require.NoError(t, db.Table("users.users").Where("id = ?", userID).Count(&count).Error)
	require.Equal(t, int64(1), count)
}

func TestPostgresUserRepositoryRollsBackWithUnitOfWork(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx, dsn := pgtest.NewDatabase(t)

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
	require.NoError(t, db.Table("users.users").Count(&count).Error)
	require.Equal(t, int64(0), count)
}

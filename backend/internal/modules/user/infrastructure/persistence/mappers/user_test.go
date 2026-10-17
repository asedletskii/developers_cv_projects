package mappers

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/user/domain"
	"developers-cv/backend/internal/modules/user/infrastructure/persistence/gormmodel"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestUserToDB(t *testing.T) {
	id := sharedvalueobjects.NewID()
	createdAt := time.Date(2026, 9, 23, 8, 0, 0, 0, time.UTC)
	user, err := domain.RestoreUser(domain.UserSnapshot{ID: id, CreatedAt: createdAt, UpdatedAt: createdAt})
	require.NoError(t, err)

	model, err := UserToDB(user)
	require.NoError(t, err)
	require.Equal(t, id.UUID(), model.ID)
	require.Equal(t, createdAt, model.CreatedAt)
	require.Equal(t, createdAt, model.UpdatedAt)
}

func TestUserToDBRejectsNil(t *testing.T) {
	_, err := UserToDB(nil)
	require.Error(t, err)
}

func TestUserToDomain(t *testing.T) {
	id := uuid.New()
	createdAt := time.Date(2026, 9, 23, 8, 0, 0, 0, time.UTC)
	user, err := UserToDomain(gormmodel.User{ID: id, CreatedAt: createdAt, UpdatedAt: createdAt})
	require.NoError(t, err)
	identity, err := sharedvalueobjects.IDFromUUID(id)
	require.NoError(t, err)
	require.Equal(t, identity, user.ID())
	require.Equal(t, createdAt, user.CreatedAt())
}

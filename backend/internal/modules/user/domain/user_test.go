package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestNewUser(t *testing.T) {
	id, err := sharedvalueobjects.IDFromUUID(uuid.New())
	require.NoError(t, err)

	createdAt := time.Date(2026, 9, 23, 8, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)

	user, err := NewUser(UserParams{
		ID:        id,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	})

	require.NoError(t, err)
	require.Equal(t, id, user.ID())
	require.Equal(t, createdAt, user.CreatedAt())
	require.Equal(t, updatedAt, user.UpdatedAt())
}

func TestNewUserGeneratesTimestamps(t *testing.T) {
	id := sharedvalueobjects.NewID()

	before := time.Now().UTC()
	user, err := NewUser(UserParams{ID: id})
	after := time.Now().UTC()

	require.NoError(t, err)
	require.False(t, user.CreatedAt().Before(before))
	require.False(t, user.CreatedAt().After(after))
	require.Equal(t, user.CreatedAt(), user.UpdatedAt())
}

func TestNewUserRejectsZeroID(t *testing.T) {
	_, err := NewUser(UserParams{})

	require.ErrorIs(t, err, ErrInvalidUser)
}

func TestNewUserRejectsInvalidTimestampOrder(t *testing.T) {
	id := sharedvalueobjects.NewID()
	createdAt := time.Date(2026, 9, 23, 8, 0, 0, 0, time.UTC)

	_, err := NewUser(UserParams{
		ID:        id,
		CreatedAt: createdAt,
		UpdatedAt: createdAt.Add(-time.Second),
	})

	require.ErrorIs(t, err, ErrInvalidUser)
}

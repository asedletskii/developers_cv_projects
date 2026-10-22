package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestNewUser(t *testing.T) {
	createdAt := time.Date(2026, 9, 23, 8, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)

	user, err := NewUser(UserParams{
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	})

	require.NoError(t, err)
	require.False(t, user.ID().IsZero())
	require.Equal(t, createdAt, user.CreatedAt())
	require.Equal(t, updatedAt, user.UpdatedAt())
}

func TestNewUserGeneratesTimestamps(t *testing.T) {
	before := time.Now().UTC()
	user, err := NewUser(UserParams{})
	after := time.Now().UTC()

	require.NoError(t, err)
	require.False(t, user.CreatedAt().Before(before))
	require.False(t, user.CreatedAt().After(after))
	require.Equal(t, user.CreatedAt(), user.UpdatedAt())
}

func TestNewUserGeneratesID(t *testing.T) {
	user, err := NewUser(UserParams{})

	require.NoError(t, err)
	require.False(t, user.ID().IsZero())
}

func TestNewUserRejectsInvalidTimestampOrder(t *testing.T) {
	createdAt := time.Date(2026, 9, 23, 8, 0, 0, 0, time.UTC)

	_, err := NewUser(UserParams{
		CreatedAt: createdAt,
		UpdatedAt: createdAt.Add(-time.Second),
	})

	require.ErrorIs(t, err, ErrInvalidUser)
}

func TestRestoreUserPreservesSnapshotTimestamps(t *testing.T) {
	id := sharedvalueobjects.NewID()
	createdAt := time.Date(2026, 9, 23, 8, 0, 0, 0, time.FixedZone("stored", 2*60*60))
	updatedAt := createdAt.Add(time.Hour)

	user, err := RestoreUser(UserSnapshot{ID: id, CreatedAt: createdAt, UpdatedAt: updatedAt})

	require.NoError(t, err)
	require.Equal(t, createdAt, user.CreatedAt())
	require.Equal(t, updatedAt, user.UpdatedAt())
	require.Equal(t, createdAt.Location(), user.CreatedAt().Location())
}

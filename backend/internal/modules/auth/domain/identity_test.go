package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestNewIdentity(t *testing.T) {
	created := time.Date(2026, 9, 30, 10, 0, 0, 0, time.FixedZone("test", 3*60*60))
	updated := created.Add(time.Hour)
	subject, err := valueobjects.NewProviderSubject("123456")
	require.NoError(t, err)
	identity, err := NewIdentity(IdentityParams{
		UserID: sharedvalueobjects.NewID(), Provider: ProviderTelegram,
		ProviderSubject: subject, CreatedAt: created, UpdatedAt: updated, LastLoginAt: updated,
	})
	require.NoError(t, err)
	require.Equal(t, created.UTC(), identity.CreatedAt())
	require.Equal(t, updated.UTC(), identity.UpdatedAt())
	require.Equal(t, subject, identity.ProviderSubject())
}

func TestNewIdentityGeneratesID(t *testing.T) {
	at := time.Date(2026, 9, 30, 10, 0, 0, 0, time.UTC)
	subject, err := valueobjects.NewProviderSubject("123456")
	require.NoError(t, err)

	identity, err := NewIdentity(IdentityParams{
		UserID: sharedvalueobjects.NewID(), Provider: ProviderTelegram, ProviderSubject: subject,
		CreatedAt: at, UpdatedAt: at, LastLoginAt: at,
	})

	require.NoError(t, err)
	require.False(t, identity.ID().IsZero())
}

func TestNewIdentityRejectsInvalidTimestampOrder(t *testing.T) {
	now := time.Date(2026, 9, 30, 10, 0, 0, 0, time.UTC)
	subject, err := valueobjects.NewProviderSubject("123456")
	require.NoError(t, err)
	_, err = NewIdentity(IdentityParams{
		UserID: sharedvalueobjects.NewID(), Provider: ProviderTelegram,
		ProviderSubject: subject, CreatedAt: now, UpdatedAt: now.Add(-time.Second), LastLoginAt: now,
	})
	require.ErrorIs(t, err, ErrInvalidIdentity)
}

func TestIdentityTouchLogin(t *testing.T) {
	created := time.Date(2026, 9, 30, 10, 0, 0, 0, time.UTC)
	subject, err := valueobjects.NewProviderSubject("123456")
	require.NoError(t, err)
	identity, err := NewIdentity(IdentityParams{
		UserID: sharedvalueobjects.NewID(), Provider: ProviderTelegram,
		ProviderSubject: subject, CreatedAt: created, UpdatedAt: created, LastLoginAt: created,
	})
	require.NoError(t, err)

	loginAt := created.Add(time.Hour)
	require.NoError(t, identity.TouchLogin(loginAt))
	require.Equal(t, loginAt, identity.LastLoginAt())
	require.Equal(t, loginAt, identity.UpdatedAt())
}

func TestIdentityTouchLoginRejectsTimeBeforeCreation(t *testing.T) {
	created := time.Date(2026, 9, 30, 10, 0, 0, 0, time.UTC)
	subject, err := valueobjects.NewProviderSubject("123456")
	require.NoError(t, err)
	identity, err := NewIdentity(IdentityParams{
		UserID: sharedvalueobjects.NewID(), Provider: ProviderTelegram,
		ProviderSubject: subject, CreatedAt: created, UpdatedAt: created, LastLoginAt: created,
	})
	require.NoError(t, err)
	require.ErrorIs(t, identity.TouchLogin(created.Add(-time.Second)), ErrInvalidIdentity)
}

func TestIdentityTouchLoginRejectsMovingBackwards(t *testing.T) {
	created := time.Date(2026, 9, 30, 10, 0, 0, 0, time.UTC)
	subject, err := valueobjects.NewProviderSubject("123456")
	require.NoError(t, err)
	identity, err := NewIdentity(IdentityParams{
		UserID: sharedvalueobjects.NewID(), Provider: ProviderTelegram, ProviderSubject: subject,
		CreatedAt: created, UpdatedAt: created.Add(2 * time.Hour), LastLoginAt: created.Add(time.Hour),
	})
	require.NoError(t, err)

	err = identity.TouchLogin(created.Add(30 * time.Minute))

	require.ErrorIs(t, err, ErrInvalidIdentity)
	require.Equal(t, created.Add(time.Hour), identity.LastLoginAt())
}

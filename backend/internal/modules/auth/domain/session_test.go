package domain

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func newSessionHash(t *testing.T) valueobjects.SessionTokenHash {
	t.Helper()
	hash, err := valueobjects.NewSessionTokenHash(bytes.Repeat([]byte{1}, valueobjects.HashSize))
	require.NoError(t, err)
	return hash
}

func newTestSession(t *testing.T) (*Session, time.Time) {
	t.Helper()
	created := time.Date(2026, 9, 30, 10, 0, 0, 0, time.UTC)
	session, err := NewSession(SessionParams{
		UserID: sharedvalueobjects.NewID(), TokenHash: newSessionHash(t),
		CreatedAt: created, LastSeenAt: created, ExpiresAt: created.Add(30 * 24 * time.Hour),
	})
	require.NoError(t, err)
	return session, created
}

func TestNewSession(t *testing.T) {
	session, created := newTestSession(t)
	require.Equal(t, created, session.CreatedAt())
	require.True(t, session.IsActive(created.Add(time.Hour), 7*24*time.Hour))
}

func TestNewSessionGeneratesIDAndRestoreCopiesOptionalValues(t *testing.T) {
	created := time.Date(2026, 9, 30, 10, 0, 0, 0, time.UTC)
	userAgent := "browser"
	revokedAt := created.Add(time.Minute)
	session, err := NewSession(SessionParams{
		UserID: sharedvalueobjects.NewID(), TokenHash: newSessionHash(t),
		CreatedAt: created, LastSeenAt: created, ExpiresAt: created.Add(time.Hour),
		RevokedAt: &revokedAt, UserAgent: &userAgent,
	})
	require.NoError(t, err)
	require.False(t, session.ID().IsZero())

	snapshot := SessionSnapshot{
		ID: session.ID(), UserID: session.UserID(), TokenHash: session.TokenHash(),
		CreatedAt: session.CreatedAt(), LastSeenAt: session.LastSeenAt(), ExpiresAt: session.ExpiresAt(),
		RevokedAt: session.RevokedAt(), UserAgent: session.UserAgent(),
	}
	restored, err := RestoreSession(snapshot)
	require.NoError(t, err)
	*snapshot.RevokedAt = created.Add(2 * time.Minute)
	*snapshot.UserAgent = "mutated"
	require.Equal(t, revokedAt, *restored.RevokedAt())
	require.Equal(t, "browser", *restored.UserAgent())
}

func TestNewSessionRejectsInvalidTimestampOrder(t *testing.T) {
	created := time.Date(2026, 9, 30, 10, 0, 0, 0, time.UTC)
	_, err := NewSession(SessionParams{
		UserID: sharedvalueobjects.NewID(), TokenHash: newSessionHash(t),
		CreatedAt: created, LastSeenAt: created.Add(-time.Second), ExpiresAt: created.Add(time.Hour),
	})
	require.ErrorIs(t, err, ErrInvalidSession)
}

func TestSessionExpirationAndRevocation(t *testing.T) {
	session, created := newTestSession(t)
	require.False(t, session.IsExpired(created.Add(29*24*time.Hour)))
	require.True(t, session.IsExpired(created.Add(30*24*time.Hour)))
	require.False(t, session.IsRevoked())
	require.NoError(t, session.Revoke(created.Add(time.Hour)))
	require.True(t, session.IsRevoked())
	require.False(t, session.IsActive(created.Add(2*time.Hour), 7*24*time.Hour))
	require.NoError(t, session.Revoke(created.Add(2*time.Hour)))
}

func TestSessionIdleExpiration(t *testing.T) {
	session, created := newTestSession(t)
	require.False(t, session.IsIdleExpired(created.Add(6*24*time.Hour), 7*24*time.Hour))
	require.True(t, session.IsIdleExpired(created.Add(7*24*time.Hour), 7*24*time.Hour))
}

func TestSessionTouchIsMonotonic(t *testing.T) {
	session, created := newTestSession(t)
	require.ErrorIs(t, session.Touch(created.Add(-time.Second)), ErrInvalidSession)
	require.NoError(t, session.Touch(created.Add(time.Hour)))
	require.Equal(t, created.Add(time.Hour), session.LastSeenAt())
}

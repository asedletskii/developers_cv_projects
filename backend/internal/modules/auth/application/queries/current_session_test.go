package queries

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/auth/application"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
)

func TestCurrentSessionRejectsNilSessionFromRepository(t *testing.T) {
	sessions := &sessionRepositoryStub{returnNil: true}
	h := NewCurrentSessionHandler(sessions, &unitOfWorkStub{}, time.Hour, 5*time.Minute)
	hash, err := valueobjects.NewSessionTokenHash(make([]byte, valueobjects.HashSize))
	require.NoError(t, err)

	_, err = h.Handle(context.Background(), hash)

	require.EqualError(t, err, "find active session: repository returned nil session")
}

func TestCurrentSessionReturnsSessionAndTouchesWhenIntervalElapsed(t *testing.T) {
	now := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	sessions := &sessionRepositoryStub{session: newQuerySession(t, now, now.Add(-10*time.Minute))}
	uow := &unitOfWorkStub{}
	h := NewCurrentSessionHandler(sessions, uow, time.Hour, 5*time.Minute)
	h.now = func() time.Time { return now }
	hash, err := valueobjects.NewSessionTokenHash(make([]byte, valueobjects.HashSize))
	require.NoError(t, err)
	view, err := h.Handle(context.Background(), hash)
	require.NoError(t, err)
	require.Equal(t, sessions.session.UserID(), view.UserID)
	require.Equal(t, now, view.LastSeenAt)
	require.Equal(t, 1, sessions.touchCalls)
	require.Equal(t, 1, uow.calls)
}

func TestCurrentSessionDoesNotTouchWithinInterval(t *testing.T) {
	now := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	sessions := &sessionRepositoryStub{session: newQuerySession(t, now, now.Add(-time.Minute))}
	h := NewCurrentSessionHandler(sessions, &unitOfWorkStub{}, time.Hour, 5*time.Minute)
	h.now = func() time.Time { return now }
	hash, err := valueobjects.NewSessionTokenHash(make([]byte, valueobjects.HashSize))
	require.NoError(t, err)
	_, err = h.Handle(context.Background(), hash)
	require.NoError(t, err)
	require.Equal(t, 0, sessions.touchCalls)
}

func TestCurrentSessionRejectsIdleExpiredSession(t *testing.T) {
	now := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	sessions := &sessionRepositoryStub{session: newQuerySession(t, now, now.Add(-2*time.Hour))}
	h := NewCurrentSessionHandler(sessions, &unitOfWorkStub{}, time.Hour, 5*time.Minute)
	h.now = func() time.Time { return now }
	hash, err := valueobjects.NewSessionTokenHash(make([]byte, valueobjects.HashSize))
	require.NoError(t, err)
	_, err = h.Handle(context.Background(), hash)
	require.ErrorIs(t, err, application.ErrSessionIdleExpired)
}

func TestCurrentSessionReturnsValidSessionWhenTouchFails(t *testing.T) {
	now := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	lastSeen := now.Add(-10 * time.Minute)
	touchErr := errors.New("database temporarily unavailable")
	sessions := &sessionRepositoryStub{session: newQuerySession(t, now, lastSeen), touchErr: touchErr}
	log := &warningLoggerStub{}
	h := NewCurrentSessionHandler(sessions, &unitOfWorkStub{}, time.Hour, 5*time.Minute, log)
	h.now = func() time.Time { return now }
	hash, err := valueobjects.NewSessionTokenHash(make([]byte, valueobjects.HashSize))
	require.NoError(t, err)

	view, err := h.Handle(context.Background(), hash)
	require.NoError(t, err)
	require.Equal(t, lastSeen, view.LastSeenAt)
	require.Equal(t, 1, sessions.touchCalls)
	require.Equal(t, 1, log.calls)
	require.Equal(t, "failed to update auth session last_seen_at", log.message)
	require.ErrorIs(t, log.err, touchErr)
}

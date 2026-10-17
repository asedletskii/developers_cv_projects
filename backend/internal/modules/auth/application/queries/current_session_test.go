package queries

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/auth/application"
	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type sessionRepoStub struct {
	session    *domain.Session
	returnNil  bool
	touchCalls int
	touchErr   error
}

func (r *sessionRepoStub) Create(context.Context, *domain.Session) error { return nil }
func (r *sessionRepoStub) FindActiveByTokenHash(context.Context, valueobjects.SessionTokenHash) (*domain.Session, error) {
	if r.returnNil {
		return nil, nil
	}
	if r.session == nil {
		return nil, domain.ErrEntityNotFound
	}
	return r.session, nil
}
func (r *sessionRepoStub) Revoke(context.Context, sharedvalueobjects.ID, time.Time) error { return nil }
func (r *sessionRepoStub) RevokeAllByUserID(context.Context, sharedvalueobjects.ID, time.Time) error {
	return nil
}
func (r *sessionRepoStub) Touch(_ context.Context, _ sharedvalueobjects.ID, at time.Time) error {
	r.touchCalls++
	if r.touchErr != nil {
		return r.touchErr
	}
	return r.session.Touch(at)
}
func (r *sessionRepoStub) DeleteExpiredBatch(context.Context, time.Time, int) (int, error) {
	return 0, nil
}

type warningLoggerStub struct {
	message string
	err     error
	calls   int
}

func (l *warningLoggerStub) Warn(message string, err error) {
	l.message = message
	l.err = err
	l.calls++
}

type queryUOW struct{ calls int }

func (u *queryUOW) Do(ctx context.Context, fn func(context.Context) error) error {
	u.calls++
	return fn(ctx)
}

func newQuerySession(t *testing.T, now, lastSeen time.Time) *domain.Session {
	t.Helper()
	hash, err := valueobjects.NewSessionTokenHash(make([]byte, valueobjects.HashSize))
	require.NoError(t, err)
	session, err := domain.NewSession(domain.SessionParams{UserID: sharedvalueobjects.NewID(), TokenHash: hash, CreatedAt: now.Add(-3 * time.Hour), LastSeenAt: lastSeen, ExpiresAt: now.Add(time.Hour)})
	require.NoError(t, err)
	return session
}

func TestCurrentSessionRejectsNilSessionFromRepository(t *testing.T) {
	sessions := &sessionRepoStub{returnNil: true}
	h := NewCurrentSessionHandler(sessions, &queryUOW{}, time.Hour, 5*time.Minute)
	hash, err := valueobjects.NewSessionTokenHash(make([]byte, valueobjects.HashSize))
	require.NoError(t, err)

	_, err = h.Handle(context.Background(), hash)

	require.EqualError(t, err, "find active session: repository returned nil session")
}

func TestCurrentSessionReturnsSessionAndTouchesWhenIntervalElapsed(t *testing.T) {
	now := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	sessions := &sessionRepoStub{session: newQuerySession(t, now, now.Add(-10*time.Minute))}
	uow := &queryUOW{}
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
	sessions := &sessionRepoStub{session: newQuerySession(t, now, now.Add(-time.Minute))}
	h := NewCurrentSessionHandler(sessions, &queryUOW{}, time.Hour, 5*time.Minute)
	h.now = func() time.Time { return now }
	hash, err := valueobjects.NewSessionTokenHash(make([]byte, valueobjects.HashSize))
	require.NoError(t, err)
	_, err = h.Handle(context.Background(), hash)
	require.NoError(t, err)
	require.Equal(t, 0, sessions.touchCalls)
}

func TestCurrentSessionRejectsIdleExpiredSession(t *testing.T) {
	now := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	sessions := &sessionRepoStub{session: newQuerySession(t, now, now.Add(-2*time.Hour))}
	h := NewCurrentSessionHandler(sessions, &queryUOW{}, time.Hour, 5*time.Minute)
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
	sessions := &sessionRepoStub{session: newQuerySession(t, now, lastSeen), touchErr: touchErr}
	log := &warningLoggerStub{}
	h := NewCurrentSessionHandler(sessions, &queryUOW{}, time.Hour, 5*time.Minute, log)
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

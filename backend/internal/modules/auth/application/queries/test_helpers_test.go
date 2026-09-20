package queries

import (
	"context"
	"testing"
	"time"

	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"

	"github.com/stretchr/testify/require"
)

type sessionRepositoryStub struct {
	session    *domain.Session
	returnNil  bool
	touchCalls int
	touchErr   error
}

func (r *sessionRepositoryStub) Create(context.Context, *domain.Session) error { return nil }
func (r *sessionRepositoryStub) FindActiveByTokenHash(context.Context, valueobjects.SessionTokenHash) (*domain.Session, error) {
	if r.returnNil {
		return nil, nil
	}
	if r.session == nil {
		return nil, domain.ErrEntityNotFound
	}
	return r.session, nil
}
func (r *sessionRepositoryStub) Revoke(context.Context, sharedvalueobjects.ID, time.Time) error {
	return nil
}
func (r *sessionRepositoryStub) RevokeAllByUserID(context.Context, sharedvalueobjects.ID, time.Time) error {
	return nil
}
func (r *sessionRepositoryStub) Touch(_ context.Context, _ sharedvalueobjects.ID, at time.Time) error {
	r.touchCalls++
	if r.touchErr != nil {
		return r.touchErr
	}
	return r.session.Touch(at)
}
func (r *sessionRepositoryStub) DeleteExpiredBatch(context.Context, time.Time, int) (int, error) {
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

type unitOfWorkStub struct{ calls int }

func (u *unitOfWorkStub) Do(ctx context.Context, fn func(context.Context) error) error {
	u.calls++
	return fn(ctx)
}

func newQuerySession(t *testing.T, now, lastSeen time.Time) *domain.Session {
	t.Helper()
	hash, err := valueobjects.NewSessionTokenHash(make([]byte, valueobjects.HashSize))
	require.NoError(t, err)
	session, err := domain.NewSession(domain.SessionParams{
		UserID: sharedvalueobjects.NewID(), TokenHash: hash, CreatedAt: now.Add(-3 * time.Hour),
		LastSeenAt: lastSeen, ExpiresAt: now.Add(time.Hour),
	})
	require.NoError(t, err)
	return session
}

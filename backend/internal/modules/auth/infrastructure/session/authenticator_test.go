package session

import (
	"context"
	"errors"
	"testing"
	"time"

	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/auth/application/ports"
	"developers-cv/backend/internal/modules/auth/application/queries"
	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	"developers-cv/backend/internal/platform/http/middleware"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type sessionRepositoryStub struct {
	session  *domain.Session
	err      error
	touchErr error
	hash     valueobjects.SessionTokenHash
}

func (r *sessionRepositoryStub) Create(context.Context, *domain.Session) error { return nil }
func (r *sessionRepositoryStub) FindActiveByTokenHash(_ context.Context, hash valueobjects.SessionTokenHash) (*domain.Session, error) {
	r.hash = hash
	if r.err != nil {
		return nil, r.err
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
	if r.touchErr != nil {
		return r.touchErr
	}
	return r.session.Touch(at)
}
func (r *sessionRepositoryStub) DeleteExpiredBatch(context.Context, time.Time, int) (int, error) {
	return 0, nil
}

type sessionUOWStub struct{ err error }

func (u *sessionUOWStub) Do(ctx context.Context, fn func(context.Context) error) error {
	if u.err != nil {
		return u.err
	}
	return fn(ctx)
}

func newAuthenticatorSession(t *testing.T, now time.Time, lastSeenOffset time.Duration) *domain.Session {
	t.Helper()
	rawToken := []byte("test-session-token")
	hashBytes := sha256.Sum256(rawToken)
	hash, err := valueobjects.NewSessionTokenHash(hashBytes[:])
	require.NoError(t, err)
	session, err := domain.NewSession(domain.SessionParams{
		UserID:     sharedvalueobjects.NewID(),
		TokenHash:  hash,
		CreatedAt:  now.Add(-3 * time.Hour),
		LastSeenAt: now.Add(lastSeenOffset),
		ExpiresAt:  now.Add(time.Hour),
	})
	require.NoError(t, err)
	return session
}

func newTestAuthenticator(repo ports.SessionRepository) *Authenticator {
	query := queries.NewCurrentSessionHandler(repo, &sessionUOWStub{}, time.Hour, 5*time.Minute)
	return NewAuthenticator(query, "__Host-session")
}

func requestWithSessionCookie(token string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:     "__Host-session",
		Value:    token,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	return req
}

func encodedTestToken() string {
	return base64.RawURLEncoding.EncodeToString([]byte("test-session-token"))
}

func TestAuthenticatorReturnsUserIDForValidSession(t *testing.T) {
	now := time.Now().UTC()
	session := newAuthenticatorSession(t, now, -time.Minute)
	repo := &sessionRepositoryStub{session: session}
	authenticator := newTestAuthenticator(repo)

	actual, err := authenticator.Authenticate(context.Background(), requestWithSessionCookie(encodedTestToken()))

	require.NoError(t, err)
	require.Equal(t, session.UserID().UUID(), actual)
	expectedHash := sha256.Sum256([]byte("test-session-token"))
	require.Equal(t, expectedHash[:], repo.hash.Bytes())
}

func TestAuthenticatorRejectsMissingCookie(t *testing.T) {
	authenticator := newTestAuthenticator(&sessionRepositoryStub{})
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	_, err := authenticator.Authenticate(context.Background(), req)

	require.ErrorIs(t, err, middleware.ErrUnauthenticated)
}

func TestAuthenticatorRejectsEmptyCookie(t *testing.T) {
	authenticator := newTestAuthenticator(&sessionRepositoryStub{})

	_, err := authenticator.Authenticate(context.Background(), requestWithSessionCookie("   "))

	require.ErrorIs(t, err, middleware.ErrUnauthenticated)
}

func TestAuthenticatorRejectsUnknownSession(t *testing.T) {
	authenticator := newTestAuthenticator(&sessionRepositoryStub{})

	_, err := authenticator.Authenticate(context.Background(), requestWithSessionCookie(encodedTestToken()))

	require.ErrorIs(t, err, middleware.ErrUnauthenticated)
}

func TestAuthenticatorPassesUnexpectedRepositoryError(t *testing.T) {
	dbErr := errors.New("database unavailable")
	authenticator := newTestAuthenticator(&sessionRepositoryStub{err: dbErr})

	_, err := authenticator.Authenticate(context.Background(), requestWithSessionCookie(encodedTestToken()))

	require.ErrorIs(t, err, dbErr)
}

func TestAuthenticatorRejectsIdleExpiredSession(t *testing.T) {
	now := time.Now().UTC()
	session := newAuthenticatorSession(t, now, -2*time.Hour)
	repo := &sessionRepositoryStub{session: session}
	authenticator := newTestAuthenticator(repo)

	_, err := authenticator.Authenticate(context.Background(), requestWithSessionCookie(encodedTestToken()))

	require.ErrorIs(t, err, middleware.ErrUnauthenticated)
}

func TestAuthenticatorRejectsAbsoluteExpiredSession(t *testing.T) {
	now := time.Now().UTC()
	hash, err := valueobjects.NewSessionTokenHash(make([]byte, valueobjects.HashSize))
	require.NoError(t, err)
	expired, err := domain.NewSession(domain.SessionParams{
		UserID:     sharedvalueobjects.NewID(),
		TokenHash:  hash,
		CreatedAt:  now.Add(-3 * time.Hour),
		LastSeenAt: now.Add(-time.Minute),
		ExpiresAt:  now.Add(-time.Second),
	})
	require.NoError(t, err)
	authenticator := newTestAuthenticator(&sessionRepositoryStub{session: expired})

	_, err = authenticator.Authenticate(context.Background(), requestWithSessionCookie(encodedTestToken()))

	require.ErrorIs(t, err, middleware.ErrUnauthenticated)
}

func TestAuthenticatorUsesConfiguredCookieName(t *testing.T) {
	now := time.Now().UTC()
	session := newAuthenticatorSession(t, now, -time.Minute)
	repo := &sessionRepositoryStub{session: session}
	query := queries.NewCurrentSessionHandler(repo, &sessionUOWStub{}, time.Hour, 5*time.Minute)
	authenticator := NewAuthenticator(query, "dev-session")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "dev-session", Value: encodedTestToken(), Secure: true, HttpOnly: true, SameSite: http.SameSiteLaxMode})

	actual, err := authenticator.Authenticate(context.Background(), req)

	require.NoError(t, err)
	require.Equal(t, session.UserID().UUID(), actual)
}

func TestAuthenticatorDoesNotAcceptNilRequest(t *testing.T) {
	authenticator := newTestAuthenticator(&sessionRepositoryStub{})

	_, err := authenticator.Authenticate(context.Background(), nil)

	require.Error(t, err)
}

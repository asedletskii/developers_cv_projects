package app

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	authports "developers-cv/backend/internal/modules/auth/application/ports"
	authqueries "developers-cv/backend/internal/modules/auth/application/queries"
	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	"developers-cv/backend/internal/modules/auth/infrastructure/session"
	"developers-cv/backend/internal/platform/config"
	"developers-cv/backend/internal/platform/http/middleware"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type runtimeSessionRepositoryStub struct{}

func (runtimeSessionRepositoryStub) Create(context.Context, *domain.Session) error { return nil }
func (runtimeSessionRepositoryStub) FindActiveByTokenHash(context.Context, valueobjects.SessionTokenHash) (*domain.Session, error) {
	return nil, domain.ErrEntityNotFound
}
func (runtimeSessionRepositoryStub) Revoke(context.Context, sharedvalueobjects.ID, time.Time) error {
	return nil
}
func (runtimeSessionRepositoryStub) RevokeAllByUserID(context.Context, sharedvalueobjects.ID, time.Time) error {
	return nil
}
func (runtimeSessionRepositoryStub) Touch(context.Context, sharedvalueobjects.ID, time.Time) error {
	return nil
}
func (runtimeSessionRepositoryStub) DeleteExpiredBatch(context.Context, time.Time, int) (int, error) {
	return 0, nil
}

var _ authports.SessionRepository = runtimeSessionRepositoryStub{}

type runtimeUOWStub struct{}

func (runtimeUOWStub) Do(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }

var _ authports.UnitOfWork = runtimeUOWStub{}

func TestNewRuntimeAuthenticatorTelegramModeIgnoresMockHeaderAndURLUserID(t *testing.T) {
	currentSession := authqueries.NewCurrentSessionHandler(runtimeSessionRepositoryStub{}, runtimeUOWStub{}, time.Hour, time.Minute)
	authenticator, err := newRuntimeAuthenticator(config.AuthConfig{Mode: "telegram"}, currentSession, nil, nil)
	require.NoError(t, err)
	require.IsType(t, &session.Authenticator{}, authenticator)

	request := httptest.NewRequest("GET", "/api/v1/lk?userId="+uuid.NewString(), nil)
	request.Header.Set("X-Mock-User-ID", uuid.NewString())

	_, err = authenticator.Authenticate(context.Background(), request)
	require.ErrorIs(t, err, middleware.ErrUnauthenticated)
}

//go:build devauth

package devauth

import (
	"context"
	"testing"
	"time"

	"net/http/httptest"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/user/application/queries"
	"developers-cv/backend/internal/modules/user/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type fakeUserReader struct {
	users map[uuid.UUID]*domain.User
	err   error
}

func (r *fakeUserReader) GetUser(_ context.Context, query queries.GetUserQuery) (*domain.User, error) {
	if r.err != nil {
		return nil, r.err
	}
	user, ok := r.users[query.ID.UUID()]
	if !ok {
		return nil, queries.ErrUserNotFound
	}
	return user, nil
}

func TestMockAuthenticatorAuthenticatesExistingUser(t *testing.T) {
	userID := uuid.New()
	id, err := sharedvalueobjects.IDFromUUID(userID)
	require.NoError(t, err)
	now := time.Now().UTC()
	user, err := domain.RestoreUser(domain.UserSnapshot{ID: id, CreatedAt: now, UpdatedAt: now})
	require.NoError(t, err)

	authenticator := NewMockAuthenticator(&fakeUserReader{
		users: map[uuid.UUID]*domain.User{userID: user},
	})
	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Set(MockUserIDHeader, userID.String())

	actual, err := authenticator.Authenticate(context.Background(), request)

	require.NoError(t, err)
	require.Equal(t, userID, actual)
}

func TestMockAuthenticatorRejectsMissingUser(t *testing.T) {
	authenticator := NewMockAuthenticator(&fakeUserReader{users: map[uuid.UUID]*domain.User{}})
	request := httptest.NewRequest("GET", "/", nil)

	_, err := authenticator.Authenticate(context.Background(), request)

	require.ErrorIs(t, err, ErrUnauthenticated)
}

func TestMockAuthenticatorRejectsUnknownUser(t *testing.T) {
	authenticator := NewMockAuthenticator(&fakeUserReader{users: map[uuid.UUID]*domain.User{}})
	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Set(MockUserIDHeader, uuid.NewString())

	_, err := authenticator.Authenticate(context.Background(), request)

	require.ErrorIs(t, err, ErrUnauthenticated)
}

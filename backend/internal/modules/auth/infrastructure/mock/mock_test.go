package mock

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/user/application/queries"
	"developers-cv/backend/internal/modules/user/domain"
	"developers-cv/backend/internal/modules/user/domain/valueobjects"
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
		return nil, domain.ErrEntityNotFound
	}
	return user, nil
}

func TestMockAuthenticatorAuthenticatesExistingUser(t *testing.T) {
	userID := uuid.New()
	id, err := valueobjects.IDFromUUID(userID)
	require.NoError(t, err)
	user, err := domain.NewUser(domain.UserParams{ID: id})
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

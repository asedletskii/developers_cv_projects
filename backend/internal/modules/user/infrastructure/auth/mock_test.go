package auth

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/user/domain"
	"developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type fakeUserRepository struct {
	users map[uuid.UUID]*domain.User
	err   error
}

func (r *fakeUserRepository) GetByID(_ context.Context, id valueobjects.ID) (*domain.User, error) {
	if r.err != nil {
		return nil, r.err
	}
	user, ok := r.users[id.UUID()]
	if !ok {
		return nil, domain.ErrEntityNotFound
	}
	return user, nil
}

func (r *fakeUserRepository) Save(_ context.Context, user *domain.User) error {
	r.users[user.ID().UUID()] = user
	return nil
}

func TestMockAuthenticatorAuthenticatesExistingUser(t *testing.T) {
	userID := uuid.New()
	id, err := valueobjects.IDFromUUID(userID)
	require.NoError(t, err)
	user, err := domain.NewUser(domain.UserParams{ID: id})
	require.NoError(t, err)

	authenticator := NewMockAuthenticator(&fakeUserRepository{
		users: map[uuid.UUID]*domain.User{userID: user},
	})
	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Set(MockUserIDHeader, userID.String())

	actual, err := authenticator.Authenticate(context.Background(), request)

	require.NoError(t, err)
	require.Equal(t, userID, actual)
}

func TestMockAuthenticatorRejectsMissingUser(t *testing.T) {
	authenticator := NewMockAuthenticator(&fakeUserRepository{users: map[uuid.UUID]*domain.User{}})
	request := httptest.NewRequest("GET", "/", nil)

	_, err := authenticator.Authenticate(context.Background(), request)

	require.ErrorIs(t, err, ErrUnauthenticated)
}

func TestMockAuthenticatorRejectsUnknownUser(t *testing.T) {
	authenticator := NewMockAuthenticator(&fakeUserRepository{users: map[uuid.UUID]*domain.User{}})
	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Set(MockUserIDHeader, uuid.NewString())

	_, err := authenticator.Authenticate(context.Background(), request)

	require.ErrorIs(t, err, ErrUnauthenticated)
}

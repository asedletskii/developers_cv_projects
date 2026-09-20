package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	usercontracts "developers-cv/backend/internal/modules/user/contracts"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestCreateUserHandlerCreatesUser(t *testing.T) {
	repo := &createUserRepositoryStub{}
	handler := NewCreateUserHandler(repo)

	id, err := handler.Handle(context.Background())

	require.NoError(t, err)
	require.False(t, id.IsZero())
	require.NotNil(t, repo.user)
	require.Equal(t, id, repo.user.ID())
	require.False(t, repo.user.CreatedAt().IsZero())
	require.Equal(t, repo.user.CreatedAt(), repo.user.UpdatedAt())
}

func TestCreateUserHandlerPropagatesRepositoryError(t *testing.T) {
	repoErr := context.Canceled
	repo := &createUserRepositoryStub{err: repoErr}
	handler := NewCreateUserHandler(repo)

	id, err := handler.Handle(context.Background())

	require.ErrorIs(t, err, repoErr)
	require.True(t, id.IsZero())
	require.NotNil(t, repo.user)
}

func TestCreateUserHandlerImplementsAuthProvisionerShape(t *testing.T) {
	var provisioner interface {
		CreateUser(context.Context) (sharedvalueobjects.ID, error)
		EnsureUserExists(context.Context, sharedvalueobjects.ID) error
	} = NewCreateUserHandler(&createUserRepositoryStub{})

	id, err := provisioner.CreateUser(context.Background())

	require.NoError(t, err)
	require.False(t, id.IsZero())
}

func TestCreateUserHandlerEnsureUserExists(t *testing.T) {
	repo := &createUserRepositoryStub{exists: true}
	handler := NewCreateUserHandler(repo)

	require.NoError(t, handler.EnsureUserExists(context.Background(), sharedvalueobjects.NewID()))

	repo.exists = false
	require.ErrorIs(t, handler.EnsureUserExists(context.Background(), sharedvalueobjects.NewID()), usercontracts.ErrUserNotFound)
	require.ErrorIs(t, handler.EnsureUserExists(context.Background(), sharedvalueobjects.ID{}), usercontracts.ErrUserNotFound)
}

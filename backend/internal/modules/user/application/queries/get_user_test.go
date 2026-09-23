package queries

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type fakeGetUserReader struct {
	user  User
	err   error
	query GetUserQuery
}

func (f *fakeGetUserReader) GetUser(_ context.Context, query GetUserQuery) (User, error) {
	f.query = query
	return f.user, f.err
}

func TestGetUserHandlerReturnsUser(t *testing.T) {
	userID := valueobjects.NewID()
	reader := &fakeGetUserReader{user: User{
		ID:       userID.String(),
		Name:     "Антон",
		Surname:  "Седлецкий",
		Position: stringPtr("Frontend"),
	}}
	handler := NewGetUserHandler(reader)

	got, err := handler.Handle(context.Background(), GetUserQuery{UserID: userID})

	require.NoError(t, err)
	require.Equal(t, reader.user, got)
	require.Equal(t, userID, reader.query.UserID)
}

func TestGetUserHandlerPropagatesNotFound(t *testing.T) {
	reader := &fakeGetUserReader{err: ErrUserNotFound}
	handler := NewGetUserHandler(reader)

	_, err := handler.Handle(context.Background(), GetUserQuery{UserID: valueobjects.NewID()})

	require.ErrorIs(t, err, ErrUserNotFound)
}

func TestGetUserHandlerRejectsEmptyUserID(t *testing.T) {
	handler := NewGetUserHandler(&fakeGetUserReader{})

	_, err := handler.Handle(context.Background(), GetUserQuery{})

	require.Error(t, err)
	require.False(t, errors.Is(err, ErrUserNotFound))
}

func stringPtr(value string) *string { return &value }

package queries

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestGetProfileHandlerReturnsProfile(t *testing.T) {
	userID := sharedvalueobjects.NewID()
	reader := &profileReaderStub{profile: ProfileView{
		ID: userID.String(), RU: LocalizedProfileView{Name: stringPtr("Антон"), Surname: stringPtr("Седлецкий"), Position: stringPtr("Frontend")},
	}}
	handler := NewGetProfileHandler(reader)

	got, err := handler.Handle(context.Background(), GetProfileQuery{UserID: userID})

	require.NoError(t, err)
	require.Equal(t, reader.profile, got)
	require.Equal(t, userID, reader.query.UserID)
}

func TestGetProfileHandlerPropagatesNotFound(t *testing.T) {
	reader := &profileReaderStub{err: ErrProfileNotFound}
	handler := NewGetProfileHandler(reader)

	_, err := handler.Handle(context.Background(), GetProfileQuery{UserID: sharedvalueobjects.NewID()})

	require.ErrorIs(t, err, ErrProfileNotFound)
}

func TestGetProfileHandlerRejectsEmptyUserID(t *testing.T) {
	handler := NewGetProfileHandler(&profileReaderStub{})

	_, err := handler.Handle(context.Background(), GetProfileQuery{})

	require.Error(t, err)
	require.False(t, errors.Is(err, ErrProfileNotFound))
}

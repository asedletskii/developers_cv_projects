package queries

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type fakeGetProfileReader struct {
	profile Profile
	err     error
	query   GetProfileQuery
}

func (f *fakeGetProfileReader) GetProfile(_ context.Context, query GetProfileQuery) (Profile, error) {
	f.query = query
	return f.profile, f.err
}

func TestGetProfileHandlerReturnsProfile(t *testing.T) {
	userID := valueobjects.NewID()
	reader := &fakeGetProfileReader{profile: Profile{
		ID:       userID.String(),
		Name:     "Антон",
		Surname:  "Седлецкий",
		Position: stringPtr("Frontend"),
	}}
	handler := NewGetProfileHandler(reader)

	got, err := handler.Handle(context.Background(), GetProfileQuery{UserID: userID})

	require.NoError(t, err)
	require.Equal(t, reader.profile, got)
	require.Equal(t, userID, reader.query.UserID)
}

func TestGetProfileHandlerPropagatesNotFound(t *testing.T) {
	reader := &fakeGetProfileReader{err: ErrProfileNotFound}
	handler := NewGetProfileHandler(reader)

	_, err := handler.Handle(context.Background(), GetProfileQuery{UserID: valueobjects.NewID()})

	require.ErrorIs(t, err, ErrProfileNotFound)
}

func TestGetProfileHandlerRejectsEmptyUserID(t *testing.T) {
	handler := NewGetProfileHandler(&fakeGetProfileReader{})

	_, err := handler.Handle(context.Background(), GetProfileQuery{})

	require.Error(t, err)
	require.False(t, errors.Is(err, ErrProfileNotFound))
}

func stringPtr(value string) *string { return &value }

package queries

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type fakeGetProfileReader struct {
	profile ProfileView
	err     error
	query   GetProfileQuery
}

func (f *fakeGetProfileReader) GetProfile(_ context.Context, query GetProfileQuery) (ProfileView, error) {
	f.query = query
	return f.profile, f.err
}

func TestGetProfileHandlerReturnsProfile(t *testing.T) {
	userID := sharedvalueobjects.NewID()
	reader := &fakeGetProfileReader{profile: ProfileView{
		ID: userID.String(), RU: LocalizedProfileView{Name: stringPtr("Антон"), Surname: stringPtr("Седлецкий"), Position: stringPtr("Frontend")},
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

	_, err := handler.Handle(context.Background(), GetProfileQuery{UserID: sharedvalueobjects.NewID()})

	require.ErrorIs(t, err, ErrProfileNotFound)
}

func TestGetProfileHandlerRejectsEmptyUserID(t *testing.T) {
	handler := NewGetProfileHandler(&fakeGetProfileReader{})

	_, err := handler.Handle(context.Background(), GetProfileQuery{})

	require.Error(t, err)
	require.False(t, errors.Is(err, ErrProfileNotFound))
}

func stringPtr(value string) *string { return &value }

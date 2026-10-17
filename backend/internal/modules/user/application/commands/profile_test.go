package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/user/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type fakeProfileRepository struct {
	profile *domain.Profile
	err     error
}

func (f *fakeProfileRepository) ReplaceProfile(_ context.Context, profile *domain.Profile) error {
	f.profile = profile
	return f.err
}

type fakeUnitOfWork struct{}

func (fakeUnitOfWork) Do(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }

func TestReplaceProfileHandlerPersistsBothLocales(t *testing.T) {
	id := sharedvalueobjects.NewID()
	name := "Дмитрий"
	enName := "Dmitry"
	repo := &fakeProfileRepository{}
	h := NewReplaceProfileHandler(repo, fakeUnitOfWork{})
	err := h.Handle(context.Background(), ReplaceProfileCommand{UserID: id, RU: LocalizedProfileInput{Name: &name}, EN: LocalizedProfileInput{Name: &enName}})
	require.NoError(t, err)
	require.Equal(t, id, repo.profile.UserID())
	require.Equal(t, &name, repo.profile.RU().Name())
	require.Equal(t, &enName, repo.profile.EN().Name())
}

func TestReplaceProfileHandlerAllowsEmptyLocales(t *testing.T) {
	repo := &fakeProfileRepository{}
	h := NewReplaceProfileHandler(repo, fakeUnitOfWork{})
	err := h.Handle(context.Background(), ReplaceProfileCommand{UserID: sharedvalueobjects.NewID()})
	require.NoError(t, err)
	require.NotNil(t, repo.profile)
}

func TestReplaceProfileHandlerRejectsInvalidLocalizedProfile(t *testing.T) {
	repo := &fakeProfileRepository{}
	h := NewReplaceProfileHandler(repo, fakeUnitOfWork{})
	err := h.Handle(context.Background(), ReplaceProfileCommand{
		UserID: sharedvalueobjects.NewID(),
		RU:     LocalizedProfileInput{Name: stringPtr("   ")},
	})
	require.ErrorIs(t, err, domain.ErrInvalidProfile)
	require.Nil(t, repo.profile)
}

func TestReplaceProfileHandlerRejectsInvalidPhotoURL(t *testing.T) {
	repo := &fakeProfileRepository{}
	h := NewReplaceProfileHandler(repo, fakeUnitOfWork{})
	err := h.Handle(context.Background(), ReplaceProfileCommand{
		UserID:   sharedvalueobjects.NewID(),
		PhotoURL: stringPtr("/profiles/dmitry.svg"),
	})
	require.ErrorIs(t, err, domain.ErrInvalidProfile)
	require.Nil(t, repo.profile)
}

func stringPtr(value string) *string { return &value }

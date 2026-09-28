package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/user/domain"
	"developers-cv/backend/internal/modules/user/domain/valueobjects"
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
	id := valueobjects.NewID()
	name := "Дмитрий"
	enName := "Dmitry"
	repo := &fakeProfileRepository{}
	h := NewReplaceProfileHandler(repo, fakeUnitOfWork{})
	err := h.Handle(context.Background(), ReplaceProfileCommand{UserID: id, RU: LocalizedProfileInput{Name: &name}, EN: LocalizedProfileInput{Name: &enName}})
	require.NoError(t, err)
	require.Equal(t, id, repo.profile.UserID)
	require.Equal(t, &name, repo.profile.RU.Name)
	require.Equal(t, &enName, repo.profile.EN.Name)
}

func TestReplaceProfileHandlerAllowsEmptyLocales(t *testing.T) {
	repo := &fakeProfileRepository{}
	h := NewReplaceProfileHandler(repo, fakeUnitOfWork{})
	err := h.Handle(context.Background(), ReplaceProfileCommand{UserID: valueobjects.NewID()})
	require.NoError(t, err)
	require.NotNil(t, repo.profile)
}

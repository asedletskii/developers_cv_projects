package commands

import (
	"context"
	"developers-cv/backend/internal/modules/resume/domain"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/ports"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type fakeTechnologiesAggregateRepository struct {
	resume  *domain.Resume
	saved   bool
	version int64
}

func (r *fakeTechnologiesAggregateRepository) Get(context.Context, ports.GetResumeParams) (*domain.Resume, error) {
	return r.resume, nil
}
func (r *fakeTechnologiesAggregateRepository) Save(_ context.Context, resume *domain.Resume) error {
	r.resume, r.version, r.saved = resume, resume.Version(), true
	return nil
}
func (r *fakeTechnologiesAggregateRepository) Create(context.Context, *domain.Resume) error {
	return nil
}
func (r *fakeTechnologiesAggregateRepository) Delete(context.Context, ports.DeleteResumeParams) error {
	return nil
}

func TestReplaceTechnologiesHandlerUsesAggregateAndVersion(t *testing.T) {
	uid := identity.NewID()
	resume := newTestResumeForCommand(t, uid)
	repo := &fakeTechnologiesAggregateRepository{resume: resume}

	err := NewReplaceTechnologiesHandler(repo, &fakeCommandUnitOfWork{}).Handle(context.Background(), ReplaceTechnologiesCommand{
		ResumeID: resume.ID(), UserID: uid, Version: 1,
		Input: TechnologiesInput{Entries: []TechnologyInput{
			{Name: "Go", Grade: ptr(domain.TechnologyGradeSenior), Position: 0},
			{Name: "Rust", Position: 1},
		}},
	})
	require.NoError(t, err)
	require.True(t, repo.saved)
	require.Equal(t, int64(1), repo.version)
	require.Len(t, repo.resume.Technologies().Entries(), 2)
}

func TestReplaceTechnologiesHandlerRejectsStaleVersion(t *testing.T) {
	uid := identity.NewID()
	resume := newTestResumeForCommand(t, uid)
	repo := &fakeTechnologiesAggregateRepository{resume: resume}

	err := NewReplaceTechnologiesHandler(repo, &fakeCommandUnitOfWork{}).Handle(context.Background(), ReplaceTechnologiesCommand{
		ResumeID: resume.ID(), UserID: uid, Version: 2,
		Input: TechnologiesInput{Entries: []TechnologyInput{{Name: "Go", Position: 0}}},
	})
	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repo.saved)
}

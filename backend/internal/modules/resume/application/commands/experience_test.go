package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type fakeExperiencesAggregateRepository struct {
	resume  *domain.Resume
	saved   bool
	version int64
}

func (r *fakeExperiencesAggregateRepository) Get(context.Context, ports.GetResumeParams) (*domain.Resume, error) {
	return r.resume, nil
}
func (r *fakeExperiencesAggregateRepository) Save(_ context.Context, resume *domain.Resume) error {
	r.resume, r.version, r.saved = resume, resume.Version(), true
	return nil
}
func (r *fakeExperiencesAggregateRepository) Create(context.Context, *domain.Resume) error {
	return nil
}
func (r *fakeExperiencesAggregateRepository) Delete(context.Context, ports.DeleteResumeParams) error {
	return nil
}

func experienceCommandItem(t *testing.T, company string) ExperienceInput {
	t.Helper()
	start, err := domain.ParseYearMonth("2024-01")
	require.NoError(t, err)
	return ExperienceInput{
		Company: company, Position: "Go Developer", StartDate: start,
		Achievements: []string{"Built API"},
		Tasks:        []string{"Implement backend", "Review PRs"},
		Stack:        []ExperienceStackItemInput{{Name: "Go", Grade: ptr(domain.TechnologyGradeSenior)}},
	}
}

func TestReplaceExperiencesHandlerUsesVersionAndAggregate(t *testing.T) {
	uid := sharedvalueobjects.NewID()
	resume := newTestResumeForCommand(t, uid)
	repo := &fakeExperiencesAggregateRepository{resume: resume}

	err := NewReplaceExperiencesHandler(repo, &fakeCommandUnitOfWork{}).Handle(context.Background(), ReplaceExperiencesCommand{
		ResumeID: resume.ID(), UserID: uid, Version: 1,
		Input: ExperiencesInput{Items: []ExperienceInput{experienceCommandItem(t, "Acme")}},
	})
	require.NoError(t, err)
	require.True(t, repo.saved)
	require.Equal(t, int64(1), repo.version)
	require.Len(t, repo.resume.Experiences(), 1)
	require.Equal(t, "Acme", repo.resume.Experiences()[0].Company())
}

func TestReplaceExperiencesHandlerRejectsStaleVersion(t *testing.T) {
	uid := sharedvalueobjects.NewID()
	resume := newTestResumeForCommand(t, uid)
	repo := &fakeExperiencesAggregateRepository{resume: resume}

	err := NewReplaceExperiencesHandler(repo, &fakeCommandUnitOfWork{}).Handle(context.Background(), ReplaceExperiencesCommand{
		ResumeID: resume.ID(), UserID: uid, Version: 2,
		Input: ExperiencesInput{Items: []ExperienceInput{experienceCommandItem(t, "Acme")}},
	})
	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repo.saved)
}

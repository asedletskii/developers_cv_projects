package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type educationRepoStub struct {
	resume  *domain.Resume
	saved   bool
	version int64
}

func (s *educationRepoStub) Get(_ context.Context, _ ports.GetResumeParams) (*domain.Resume, error) {
	return s.resume, nil
}
func (s *educationRepoStub) Save(_ context.Context, r *domain.Resume) error {
	s.saved = true
	s.version = r.Version()
	return nil
}
func (s *educationRepoStub) Create(context.Context, *domain.Resume) error { return nil }
func (s *educationRepoStub) Delete(context.Context, ports.DeleteResumeParams) error {
	return nil
}

type educationUOW struct{}

func (educationUOW) Do(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }
func TestReplaceEducationsHandlerRejectsStaleVersion(t *testing.T) {
	r, err := domain.NewResume(domain.ResumeParams{ID: sharedvalueobjects.NewID(), UserID: sharedvalueobjects.NewID(), Title: "Resume", Language: domain.ResumeLanguageEN, Version: 2})
	require.NoError(t, err)
	repo := &educationRepoStub{resume: r}
	h := NewReplaceEducationsHandler(repo, educationUOW{})
	err = h.Handle(context.Background(), ReplaceEducationsCommand{ResumeID: r.ID(), UserID: r.UserID(), Version: 1, Input: EducationsInput{}})
	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repo.saved)
}
func TestReplaceEducationsHandlerUsesAggregateRepository(t *testing.T) {
	r, err := domain.NewResume(domain.ResumeParams{ID: sharedvalueobjects.NewID(), UserID: sharedvalueobjects.NewID(), Title: "Resume", Language: domain.ResumeLanguageEN})
	require.NoError(t, err)
	repo := &educationRepoStub{resume: r}
	h := NewReplaceEducationsHandler(repo, educationUOW{})
	err = h.Handle(context.Background(), ReplaceEducationsCommand{ResumeID: r.ID(), UserID: r.UserID(), Version: 1, Input: EducationsInput{Items: []EducationInput{{Location: "Berlin", University: "Uni", Grade: domain.EducationLevelMaster, ProgramName: "CS", Faculty: "Engineering"}}}})
	require.NoError(t, err)
	require.True(t, repo.saved)
	require.Len(t, r.Educations(), 1)
}

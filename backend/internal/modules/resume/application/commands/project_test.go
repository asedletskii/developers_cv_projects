package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type projectRepoStub struct {
	resume *domain.Resume
	saved  bool
}

func (s *projectRepoStub) Get(context.Context, ports.GetResumeParams) (*domain.Resume, error) {
	return s.resume, nil
}
func (s *projectRepoStub) Save(context.Context, *domain.Resume) error {
	s.saved = true
	return nil
}
func (s *projectRepoStub) Create(context.Context, *domain.Resume) error           { return nil }
func (s *projectRepoStub) Delete(context.Context, ports.DeleteResumeParams) error { return nil }

type projectUOWStub struct{}

func (projectUOWStub) Do(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }

func TestReplaceProjectsHandlerUsesAggregateRepository(t *testing.T) {
	resume, err := domain.NewResume(domain.ResumeParams{
		ID: identity.NewID(), UserID: identity.NewID(),
		Title: "Resume", Language: domain.ResumeLanguageEN,
	})
	require.NoError(t, err)

	repo := &projectRepoStub{resume: resume}
	err = NewReplaceProjectsHandler(repo, projectUOWStub{}).Handle(context.Background(), ReplaceProjectsCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 1,
		Input: ProjectsInput{Items: []ProjectInput{{Name: "Project"}}},
	})
	require.NoError(t, err)
	require.True(t, repo.saved)
	require.Len(t, resume.Projects(), 1)
}

func TestReplaceProjectsHandlerRejectsStaleVersion(t *testing.T) {
	resume, err := domain.NewResume(domain.ResumeParams{
		ID: identity.NewID(), UserID: identity.NewID(),
		Title: "Resume", Language: domain.ResumeLanguageEN,
	})
	require.NoError(t, err)

	repo := &projectRepoStub{resume: resume}
	err = NewReplaceProjectsHandler(repo, projectUOWStub{}).Handle(context.Background(), ReplaceProjectsCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 2,
		Input: ProjectsInput{Items: []ProjectInput{{Name: "Project"}}},
	})
	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repo.saved)
}

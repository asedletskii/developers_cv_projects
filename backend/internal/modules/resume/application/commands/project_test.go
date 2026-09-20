package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestReplaceProjectsHandlerUsesAggregateRepository(t *testing.T) {
	resume := newCommandTestResume(t, sharedvalueobjects.NewID(), "Resume", 1)

	repo := &commandRepositoryStub{resume: resume}
	err := NewReplaceProjectsHandler(repo, &unitOfWorkStub{}).Handle(context.Background(), ReplaceProjectsCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 1,
		Input: ProjectsInput{Items: []ProjectInput{{Name: "Project"}}},
	})
	require.NoError(t, err)
	require.True(t, repo.saveCalled)
	require.Len(t, resume.Projects(), 1)
}

func TestReplaceProjectsHandlerRejectsStaleVersion(t *testing.T) {
	resume := newCommandTestResume(t, sharedvalueobjects.NewID(), "Resume", 1)

	repo := &commandRepositoryStub{resume: resume}
	err := NewReplaceProjectsHandler(repo, &unitOfWorkStub{}).Handle(context.Background(), ReplaceProjectsCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 2,
		Input: ProjectsInput{Items: []ProjectInput{{Name: "Project"}}},
	})
	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repo.saveCalled)
}

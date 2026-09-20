package commands

import (
	"context"
	"developers-cv/backend/internal/modules/resume/domain"
	"testing"

	"github.com/stretchr/testify/require"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestReplaceTechnologiesHandlerUsesAggregateAndVersion(t *testing.T) {
	uid := sharedvalueobjects.NewID()
	resume := newCommandTestResume(t, uid, "Resume", 1)
	repo := &commandRepositoryStub{resume: resume}

	err := NewReplaceTechnologiesHandler(repo, &unitOfWorkStub{}).Handle(context.Background(), ReplaceTechnologiesCommand{
		ResumeID: resume.ID(), UserID: uid, Version: 1,
		Input: TechnologiesInput{Entries: []TechnologyInput{
			{Name: "Go", Position: 0},
			{Name: "Rust", Position: 1},
		}},
	})
	require.NoError(t, err)
	require.True(t, repo.saveCalled)
	require.Equal(t, int64(1), repo.version)
	require.Len(t, repo.resume.Technologies().Entries(), 2)
}

func TestReplaceTechnologiesHandlerRejectsStaleVersion(t *testing.T) {
	uid := sharedvalueobjects.NewID()
	resume := newCommandTestResume(t, uid, "Resume", 1)
	repo := &commandRepositoryStub{resume: resume}

	err := NewReplaceTechnologiesHandler(repo, &unitOfWorkStub{}).Handle(context.Background(), ReplaceTechnologiesCommand{
		ResumeID: resume.ID(), UserID: uid, Version: 2,
		Input: TechnologiesInput{Entries: []TechnologyInput{{Name: "Go", Position: 0}}},
	})
	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repo.saveCalled)
}

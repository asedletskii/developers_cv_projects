package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestReplaceActivitiesHandlerUsesAggregateRepository(t *testing.T) {
	resume, err := domain.NewResume(domain.ResumeParams{UserID: sharedvalueobjects.NewID(), Title: "Resume", Language: domain.ResumeLanguageEN})
	require.NoError(t, err)
	repo := &commandRepositoryStub{resume: resume}
	err = NewReplaceActivitiesHandler(repo, &unitOfWorkStub{}).Handle(context.Background(), ReplaceActivitiesCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 1,
		Input: ActivitiesInput{Items: []ActivityInput{{Type: domain.ActivityTypeOpenSource, Title: "OSS"}}},
	})
	require.NoError(t, err)
	require.True(t, repo.saveCalled)
	require.Len(t, resume.Activities(), 1)
}

func TestReplaceActivitiesHandlerRejectsStaleVersion(t *testing.T) {
	resume, err := domain.NewResume(domain.ResumeParams{UserID: sharedvalueobjects.NewID(), Title: "Resume", Language: domain.ResumeLanguageEN})
	require.NoError(t, err)
	repo := &commandRepositoryStub{resume: resume}
	err = NewReplaceActivitiesHandler(repo, &unitOfWorkStub{}).Handle(context.Background(), ReplaceActivitiesCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 2,
		Input: ActivitiesInput{Items: []ActivityInput{{Type: domain.ActivityTypeOpenSource, Title: "OSS"}}},
	})
	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repo.saveCalled)
}

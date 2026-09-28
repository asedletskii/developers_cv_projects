package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type activityRepoStub struct {
	resume *domain.Resume
	saved  bool
}

func (s *activityRepoStub) Get(context.Context, ports.GetResumeParams) (*domain.Resume, error) {
	return s.resume, nil
}
func (s *activityRepoStub) Save(context.Context, *domain.Resume) error {
	s.saved = true
	return nil
}
func (s *activityRepoStub) Create(context.Context, *domain.Resume) error { return nil }
func (s *activityRepoStub) Delete(context.Context, ports.DeleteResumeParams) error {
	return nil
}

type activityUOWStub struct{}

func (activityUOWStub) Do(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }

func TestReplaceActivitiesHandlerUsesAggregateRepository(t *testing.T) {
	resume, err := domain.NewResume(domain.ResumeParams{ID: valueobjects.NewID(), UserID: useridentity.NewID(), Title: "Resume", Language: enums.ResumeLanguageEN})
	require.NoError(t, err)
	repo := &activityRepoStub{resume: resume}
	err = NewReplaceActivitiesHandler(repo, activityUOWStub{}).Handle(context.Background(), ReplaceActivitiesCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 1,
		Input: ActivitiesInput{Items: []ActivityInput{{Type: enums.ActivityTypeOpenSource, Title: "OSS"}}},
	})
	require.NoError(t, err)
	require.True(t, repo.saved)
	require.Len(t, resume.Activities(), 1)
}

func TestReplaceActivitiesHandlerRejectsStaleVersion(t *testing.T) {
	resume, err := domain.NewResume(domain.ResumeParams{ID: valueobjects.NewID(), UserID: useridentity.NewID(), Title: "Resume", Language: enums.ResumeLanguageEN})
	require.NoError(t, err)
	repo := &activityRepoStub{resume: resume}
	err = NewReplaceActivitiesHandler(repo, activityUOWStub{}).Handle(context.Background(), ReplaceActivitiesCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 2,
		Input: ActivitiesInput{Items: []ActivityInput{{Type: enums.ActivityTypeOpenSource, Title: "OSS"}}},
	})
	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repo.saved)
}

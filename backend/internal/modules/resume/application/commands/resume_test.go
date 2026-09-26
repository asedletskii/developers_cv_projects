package commands

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type fakeCreateResumeRepository struct {
	called bool
	resume *domain.Resume
}

func (r *fakeCreateResumeRepository) CreateResume(_ context.Context, cmd ports.CreateResumeRepositoryCommand) error {
	r.called = true
	r.resume = cmd.Resume
	return nil
}

type fakeCreateOverviewRepository struct {
	called   bool
	resumeID valueobjects.ID
	userID   useridentity.ID
}

func (r *fakeCreateOverviewRepository) CreateOverview(_ context.Context, cmd ports.CreateOverviewRepositoryCommand) error {
	r.called = true
	r.resumeID = cmd.ResumeID
	r.userID = cmd.UserID
	return nil
}

func TestCreateResumeHandlerCreatesResumeWithTitleAndLanguage(t *testing.T) {
	resumeRepository := &fakeCreateResumeRepository{}
	overviewRepository := &fakeCreateOverviewRepository{}
	uow := &fakeResumeUnitOfWork{}
	handler := NewCreateResumeHandler(resumeRepository, overviewRepository, uow)
	userID := useridentity.NewID()

	resumeID, err := handler.Handle(context.Background(), CreateResumeCommand{
		UserID:   userID,
		Title:    "  Senior Go Developer  ",
		Language: "en",
	})

	require.NoError(t, err)
	require.False(t, resumeID.IsZero())
	require.True(t, uow.called)
	require.True(t, resumeRepository.called)
	require.True(t, overviewRepository.called)
	require.Equal(t, resumeID, resumeRepository.resume.ID())
	require.Equal(t, userID, resumeRepository.resume.UserID())
	require.Equal(t, "Senior Go Developer", resumeRepository.resume.Title())
	require.Equal(t, "en", string(resumeRepository.resume.Language()))
	require.Equal(t, resumeID, overviewRepository.resumeID)
	require.Equal(t, userID, overviewRepository.userID)
}

type fakeReplaceResumeTitleRepository struct {
	called bool
	cmd    ports.ReplaceResumeTitleRepositoryCommand
	err    error
}

func (r *fakeReplaceResumeTitleRepository) ReplaceResumeTitle(_ context.Context, cmd ports.ReplaceResumeTitleRepositoryCommand) error {
	r.called = true
	r.cmd = cmd
	return r.err
}

type fakeResumeUnitOfWork struct {
	called bool
}

func (u *fakeResumeUnitOfWork) Do(ctx context.Context, fn func(context.Context) error) error {
	u.called = true
	return fn(ctx)
}

func TestReplaceResumeTitleHandlerReplacesOnlyTitle(t *testing.T) {
	repository := &fakeReplaceResumeTitleRepository{}
	uow := &fakeResumeUnitOfWork{}
	handler := NewReplaceResumeTitleHandler(repository, uow)

	resumeID := valueobjects.NewID()
	userID := useridentity.NewID()
	err := handler.Handle(context.Background(), ReplaceResumeTitleCommand{
		ResumeID: resumeID,
		UserID:   userID,
		Title:    "Senior Go Developer",
	})

	require.NoError(t, err)
	require.True(t, uow.called)
	require.True(t, repository.called)
	require.Equal(t, resumeID, repository.cmd.ResumeID)
	require.Equal(t, userID, repository.cmd.UserID)
	require.Equal(t, "Senior Go Developer", repository.cmd.Title)
}

func TestReplaceResumeTitleHandlerTrimsTitle(t *testing.T) {
	repository := &fakeReplaceResumeTitleRepository{}
	uow := &fakeResumeUnitOfWork{}
	handler := NewReplaceResumeTitleHandler(repository, uow)

	err := handler.Handle(context.Background(), ReplaceResumeTitleCommand{
		ResumeID: valueobjects.NewID(),
		UserID:   useridentity.NewID(),
		Title:    "  Senior Go Developer  ",
	})

	require.NoError(t, err)
	require.Equal(t, "Senior Go Developer", repository.cmd.Title)
}

func TestReplaceResumeTitleHandlerRejectsInvalidIdentity(t *testing.T) {
	repository := &fakeReplaceResumeTitleRepository{}
	uow := &fakeResumeUnitOfWork{}
	handler := NewReplaceResumeTitleHandler(repository, uow)
	userID := useridentity.NewID()

	tests := []ReplaceResumeTitleCommand{
		{UserID: userID, Title: "Senior Go Developer"},
		{ResumeID: valueobjects.NewID(), Title: "Senior Go Developer"},
	}

	for _, command := range tests {
		err := handler.Handle(context.Background(), command)
		require.Error(t, err)
		require.False(t, repository.called)
	}
}

func TestReplaceResumeTitleHandlerRejectsInvalidTitle(t *testing.T) {
	repository := &fakeReplaceResumeTitleRepository{}
	uow := &fakeResumeUnitOfWork{}
	handler := NewReplaceResumeTitleHandler(repository, uow)

	err := handler.Handle(context.Background(), ReplaceResumeTitleCommand{
		ResumeID: valueobjects.NewID(),
		UserID:   useridentity.NewID(),
		Title:    "   ",
	})

	require.Error(t, err)
	require.ErrorIs(t, err, domain.ErrInvalidResume)
	require.False(t, repository.called)
}

func TestReplaceResumeTitleHandlerPropagatesRepositoryError(t *testing.T) {
	repositoryErr := errors.New("repository failure")
	repository := &fakeReplaceResumeTitleRepository{err: repositoryErr}
	uow := &fakeResumeUnitOfWork{}
	handler := NewReplaceResumeTitleHandler(repository, uow)

	err := handler.Handle(context.Background(), ReplaceResumeTitleCommand{
		ResumeID: valueobjects.NewID(),
		UserID:   useridentity.NewID(),
		Title:    "Senior Go Developer",
	})

	require.ErrorIs(t, err, repositoryErr)
}

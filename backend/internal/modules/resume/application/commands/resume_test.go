package commands

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type fakeResumeRepository struct {
	called   bool
	resume   *domain.Resume
	saved    *domain.Resume
	expected int64
	err      error
}

func (r *fakeResumeRepository) Get(_ context.Context, _ ports.GetResumeParams) (*domain.Resume, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.resume == nil {
		return nil, errors.New("resume is nil")
	}
	return r.resume, nil
}

func (r *fakeResumeRepository) Save(_ context.Context, resume *domain.Resume) error {
	r.called = true
	r.saved = resume
	r.expected = resume.Version()
	return r.err
}

func (r *fakeResumeRepository) Create(_ context.Context, resume *domain.Resume) error {
	r.called = true
	r.resume = resume
	return r.err
}

func (r *fakeResumeRepository) Delete(_ context.Context, _ ports.DeleteResumeParams) error {
	r.called = true
	return r.err
}

func TestCreateResumeHandlerCreatesResumeWithOverview(t *testing.T) {
	repository := &fakeResumeRepository{}
	uow := &fakeCommandUnitOfWork{}
	handler := NewCreateResumeHandler(repository, uow)
	userID := useridentity.NewID()

	resumeID, err := handler.Handle(context.Background(), CreateResumeCommand{UserID: userID, Title: "  Senior Go Developer  ", Language: enums.ResumeLanguage("en")})
	require.NoError(t, err)
	require.False(t, resumeID.IsZero())
	require.True(t, uow.called)
	require.True(t, repository.called)
	require.Equal(t, resumeID, repository.resume.ID())
	require.Equal(t, userID, repository.resume.UserID())
	require.Equal(t, "Senior Go Developer", repository.resume.Title())
	require.Equal(t, "en", string(repository.resume.Language()))
	require.NotNil(t, repository.resume.Overview())
}

func newTestResume(t *testing.T, title string, version int64) *domain.Resume {
	t.Helper()
	resume, err := domain.NewResume(domain.ResumeParams{
		ID: valueobjects.NewID(), UserID: useridentity.NewID(), Title: title,
		Language: enums.ResumeLanguage("en"), Version: version,
	})
	require.NoError(t, err)
	return resume
}

func TestReplaceResumeTitleHandlerReplacesTitleThroughAggregate(t *testing.T) {
	resume := newTestResume(t, "Old title", 4)
	repository := &fakeResumeRepository{resume: resume}
	uow := &fakeCommandUnitOfWork{}
	handler := NewReplaceResumeTitleHandler(repository, uow)

	err := handler.Handle(context.Background(), ReplaceResumeTitleCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 4, Title: "  Senior Go Developer  ",
	})

	require.NoError(t, err)
	require.True(t, uow.called)
	require.True(t, repository.called)
	require.Same(t, resume, repository.saved)
	require.Equal(t, int64(4), repository.expected)
	require.Equal(t, "Senior Go Developer", resume.Title())
}

func TestReplaceResumeTitleHandlerRejectsStaleVersion(t *testing.T) {
	resume := newTestResume(t, "Old title", 5)
	repository := &fakeResumeRepository{resume: resume}
	handler := NewReplaceResumeTitleHandler(repository, &fakeCommandUnitOfWork{})

	err := handler.Handle(context.Background(), ReplaceResumeTitleCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 4, Title: "New title",
	})

	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repository.called)
	require.Equal(t, "Old title", resume.Title())
}

func TestReplaceResumeTitleHandlerRejectsInvalidIdentityOrVersion(t *testing.T) {
	resume := newTestResume(t, "Old title", 1)
	repository := &fakeResumeRepository{resume: resume}
	handler := NewReplaceResumeTitleHandler(repository, &fakeCommandUnitOfWork{})

	tests := []ReplaceResumeTitleCommand{
		{UserID: resume.UserID(), Version: 1, Title: "New"},
		{ResumeID: resume.ID(), Version: 1, Title: "New"},
		{ResumeID: resume.ID(), UserID: resume.UserID(), Title: "New"},
	}

	for _, command := range tests {
		err := handler.Handle(context.Background(), command)
		require.Error(t, err)
		require.False(t, repository.called)
	}
}

func TestReplaceResumeTitleHandlerRejectsInvalidTitle(t *testing.T) {
	resume := newTestResume(t, "Old title", 1)
	repository := &fakeResumeRepository{resume: resume}
	handler := NewReplaceResumeTitleHandler(repository, &fakeCommandUnitOfWork{})

	err := handler.Handle(context.Background(), ReplaceResumeTitleCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 1, Title: "   ",
	})

	require.ErrorIs(t, err, domain.ErrInvalidResume)
	require.False(t, repository.called)
}

func TestReplaceResumeTitleHandlerPropagatesRepositoryError(t *testing.T) {
	repositoryErr := errors.New("repository failure")
	resume := newTestResume(t, "Old title", 1)
	repository := &fakeResumeRepository{resume: resume, err: repositoryErr}
	handler := NewReplaceResumeTitleHandler(repository, &fakeCommandUnitOfWork{})

	err := handler.Handle(context.Background(), ReplaceResumeTitleCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 1, Title: "New title",
	})

	require.ErrorIs(t, err, repositoryErr)
}

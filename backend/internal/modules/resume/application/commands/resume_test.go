package commands

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestDeleteResumeHandlerDeletesExpectedVersion(t *testing.T) {
	resume := newCommandTestResume(t, sharedvalueobjects.NewID(), "CV", 7)
	repository := &commandRepositoryStub{}
	handler := NewDeleteResumeHandler(repository, &unitOfWorkStub{})

	err := handler.Handle(context.Background(), DeleteResumeCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 7,
	})
	require.NoError(t, err)
	require.True(t, repository.called)
	require.Equal(t, resume.ID(), repository.deleted.ResumeID)
	require.Equal(t, resume.UserID(), repository.deleted.UserID)
	require.Equal(t, int64(7), repository.deleted.Version)
}

func TestDeleteResumeHandlerRejectsInvalidVersion(t *testing.T) {
	resume := newCommandTestResume(t, sharedvalueobjects.NewID(), "CV", 7)
	repository := &commandRepositoryStub{}
	handler := NewDeleteResumeHandler(repository, &unitOfWorkStub{})

	err := handler.Handle(context.Background(), DeleteResumeCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 0,
	})

	require.Error(t, err)
	require.False(t, repository.called)
}

func TestDeleteResumeHandlerPropagatesConcurrencyConflict(t *testing.T) {
	resume := newCommandTestResume(t, sharedvalueobjects.NewID(), "CV", 7)
	repository := &commandRepositoryStub{err: domain.ErrConcurrencyConflict}
	handler := NewDeleteResumeHandler(repository, &unitOfWorkStub{})

	err := handler.Handle(context.Background(), DeleteResumeCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 6,
	})

	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.True(t, repository.called)
}

func TestCreateResumeHandlerCreatesResumeWithOverview(t *testing.T) {
	repository := &commandRepositoryStub{}
	uow := &unitOfWorkStub{}
	handler := NewCreateResumeHandler(repository, uow)
	userID := sharedvalueobjects.NewID()

	resumeID, err := handler.Handle(context.Background(), CreateResumeCommand{UserID: userID, Title: "  Senior Go Developer  ", Language: domain.ResumeLanguage("en")})
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

func TestReplaceResumeTitleHandlerUpdatesOnlyResumeRoot(t *testing.T) {
	resume := newCommandTestResume(t, sharedvalueobjects.NewID(), "Old title", 4)
	repository := &commandRepositoryStub{resume: resume}
	uow := &unitOfWorkStub{}
	handler := NewReplaceResumeTitleHandler(repository, uow)

	err := handler.Handle(context.Background(), ReplaceResumeTitleCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 4, Title: "  Senior Go Developer  ",
	})

	require.NoError(t, err)
	require.True(t, uow.called)
	require.True(t, repository.called)
	require.Same(t, resume, repository.saved)
	require.Equal(t, int64(4), repository.version)
	require.Equal(t, "Senior Go Developer", resume.Title())
}

func TestReplaceResumeTitleHandlerRejectsStaleVersion(t *testing.T) {
	resume := newCommandTestResume(t, sharedvalueobjects.NewID(), "Old title", 5)
	repository := &commandRepositoryStub{resume: resume}
	handler := NewReplaceResumeTitleHandler(repository, &unitOfWorkStub{})

	err := handler.Handle(context.Background(), ReplaceResumeTitleCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 4, Title: "New title",
	})

	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repository.called)
	require.Equal(t, "Old title", resume.Title())
}

func TestReplaceResumeTitleHandlerRejectsInvalidIdentityOrVersion(t *testing.T) {
	resume := newCommandTestResume(t, sharedvalueobjects.NewID(), "Old title", 1)
	repository := &commandRepositoryStub{resume: resume}
	handler := NewReplaceResumeTitleHandler(repository, &unitOfWorkStub{})

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
	resume := newCommandTestResume(t, sharedvalueobjects.NewID(), "Old title", 1)
	repository := &commandRepositoryStub{resume: resume}
	handler := NewReplaceResumeTitleHandler(repository, &unitOfWorkStub{})

	err := handler.Handle(context.Background(), ReplaceResumeTitleCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 1, Title: "   ",
	})

	require.ErrorIs(t, err, domain.ErrInvalidResume)
	require.False(t, repository.called)
}

func TestReplaceResumeTitleHandlerPropagatesRepositoryError(t *testing.T) {
	repositoryErr := errors.New("repository failure")
	resume := newCommandTestResume(t, sharedvalueobjects.NewID(), "Old title", 1)
	repository := &commandRepositoryStub{resume: resume, err: repositoryErr}
	handler := NewReplaceResumeTitleHandler(repository, &unitOfWorkStub{})

	err := handler.Handle(context.Background(), ReplaceResumeTitleCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 1, Title: "New title",
	})

	require.ErrorIs(t, err, repositoryErr)
}

func TestChangeResumeStatusHandlerPublishesDraft(t *testing.T) {
	resume := newCommandTestResume(t, sharedvalueobjects.NewID(), "CV", 4)
	repository := &commandRepositoryStub{resume: resume, advanceSectionVersion: true}
	handler := NewChangeResumeStatusHandler(repository, &unitOfWorkStub{})

	err := handler.Handle(context.Background(), ChangeResumeStatusCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 4, Status: domain.ResumeStatusOpened,
	})

	require.NoError(t, err)
	require.Equal(t, domain.ResumeStatusOpened, resume.Status())
	require.True(t, repository.called)
}

func TestChangeResumeStatusHandlerSupportsCloseAndReopen(t *testing.T) {
	resume := newCommandTestResume(t, sharedvalueobjects.NewID(), "CV", 4)
	require.NoError(t, resume.ReplaceStatus(domain.ResumeStatusOpened))
	repository := &commandRepositoryStub{resume: resume, advanceSectionVersion: true}
	handler := NewChangeResumeStatusHandler(repository, &unitOfWorkStub{})

	require.NoError(t, handler.Handle(context.Background(), ChangeResumeStatusCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 4, Status: domain.ResumeStatusClosed,
	}))
	require.Equal(t, domain.ResumeStatusClosed, resume.Status())

	resume = newCommandTestResume(t, sharedvalueobjects.NewID(), "CV", 5)
	require.NoError(t, resume.ReplaceStatus(domain.ResumeStatusOpened))
	repository = &commandRepositoryStub{resume: resume, advanceSectionVersion: true}
	handler = NewChangeResumeStatusHandler(repository, &unitOfWorkStub{})
	require.NoError(t, handler.Handle(context.Background(), ChangeResumeStatusCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 5, Status: domain.ResumeStatusClosed,
	}))
	require.NoError(t, handler.Handle(context.Background(), ChangeResumeStatusCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 6, Status: domain.ResumeStatusOpened,
	}))
	require.Equal(t, domain.ResumeStatusOpened, repository.resume.Status())
}

func TestChangeResumeStatusHandlerRejectsDraftTarget(t *testing.T) {
	resume := newCommandTestResume(t, sharedvalueobjects.NewID(), "CV", 1)
	repository := &commandRepositoryStub{resume: resume}
	handler := NewChangeResumeStatusHandler(repository, &unitOfWorkStub{})

	err := handler.Handle(context.Background(), ChangeResumeStatusCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 1, Status: domain.ResumeStatusDraft,
	})

	require.ErrorIs(t, err, domain.ErrInvalidResumeStatusTransition)
	require.False(t, repository.called)
}

func TestChangeResumeStatusHandlerRejectsStaleVersion(t *testing.T) {
	resume := newCommandTestResume(t, sharedvalueobjects.NewID(), "CV", 5)
	repository := &commandRepositoryStub{resume: resume}
	handler := NewChangeResumeStatusHandler(repository, &unitOfWorkStub{})

	err := handler.Handle(context.Background(), ChangeResumeStatusCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 4, Status: domain.ResumeStatusOpened,
	})

	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repository.called)
	require.Equal(t, domain.ResumeStatusDraft, resume.Status())
}

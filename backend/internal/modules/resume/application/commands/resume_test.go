package commands

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type fakeResumeRepository struct {
	called   bool
	resume   *domain.Resume
	saved    *domain.Resume
	expected int64
	deleted  ports.DeleteResumeParams
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

func (r *fakeResumeRepository) SaveSection(ctx context.Context, resume *domain.Resume, _ ports.ResumeMutationSection) error {
	if err := r.Save(ctx, resume); err != nil {
		return err
	}
	restored := domain.RestoreResume(domain.ResumeSnapshot{
		ID: resume.ID(), UserID: resume.UserID(), Title: resume.Title(),
		Language: resume.Language(), Status: resume.Status(), Version: resume.Version() + 1,
		Overview: resume.Overview(), Profile: resume.Profile(), Contacts: resume.Contacts(),
		Skills: resume.Skills(), SpokenLanguages: resume.SpokenLanguages(), Technologies: resume.Technologies(),
		Experiences: resume.Experiences(), Educations: resume.Educations(), Courses: resume.Courses(),
		Activities: resume.Activities(), Projects: resume.Projects(),
	})
	r.resume = restored
	return nil
}

func (r *fakeResumeRepository) Create(_ context.Context, resume *domain.Resume) error {
	r.called = true
	r.resume = resume
	return r.err
}

func (r *fakeResumeRepository) Delete(_ context.Context, params ports.DeleteResumeParams) error {
	r.called = true
	r.deleted = params
	return r.err
}

func TestDeleteResumeHandlerDeletesExpectedVersion(t *testing.T) {
	resume := newTestResume(t, "CV", 7)
	repository := &fakeResumeRepository{}
	handler := NewDeleteResumeHandler(repository, &fakeCommandUnitOfWork{})

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
	resume := newTestResume(t, "CV", 7)
	repository := &fakeResumeRepository{}
	handler := NewDeleteResumeHandler(repository, &fakeCommandUnitOfWork{})

	err := handler.Handle(context.Background(), DeleteResumeCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 0,
	})

	require.Error(t, err)
	require.False(t, repository.called)
}

func TestDeleteResumeHandlerPropagatesConcurrencyConflict(t *testing.T) {
	resume := newTestResume(t, "CV", 7)
	repository := &fakeResumeRepository{err: domain.ErrConcurrencyConflict}
	handler := NewDeleteResumeHandler(repository, &fakeCommandUnitOfWork{})

	err := handler.Handle(context.Background(), DeleteResumeCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 6,
	})

	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.True(t, repository.called)
}

func TestCreateResumeHandlerCreatesResumeWithOverview(t *testing.T) {
	repository := &fakeResumeRepository{}
	uow := &fakeCommandUnitOfWork{}
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

func newTestResume(t *testing.T, title string, version int64) *domain.Resume {
	t.Helper()
	resume, err := domain.NewResume(domain.ResumeParams{
		UserID: sharedvalueobjects.NewID(), Title: title,
		Language: domain.ResumeLanguageEN,
	})
	require.NoError(t, err)
	if version == 1 {
		return resume
	}
	restored := domain.RestoreResume(domain.ResumeSnapshot{
		ID: resume.ID(), UserID: resume.UserID(), Title: resume.Title(),
		Language: resume.Language(), Status: resume.Status(), Version: version,
	})
	return restored
}

func TestReplaceResumeTitleHandlerUpdatesOnlyResumeRoot(t *testing.T) {
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

func TestChangeResumeStatusHandlerPublishesDraft(t *testing.T) {
	resume := newTestResume(t, "CV", 4)
	repository := &fakeResumeRepository{resume: resume}
	handler := NewChangeResumeStatusHandler(repository, &fakeCommandUnitOfWork{})

	err := handler.Handle(context.Background(), ChangeResumeStatusCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 4, Status: domain.ResumeStatusOpened,
	})

	require.NoError(t, err)
	require.Equal(t, domain.ResumeStatusOpened, resume.Status())
	require.True(t, repository.called)
}

func TestChangeResumeStatusHandlerSupportsCloseAndReopen(t *testing.T) {
	resume := newTestResume(t, "CV", 4)
	require.NoError(t, resume.ReplaceStatus(domain.ResumeStatusOpened))
	repository := &fakeResumeRepository{resume: resume}
	handler := NewChangeResumeStatusHandler(repository, &fakeCommandUnitOfWork{})

	require.NoError(t, handler.Handle(context.Background(), ChangeResumeStatusCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 4, Status: domain.ResumeStatusClosed,
	}))
	require.Equal(t, domain.ResumeStatusClosed, resume.Status())

	resume = newTestResume(t, "CV", 5)
	require.NoError(t, resume.ReplaceStatus(domain.ResumeStatusOpened))
	repository = &fakeResumeRepository{resume: resume}
	handler = NewChangeResumeStatusHandler(repository, &fakeCommandUnitOfWork{})
	require.NoError(t, handler.Handle(context.Background(), ChangeResumeStatusCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 5, Status: domain.ResumeStatusClosed,
	}))
	require.NoError(t, handler.Handle(context.Background(), ChangeResumeStatusCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 6, Status: domain.ResumeStatusOpened,
	}))
	require.Equal(t, domain.ResumeStatusOpened, repository.resume.Status())
}

func TestChangeResumeStatusHandlerRejectsDraftTarget(t *testing.T) {
	resume := newTestResume(t, "CV", 1)
	repository := &fakeResumeRepository{resume: resume}
	handler := NewChangeResumeStatusHandler(repository, &fakeCommandUnitOfWork{})

	err := handler.Handle(context.Background(), ChangeResumeStatusCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 1, Status: domain.ResumeStatusDraft,
	})

	require.ErrorIs(t, err, domain.ErrInvalidResumeStatusTransition)
	require.False(t, repository.called)
}

func TestChangeResumeStatusHandlerRejectsStaleVersion(t *testing.T) {
	resume := newTestResume(t, "CV", 5)
	repository := &fakeResumeRepository{resume: resume}
	handler := NewChangeResumeStatusHandler(repository, &fakeCommandUnitOfWork{})

	err := handler.Handle(context.Background(), ChangeResumeStatusCommand{
		ResumeID: resume.ID(), UserID: resume.UserID(), Version: 4, Status: domain.ResumeStatusOpened,
	})

	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repository.called)
	require.Equal(t, domain.ResumeStatusDraft, resume.Status())
}

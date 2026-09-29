package commands

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumevalueobjects "developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/modules/resume_link/domain/valueobjects"
	uservalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type fakeUnitOfWork struct{ called bool }

func (u *fakeUnitOfWork) Do(ctx context.Context, fn func(context.Context) error) error {
	u.called = true
	return fn(ctx)
}

type fakeResumeReader struct {
	view   resumequeries.ResumeView
	err    error
	called bool
	query  resumequeries.GetResumeQuery
}

func (r *fakeResumeReader) GetOwned(_ context.Context, query resumequeries.GetResumeQuery) (resumequeries.ResumeView, error) {
	r.called = true
	r.query = query
	if r.err != nil {
		return resumequeries.ResumeView{}, r.err
	}
	return r.view, nil
}

type fakeRepository struct {
	link   *domain.ResumeLink
	err    error
	called bool
}

func (r *fakeRepository) Create(_ context.Context, link *domain.ResumeLink) error {
	r.called = true
	r.link = link
	return r.err
}

func (r *fakeRepository) GetByID(context.Context, valueobjects.ID) (*domain.ResumeLink, error) {
	return nil, r.err
}
func (r *fakeRepository) ReplaceStatus(context.Context, valueobjects.ID, bool, int64) error {
	return r.err
}

func (r *fakeRepository) Delete(context.Context, valueobjects.ID) error              { return r.err }
func (r *fakeRepository) RecordEvent(context.Context, *domain.ResumeLinkEvent) error { return r.err }

func validCreateCommand() CreateResumeLinkCommand {
	return CreateResumeLinkCommand{UserID: uservalueobjects.NewID(), ResumeID: resumevalueobjects.NewID(), Name: "  HH  "}
}

func TestCreateResumeLinkHandlerCreatesLinkAfterOwnershipValidation(t *testing.T) {
	reader := &fakeResumeReader{}
	repository := &fakeRepository{}
	uow := &fakeUnitOfWork{}
	handler := NewCreateResumeLinkHandler(repository, reader, uow)
	command := validCreateCommand()

	id, err := handler.Handle(context.Background(), command)

	require.NoError(t, err)
	require.False(t, id.IsZero())
	require.True(t, uow.called)
	require.True(t, reader.called)
	require.Equal(t, command.UserID, reader.query.UserID)
	require.Equal(t, command.ResumeID, reader.query.ResumeID)
	require.True(t, repository.called)
	require.Equal(t, id, repository.link.ID())
	require.Equal(t, command.ResumeID, repository.link.ResumeID())
	require.Equal(t, "HH", repository.link.Name())
	require.True(t, repository.link.IsActive())
}

func TestCreateResumeLinkHandlerPropagatesResumeNotFound(t *testing.T) {
	reader := &fakeResumeReader{err: resumequeries.ErrResumeNotFound}
	repository := &fakeRepository{}
	handler := NewCreateResumeLinkHandler(repository, reader, &fakeUnitOfWork{})

	_, err := handler.Handle(context.Background(), validCreateCommand())

	require.ErrorIs(t, err, resumequeries.ErrResumeNotFound)
	require.False(t, repository.called)
}

func TestCreateResumeLinkHandlerRejectsInvalidName(t *testing.T) {
	reader := &fakeResumeReader{}
	repository := &fakeRepository{}
	handler := NewCreateResumeLinkHandler(repository, reader, &fakeUnitOfWork{})
	command := validCreateCommand()
	command.Name = "   "

	_, err := handler.Handle(context.Background(), command)

	require.ErrorIs(t, err, domain.ErrInvalidResumeLinkName)
	require.False(t, reader.called)
	require.False(t, repository.called)
}

func TestCreateResumeLinkHandlerPropagatesRepositoryFailure(t *testing.T) {
	repositoryErr := errors.New("repository failure")
	repository := &fakeRepository{err: repositoryErr}
	handler := NewCreateResumeLinkHandler(repository, &fakeResumeReader{}, &fakeUnitOfWork{})

	_, err := handler.Handle(context.Background(), validCreateCommand())

	require.ErrorIs(t, err, repositoryErr)
}

func TestCreateResumeLinkHandlerRejectsInvalidIdentifiers(t *testing.T) {
	handler := NewCreateResumeLinkHandler(&fakeRepository{}, &fakeResumeReader{}, &fakeUnitOfWork{})
	command := validCreateCommand()
	command.UserID = uservalueobjects.ID{}

	_, err := handler.Handle(context.Background(), command)

	require.Error(t, err)
}

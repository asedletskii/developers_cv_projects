package commands

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestCreateResumeLinkHandlerCreatesLinkAfterOwnershipValidation(t *testing.T) {
	reader := &resumeReaderStub{}
	repository := &resumeLinkRepositoryStub{}
	uow := &unitOfWorkStub{}
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
	require.True(t, repository.lockCalled)
	require.True(t, repository.countCalled)
}

func TestCreateResumeLinkHandlerRejectsAtCapacityBeforeCreate(t *testing.T) {
	repository := &resumeLinkRepositoryStub{count: domain.MaxResumeLinksPerResume}
	handler := NewCreateResumeLinkHandler(repository, &resumeReaderStub{}, &unitOfWorkStub{})

	_, err := handler.Handle(context.Background(), validCreateCommand())

	require.ErrorIs(t, err, domain.ErrResumeLinkLimitExceeded)
	require.True(t, repository.lockCalled)
	require.True(t, repository.countCalled)
	require.False(t, repository.called)
}

func TestCreateResumeLinkHandlerPropagatesCapacityReadFailure(t *testing.T) {
	repositoryErr := errors.New("capacity query failure")
	repository := &resumeLinkRepositoryStub{countErr: repositoryErr}
	handler := NewCreateResumeLinkHandler(repository, &resumeReaderStub{}, &unitOfWorkStub{})

	_, err := handler.Handle(context.Background(), validCreateCommand())

	require.ErrorIs(t, err, repositoryErr)
	require.False(t, repository.called)
}

func TestCreateResumeLinkHandlerPropagatesResumeNotFound(t *testing.T) {
	reader := &resumeReaderStub{err: resumequeries.ErrResumeNotFound}
	repository := &resumeLinkRepositoryStub{}
	handler := NewCreateResumeLinkHandler(repository, reader, &unitOfWorkStub{})

	_, err := handler.Handle(context.Background(), validCreateCommand())

	require.ErrorIs(t, err, resumequeries.ErrResumeNotFound)
	require.False(t, repository.called)
}

func TestCreateResumeLinkHandlerRejectsInvalidName(t *testing.T) {
	reader := &resumeReaderStub{}
	repository := &resumeLinkRepositoryStub{}
	handler := NewCreateResumeLinkHandler(repository, reader, &unitOfWorkStub{})
	command := validCreateCommand()
	command.Name = "   "

	_, err := handler.Handle(context.Background(), command)

	require.ErrorIs(t, err, domain.ErrInvalidResumeLinkName)
	require.False(t, reader.called)
	require.False(t, repository.called)
}

func TestCreateResumeLinkHandlerPropagatesRepositoryFailure(t *testing.T) {
	repositoryErr := errors.New("repository failure")
	repository := &resumeLinkRepositoryStub{err: repositoryErr}
	handler := NewCreateResumeLinkHandler(repository, &resumeReaderStub{}, &unitOfWorkStub{})

	_, err := handler.Handle(context.Background(), validCreateCommand())

	require.ErrorIs(t, err, repositoryErr)
}

func TestCreateResumeLinkHandlerRejectsInvalidIdentifiers(t *testing.T) {
	handler := NewCreateResumeLinkHandler(&resumeLinkRepositoryStub{}, &resumeReaderStub{}, &unitOfWorkStub{})
	command := validCreateCommand()
	command.UserID = sharedvalueobjects.ID{}

	_, err := handler.Handle(context.Background(), command)

	require.Error(t, err)
}

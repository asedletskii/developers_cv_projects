package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
	application "developers-cv/backend/internal/modules/resume_link/application"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestDeleteResumeLinkHandlerChecksOwnershipBeforeDelete(t *testing.T) {
	userID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	linkID := sharedvalueobjects.NewID()
	link, err := domain.RestoreResumeLink(domain.ResumeLinkSnapshot{ID: linkID, ResumeID: resumeID, Name: "HH", Status: domain.ResumeLinkStatusActive, Version: 1})
	require.NoError(t, err)

	repo := &resumeLinkRepositoryStub{link: link}
	ownership := &resumeOwnershipReaderStub{}
	err = NewDeleteResumeLinkHandler(repo, ownership, &analyticsCleanupStub{}, &unitOfWorkStub{}).Handle(context.Background(), DeleteResumeLinkCommand{UserID: userID, LinkID: linkID})
	require.NoError(t, err)
	require.True(t, ownership.lockCalled)
	require.Equal(t, linkID, repo.deleted)
	require.Equal(t, resumeID, ownership.query.ResumeID)
	require.Equal(t, userID, ownership.query.UserID)
}

func TestDeleteResumeLinkHandlerDoesNotDeleteForeignLink(t *testing.T) {
	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ResumeID: sharedvalueobjects.NewID(), Name: "HH"})
	require.NoError(t, err)
	repo := &resumeLinkRepositoryStub{link: link}
	ownership := &resumeOwnershipReaderStub{err: resumecontracts.ErrResumeNotFound}

	err = NewDeleteResumeLinkHandler(repo, ownership, &analyticsCleanupStub{}, &unitOfWorkStub{}).Handle(context.Background(), DeleteResumeLinkCommand{UserID: sharedvalueobjects.NewID(), LinkID: link.ID()})
	require.ErrorIs(t, err, application.ErrResumeNotFound)
	require.True(t, ownership.lockCalled)
	require.True(t, repo.deleted.IsZero())
}

func TestDeleteResumeLinkHandlerPropagatesMissingLink(t *testing.T) {
	repo := &resumeLinkRepositoryStub{getErr: domain.ErrEntityNotFound}
	err := NewDeleteResumeLinkHandler(repo, &resumeOwnershipReaderStub{}, &analyticsCleanupStub{}, &unitOfWorkStub{}).Handle(context.Background(), DeleteResumeLinkCommand{UserID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID()})
	require.ErrorIs(t, err, application.ErrResumeLinkNotFound)
	require.ErrorIs(t, err, domain.ErrEntityNotFound)
}

func TestDeleteResumeLinkHandlerNormalizesDeleteRaceNotFound(t *testing.T) {
	link := restoreCommandTestLink(t, domain.ResumeLinkStatusActive, 1)
	repo := &resumeLinkRepositoryStub{link: link, deleteErr: domain.ErrEntityNotFound}

	err := NewDeleteResumeLinkHandler(repo, &resumeOwnershipReaderStub{}, &analyticsCleanupStub{}, &unitOfWorkStub{}).Handle(context.Background(), DeleteResumeLinkCommand{
		UserID: sharedvalueobjects.NewID(), LinkID: link.ID(),
	})

	require.ErrorIs(t, err, application.ErrResumeLinkNotFound)
	require.ErrorIs(t, err, domain.ErrEntityNotFound)
}

func TestDeleteResumeLinkHandlerRejectsNilLinkFromRepository(t *testing.T) {
	repo := &resumeLinkRepositoryStub{returnNil: true}
	err := NewDeleteResumeLinkHandler(repo, &resumeOwnershipReaderStub{}, &analyticsCleanupStub{}, &unitOfWorkStub{}).Handle(context.Background(), DeleteResumeLinkCommand{
		UserID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(),
	})

	require.EqualError(t, err, "delete resume link: repository returned nil link")
	require.True(t, repo.deleted.IsZero())
}

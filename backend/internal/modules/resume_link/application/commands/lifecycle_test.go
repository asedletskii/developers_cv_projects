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

func TestReplaceStatusResumeLinkDeactivatesOwnedLink(t *testing.T) {
	link := restoreCommandTestLink(t, domain.ResumeLinkStatusActive, 2)
	repo := &resumeLinkRepositoryStub{link: link}
	err := NewReplaceStatusResumeLinkHandler(repo, &resumeOwnershipReaderStub{}, &unitOfWorkStub{}).Handle(context.Background(), ReplaceStatusResumeLinkCommand{
		UserID: sharedvalueobjects.NewID(), LinkID: link.ID(), Active: false, Version: 2,
	})
	require.NoError(t, err)
	require.True(t, repo.replaced)
	require.False(t, repo.active)
	require.Equal(t, int64(2), repo.version)
	require.False(t, link.IsActive())
}

func TestReplaceStatusResumeLinkReactivatesLink(t *testing.T) {
	link := restoreCommandTestLink(t, domain.ResumeLinkStatusInactive, 5)
	repo := &resumeLinkRepositoryStub{link: link}
	err := NewReplaceStatusResumeLinkHandler(repo, &resumeOwnershipReaderStub{}, &unitOfWorkStub{}).Handle(context.Background(), ReplaceStatusResumeLinkCommand{
		UserID: sharedvalueobjects.NewID(), LinkID: link.ID(), Active: true, Version: 5,
	})
	require.NoError(t, err)
	require.True(t, repo.active)
	require.True(t, link.IsActive())
}

func TestReplaceStatusResumeLinkChecksOwnership(t *testing.T) {
	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ResumeID: sharedvalueobjects.NewID(), Name: "HH"})
	require.NoError(t, err)
	repo := &resumeLinkRepositoryStub{link: link}
	ownership := &resumeOwnershipReaderStub{err: resumecontracts.ErrResumeNotFound}
	err = NewReplaceStatusResumeLinkHandler(repo, ownership, &unitOfWorkStub{}).Handle(context.Background(), ReplaceStatusResumeLinkCommand{
		UserID: sharedvalueobjects.NewID(), LinkID: link.ID(), Active: false, Version: 1,
	})
	require.ErrorIs(t, err, application.ErrResumeNotFound)
	require.True(t, ownership.lockCalled)
	require.False(t, repo.replaced)
}

func TestReplaceStatusResumeLinkRejectsStaleVersion(t *testing.T) {
	link := restoreCommandTestLink(t, domain.ResumeLinkStatusActive, 2)
	repo := &resumeLinkRepositoryStub{link: link}
	err := NewReplaceStatusResumeLinkHandler(repo, &resumeOwnershipReaderStub{}, &unitOfWorkStub{}).Handle(context.Background(), ReplaceStatusResumeLinkCommand{
		UserID: sharedvalueobjects.NewID(), LinkID: link.ID(), Active: false, Version: 1,
	})
	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repo.replaced)
}

func TestReplaceStatusResumeLinkRejectsMissingLink(t *testing.T) {
	repo := &resumeLinkRepositoryStub{getErr: domain.ErrEntityNotFound}
	err := NewReplaceStatusResumeLinkHandler(repo, &resumeOwnershipReaderStub{}, &unitOfWorkStub{}).Handle(context.Background(), ReplaceStatusResumeLinkCommand{
		UserID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(), Active: false, Version: 1,
	})
	require.ErrorIs(t, err, application.ErrResumeLinkNotFound)
	require.ErrorIs(t, err, domain.ErrEntityNotFound)
}

func TestReplaceStatusResumeLinkNormalizesReplaceRaceNotFound(t *testing.T) {
	link := restoreCommandTestLink(t, domain.ResumeLinkStatusActive, 1)
	repo := &resumeLinkRepositoryStub{link: link, replaceErr: domain.ErrEntityNotFound}

	err := NewReplaceStatusResumeLinkHandler(repo, &resumeOwnershipReaderStub{}, &unitOfWorkStub{}).Handle(context.Background(), ReplaceStatusResumeLinkCommand{
		UserID: sharedvalueobjects.NewID(), LinkID: link.ID(), Active: false, Version: 1,
	})

	require.ErrorIs(t, err, application.ErrResumeLinkNotFound)
	require.ErrorIs(t, err, domain.ErrEntityNotFound)
}

func TestReplaceStatusResumeLinkRejectsNilLinkFromRepository(t *testing.T) {
	repo := &resumeLinkRepositoryStub{}
	err := NewReplaceStatusResumeLinkHandler(repo, &resumeOwnershipReaderStub{}, &unitOfWorkStub{}).Handle(context.Background(), ReplaceStatusResumeLinkCommand{
		UserID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(), Active: false, Version: 1,
	})

	require.EqualError(t, err, "replace resume link status: repository returned nil link")
}

func TestReplaceStatusResumeLinkNoOpKeepsVersionAndSkipsPersistence(t *testing.T) {
	link := restoreCommandTestLink(t, domain.ResumeLinkStatusActive, 7)
	repo := &resumeLinkRepositoryStub{link: link}

	err := NewReplaceStatusResumeLinkHandler(repo, &resumeOwnershipReaderStub{}, &unitOfWorkStub{}).Handle(context.Background(), ReplaceStatusResumeLinkCommand{
		UserID: sharedvalueobjects.NewID(), LinkID: link.ID(), Active: true, Version: 7,
	})

	require.NoError(t, err)
	require.False(t, repo.replaced)
	require.True(t, link.IsActive())
	require.Equal(t, int64(7), link.Version())
}

func TestReplaceStatusResumeLinkNoOpStillRejectsStaleVersion(t *testing.T) {
	link := restoreCommandTestLink(t, domain.ResumeLinkStatusActive, 7)
	repo := &resumeLinkRepositoryStub{link: link}

	err := NewReplaceStatusResumeLinkHandler(repo, &resumeOwnershipReaderStub{}, &unitOfWorkStub{}).Handle(context.Background(), ReplaceStatusResumeLinkCommand{
		UserID: sharedvalueobjects.NewID(), LinkID: link.ID(), Active: true, Version: 6,
	})

	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repo.replaced)
}

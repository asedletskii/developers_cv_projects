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

func TestRenameResumeLinkRenamesOwnedLink(t *testing.T) {
	link := restoreCommandTestLink(t, domain.ResumeLinkStatusActive, 2)
	repo := &resumeLinkRepositoryStub{link: link}

	err := NewRenameResumeLinkHandler(repo, &resumeOwnershipReaderStub{}, &unitOfWorkStub{}).Handle(
		context.Background(),
		RenameResumeLinkCommand{
			UserID: sharedvalueobjects.NewID(), LinkID: link.ID(), Name: "New Name", Version: 2,
		},
	)

	require.NoError(t, err)
	require.True(t, repo.renamed)
	require.Equal(t, "New Name", repo.renameName)
	require.Equal(t, int64(2), repo.renameVersion)
	require.Equal(t, "New Name", link.Name())
}

func TestRenameResumeLinkChecksOwnership(t *testing.T) {
	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ResumeID: sharedvalueobjects.NewID(), Name: "HH"})
	require.NoError(t, err)
	repo := &resumeLinkRepositoryStub{link: link}
	ownership := &resumeOwnershipReaderStub{err: resumecontracts.ErrResumeNotFound}

	err = NewRenameResumeLinkHandler(repo, ownership, &unitOfWorkStub{}).Handle(
		context.Background(),
		RenameResumeLinkCommand{
			UserID: sharedvalueobjects.NewID(), LinkID: link.ID(), Name: "New", Version: 1,
		},
	)

	require.ErrorIs(t, err, application.ErrResumeNotFound)
	require.True(t, ownership.lockCalled)
	require.False(t, repo.renamed)
}

func TestRenameResumeLinkRejectsStaleVersion(t *testing.T) {
	link := restoreCommandTestLink(t, domain.ResumeLinkStatusActive, 2)
	repo := &resumeLinkRepositoryStub{link: link}

	err := NewRenameResumeLinkHandler(repo, &resumeOwnershipReaderStub{}, &unitOfWorkStub{}).Handle(
		context.Background(),
		RenameResumeLinkCommand{
			UserID: sharedvalueobjects.NewID(), LinkID: link.ID(), Name: "New", Version: 1,
		},
	)

	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repo.renamed)
}

func TestRenameResumeLinkRejectsMissingLink(t *testing.T) {
	repo := &resumeLinkRepositoryStub{getErr: domain.ErrEntityNotFound}

	err := NewRenameResumeLinkHandler(repo, &resumeOwnershipReaderStub{}, &unitOfWorkStub{}).Handle(
		context.Background(),
		RenameResumeLinkCommand{
			UserID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(), Name: "New", Version: 1,
		},
	)

	require.ErrorIs(t, err, application.ErrResumeLinkNotFound)
	require.False(t, repo.renamed)
}

func TestRenameResumeLinkValidatesName(t *testing.T) {
	link := restoreCommandTestLink(t, domain.ResumeLinkStatusActive, 1)
	repo := &resumeLinkRepositoryStub{link: link}

	err := NewRenameResumeLinkHandler(repo, &resumeOwnershipReaderStub{}, &unitOfWorkStub{}).Handle(
		context.Background(),
		RenameResumeLinkCommand{
			UserID: sharedvalueobjects.NewID(), LinkID: link.ID(), Name: "   ", Version: 1,
		},
	)

	require.ErrorIs(t, err, domain.ErrInvalidResumeLinkName)
	require.False(t, repo.renamed)
}

func TestRenameResumeLinkRejectsNilLinkFromRepository(t *testing.T) {
	repo := &resumeLinkRepositoryStub{returnNil: true}

	err := NewRenameResumeLinkHandler(repo, &resumeOwnershipReaderStub{}, &unitOfWorkStub{}).Handle(
		context.Background(),
		RenameResumeLinkCommand{
			UserID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(), Name: "New", Version: 1,
		},
	)

	require.EqualError(t, err, "rename resume link: repository returned nil link")
	require.False(t, repo.renamed)
}

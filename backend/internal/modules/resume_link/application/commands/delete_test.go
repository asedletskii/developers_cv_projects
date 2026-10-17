package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type deleteRepositoryStub struct {
	link      *domain.ResumeLink
	deleted   sharedvalueobjects.ID
	getErr    error
	returnNil bool
	delErr    error
}

func (s *deleteRepositoryStub) Create(context.Context, *domain.ResumeLink) error { return nil }
func (s *deleteRepositoryStub) LockResumeForLinkCreation(context.Context, sharedvalueobjects.ID) error {
	return nil
}
func (s *deleteRepositoryStub) CountByResumeID(context.Context, sharedvalueobjects.ID) (int64, error) {
	return 0, nil
}
func (s *deleteRepositoryStub) GetByID(context.Context, sharedvalueobjects.ID) (*domain.ResumeLink, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	if s.returnNil {
		return nil, nil
	}
	return s.link, nil
}
func (s *deleteRepositoryStub) ReplaceStatus(context.Context, sharedvalueobjects.ID, bool, int64) error {
	return nil
}

func (s *deleteRepositoryStub) Delete(_ context.Context, id sharedvalueobjects.ID) error {
	s.deleted = id
	return s.delErr
}

type deleteOwnershipStub struct {
	err error
	q   resumequeries.GetResumeQuery
}

func (s *deleteOwnershipStub) GetOwned(_ context.Context, q resumequeries.GetResumeQuery) (resumequeries.ResumeView, error) {
	s.q = q
	if s.err != nil {
		return resumequeries.ResumeView{}, s.err
	}
	return resumequeries.ResumeView{ID: q.ResumeID.String()}, nil
}

func (s *deleteOwnershipStub) OwnsResume(_ context.Context, query resumequeries.GetResumeQuery) (bool, error) {
	s.q = query
	return s.err == nil, s.err
}

func TestDeleteResumeLinkHandlerChecksOwnershipBeforeDelete(t *testing.T) {
	userID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	linkID := sharedvalueobjects.NewID()
	link, err := domain.RestoreResumeLink(domain.ResumeLinkSnapshot{ID: linkID, ResumeID: resumeID, Name: "HH", Status: domain.ResumeLinkStatusActive, Version: 1})
	require.NoError(t, err)

	repo := &deleteRepositoryStub{link: link}
	ownership := &deleteOwnershipStub{}
	err = NewDeleteResumeLinkHandler(repo, ownership, &fakeUnitOfWork{}).Handle(context.Background(), DeleteResumeLinkCommand{UserID: userID, LinkID: linkID})
	require.NoError(t, err)
	require.Equal(t, linkID, repo.deleted)
	require.Equal(t, resumeID, ownership.q.ResumeID)
	require.Equal(t, userID, ownership.q.UserID)
}

func TestDeleteResumeLinkHandlerDoesNotDeleteForeignLink(t *testing.T) {
	link, err := domain.NewResumeLink(domain.ResumeLinkParams{ResumeID: sharedvalueobjects.NewID(), Name: "HH"})
	require.NoError(t, err)
	repo := &deleteRepositoryStub{link: link}
	ownership := &deleteOwnershipStub{err: resumequeries.ErrResumeNotFound}

	err = NewDeleteResumeLinkHandler(repo, ownership, &fakeUnitOfWork{}).Handle(context.Background(), DeleteResumeLinkCommand{UserID: sharedvalueobjects.NewID(), LinkID: link.ID()})
	require.ErrorIs(t, err, resumequeries.ErrResumeNotFound)
	require.True(t, repo.deleted.IsZero())
}

func TestDeleteResumeLinkHandlerPropagatesMissingLink(t *testing.T) {
	repo := &deleteRepositoryStub{getErr: domain.ErrEntityNotFound}
	err := NewDeleteResumeLinkHandler(repo, &deleteOwnershipStub{}, &fakeUnitOfWork{}).Handle(context.Background(), DeleteResumeLinkCommand{UserID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID()})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)
}

func TestDeleteResumeLinkHandlerRejectsNilLinkFromRepository(t *testing.T) {
	repo := &deleteRepositoryStub{returnNil: true}
	err := NewDeleteResumeLinkHandler(repo, &deleteOwnershipStub{}, &fakeUnitOfWork{}).Handle(context.Background(), DeleteResumeLinkCommand{
		UserID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(),
	})

	require.EqualError(t, err, "delete resume link: repository returned nil link")
	require.True(t, repo.deleted.IsZero())
}

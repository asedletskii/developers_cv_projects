package http

import (
	"context"
	"errors"

	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type resumeOwnershipReaderStub struct {
	err error
}

func (s resumeOwnershipReaderStub) Handle(context.Context, resumecontracts.OwnedResumeSectionsQuery) ([]resumecontracts.ResumeSectionSettings, error) {
	if s.err != nil {
		return nil, s.err
	}
	return nil, nil
}

func (s resumeOwnershipReaderStub) OwnsResume(context.Context, resumecontracts.ResumeOwnershipQuery) (bool, error) {
	if errors.Is(s.err, resumecontracts.ErrResumeNotFound) {
		return false, nil
	}
	if s.err != nil {
		return false, s.err
	}
	return true, nil
}
func (s resumeOwnershipReaderStub) LockOwnedResumeForUpdate(_ context.Context, query resumecontracts.OwnedResumeLock) (bool, error) {
	if errors.Is(s.err, resumecontracts.ErrResumeNotFound) {
		return false, nil
	}
	if s.err != nil {
		return false, s.err
	}
	return true, nil
}

func (s resumeOwnershipReaderStub) LockPublicationState(_ context.Context, query resumecontracts.ResumePublicationQuery) (resumecontracts.ResumePublicationState, error) {
	if errors.Is(s.err, resumecontracts.ErrResumeNotFound) {
		return resumecontracts.ResumePublicationState{}, resumecontracts.ErrResumeNotFound
	}
	if s.err != nil {
		return resumecontracts.ResumePublicationState{}, s.err
	}
	return resumecontracts.ResumePublicationState{Exists: true, Published: true, Status: "opened"}, nil
}

type resumeLinkPublicReaderStub struct {
	view  queries.ResumeLinkPublicView
	items []queries.ResumeLinkListItem
	err   error
}

func (s resumeLinkPublicReaderStub) GetByID(context.Context, queries.PublicResumeLinkQuery) (queries.ResumeLinkPublicView, error) {
	return s.view, s.err
}

func (s resumeLinkPublicReaderStub) LockByID(context.Context, queries.PublicResumeLinkQuery) (queries.ResumeLinkPublicView, error) {
	return s.view, s.err
}

func (s resumeLinkPublicReaderStub) ListByResumeID(context.Context, queries.ListResumeLinksQuery) ([]queries.ResumeLinkListItem, error) {
	return s.items, s.err
}

type analyticsCleanupStub struct {
	called bool
	err    error
}

func (s *analyticsCleanupStub) DeleteByLinkIDs(context.Context, []sharedvalueobjects.ID) error {
	s.called = true
	return s.err
}

type resumeLinkRepositoryStub struct {
	err            error
	getErr         error
	delErr         error
	replaceErr     error
	renameErr      error
	link           *domain.ResumeLink
	deleted        bool
	replaced       bool
	renamed        bool
	replaceActive  bool
	replaceVersion int64
	renameName     string
	renameVersion  int64
}

func (s *resumeLinkRepositoryStub) Create(context.Context, *domain.ResumeLink) error {
	return s.err
}

func (*resumeLinkRepositoryStub) CountByResumeID(context.Context, sharedvalueobjects.ID) (int64, error) {
	return 0, nil
}

func (s *resumeLinkRepositoryStub) GetByID(context.Context, sharedvalueobjects.ID) (*domain.ResumeLink, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	if s.err != nil {
		return nil, s.err
	}
	return s.link, nil
}

func (s *resumeLinkRepositoryStub) GetByIDForUpdate(context.Context, sharedvalueobjects.ID) (*domain.ResumeLink, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	if s.err != nil {
		return nil, s.err
	}
	return s.link, nil
}

func (s *resumeLinkRepositoryStub) Rename(_ context.Context, _ sharedvalueobjects.ID, name string, version int64) error {
	s.renamed = true
	s.renameName = name
	s.renameVersion = version
	if s.renameErr != nil {
		return s.renameErr
	}
	return s.err
}

func (s *resumeLinkRepositoryStub) ReplaceStatus(_ context.Context, _ sharedvalueobjects.ID, active bool, version int64) error {
	s.replaced = true
	s.replaceActive = active
	s.replaceVersion = version
	if s.replaceErr != nil {
		return s.replaceErr
	}
	return s.err
}

func (s *resumeLinkRepositoryStub) Delete(context.Context, sharedvalueobjects.ID) error {
	s.deleted = true
	if s.delErr != nil {
		return s.delErr
	}
	return s.err
}

type unitOfWorkStub struct{}

func (unitOfWorkStub) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

var _ resumecontracts.ResumeOwnershipChecker = resumeOwnershipReaderStub{}
var _ ports.ResumeLinkRepository = (*resumeLinkRepositoryStub)(nil)

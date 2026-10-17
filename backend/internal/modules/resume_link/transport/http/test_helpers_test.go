package http

import (
	"context"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type resumeOwnershipReaderStub struct {
	view resumequeries.ResumeView
	err  error
}

func (s resumeOwnershipReaderStub) GetOwned(context.Context, resumequeries.GetResumeQuery) (resumequeries.ResumeView, error) {
	if s.err != nil {
		return resumequeries.ResumeView{}, s.err
	}
	if s.view.ID != "" {
		return s.view, nil
	}
	return resumequeries.ResumeView{ID: "resume"}, nil
}

func (s resumeOwnershipReaderStub) GetOwnedOverview(context.Context, resumequeries.GetResumeQuery) (resumequeries.ResumeOverviewView, error) {
	if s.err != nil {
		return resumequeries.ResumeOverviewView{}, s.err
	}
	return resumequeries.ResumeOverviewView{}, nil
}

func (s resumeOwnershipReaderStub) OwnsResume(context.Context, resumequeries.GetResumeQuery) (bool, error) {
	return s.err == nil, s.err
}

type resumeLinkPublicReaderStub struct {
	view queries.ResumeLinkPublicView
	err  error
}

func (s resumeLinkPublicReaderStub) GetByID(context.Context, queries.PublicResumeLinkQuery) (queries.ResumeLinkPublicView, error) {
	return s.view, s.err
}

type resumeLinkRepositoryStub struct {
	err            error
	getErr         error
	delErr         error
	replaceErr     error
	link           *domain.ResumeLink
	deleted        bool
	replaced       bool
	replaceActive  bool
	replaceVersion int64
}

func (s *resumeLinkRepositoryStub) Create(context.Context, *domain.ResumeLink) error {
	return s.err
}

func (*resumeLinkRepositoryStub) LockResumeForLinkCreation(context.Context, sharedvalueobjects.ID) error {
	return nil
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

var _ ports.ResumeOwnershipReader = resumeOwnershipReaderStub{}
var _ ports.ResumeLinkRepository = (*resumeLinkRepositoryStub)(nil)

package queries

import (
	"context"

	analyticsports "developers-cv/backend/internal/modules/analytics/application/ports"
	analyticsdomain "developers-cv/backend/internal/modules/analytics/domain"
	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ownershipReaderStub struct {
	err            error
	viewID         string
	query          resumequeries.GetResumeQuery
	calls          int
	getCalls       int
	ownershipCalls int
}

func (s *ownershipReaderStub) GetOwned(_ context.Context, query resumequeries.GetResumeQuery) (resumequeries.ResumeView, error) {
	s.calls++
	s.getCalls++
	s.query = query
	if s.err != nil {
		return resumequeries.ResumeView{}, s.err
	}
	viewID := s.viewID
	if viewID == "" {
		viewID = query.ResumeID.String()
	}
	return resumequeries.ResumeView{ID: viewID}, nil
}
func (s *ownershipReaderStub) GetOwnedOverview(_ context.Context, query resumequeries.GetResumeQuery) (resumequeries.ResumeOverviewView, error) {
	s.calls++
	s.getCalls++
	s.query = query
	if s.err != nil {
		return resumequeries.ResumeOverviewView{}, s.err
	}
	return resumequeries.ResumeOverviewView{}, nil
}
func (s *ownershipReaderStub) OwnsResume(_ context.Context, query resumequeries.GetResumeQuery) (bool, error) {
	s.calls++
	s.ownershipCalls++
	s.query = query
	if s.err != nil {
		return false, s.err
	}
	return s.viewID == "" || s.viewID == query.ResumeID.String(), nil
}

type resumeLinkReaderStub struct {
	view  ResumeLinkPublicView
	err   error
	query PublicResumeLinkQuery
}

func (s *resumeLinkReaderStub) GetByID(_ context.Context, query PublicResumeLinkQuery) (ResumeLinkPublicView, error) {
	s.query = query
	return s.view, s.err
}

type listReaderStub struct {
	items []ResumeLinkListItem
	err   error
	query ListResumeLinksQuery
}

func (s *listReaderStub) ListByResumeID(_ context.Context, query ListResumeLinksQuery) ([]ResumeLinkListItem, error) {
	s.query = query
	return s.items, s.err
}

type publicResumeQueryStub struct {
	view  resumequeries.PublicResumeView
	err   error
	query resumequeries.PublicResumeQuery
}

func (s *publicResumeQueryStub) Handle(_ context.Context, query resumequeries.PublicResumeQuery) (resumequeries.PublicResumeView, error) {
	s.query = query
	return s.view, s.err
}

type publicVisitCreatorStub struct {
	identity analyticsports.VisitorIdentityResult
	err      error
	input    analyticsports.CreatePublicVisitInput
	calls    int
}

func (s *publicVisitCreatorStub) Handle(_ context.Context, input analyticsports.CreatePublicVisitInput) (*analyticsdomain.View, analyticsports.VisitorIdentityResult, error) {
	s.calls++
	s.input = input
	return nil, s.identity, s.err
}

type contactViewResumeQueryStub struct {
	view  resumequeries.ContactView
	err   error
	query resumequeries.ContactViewQuery
}

func (s *contactViewResumeQueryStub) Handle(_ context.Context, query resumequeries.ContactViewQuery) (resumequeries.ContactView, error) {
	s.query = query
	return s.view, s.err
}

type contactViewResumeDownloadQueryStub struct {
	view  resumequeries.DownloadContactsView
	err   error
	query resumequeries.DownloadContactsQuery
}

func (s *contactViewResumeDownloadQueryStub) Handle(_ context.Context, query resumequeries.DownloadContactsQuery) (resumequeries.DownloadContactsView, error) {
	s.query = query
	return s.view, s.err
}

type contactViewAnalyticsRecorderStub struct {
	err       error
	calls     int
	initiator string
	viewID    sharedvalueobjects.ID
	linkID    sharedvalueobjects.ID
}

func (s *contactViewAnalyticsRecorderStub) RecordContactView(_ context.Context, viewID, linkID sharedvalueobjects.ID, initiator string) error {
	s.calls++
	s.initiator = "contact"
	s.viewID = viewID
	s.linkID = linkID
	return s.err
}
func (s *contactViewAnalyticsRecorderStub) RecordDownload(_ context.Context, viewID, linkID sharedvalueobjects.ID) error {
	s.calls++
	s.initiator = "download"
	s.viewID = viewID
	s.linkID = linkID
	return s.err
}

func newPublicAnalyticsTestHandler(
	owner ports.ResumeOwnershipReader,
	visit analyticsports.PublicVisitCreator,
	reporter PublicAnalyticsErrorReporter,
) (*PublicResumeLinkQueryHandler, PublicResumeLinkQuery, *publicResumeQueryStub) {
	linkID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	links := &resumeLinkReaderStub{view: ResumeLinkPublicView{
		ID: linkID.String(), ResumeID: resumeID.String(), Status: string(domain.ResumeLinkStatusActive),
	}}
	resume := &publicResumeQueryStub{view: resumequeries.PublicResumeView{Title: "Public CV"}}
	handler := NewPublicResumeLinkQueryHandlerWithVisit(NewActivePublicResumeLinkResolver(links), resume, owner, visit, reporter)
	query := PublicResumeLinkQuery{
		LinkID: linkID,
		Visit: &analyticsports.CreatePublicVisitInput{
			LinkID:  linkID,
			Visitor: analyticsports.PublicVisitInput{OwnerID: sharedvalueobjects.NewID()},
		},
	}
	return handler, query, resume
}

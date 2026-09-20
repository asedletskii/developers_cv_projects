package queries

import (
	"context"
	"errors"

	analyticscontracts "developers-cv/backend/internal/modules/analytics/contracts"
	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ownershipReaderStub struct {
	err            error
	viewID         string
	query          resumecontracts.ResumeOwnershipQuery
	calls          int
	ownershipCalls int
}

func (s *ownershipReaderStub) Handle(_ context.Context, query resumecontracts.OwnedResumeSectionsQuery) ([]resumecontracts.ResumeSectionSettings, error) {
	s.calls++
	s.query = resumecontracts.ResumeOwnershipQuery(query)
	if s.err != nil {
		return nil, s.err
	}
	return nil, nil
}

func (s *ownershipReaderStub) OwnsResume(_ context.Context, query resumecontracts.ResumeOwnershipQuery) (bool, error) {
	s.calls++
	s.ownershipCalls++
	s.query = query
	if errors.Is(s.err, resumecontracts.ErrResumeNotFound) {
		return false, nil
	}
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
	view  resumecontracts.PublicResumeView
	err   error
	query resumecontracts.PublicResumeQuery
}

func (s *publicResumeQueryStub) Handle(_ context.Context, query resumecontracts.PublicResumeQuery) (resumecontracts.PublicResumeView, error) {
	s.query = query
	return s.view, s.err
}

type publicVisitCreatorStub struct {
	identity analyticscontracts.VisitorIdentityResult
	err      error
	input    analyticscontracts.CreatePublicVisitInput
	calls    int
}

func (s *publicVisitCreatorStub) Handle(_ context.Context, input analyticscontracts.CreatePublicVisitInput) (analyticscontracts.VisitorIdentityResult, error) {
	s.calls++
	s.input = input
	return s.identity, s.err
}

type contactViewResumeQueryStub struct {
	view  resumecontracts.ContactView
	err   error
	query resumecontracts.ContactViewQuery
}

func (s *contactViewResumeQueryStub) Handle(_ context.Context, query resumecontracts.ContactViewQuery) (resumecontracts.ContactView, error) {
	s.query = query
	return s.view, s.err
}

type contactViewResumeDownloadQueryStub struct {
	view  resumecontracts.DownloadContactsView
	err   error
	query resumecontracts.DownloadContactsQuery
}

func (s *contactViewResumeDownloadQueryStub) Handle(_ context.Context, query resumecontracts.DownloadContactsQuery) (resumecontracts.DownloadContactsView, error) {
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
	owner resumecontracts.ResumeOwnershipChecker,
	visit analyticscontracts.PublicVisitCreator,
	reporter PublicAnalyticsErrorReporter,
) (*PublicResumeLinkQueryHandler, PublicResumeLinkQuery, *publicResumeQueryStub) {
	linkID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	links := &resumeLinkReaderStub{view: ResumeLinkPublicView{
		ID: linkID.String(), ResumeID: resumeID.String(), Status: string(domain.ResumeLinkStatusActive),
	}}
	resume := &publicResumeQueryStub{view: resumecontracts.PublicResumeView{
		Title:        "Public CV",
		Technologies: []resumecontracts.PublicTechnologyView{{Name: "Go"}},
	}}
	handler := NewPublicResumeLinkQueryHandlerWithVisit(NewActivePublicResumeLinkResolver(links), resume, owner, visit, reporter)
	query := PublicResumeLinkQuery{
		LinkID: linkID,
		Visit: &analyticscontracts.CreatePublicVisitInput{
			LinkID:  linkID,
			Visitor: analyticscontracts.PublicVisitInput{OwnerID: sharedvalueobjects.NewID()},
		},
	}
	return handler, query, resume
}

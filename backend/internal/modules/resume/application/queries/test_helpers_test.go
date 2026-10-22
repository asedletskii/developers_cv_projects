package queries

import (
	"context"

	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type contactViewReaderStub struct {
	view  resumecontracts.ContactView
	err   error
	query resumecontracts.ContactViewQuery
}

func (s *contactViewReaderStub) GetContactForView(_ context.Context, query resumecontracts.ContactViewQuery) (resumecontracts.ContactView, error) {
	s.query = query
	if s.err != nil {
		return resumecontracts.ContactView{}, s.err
	}
	return s.view, nil
}

type downloadContactsReaderStub struct {
	view  resumecontracts.DownloadContactsView
	err   error
	query resumecontracts.DownloadContactsQuery
}

func (s *downloadContactsReaderStub) GetContactsForDownload(_ context.Context, query resumecontracts.DownloadContactsQuery) (resumecontracts.DownloadContactsView, error) {
	s.query = query
	return s.view, s.err
}

type listResumesReaderStub struct {
	items  []ResumeListItem
	err    error
	userID sharedvalueobjects.ID
}

func (s *listResumesReaderStub) ListResumes(_ context.Context, query ListResumesQuery) ([]ResumeListItem, error) {
	s.userID = query.UserID
	return s.items, s.err
}

type publicResumeReaderStub struct {
	view ResumeView
	err  error
}

func (s publicResumeReaderStub) GetPublic(context.Context, resumecontracts.PublicResumeQuery) (ResumeView, error) {
	if s.err != nil {
		return ResumeView{}, s.err
	}
	return s.view, nil
}

type resumeReaderStub struct{ view ResumeView }

func (s resumeReaderStub) GetOwned(context.Context, GetResumeQuery) (ResumeView, error) {
	return s.view, nil
}

type resumeReaderErrorStub struct{ err error }

func (s resumeReaderErrorStub) GetOwned(context.Context, GetResumeQuery) (ResumeView, error) {
	return ResumeView{}, s.err
}

func stringPtr(value string) *string {
	return &value
}

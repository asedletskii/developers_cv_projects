package queries

import (
	"context"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type contactViewReaderStub struct {
	view  ContactView
	err   error
	query ContactViewQuery
}

func (s *contactViewReaderStub) GetContactForView(_ context.Context, query ContactViewQuery) (ContactView, error) {
	s.query = query
	if s.err != nil {
		return ContactView{}, s.err
	}
	return s.view, nil
}

type downloadContactsReaderStub struct {
	view  DownloadContactsView
	err   error
	query DownloadContactsQuery
}

func (s *downloadContactsReaderStub) GetContactsForDownload(_ context.Context, query DownloadContactsQuery) (DownloadContactsView, error) {
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

func (s publicResumeReaderStub) GetPublic(context.Context, PublicResumeQuery) (PublicResumeView, error) {
	if s.err != nil {
		return PublicResumeView{}, s.err
	}
	return ToPublicResumeView(s.view), nil
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

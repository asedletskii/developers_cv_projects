package queries

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type contactViewLinkReaderStub struct {
	view  ResumeLinkPublicView
	err   error
	query PublicResumeLinkQuery
}

func (s *contactViewLinkReaderStub) GetByID(_ context.Context, query PublicResumeLinkQuery) (ResumeLinkPublicView, error) {
	s.query = query
	return s.view, s.err
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

func (s *contactViewResumeDownloadQueryStub) Handle(
	_ context.Context,
	query resumequeries.DownloadContactsQuery,
) (resumequeries.DownloadContactsView, error) {
	s.query = query
	return s.view, s.err
}

func TestContactViewThroughLinkResolvesActiveLink(t *testing.T) {
	linkID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	links := &contactViewLinkReaderStub{view: ResumeLinkPublicView{
		ID: linkID.String(), ResumeID: resumeID.String(), Status: string(domain.ResumeLinkStatusActive),
	}}
	resume := &contactViewResumeQueryStub{view: resumequeries.ContactView{Contact: resumedomain.ContactKeyEmail, Value: "person@example.com"}}

	view, err := NewContactViewThroughLinkQueryHandler(links, resume).Handle(context.Background(), ContactViewThroughLinkQuery{
		LinkID:  linkID,
		Contact: resumedomain.ContactKeyEmail,
	})

	require.NoError(t, err)
	require.Equal(t, resumedomain.ContactKeyEmail, *view.Contact)
	require.Equal(t, "person@example.com", *view.Value)
	require.Equal(t, linkID, links.query.LinkID)
	require.Equal(t, resumeID, resume.query.ResumeID)
	require.Equal(t, resumedomain.ContactKeyEmail, resume.query.Contact)
}

func TestContactViewThroughLinkReturnsAllContactsForDownload(t *testing.T) {
	linkID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	viewID := sharedvalueobjects.NewID()
	links := &contactViewLinkReaderStub{view: ResumeLinkPublicView{
		ID: linkID.String(), ResumeID: resumeID.String(), Status: string(domain.ResumeLinkStatusActive),
	}}
	phone, email := "+33123456789", "person@example.com"
	contacts := &contactViewResumeDownloadQueryStub{
		view: resumequeries.DownloadContactsView{Phone: &phone, Email: &email},
	}
	analytics := &contactViewAnalyticsRecorderStub{}
	handler := NewContactViewThroughLinkQueryHandlerWithAnalytics(
		links,
		&contactViewResumeQueryStub{},
		contacts,
		analytics,
		nil,
	)

	result, err := handler.Handle(context.Background(), ContactViewThroughLinkQuery{
		LinkID: linkID, Initiator: ContactViewInitiatorDownload, ViewID: viewID,
	})

	require.NoError(t, err)
	require.Nil(t, result.Contact)
	require.Nil(t, result.Value)
	require.NotNil(t, result.Contacts)
	require.Equal(t, []ContactViewThroughLinkContact{
		{Contact: resumedomain.ContactKeyEmail, Value: email},
		{Contact: resumedomain.ContactKeyPhone, Value: phone},
	}, *result.Contacts)
	require.Equal(t, resumeID, contacts.query.ResumeID)
	require.Empty(t, analytics.calls)
}

type contactViewAnalyticsRecorderStub struct {
	err       error
	calls     int
	initiator string
	viewID    sharedvalueobjects.ID
	linkID    sharedvalueobjects.ID
}

func (s *contactViewAnalyticsRecorderStub) RecordContactView(
	_ context.Context,
	viewID, linkID sharedvalueobjects.ID,
	_ string,
) error {
	s.calls++
	s.initiator = "contact"
	s.viewID = viewID
	s.linkID = linkID
	return s.err
}

func (s *contactViewAnalyticsRecorderStub) RecordDownload(
	_ context.Context,
	viewID, linkID sharedvalueobjects.ID,
) error {
	s.calls++
	s.initiator = "download"
	s.viewID = viewID
	s.linkID = linkID
	return s.err
}

func TestContactViewThroughLinkWithoutAnalyticsReturnsContact(t *testing.T) {
	linkID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	links := &contactViewLinkReaderStub{view: ResumeLinkPublicView{
		ID: linkID.String(), ResumeID: resumeID.String(), Status: string(domain.ResumeLinkStatusActive),
	}}
	resume := &contactViewResumeQueryStub{view: resumequeries.ContactView{
		Contact: resumedomain.ContactKeyEmail, Value: "person@example.com",
	}}
	handler := NewContactViewThroughLinkQueryHandler(links, resume)

	view, err := handler.Handle(context.Background(), ContactViewThroughLinkQuery{
		LinkID: linkID, Contact: resumedomain.ContactKeyEmail, ViewID: sharedvalueobjects.NewID(),
	})

	require.NoError(t, err)
	require.Equal(t, "person@example.com", *view.Value)
}

func TestContactViewThroughLinkRejectsMissingRequiredDependenciesWithoutPanic(t *testing.T) {
	query := ContactViewThroughLinkQuery{LinkID: sharedvalueobjects.NewID(), Contact: resumedomain.ContactKeyEmail}

	tests := []struct {
		name    string
		handler *ContactViewThroughLinkQueryHandler
	}{
		{name: "nil handler", handler: nil},
		{name: "nil link reader", handler: NewContactViewThroughLinkQueryHandler(nil, &contactViewResumeQueryStub{})},
		{name: "nil resume query", handler: NewContactViewThroughLinkQueryHandler(&contactViewLinkReaderStub{}, nil)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			require.NotPanics(t, func() {
				_, err = tt.handler.Handle(context.Background(), query)
			})
			require.ErrorContains(t, err, "dependencies are nil")
		})
	}
}

func TestContactViewThroughLinkReportsAnalyticsFailureWithoutFailingContactView(t *testing.T) {
	linkID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	viewID := sharedvalueobjects.NewID()
	analyticsErr := errors.New("analytics storage unavailable")
	links := &contactViewLinkReaderStub{view: ResumeLinkPublicView{
		ID: linkID.String(), ResumeID: resumeID.String(), Status: string(domain.ResumeLinkStatusActive),
	}}
	resume := &contactViewResumeQueryStub{view: resumequeries.ContactView{Contact: resumedomain.ContactKeyEmail, Value: "person@example.com"}}
	analytics := &contactViewAnalyticsRecorderStub{err: analyticsErr}
	var reportedErr error

	handler := NewContactViewThroughLinkQueryHandlerWithAnalytics(
		links,
		resume,
		nil,
		analytics,
		func(err error) { reportedErr = err },
	)
	view, err := handler.Handle(context.Background(), ContactViewThroughLinkQuery{
		LinkID: linkID, Contact: resumedomain.ContactKeyEmail, ViewID: viewID,
	})

	require.NoError(t, err)
	require.Equal(t, "person@example.com", *view.Value)
	require.Equal(t, 1, analytics.calls)
	require.ErrorIs(t, reportedErr, analyticsErr)
}

func TestContactViewThroughLinkRejectsInactiveLink(t *testing.T) {
	resume := &contactViewResumeQueryStub{}
	links := &contactViewLinkReaderStub{view: ResumeLinkPublicView{
		ResumeID: sharedvalueobjects.NewID().String(), Status: string(domain.ResumeLinkStatusInactive),
	}}

	_, err := NewContactViewThroughLinkQueryHandler(links, resume).Handle(context.Background(), ContactViewThroughLinkQuery{
		LinkID: sharedvalueobjects.NewID(), Contact: resumedomain.ContactKeyEmail,
	})

	require.ErrorIs(t, err, domain.ErrResumeLinkInactive)
	require.True(t, resume.query.ResumeID.IsZero())
}

func TestContactViewThroughLinkPropagatesMissingLink(t *testing.T) {
	resume := &contactViewResumeQueryStub{}
	links := &contactViewLinkReaderStub{err: ErrResumeLinkNotFound}

	_, err := NewContactViewThroughLinkQueryHandler(links, resume).Handle(context.Background(), ContactViewThroughLinkQuery{
		LinkID: sharedvalueobjects.NewID(), Contact: resumedomain.ContactKeyEmail,
	})

	require.ErrorIs(t, err, ErrResumeLinkNotFound)
	require.True(t, resume.query.ResumeID.IsZero())
}

func TestContactViewThroughLinkDoesNotRecordFailedContact(t *testing.T) {
	linkID := sharedvalueobjects.NewID()
	links := &contactViewLinkReaderStub{view: ResumeLinkPublicView{
		ResumeID: sharedvalueobjects.NewID().String(), Status: string(domain.ResumeLinkStatusActive),
	}}
	resume := &contactViewResumeQueryStub{err: resumequeries.ErrContactNotHidden}

	_, err := NewContactViewThroughLinkQueryHandler(links, resume).Handle(context.Background(), ContactViewThroughLinkQuery{
		LinkID: linkID, Contact: resumedomain.ContactKeyEmail,
	})

	require.ErrorIs(t, err, resumequeries.ErrContactNotHidden)
}

func TestContactViewThroughLinkPropagatesResumeContactRules(t *testing.T) {
	want := resumequeries.ErrContactNotHidden
	resume := &contactViewResumeQueryStub{err: want}
	links := &contactViewLinkReaderStub{view: ResumeLinkPublicView{
		ResumeID: sharedvalueobjects.NewID().String(), Status: string(domain.ResumeLinkStatusActive),
	}}

	_, err := NewContactViewThroughLinkQueryHandler(links, resume).Handle(context.Background(), ContactViewThroughLinkQuery{
		LinkID: sharedvalueobjects.NewID(), Contact: resumedomain.ContactKeyEmail,
	})

	require.ErrorIs(t, err, want)
}

func TestContactViewThroughLinkTranslatesDomainNotFound(t *testing.T) {
	links := &contactViewLinkReaderStub{err: domain.ErrEntityNotFound}
	handler := NewContactViewThroughLinkQueryHandler(links, &contactViewResumeQueryStub{})
	_, err := handler.Handle(context.Background(), ContactViewThroughLinkQuery{
		LinkID:  sharedvalueobjects.NewID(),
		Contact: resumedomain.ContactKeyEmail,
	})
	require.ErrorIs(t, err, ErrResumeLinkNotFound)
}

package queries

import (
	"context"
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

func TestContactViewThroughLinkResolvesActiveLink(t *testing.T) {
	linkID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	links := &contactViewLinkReaderStub{view: ResumeLinkPublicView{
		ID: linkID.String(), ResumeID: resumeID.String(), Status: string(domain.ResumeLinkStatusActive),
	}}
	resume := &contactViewResumeQueryStub{view: resumequeries.ContactView{Contact: resumedomain.ContactKeyEmail, Value: "person@example.com"}}
	events := &linkEventRecorderStub{}

	view, err := NewContactViewThroughLinkQueryHandler(links, resume, events).Handle(context.Background(), ContactViewThroughLinkQuery{
		LinkID:  linkID,
		Contact: resumedomain.ContactKeyEmail,
	})

	require.NoError(t, err)
	require.Equal(t, resumedomain.ContactKeyEmail, view.Contact)
	require.Equal(t, "person@example.com", view.Value)
	require.Equal(t, linkID, links.query.LinkID)
	require.Equal(t, resumeID, resume.query.ResumeID)
	require.Equal(t, resumedomain.ContactKeyEmail, resume.query.Contact)
	require.NotNil(t, events.event)
	require.Equal(t, linkID, events.event.LinkID())
	require.Equal(t, domain.EventTypeContactView, events.event.Type())
}

func TestContactViewThroughLinkIgnoresAnalyticsFailure(t *testing.T) {
	linkID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	links := &contactViewLinkReaderStub{view: ResumeLinkPublicView{
		ID: linkID.String(), ResumeID: resumeID.String(), Status: string(domain.ResumeLinkStatusActive),
	}}
	resume := &contactViewResumeQueryStub{view: resumequeries.ContactView{Contact: resumedomain.ContactKeyEmail, Value: "person@example.com"}}
	events := &linkEventRecorderStub{err: context.Canceled}

	view, err := NewContactViewThroughLinkQueryHandler(links, resume, events).Handle(context.Background(), ContactViewThroughLinkQuery{
		LinkID: linkID, Contact: resumedomain.ContactKeyEmail,
	})

	require.NoError(t, err)
	require.Equal(t, "person@example.com", view.Value)
	require.NotNil(t, events.event)
}

func TestContactViewThroughLinkRejectsInactiveLink(t *testing.T) {
	resume := &contactViewResumeQueryStub{}
	links := &contactViewLinkReaderStub{view: ResumeLinkPublicView{
		ResumeID: sharedvalueobjects.NewID().String(), Status: string(domain.ResumeLinkStatusInactive),
	}}

	_, err := NewContactViewThroughLinkQueryHandler(links, resume, &linkEventRecorderStub{}).Handle(context.Background(), ContactViewThroughLinkQuery{
		LinkID: sharedvalueobjects.NewID(), Contact: resumedomain.ContactKeyEmail,
	})

	require.ErrorIs(t, err, domain.ErrResumeLinkInactive)
	require.True(t, resume.query.ResumeID.IsZero())
}

func TestContactViewThroughLinkPropagatesMissingLink(t *testing.T) {
	resume := &contactViewResumeQueryStub{}
	links := &contactViewLinkReaderStub{err: ErrResumeLinkNotFound}

	_, err := NewContactViewThroughLinkQueryHandler(links, resume, &linkEventRecorderStub{}).Handle(context.Background(), ContactViewThroughLinkQuery{
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
	events := &linkEventRecorderStub{}

	_, err := NewContactViewThroughLinkQueryHandler(links, resume, events).Handle(context.Background(), ContactViewThroughLinkQuery{
		LinkID: linkID, Contact: resumedomain.ContactKeyEmail,
	})

	require.ErrorIs(t, err, resumequeries.ErrContactNotHidden)
	require.Nil(t, events.event)
}

func TestContactViewThroughLinkPropagatesResumeContactRules(t *testing.T) {
	want := resumequeries.ErrContactNotHidden
	resume := &contactViewResumeQueryStub{err: want}
	links := &contactViewLinkReaderStub{view: ResumeLinkPublicView{
		ResumeID: sharedvalueobjects.NewID().String(), Status: string(domain.ResumeLinkStatusActive),
	}}

	_, err := NewContactViewThroughLinkQueryHandler(links, resume, &linkEventRecorderStub{}).Handle(context.Background(), ContactViewThroughLinkQuery{
		LinkID: sharedvalueobjects.NewID(), Contact: resumedomain.ContactKeyEmail,
	})

	require.ErrorIs(t, err, want)
}

func TestContactViewThroughLinkTranslatesDomainNotFound(t *testing.T) {
	links := &contactViewLinkReaderStub{err: domain.ErrEntityNotFound}
	handler := NewContactViewThroughLinkQueryHandler(links, &contactViewResumeQueryStub{}, &linkEventRecorderStub{})
	_, err := handler.Handle(context.Background(), ContactViewThroughLinkQuery{LinkID: sharedvalueobjects.NewID(), Contact: resumedomain.ContactKeyEmail})
	require.ErrorIs(t, err, ErrResumeLinkNotFound)
}

package queries

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type publicLinkReaderStub struct {
	view  ResumeLinkPublicView
	err   error
	query PublicResumeLinkQuery
}

func (s *publicLinkReaderStub) GetByID(_ context.Context, query PublicResumeLinkQuery) (ResumeLinkPublicView, error) {
	s.query = query
	return s.view, s.err
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

func TestPublicResumeLinkQueryResolvesActiveLinkToResume(t *testing.T) {
	linkID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	links := &publicLinkReaderStub{view: ResumeLinkPublicView{ID: linkID.String(), ResumeID: resumeID.String(), Status: string(domain.ResumeLinkStatusActive)}}
	resume := &publicResumeQueryStub{view: resumequeries.PublicResumeView{Title: "CV"}}
	events := &linkEventRecorderStub{}

	view, err := NewPublicResumeLinkQueryHandler(links, resume, events).Handle(context.Background(), PublicResumeLinkQuery{LinkID: linkID})

	require.NoError(t, err)
	require.Equal(t, "CV", view.Title)
	require.Equal(t, linkID, links.query.LinkID)
	require.Equal(t, resumeID, resume.query.ID)
	require.NotNil(t, events.event)
	require.Equal(t, linkID, events.event.LinkID())
	require.Equal(t, domain.EventTypeView, events.event.Type())
}

func TestPublicResumeLinkQueryIgnoresAnalyticsFailure(t *testing.T) {
	linkID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	links := &publicLinkReaderStub{view: ResumeLinkPublicView{ID: linkID.String(), ResumeID: resumeID.String(), Status: string(domain.ResumeLinkStatusActive)}}
	resume := &publicResumeQueryStub{view: resumequeries.PublicResumeView{Title: "CV"}}
	events := &linkEventRecorderStub{err: context.Canceled}

	view, err := NewPublicResumeLinkQueryHandler(links, resume, events).Handle(context.Background(), PublicResumeLinkQuery{LinkID: linkID})

	require.NoError(t, err)
	require.Equal(t, "CV", view.Title)
	require.NotNil(t, events.event)
}

func TestPublicResumeLinkQueryDoesNotRecordInactiveLink(t *testing.T) {
	linkID := sharedvalueobjects.NewID()
	links := &publicLinkReaderStub{view: ResumeLinkPublicView{
		ResumeID: sharedvalueobjects.NewID().String(), Status: string(domain.ResumeLinkStatusInactive),
	}}
	resume := &publicResumeQueryStub{}
	events := &linkEventRecorderStub{}

	_, err := NewPublicResumeLinkQueryHandler(links, resume, events).Handle(context.Background(), PublicResumeLinkQuery{LinkID: linkID})

	require.ErrorIs(t, err, domain.ErrResumeLinkInactive)
	require.Nil(t, events.event)
}

func TestPublicResumeLinkQueryRejectsInactiveLink(t *testing.T) {
	resume := &publicResumeQueryStub{}
	links := &publicLinkReaderStub{view: ResumeLinkPublicView{ResumeID: sharedvalueobjects.NewID().String(), Status: string(domain.ResumeLinkStatusInactive)}}

	_, err := NewPublicResumeLinkQueryHandler(links, resume, &linkEventRecorderStub{}).Handle(context.Background(), PublicResumeLinkQuery{LinkID: sharedvalueobjects.NewID()})

	require.ErrorIs(t, err, domain.ErrResumeLinkInactive)
	require.True(t, resume.query.ID.IsZero())
}

func TestPublicResumeLinkQueryDoesNotRecordFailedLookup(t *testing.T) {
	links := &publicLinkReaderStub{err: ErrResumeLinkNotFound}
	resume := &publicResumeQueryStub{}
	events := &linkEventRecorderStub{}

	_, err := NewPublicResumeLinkQueryHandler(links, resume, events).Handle(context.Background(), PublicResumeLinkQuery{LinkID: sharedvalueobjects.NewID()})

	require.ErrorIs(t, err, ErrResumeLinkNotFound)
	require.Nil(t, events.event)
}

func TestPublicResumeLinkQueryPropagatesMissingLink(t *testing.T) {
	resume := &publicResumeQueryStub{}
	links := &publicLinkReaderStub{err: ErrResumeLinkNotFound}

	_, err := NewPublicResumeLinkQueryHandler(links, resume, &linkEventRecorderStub{}).Handle(context.Background(), PublicResumeLinkQuery{LinkID: sharedvalueobjects.NewID()})

	require.ErrorIs(t, err, ErrResumeLinkNotFound)
	require.True(t, resume.query.ID.IsZero())
}

func TestPublicResumeLinkQueryRejectsInvalidResolvedResumeID(t *testing.T) {
	links := &publicLinkReaderStub{view: ResumeLinkPublicView{
		ResumeID: "not-a-uuid",
		Status:   string(domain.ResumeLinkStatusActive),
	}}
	resume := &publicResumeQueryStub{}

	_, err := NewPublicResumeLinkQueryHandler(links, resume, &linkEventRecorderStub{}).Handle(
		context.Background(),
		PublicResumeLinkQuery{LinkID: sharedvalueobjects.NewID()},
	)

	require.Error(t, err)
	require.True(t, resume.query.ID.IsZero())
}

func TestPublicResumeLinkQueryTranslatesDomainNotFound(t *testing.T) {
	links := &publicLinkReaderStub{err: domain.ErrEntityNotFound}
	handler := NewPublicResumeLinkQueryHandler(links, &publicResumeQueryStub{}, &linkEventRecorderStub{})
	_, err := handler.Handle(context.Background(), PublicResumeLinkQuery{LinkID: sharedvalueobjects.NewID()})
	require.ErrorIs(t, err, ErrResumeLinkNotFound)
}

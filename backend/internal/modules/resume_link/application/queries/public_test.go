package queries

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/shared/domain/valueobjects"
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
	view  resumequeries.ResumeView
	err   error
	query resumequeries.PublicResumeQuery
}

func (s *publicResumeQueryStub) Handle(_ context.Context, query resumequeries.PublicResumeQuery) (resumequeries.ResumeView, error) {
	s.query = query
	return s.view, s.err
}

func TestPublicResumeLinkQueryResolvesActiveLinkToResume(t *testing.T) {
	linkID := valueobjects.NewID()
	resumeID := valueobjects.NewID()
	links := &publicLinkReaderStub{view: ResumeLinkPublicView{ID: linkID.String(), ResumeID: resumeID.String(), Status: domain.ResumeLinkStatusActive}}
	resume := &publicResumeQueryStub{view: resumequeries.ResumeView{ID: resumeID.String()}}
	events := &linkEventRecorderStub{}

	view, err := NewPublicResumeLinkQueryHandler(links, resume, events).Handle(context.Background(), PublicResumeLinkQuery{LinkID: linkID})

	require.NoError(t, err)
	require.Equal(t, resumeID.String(), view.ID)
	require.Equal(t, linkID, links.query.LinkID)
	require.Equal(t, resumeID, resume.query.ID)
	require.NotNil(t, events.event)
	require.Equal(t, linkID, events.event.LinkID())
	require.Equal(t, domain.EventTypeView, events.event.Type())
}

func TestPublicResumeLinkQueryIgnoresAnalyticsFailure(t *testing.T) {
	linkID := valueobjects.NewID()
	resumeID := valueobjects.NewID()
	links := &publicLinkReaderStub{view: ResumeLinkPublicView{ID: linkID.String(), ResumeID: resumeID.String(), Status: domain.ResumeLinkStatusActive}}
	resume := &publicResumeQueryStub{view: resumequeries.ResumeView{ID: resumeID.String()}}
	events := &linkEventRecorderStub{err: context.Canceled}

	view, err := NewPublicResumeLinkQueryHandler(links, resume, events).Handle(context.Background(), PublicResumeLinkQuery{LinkID: linkID})

	require.NoError(t, err)
	require.Equal(t, resumeID.String(), view.ID)
	require.NotNil(t, events.event)
}

func TestPublicResumeLinkQueryDoesNotRecordInactiveLink(t *testing.T) {
	linkID := valueobjects.NewID()
	links := &publicLinkReaderStub{view: ResumeLinkPublicView{
		ResumeID: valueobjects.NewID().String(), Status: domain.ResumeLinkStatusInactive,
	}}
	resume := &publicResumeQueryStub{}
	events := &linkEventRecorderStub{}

	_, err := NewPublicResumeLinkQueryHandler(links, resume, events).Handle(context.Background(), PublicResumeLinkQuery{LinkID: linkID})

	require.ErrorIs(t, err, domain.ErrResumeLinkInactive)
	require.Nil(t, events.event)
}

func TestPublicResumeLinkQueryRejectsInactiveLink(t *testing.T) {
	resume := &publicResumeQueryStub{}
	links := &publicLinkReaderStub{view: ResumeLinkPublicView{ResumeID: valueobjects.NewID().String(), Status: domain.ResumeLinkStatusInactive}}

	_, err := NewPublicResumeLinkQueryHandler(links, resume, &linkEventRecorderStub{}).Handle(context.Background(), PublicResumeLinkQuery{LinkID: valueobjects.NewID()})

	require.ErrorIs(t, err, domain.ErrResumeLinkInactive)
	require.True(t, resume.query.ID.IsZero())
}

func TestPublicResumeLinkQueryDoesNotRecordFailedLookup(t *testing.T) {
	links := &publicLinkReaderStub{err: domain.ErrEntityNotFound}
	resume := &publicResumeQueryStub{}
	events := &linkEventRecorderStub{}

	_, err := NewPublicResumeLinkQueryHandler(links, resume, events).Handle(context.Background(), PublicResumeLinkQuery{LinkID: valueobjects.NewID()})

	require.ErrorIs(t, err, domain.ErrEntityNotFound)
	require.Nil(t, events.event)
}

func TestPublicResumeLinkQueryPropagatesMissingLink(t *testing.T) {
	resume := &publicResumeQueryStub{}
	links := &publicLinkReaderStub{err: domain.ErrEntityNotFound}

	_, err := NewPublicResumeLinkQueryHandler(links, resume, &linkEventRecorderStub{}).Handle(context.Background(), PublicResumeLinkQuery{LinkID: valueobjects.NewID()})

	require.ErrorIs(t, err, domain.ErrEntityNotFound)
	require.True(t, resume.query.ID.IsZero())
}

func TestPublicResumeLinkQueryRejectsInvalidResolvedResumeID(t *testing.T) {
	links := &publicLinkReaderStub{view: ResumeLinkPublicView{
		ResumeID: "not-a-uuid",
		Status:   domain.ResumeLinkStatusActive,
	}}
	resume := &publicResumeQueryStub{}

	_, err := NewPublicResumeLinkQueryHandler(links, resume, &linkEventRecorderStub{}).Handle(
		context.Background(),
		PublicResumeLinkQuery{LinkID: valueobjects.NewID()},
	)

	require.Error(t, err)
	require.True(t, resume.query.ID.IsZero())
}

package queries

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumevalueobjects "developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume_link/domain/valueobjects"
	uservalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type statisticsLinkReaderStub struct {
	view ResumeLinkPublicView
	err  error
}

func (s statisticsLinkReaderStub) GetByID(context.Context, PublicResumeLinkQuery) (ResumeLinkPublicView, error) {
	return s.view, s.err
}

type statisticsOwnershipStub struct {
	err error
}

func (s statisticsOwnershipStub) GetOwned(context.Context, resumequeries.GetResumeQuery) (resumequeries.ResumeView, error) {
	if s.err != nil {
		return resumequeries.ResumeView{}, s.err
	}
	return resumequeries.ResumeView{ID: "resume"}, nil
}

type statisticsReaderStub struct {
	result ResumeLinkStatistics
	err    error
	query  GetResumeLinkStatisticsQuery
}

func (s *statisticsReaderStub) GetStatistics(_ context.Context, query GetResumeLinkStatisticsQuery) (ResumeLinkStatistics, error) {
	s.query = query
	return s.result, s.err
}

func TestGetResumeLinkStatisticsChecksOwnershipBeforeAggregation(t *testing.T) {
	userID := uservalueobjects.NewID()
	linkID := valueobjects.NewID()
	resumeID := resumevalueobjects.NewID()
	reader := &statisticsReaderStub{result: ResumeLinkStatistics{Views: 2, ContactViews: 1}}

	handler := NewGetResumeLinkStatisticsQueryHandler(
		statisticsLinkReaderStub{view: ResumeLinkPublicView{
			ID: linkID.String(), ResumeID: resumeID.String(), Status: "inactive",
		}},
		statisticsOwnershipStub{},
		reader,
	)

	result, err := handler.Handle(context.Background(), GetResumeLinkStatisticsQuery{
		UserID: userID,
		LinkID: linkID,
	})
	require.NoError(t, err)
	require.Equal(t, int64(2), result.Views)
	require.Equal(t, int64(1), result.ContactViews)
	require.Equal(t, linkID, reader.query.LinkID)
	require.Equal(t, userID, reader.query.UserID)
}

func TestGetResumeLinkStatisticsRejectsForeignResume(t *testing.T) {
	reader := &statisticsReaderStub{}
	handler := NewGetResumeLinkStatisticsQueryHandler(
		statisticsLinkReaderStub{view: ResumeLinkPublicView{
			ID: valueobjects.NewID().String(), ResumeID: resumevalueobjects.NewID().String(), Status: "active",
		}},
		statisticsOwnershipStub{err: resumequeries.ErrResumeNotFound},
		reader,
	)

	_, err := handler.Handle(context.Background(), GetResumeLinkStatisticsQuery{
		UserID: uservalueobjects.NewID(),
		LinkID: valueobjects.NewID(),
	})
	require.ErrorIs(t, err, resumequeries.ErrResumeNotFound)
	require.True(t, reader.query.LinkID.IsZero())
}

func TestGetResumeLinkStatisticsPropagatesMissingLink(t *testing.T) {
	reader := &statisticsReaderStub{}
	handler := NewGetResumeLinkStatisticsQueryHandler(
		statisticsLinkReaderStub{err: errors.New("missing")},
		statisticsOwnershipStub{},
		reader,
	)

	_, err := handler.Handle(context.Background(), GetResumeLinkStatisticsQuery{
		UserID: uservalueobjects.NewID(),
		LinkID: valueobjects.NewID(),
	})
	require.EqualError(t, err, "missing")
	require.True(t, reader.query.LinkID.IsZero())
}

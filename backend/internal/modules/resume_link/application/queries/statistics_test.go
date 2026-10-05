package queries

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
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

func (s *statisticsReaderStub) GetLinkStatistics(_ context.Context, query LinkAnalyticsQuery) (ResumeLinkStatistics, error) {
	s.query = query
	return s.result, s.err
}

func (s *statisticsReaderStub) GetDaily(context.Context, LinkAnalyticsQuery) (LinkDaily, error) {
	return LinkDaily{}, s.err
}
func (s *statisticsReaderStub) GetVisits(context.Context, LinkVisitsQuery) (LinkVisitsPage, error) {
	return LinkVisitsPage{}, s.err
}
func (s *statisticsReaderStub) GetSections(context.Context, LinkSectionsQuery) (SectionsStatistics, error) {
	return SectionsStatistics{}, s.err
}
func (s *statisticsReaderStub) GetTechnologies(context.Context, LinkAnalyticsQuery) (TechnologiesStatistics, error) {
	return TechnologiesStatistics{}, s.err
}
func (s *statisticsReaderStub) GetResumeStatistics(context.Context, ResumeAnalyticsQuery) (ResumeStatistics, error) {
	return ResumeStatistics{}, s.err
}

func TestLinkStatisticsNormalizeResponse(t *testing.T) {
	var result LinkStatistics
	result.NormalizeResponse()
	require.Equal(t, map[string]int64{"desktop": 0, "mobile": 0, "tablet": 0, "unknown": 0}, result.Devices)
	require.NotNil(t, result.Referrers)
	require.Empty(t, result.Referrers)
}

func TestGetResumeLinkStatisticsChecksOwnershipBeforeAggregation(t *testing.T) {
	userID := sharedvalueobjects.NewID()
	linkID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
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
			ID: sharedvalueobjects.NewID().String(), ResumeID: sharedvalueobjects.NewID().String(), Status: "active",
		}},
		statisticsOwnershipStub{err: resumequeries.ErrResumeNotFound},
		reader,
	)

	_, err := handler.Handle(context.Background(), GetResumeLinkStatisticsQuery{
		UserID: sharedvalueobjects.NewID(),
		LinkID: sharedvalueobjects.NewID(),
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
		UserID: sharedvalueobjects.NewID(),
		LinkID: sharedvalueobjects.NewID(),
	})
	require.EqualError(t, err, "missing")
	require.True(t, reader.query.LinkID.IsZero())
}

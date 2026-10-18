package queries

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestLinkStatisticsNormalizeResponse(t *testing.T) {
	var result LinkStatistics
	result.NormalizeResponse()
	require.Equal(t, map[string]int64{"desktop": 0, "mobile": 0, "tablet": 0, "unknown": 0}, result.Devices)
	require.NotNil(t, result.Referrers)
	require.Empty(t, result.Referrers)
}

func TestGetLinkStatisticsChecksOwnershipBeforeAggregation(t *testing.T) {
	userID := sharedvalueobjects.NewID()
	linkID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	reader := &statisticsReaderStub{result: LinkStatistics{Views: 2, ContactViews: 1}}

	access := &analyticsAccessReaderStub{resumeID: resumeID}
	handler := NewLinkStatisticsHandler(access, reader)

	result, err := handler.Handle(context.Background(), LinkAnalyticsQuery{
		UserID: userID,
		LinkID: linkID,
		Period: validAnalyticsTestPeriod(t),
	})
	require.NoError(t, err)
	require.Equal(t, int64(2), result.Views)
	require.Equal(t, int64(1), result.ContactViews)
	require.Equal(t, linkID, reader.query.LinkID)
	require.Equal(t, userID, reader.query.UserID)
}

func TestGetLinkStatisticsRejectsForeignResume(t *testing.T) {
	reader := &statisticsReaderStub{}
	expected := errors.New("resume not found")
	handler := NewLinkStatisticsHandler(&analyticsAccessReaderStub{err: expected}, reader)

	_, err := handler.Handle(context.Background(), LinkAnalyticsQuery{
		UserID: sharedvalueobjects.NewID(),
		LinkID: sharedvalueobjects.NewID(),
		Period: validAnalyticsTestPeriod(t),
	})
	require.ErrorIs(t, err, expected)
	require.True(t, reader.query.LinkID.IsZero())
}

func TestGetLinkStatisticsPropagatesMissingLink(t *testing.T) {
	reader := &statisticsReaderStub{}
	handler := NewLinkStatisticsHandler(&analyticsAccessReaderStub{err: errors.New("missing")}, reader)

	_, err := handler.Handle(context.Background(), LinkAnalyticsQuery{
		UserID: sharedvalueobjects.NewID(),
		LinkID: sharedvalueobjects.NewID(),
		Period: validAnalyticsTestPeriod(t),
	})
	require.EqualError(t, err, "missing")
	require.True(t, reader.query.LinkID.IsZero())
}

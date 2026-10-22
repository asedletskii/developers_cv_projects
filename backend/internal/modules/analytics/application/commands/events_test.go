package commands

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/analytics/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestIngestAnalyticsBatchDropsInvalidIndividualEvents(t *testing.T) {
	repo := &analyticsBatchRepositoryStub{}
	observer := &analyticsInputObserverStub{}
	h := NewIngestAnalyticsBatchHandler(repo, &unitOfWorkStub{}, &analyticsSubjectReaderStub{}, observer)
	linkID := sharedvalueobjects.NewID()
	viewID := sharedvalueobjects.NewID()

	err := h.Handle(context.Background(), AnalyticsBatchInput{
		ViewID: viewID, LinkID: linkID, Seq: 1,
		Events: []AnalyticsEventInput{
			{ClientEventID: uuid.New(), Type: domain.EventTypeDownload},
			{ClientEventID: uuid.New(), Type: domain.AnalyticsEventType("view")},
			{Type: domain.EventTypeDownload},
		},
	})
	require.NoError(t, err)
	require.True(t, repo.called)
	require.Equal(t, 1, repo.events)
	require.Equal(t, 1, observer.dropped["event:invalid_payload_or_type"])
	require.Equal(t, 1, observer.dropped["event:missing_client_id"])
	require.Equal(t, 1, observer.outcomes["accepted:success"])
}

func TestIngestAnalyticsBatchValidatesTechnologyAgainstViewSnapshot(t *testing.T) {
	repo := &analyticsBatchRepositoryStub{technologies: []string{"Go", "PostgreSQL"}}
	observer := &analyticsInputObserverStub{}
	h := NewIngestAnalyticsBatchHandler(repo, &unitOfWorkStub{}, &analyticsSubjectReaderStub{}, observer)

	err := h.Handle(context.Background(), AnalyticsBatchInput{
		ViewID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(), Seq: 1,
		Events: []AnalyticsEventInput{
			{ClientEventID: uuid.New(), Type: domain.EventTypeTechFilter, Payload: map[string]string{"technology": "go", "action": "select"}},
			{ClientEventID: uuid.New(), Type: domain.EventTypeTechFilter, Payload: map[string]string{"technology": "MadeUp", "action": "select"}},
		},
	})
	require.NoError(t, err)
	require.Equal(t, 1, repo.events)
	require.Equal(t, "Go", repo.received[0].Payload()["technology"])
	require.Equal(t, 1, observer.dropped["event:technology_not_in_resume"])
}

func TestIngestAnalyticsBatchTreatsFinalRetryAsIdempotent(t *testing.T) {
	repo := &analyticsBatchRepositoryStub{lastSeq: 5}
	observer := &analyticsInputObserverStub{}
	h := NewIngestAnalyticsBatchHandler(repo, &unitOfWorkStub{}, &analyticsSubjectReaderStub{}, observer)

	err := h.Handle(context.Background(), AnalyticsBatchInput{
		ViewID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(), Seq: 5, Final: true,
		Events: []AnalyticsEventInput{{ClientEventID: uuid.New(), Type: domain.EventTypeDownload}},
		Clicks: []AnalyticsClickInput{{Section: string(domain.AnalyticsSectionAbout), XPct: 100, YPct: 100, ViewportW: 800}},
	})

	require.NoError(t, err)
	require.True(t, repo.called)
	require.Zero(t, repo.events)
	require.Empty(t, repo.clicks)
	require.Equal(t, 1, observer.outcomes["ignored:same_sequence"])
}

func TestIngestAnalyticsBatchRejectsBatchLimits(t *testing.T) {
	observer := &analyticsInputObserverStub{}
	h := NewIngestAnalyticsBatchHandler(&analyticsBatchRepositoryStub{}, &unitOfWorkStub{}, &analyticsSubjectReaderStub{}, observer)
	input := AnalyticsBatchInput{ViewID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(), Seq: 1, Clicks: make([]AnalyticsClickInput, MaxAnalyticsClicksPerBatch+1)}
	require.ErrorIs(t, h.Handle(context.Background(), input), domain.ErrInvalidAnalyticsEvent)
	require.Equal(t, 1, observer.outcomes["rejected:batch_limit"])
}

func TestIngestAnalyticsBatchObservesPersistenceFailure(t *testing.T) {
	repo := &analyticsBatchRepositoryStub{err: errors.New("database unavailable")}
	observer := &analyticsInputObserverStub{}
	h := NewIngestAnalyticsBatchHandler(repo, &unitOfWorkStub{}, &analyticsSubjectReaderStub{}, observer)

	err := h.Handle(context.Background(), AnalyticsBatchInput{
		ViewID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(), Seq: 1,
	})

	require.Error(t, err)
	require.Equal(t, 1, observer.outcomes["failure:save_failed"])
}

func TestIngestAnalyticsBatchUsesClientTimestampOnlyForOrdering(t *testing.T) {
	repo := &analyticsBatchRepositoryStub{}
	handler := NewIngestAnalyticsBatchHandler(repo, &unitOfWorkStub{}, &analyticsSubjectReaderStub{}, nil)
	serverTime := time.Date(2026, 10, 18, 10, 0, 0, 0, time.UTC)
	handler.clock = func() time.Time { return serverTime }

	err := handler.Handle(context.Background(), AnalyticsBatchInput{
		ViewID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(), Seq: 1,
		Events: []AnalyticsEventInput{
			{ClientEventID: uuid.New(), Type: domain.EventTypeDownload, ClientTimestampMS: 200},
			{ClientEventID: uuid.New(), Type: domain.EventTypePrint, ClientTimestampMS: 100},
		},
		Clicks: []AnalyticsClickInput{
			{Section: string(domain.AnalyticsSectionAbout), XPct: 200, YPct: 300, ViewportW: 800, ClientTimestampMS: 200},
			{Section: string(domain.AnalyticsSectionAbout), XPct: 100, YPct: 200, ViewportW: 800, ClientTimestampMS: 100},
		},
	})

	require.NoError(t, err)
	require.Len(t, repo.received, 2)
	require.Equal(t, domain.EventTypePrint, repo.received[0].Type())
	require.Equal(t, serverTime, repo.received[0].OccurredAt())
	require.Equal(t, domain.EventTypeDownload, repo.received[1].Type())
	require.Equal(t, serverTime.Add(time.Nanosecond), repo.received[1].OccurredAt())
	require.Len(t, repo.clicks, 2)
	require.Equal(t, 100, repo.clicks[0].XPct())
	require.Equal(t, serverTime, repo.clicks[0].OccurredAt())
	require.Equal(t, 200, repo.clicks[1].XPct())
	require.Equal(t, serverTime.Add(time.Nanosecond), repo.clicks[1].OccurredAt())
}

package commands

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestIngestAnalyticsBatchDropsInvalidIndividualEvents(t *testing.T) {
	repo := &analyticsBatchRepositoryStub{}
	observer := &analyticsInputObserverStub{}
	h := NewIngestAnalyticsBatchHandler(repo, &unitOfWorkStub{}, observer)
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
	h := NewIngestAnalyticsBatchHandler(repo, &unitOfWorkStub{}, observer)

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

func TestIngestAnalyticsBatchRejectsBatchLimits(t *testing.T) {
	observer := &analyticsInputObserverStub{}
	h := NewIngestAnalyticsBatchHandler(&analyticsBatchRepositoryStub{}, &unitOfWorkStub{}, observer)
	input := AnalyticsBatchInput{ViewID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(), Seq: 1, Clicks: make([]AnalyticsClickInput, MaxAnalyticsClicksPerBatch+1)}
	require.ErrorIs(t, h.Handle(context.Background(), input), domain.ErrInvalidAnalyticsEvent)
	require.Equal(t, 1, observer.outcomes["rejected:batch_limit"])
}

func TestIngestAnalyticsBatchObservesPersistenceFailure(t *testing.T) {
	repo := &analyticsBatchRepositoryStub{err: errors.New("database unavailable")}
	observer := &analyticsInputObserverStub{}
	h := NewIngestAnalyticsBatchHandler(repo, &unitOfWorkStub{}, observer)

	err := h.Handle(context.Background(), AnalyticsBatchInput{
		ViewID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(), Seq: 1,
	})

	require.Error(t, err)
	require.Equal(t, 1, observer.outcomes["failure:save_failed"])
}

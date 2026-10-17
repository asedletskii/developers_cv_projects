package commands

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type analyticsBatchRepoStub struct {
	called       bool
	err          error
	events       int
	received     []*domain.AnalyticsEvent
	technologies []string
}

type analyticsInputObserverStub struct {
	dropped  map[string]int
	outcomes map[string]int
}

func (s *analyticsInputObserverStub) RecordAnalyticsBatchOutcome(outcome, reason string) {
	if s.outcomes == nil {
		s.outcomes = make(map[string]int)
	}
	s.outcomes[outcome+":"+reason]++
}

func (s *analyticsInputObserverStub) RecordDroppedAnalyticsInput(kind, reason string) {
	if s.dropped == nil {
		s.dropped = make(map[string]int)
	}
	s.dropped[kind+":"+reason]++
}

func (s *analyticsBatchRepoStub) LoadAnalyticsBatchForUpdate(_ context.Context, viewID, linkID sharedvalueobjects.ID, _, _, _ bool) (domain.AnalyticsBatchSnapshot, error) {
	layout, err := domain.NewLayoutSnapshot(
		[]resumedomain.SectionID{resumedomain.SectionAbout},
		map[resumedomain.SectionID]string{resumedomain.SectionAbout: "0123456789abcdef"},
		s.technologies,
	)
	if err != nil {
		return domain.AnalyticsBatchSnapshot{}, err
	}
	view, err := domain.RestoreView(domain.ViewSnapshot{
		ID: viewID, LinkID: linkID, VisitorID: sharedvalueobjects.NewID(), IPHash: make([]byte, 32),
		Device: domain.DeviceDesktop, Browser: domain.BrowserChrome, OS: domain.OSLinux,
		Traffic: domain.TrafficHuman, Layout: layout, StartedAt: time.Now().UTC(), CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		return domain.AnalyticsBatchSnapshot{}, err
	}
	return domain.AnalyticsBatchSnapshot{View: view, LinkStatus: domain.ResumeLinkStatusActive, ResumeStatus: resumedomain.ResumeStatusOpened}, nil
}

func (s *analyticsBatchRepoStub) SaveAnalyticsBatch(_ context.Context, _ *domain.View, events []*domain.AnalyticsEvent, _ []*domain.AnalyticsClick) error {
	s.called = true
	s.events = len(events)
	s.received = append([]*domain.AnalyticsEvent(nil), events...)
	return s.err
}

type analyticsBatchUOWStub struct{}

func (analyticsBatchUOWStub) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func TestIngestAnalyticsBatchDropsInvalidIndividualEvents(t *testing.T) {
	repo := &analyticsBatchRepoStub{}
	observer := &analyticsInputObserverStub{}
	h := NewIngestAnalyticsBatchHandler(repo, analyticsBatchUOWStub{}, observer)
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
	repo := &analyticsBatchRepoStub{technologies: []string{"Go", "PostgreSQL"}}
	observer := &analyticsInputObserverStub{}
	h := NewIngestAnalyticsBatchHandler(repo, analyticsBatchUOWStub{}, observer)

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
	h := NewIngestAnalyticsBatchHandler(&analyticsBatchRepoStub{}, analyticsBatchUOWStub{}, observer)
	input := AnalyticsBatchInput{ViewID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(), Seq: 1, Clicks: make([]AnalyticsClickInput, MaxAnalyticsClicksPerBatch+1)}
	require.ErrorIs(t, h.Handle(context.Background(), input), domain.ErrInvalidAnalyticsEvent)
	require.Equal(t, 1, observer.outcomes["rejected:batch_limit"])
}

func TestIngestAnalyticsBatchObservesPersistenceFailure(t *testing.T) {
	repo := &analyticsBatchRepoStub{err: errors.New("database unavailable")}
	observer := &analyticsInputObserverStub{}
	h := NewIngestAnalyticsBatchHandler(repo, analyticsBatchUOWStub{}, observer)

	err := h.Handle(context.Background(), AnalyticsBatchInput{
		ViewID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(), Seq: 1,
	})

	require.Error(t, err)
	require.Equal(t, 1, observer.outcomes["failure:save_failed"])
}

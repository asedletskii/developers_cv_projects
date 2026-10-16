package commands

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type analyticsBatchRepoStub struct {
	called bool
	err    error
	events int
}

func (s *analyticsBatchRepoStub) ApplyBatch(_ context.Context, _, _ sharedvalueobjects.ID, _ domain.ViewStateUpdate, events []*domain.AnalyticsEvent, _ []*domain.AnalyticsClick) error {
	s.called = true
	s.events = len(events)
	return s.err
}

type analyticsViewRepoStub struct{}

func (analyticsViewRepoStub) CreateView(context.Context, *domain.View) error { return nil }
func (analyticsViewRepoStub) UpdateViewState(context.Context, sharedvalueobjects.ID, domain.ViewStateUpdate) error {
	return nil
}
func (analyticsViewRepoStub) GetView(_ context.Context, id sharedvalueobjects.ID) (*domain.View, error) {
	layout, err := domain.NewLayoutSnapshot([]resumedomain.SectionID{resumedomain.SectionAbout}, map[resumedomain.SectionID]string{resumedomain.SectionAbout: "0123456789abcdef"})
	if err != nil {
		return nil, err
	}
	return domain.NewView(domain.ViewParams{ID: id, LinkID: sharedvalueobjects.NewID(), VisitorID: sharedvalueobjects.NewID(), IPHash: make([]byte, 32), Device: domain.DeviceDesktop, Browser: domain.BrowserChrome, OS: domain.OSLinux, Traffic: domain.TrafficHuman, Layout: layout, StartedAt: time.Now().UTC()})
}

type analyticsBatchUOWStub struct{}

func (analyticsBatchUOWStub) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

var _ ports.AnalyticsBatchRepository = (*analyticsBatchRepoStub)(nil)
var _ ports.UnitOfWork = analyticsBatchUOWStub{}

func TestIngestAnalyticsBatchDropsInvalidIndividualEvents(t *testing.T) {
	repo := &analyticsBatchRepoStub{}
	h := NewIngestAnalyticsBatchHandler(repo, analyticsViewRepoStub{}, analyticsBatchUOWStub{})
	linkID := sharedvalueobjects.NewID()
	viewID := sharedvalueobjects.NewID()

	err := h.Handle(context.Background(), AnalyticsBatchInput{
		ViewID: viewID, LinkID: linkID, Seq: 1,
		Events: []AnalyticsEventInput{
			{ClientEventID: uuid.New(), Type: domain.EventTypeDownload},
			{ClientEventID: uuid.New(), Type: domain.AnalyticsEventType("view")},
		},
	})
	require.NoError(t, err)
	require.True(t, repo.called)
	require.Equal(t, 1, repo.events)
}

func TestIngestAnalyticsBatchRejectsBatchLimits(t *testing.T) {
	h := NewIngestAnalyticsBatchHandler(&analyticsBatchRepoStub{}, analyticsViewRepoStub{}, analyticsBatchUOWStub{})
	input := AnalyticsBatchInput{ViewID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(), Seq: 1, Clicks: make([]AnalyticsClickInput, MaxAnalyticsClicksPerBatch+1)}
	require.ErrorIs(t, h.Handle(context.Background(), input), domain.ErrInvalidAnalyticsEvent)
}

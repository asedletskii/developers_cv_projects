package domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/analytics/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func newPolicyEvent(t *testing.T, eventType domain.AnalyticsEventType, occurredAt time.Time) *domain.AnalyticsEvent {
	t.Helper()
	clientEventID := sharedvalueobjects.NewID()
	event, err := domain.NewAnalyticsEvent(domain.AnalyticsEventParams{
		ViewID: sharedvalueobjects.NewID(),
		LinkID: sharedvalueobjects.NewID(), ClientEventID: &clientEventID,
		Type: eventType, OccurredAt: occurredAt,
	})
	require.NoError(t, err)
	return event
}

func newPolicyClick(t *testing.T, occurredAt time.Time) *domain.AnalyticsClick {
	t.Helper()
	click, err := domain.NewAnalyticsClick(domain.AnalyticsClickParams{
		ViewID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(),
		Section: domain.AnalyticsSectionAbout, XPct: 100, YPct: 200, ViewportW: 1280,
		LayoutHash: "0123456789abcdef", Target: "button", OccurredAt: occurredAt,
	})
	require.NoError(t, err)
	return click
}

func TestSelectAnalyticsEventsForStorageAppliesPerViewLimitBeforePrintDeduplication(t *testing.T) {
	now := time.Date(2026, 10, 15, 12, 0, 0, 0, time.UTC)
	firstPrint := newPolicyEvent(t, domain.EventTypePrint, now)
	secondDownload := newPolicyEvent(t, domain.EventTypeDownload, now.Add(time.Second))

	selected := domain.SelectAnalyticsEventsForStorage(
		[]*domain.AnalyticsEvent{firstPrint, secondDownload},
		domain.MaxAnalyticsEventsPerView-1,
		nil,
	)
	require.Len(t, selected, 1)
	require.Same(t, firstPrint, selected[0])
}

func TestSelectAnalyticsEventsForStorageDeduplicatesPrintsWithinAndAcrossBatches(t *testing.T) {
	now := time.Date(2026, 10, 15, 12, 0, 0, 0, time.UTC)
	firstPrint := newPolicyEvent(t, domain.EventTypePrint, now)
	secondPrint := newPolicyEvent(t, domain.EventTypePrint, now.Add(time.Second))
	download := newPolicyEvent(t, domain.EventTypeDownload, now.Add(2*time.Second))

	selected := domain.SelectAnalyticsEventsForStorage(
		[]*domain.AnalyticsEvent{firstPrint, secondPrint, download}, 0, nil,
	)
	require.Len(t, selected, 2)
	require.Same(t, firstPrint, selected[0])
	require.Same(t, download, selected[1])

	selected = domain.SelectAnalyticsEventsForStorage(
		[]*domain.AnalyticsEvent{secondPrint, download}, 0, []time.Time{now},
	)
	require.Len(t, selected, 1)
	require.Same(t, download, selected[0])
}

func TestSelectAnalyticsEventsForStorageHonorsCapacityAndEmptyRemainingCapacity(t *testing.T) {
	now := time.Now().UTC()
	first := newPolicyEvent(t, domain.EventTypeDownload, now)
	second := newPolicyEvent(t, domain.EventTypeDownload, now.Add(time.Second))

	selected := domain.SelectAnalyticsEventsForStorage(
		[]*domain.AnalyticsEvent{first, second}, domain.MaxAnalyticsEventsPerView-1, nil,
	)
	require.Equal(t, []*domain.AnalyticsEvent{first}, selected)
	require.Empty(t, domain.SelectAnalyticsEventsForStorage(
		[]*domain.AnalyticsEvent{first}, domain.MaxAnalyticsEventsPerView, nil,
	))
}

func TestSelectAnalyticsClicksForStorageHonorsPerViewLimit(t *testing.T) {
	now := time.Now().UTC()
	first := newPolicyClick(t, now)
	second := newPolicyClick(t, now.Add(time.Second))
	third := newPolicyClick(t, now.Add(2*time.Second))
	clicks := []*domain.AnalyticsClick{first, second, third}
	selected := domain.SelectAnalyticsClicksForStorage(clicks, domain.MaxAnalyticsClicksPerView-2)
	require.Len(t, selected, 2)
	require.Same(t, first, selected[0])
	require.Same(t, second, selected[1])
	require.Empty(t, domain.SelectAnalyticsClicksForStorage(clicks, domain.MaxAnalyticsClicksPerView))
}

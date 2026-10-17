package domain

import (
	"time"

	resumedomain "developers-cv/backend/internal/modules/resume/domain"
)

type AnalyticsBatchSnapshot struct {
	View             *View
	LinkStatus       ResumeLinkStatus
	ResumeStatus     resumedomain.ResumeStatus
	StoredEventCount int64
	StoredClickCount int64
	ExistingPrints   []time.Time
}

const MaxAnalyticsClicksPerView = 1000

func SelectAnalyticsEventsForStorage(
	events []*AnalyticsEvent,
	storedCount int64,
	existingPrints []time.Time,
) []*AnalyticsEvent {
	if storedCount < 0 || storedCount >= MaxAnalyticsEventsPerView || len(events) == 0 {
		return nil
	}

	remaining := int64(MaxAnalyticsEventsPerView) - storedCount
	if int64(len(events)) > remaining {
		events = events[:int(remaining)]
	}

	selected := make([]*AnalyticsEvent, 0, len(events))
	printSeenInBatch := false
	for _, event := range events {
		if event == nil {
			continue
		}
		if event.Type() == EventTypePrint {
			if printSeenInBatch {
				continue
			}
			printSeenInBatch = true
			if hasRecentStoredPrint(event.OccurredAt(), existingPrints) {
				continue
			}
		}
		selected = append(selected, event)
	}
	return selected
}

func SelectAnalyticsClicksForStorage(clicks []*AnalyticsClick, storedCount int64) []*AnalyticsClick {
	if storedCount < 0 || storedCount >= MaxAnalyticsClicksPerView || len(clicks) == 0 {
		return nil
	}
	remaining := int64(MaxAnalyticsClicksPerView) - storedCount
	if int64(len(clicks)) > remaining {
		clicks = clicks[:int(remaining)]
	}
	return clicks
}

func hasRecentStoredPrint(occurredAt time.Time, existingPrints []time.Time) bool {
	threshold := occurredAt.Add(-2 * time.Second)
	for _, existing := range existingPrints {
		if !existing.Before(threshold) {
			return true
		}
	}
	return false
}

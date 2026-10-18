package ports

import (
	"context"
	"time"

	"developers-cv/backend/internal/modules/analytics/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type LinkAnalyticsContext struct {
	LinkID          sharedvalueobjects.ID
	ResumeID        sharedvalueobjects.ID
	LinkActive      bool
	ResumePublished bool
}

type AnalyticsBatchSnapshot struct {
	View             *domain.View
	LinkContext      LinkAnalyticsContext
	StoredEventCount int64
	StoredClickCount int64
	ExistingPrints   []time.Time
}

type ViewRepository interface {
	CreateView(context.Context, *domain.View) error
	GetView(context.Context, sharedvalueobjects.ID) (*domain.View, error)
	UpdateViewState(context.Context, sharedvalueobjects.ID, domain.ViewStateUpdate) error
}

type EventRepository interface {
	RecordEvents(context.Context, []*domain.AnalyticsEvent) error
}

type AnalyticsBatchRepository interface {
	LoadAnalyticsBatchForUpdate(context.Context, sharedvalueobjects.ID, sharedvalueobjects.ID, bool, bool, bool) (AnalyticsBatchSnapshot, error)
	SaveAnalyticsBatch(context.Context, *domain.View, []*domain.AnalyticsEvent, []*domain.AnalyticsClick) error
}

type AnalyticsObserver interface {
	RecordDroppedAnalyticsInput(kind, reason string)
	RecordAnalyticsBatchOutcome(outcome, reason string)
}

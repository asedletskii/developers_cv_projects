package ports

import (
	"context"

	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ViewRepository interface {
	CreateView(context.Context, *domain.View) error
	GetView(context.Context, sharedvalueobjects.ID) (*domain.View, error)
	UpdateViewState(context.Context, sharedvalueobjects.ID, domain.ViewStateUpdate) error
}

type AnalyticsEventRepository interface {
	RecordEvents(context.Context, []*domain.AnalyticsEvent) error
}

type ClickRepository interface {
	RecordClicks(context.Context, []*domain.AnalyticsClick) error
}

type AnalyticsBatchRepository interface {
	LoadAnalyticsBatchForUpdate(context.Context, sharedvalueobjects.ID, sharedvalueobjects.ID, bool, bool, bool) (domain.AnalyticsBatchSnapshot, error)
	SaveAnalyticsBatch(context.Context, *domain.View, []*domain.AnalyticsEvent, []*domain.AnalyticsClick) error
}

type AnalyticsInputObserver interface {
	RecordDroppedAnalyticsInput(kind, reason string)
}

type AnalyticsBatchOutcomeObserver interface {
	RecordAnalyticsBatchOutcome(outcome, reason string)
}

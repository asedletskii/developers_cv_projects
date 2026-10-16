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
	ApplyBatch(context.Context, sharedvalueobjects.ID, sharedvalueobjects.ID, domain.ViewStateUpdate, []*domain.AnalyticsEvent, []*domain.AnalyticsClick) error
}

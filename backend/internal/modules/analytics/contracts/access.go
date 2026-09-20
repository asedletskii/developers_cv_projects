package contracts

import (
	"context"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type AnalyticsAccessReader interface {
	ResolveLink(context.Context, sharedvalueobjects.ID) (sharedvalueobjects.ID, error)
	EnsureResumeOwned(context.Context, sharedvalueobjects.ID, sharedvalueobjects.ID) error
	ResumeSectionOrder(context.Context, sharedvalueobjects.ID, sharedvalueobjects.ID) ([]string, error)
	ResumeAnalyticsLayout(context.Context, sharedvalueobjects.ID, sharedvalueobjects.ID) (ResumeAnalyticsLayout, error)
}

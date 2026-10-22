package ports

import (
	"context"
	"time"

	analyticscontracts "developers-cv/backend/internal/modules/analytics/contracts"
	"developers-cv/backend/internal/modules/analytics/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type VisitorIdentity interface {
	Identify(now time.Time, input analyticscontracts.PublicVisitInput) (analyticscontracts.VisitorIdentityResult, error)
}

type TrafficClassifier interface {
	Classify(context.Context, TrafficClassificationInput) (TrafficClassificationResult, error)
}

type TrafficClassificationInput struct {
	LinkID    sharedvalueobjects.ID
	OwnerID   sharedvalueobjects.ID
	UserAgent string
	IPHash    []byte
	Now       time.Time
}

type TrafficClassificationResult struct {
	Traffic domain.TrafficType
	Reason  string
	Device  domain.DeviceType
	Browser domain.BrowserType
	OS      domain.OSType
}

type LayoutSnapshotBuilder interface {
	Build(analyticscontracts.ResumeAnalyticsLayout) (domain.LayoutSnapshot, error)
}

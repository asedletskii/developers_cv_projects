package ports

import (
	"context"
	"net/http"
	"time"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type PublicVisitInput struct {
	IP        string
	UserAgent string
	Referrer  string
	VisitorID sharedvalueobjects.ID
	ResumeID  string
	DNT       bool
	GPC       bool
	OwnerID   sharedvalueobjects.ID
	HasCookie bool
}

type VisitorIdentity interface {
	Identify(now time.Time, input PublicVisitInput) (VisitorIdentityResult, error)
}

type VisitorIdentityResult struct {
	ViewID    sharedvalueobjects.ID
	VisitorID sharedvalueobjects.ID
	IPHash    []byte
	SetCookie bool
	Cookie    *http.Cookie
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
	Build(resumequeries.PublicResumeView) (domain.LayoutSnapshot, error)
}

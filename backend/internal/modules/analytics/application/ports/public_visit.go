package ports

import (
	"context"
	"net/http"
	"time"

	"developers-cv/backend/internal/modules/analytics/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

const (
	ProductionVisitorCookieName = "__Secure-cv_vid"
	LocalVisitorCookieName      = "cv_vid"
)

type PublicVisitInput struct {
	IP        string
	UserAgent string
	Referrer  string
	VisitorID sharedvalueobjects.ID
	OwnerID   sharedvalueobjects.ID
	HasCookie bool
}

type CreatePublicVisitInput struct {
	LinkID       sharedvalueobjects.ID
	ResumeID     sharedvalueobjects.ID
	ResumeLayout ResumeAnalyticsLayout
	Visitor      PublicVisitInput
}

type ResumeAnalyticsLayout struct {
	Sections     []ResumeAnalyticsSection
	Technologies []string
}

type ResumeAnalyticsSection struct {
	ID   string
	Hash string
}

type PublicVisitCreator interface {
	Handle(context.Context, CreatePublicVisitInput) (*domain.View, VisitorIdentityResult, error)
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
	Build(ResumeAnalyticsLayout) (domain.LayoutSnapshot, error)
}

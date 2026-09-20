package contracts

import (
	"context"
	"net/http"

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

type ResumeAnalyticsSection struct {
	ID   string
	Hash string
}

type ResumeAnalyticsLayout struct {
	Sections     []ResumeAnalyticsSection
	Technologies []string
}

type CreatePublicVisitInput struct {
	LinkID       sharedvalueobjects.ID
	ResumeID     sharedvalueobjects.ID
	ResumeLayout ResumeAnalyticsLayout
	Visitor      PublicVisitInput
}

type VisitorIdentityResult struct {
	ViewID    sharedvalueobjects.ID
	VisitorID sharedvalueobjects.ID
	IPHash    []byte
	SetCookie bool
	Cookie    *http.Cookie
}

type PublicVisitCreator interface {
	Handle(context.Context, CreatePublicVisitInput) (VisitorIdentityResult, error)
}

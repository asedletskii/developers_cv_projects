package contracts

import (
	"context"
	"time"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type LinkLifecycle struct {
	LinkID          sharedvalueobjects.ID
	ResumeID        sharedvalueobjects.ID
	LinkActive      bool
	ResumePublished bool
}

type ResumeLinkRef struct {
	ID        sharedvalueobjects.ID
	Name      string
	Status    string
	CreatedAt time.Time
}

type AnalyticsSubjectReader interface {
	LoadLinkLifecycle(context.Context, sharedvalueobjects.ID) (LinkLifecycle, error)
	ListResumeLinks(context.Context, sharedvalueobjects.ID) ([]ResumeLinkRef, error)
}

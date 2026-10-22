package contracts

import (
	"context"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type OwnedResumeSectionsQuery struct {
	ResumeID sharedvalueobjects.ID
	UserID   sharedvalueobjects.ID
}

type ResumeSectionSettings struct {
	ID        string
	Enabled   bool
	Collapsed bool
}

type OwnedResumeSectionsQueryHandler interface {
	Handle(context.Context, OwnedResumeSectionsQuery) ([]ResumeSectionSettings, error)
}

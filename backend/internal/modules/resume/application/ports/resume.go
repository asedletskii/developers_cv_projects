package ports

import (
	"context"

	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type GetResumeParams struct {
	ResumeID sharedvalueobjects.ID
	UserID   sharedvalueobjects.ID
}

type DeleteResumeParams struct {
	ResumeID sharedvalueobjects.ID
	UserID   sharedvalueobjects.ID
	Version  int64
}

type ResumeRepository interface {
	Get(context.Context, GetResumeParams) (*domain.Resume, error)
	Save(context.Context, *domain.Resume) error
	Create(context.Context, *domain.Resume) error
	Delete(context.Context, DeleteResumeParams) error
}

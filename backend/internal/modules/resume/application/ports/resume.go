package ports

import (
	"context"

	"developers-cv/backend/internal/modules/resume/domain"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type GetResumeParams struct {
	ResumeID identity.ID
	UserID   identity.ID
}

type DeleteResumeParams struct {
	ResumeID identity.ID
	UserID   identity.ID
	Version  int64
}

type ResumeRepository interface {
	Get(context.Context, GetResumeParams) (*domain.Resume, error)
	Save(context.Context, *domain.Resume) error
	Create(context.Context, *domain.Resume) error
	Delete(context.Context, DeleteResumeParams) error
}

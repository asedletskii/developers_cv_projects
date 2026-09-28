package ports

import (
	"context"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type GetResumeParams struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}

type DeleteResumeParams struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}

type ResumeRepository interface {
	Get(context.Context, GetResumeParams) (*domain.Resume, error)
	Save(context.Context, *domain.Resume) error
	Create(context.Context, *domain.Resume) error
	Delete(context.Context, DeleteResumeParams) error
}

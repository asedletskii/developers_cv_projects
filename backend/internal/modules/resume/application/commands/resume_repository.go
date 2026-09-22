package commands

import (
	"context"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

type ResumeRepository interface {
	GetByID(ctx context.Context, id valueobjects.ID) (*domain.Resume, error)
	Save(ctx context.Context, resume *domain.Resume) error
}

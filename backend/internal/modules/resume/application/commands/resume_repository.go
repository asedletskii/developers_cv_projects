package commands

import (
	"context"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type ResumeRepository interface {
	GetByID(ctx context.Context, userID useridentity.ID, id valueobjects.ID) (*domain.Resume, error)
	Save(ctx context.Context, resume *domain.Resume) error
}

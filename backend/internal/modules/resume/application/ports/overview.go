package ports

import (
	"context"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type CreateOverviewRepositoryCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
	Overview *domain.ResumeOverview
}
type GetOverviewRepositoryCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type UpdateOverviewRepositoryCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
	Overview *domain.ResumeOverview
}
type CreateOverviewRepository interface {
	CreateOverview(context.Context, CreateOverviewRepositoryCommand) error
}
type GetOverviewRepository interface {
	GetOverview(context.Context, GetOverviewRepositoryCommand) (*domain.ResumeOverview, error)
}
type UpdateOverviewRepository interface {
	UpdateOverview(context.Context, UpdateOverviewRepositoryCommand) error
}

package ports

import (
	"context"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type CreateEducationRepositoryCommand struct {
	Education *domain.Education
	ResumeID  valueobjects.ID
	UserID    useridentity.ID
}

type GetEducationRepositoryCommand struct {
	ResumeID, EducationID valueobjects.ID
	UserID                useridentity.ID
}

type UpdateEducationRepositoryCommand struct {
	Education *domain.Education
	ResumeID  valueobjects.ID
	UserID    useridentity.ID
}

type DeleteEducationRepositoryCommand struct {
	ResumeID, EducationID valueobjects.ID
	UserID                useridentity.ID
}

type CreateEducationRepository interface {
	CreateEducation(context.Context, CreateEducationRepositoryCommand) error
}

type GetEducationRepository interface {
	GetEducation(context.Context, GetEducationRepositoryCommand) (*domain.Education, error)
}

type UpdateEducationRepository interface {
	UpdateEducation(context.Context, UpdateEducationRepositoryCommand) error
}

type DeleteEducationRepository interface {
	DeleteEducation(context.Context, DeleteEducationRepositoryCommand) error
}

package ports

import (
	"context"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type CreateSkillsRepositoryCommand struct {
	Skills *domain.Skills
	UserID useridentity.ID
}

type GetSkillsRepositoryCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}

type UpdateSkillsRepositoryCommand struct {
	Skills *domain.Skills
	UserID useridentity.ID
}

type CreateSkillsRepository interface {
	CreateSkills(context.Context, CreateSkillsRepositoryCommand) error
}

type GetSkillsRepository interface {
	GetSkills(context.Context, GetSkillsRepositoryCommand) (*domain.Skills, error)
}

type UpdateSkillsRepository interface {
	UpdateSkills(context.Context, UpdateSkillsRepositoryCommand) error
}

package commands

import (
	"context"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

type PersonRepository interface {
	GetByID(ctx context.Context, id valueobjects.ID) (*domain.Person, error)
	Save(ctx context.Context, person *domain.Person) error
}

package application

import (
	"context"

	"developers-cv/backend/internal/modules/user/domain"
	"developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type UserRepository interface {
	GetByID(ctx context.Context, id valueobjects.ID) (*domain.User, error)
	Save(ctx context.Context, user *domain.User) error
}

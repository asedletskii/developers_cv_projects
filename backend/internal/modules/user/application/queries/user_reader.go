package queries

import (
	"context"

	"developers-cv/backend/internal/modules/user/domain"
	"developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type UserReader interface {
	GetByID(ctx context.Context, id valueobjects.ID) (*domain.User, error)
}

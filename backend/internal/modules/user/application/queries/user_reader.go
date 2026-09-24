package queries

import (
	"context"

	"developers-cv/backend/internal/modules/user/domain"
	"developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type GetUserQuery struct {
	ID valueobjects.ID
}

type UserReader interface {
	GetUser(ctx context.Context, query GetUserQuery) (*domain.User, error)
}

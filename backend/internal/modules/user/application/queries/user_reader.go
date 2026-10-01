package queries

import (
	"context"

	"developers-cv/backend/internal/modules/user/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type GetUserQuery struct {
	ID sharedvalueobjects.ID
}

type UserReader interface {
	GetUser(ctx context.Context, query GetUserQuery) (*domain.User, error)
}

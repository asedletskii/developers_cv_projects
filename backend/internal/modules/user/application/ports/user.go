package ports

import (
	"context"

	"developers-cv/backend/internal/modules/user/domain"
)

type UserRepository interface {
	Create(context.Context, *domain.User) error
}

type ProfileRepository interface {
	ReplaceProfile(context.Context, *domain.Profile) error
}

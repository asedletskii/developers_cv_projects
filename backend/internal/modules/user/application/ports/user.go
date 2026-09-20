package ports

import (
	"context"

	"developers-cv/backend/internal/modules/user/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type UserRepository interface {
	Create(context.Context, *domain.User) error
	Exists(context.Context, sharedvalueobjects.ID) (bool, error)
}

type UserDeleteRepository interface {
	LockForDelete(context.Context, sharedvalueobjects.ID) (bool, error)
	Delete(context.Context, sharedvalueobjects.ID) error
}

type ProfileRepository interface {
	ReplaceProfile(context.Context, *domain.Profile) error
}

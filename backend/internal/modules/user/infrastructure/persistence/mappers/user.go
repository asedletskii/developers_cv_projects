package mappers

import (
	"fmt"

	"developers-cv/backend/internal/modules/user/domain"
	"developers-cv/backend/internal/modules/user/infrastructure/persistence/gormmodel"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func UserToDB(user *domain.User) (gormmodel.User, error) {
	if user == nil || user.ID().IsZero() {
		return gormmodel.User{}, fmt.Errorf("user identity is required")
	}
	return gormmodel.User{
		ID:        user.ID().UUID(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}, nil
}

func UserToDomain(model gormmodel.User) (*domain.User, error) {
	id, err := sharedvalueobjects.IDFromUUID(model.ID)
	if err != nil {
		return nil, fmt.Errorf("map user id: %w", err)
	}
	return domain.NewUser(domain.UserParams{
		ID:        id,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	})
}

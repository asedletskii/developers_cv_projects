package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/user/domain"
	"developers-cv/backend/internal/modules/user/domain/valueobjects"
	"developers-cv/backend/internal/modules/user/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/platform/database"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(ctx context.Context, id valueobjects.ID) (*domain.User, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("get user: repository dependencies are nil")
	}
	if id.IsZero() {
		return nil, fmt.Errorf("get user: %w", domain.ErrInvalidUser)
	}

	var model gormmodel.User
	db := database.DBFromContext(ctx, r.db)
	if err := db.WithContext(ctx).First(&model, "id = ?", id.UUID()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, fmt.Errorf("get user: %w", err)
	}
	return fromModel(model)
}

func (r *UserRepository) Save(ctx context.Context, user *domain.User) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("save user: repository dependencies are nil")
	}
	if user == nil {
		return fmt.Errorf("save user: %w", domain.ErrInvalidUser)
	}

	db := database.DBFromContext(ctx, r.db)
	model := toModel(user)
	if err := db.WithContext(ctx).Save(&model).Error; err != nil {
		return fmt.Errorf("save user: %w", err)
	}
	return nil
}

func toModel(user *domain.User) gormmodel.User {
	return gormmodel.User{
		ID:        user.ID().UUID(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}
}

func fromModel(model gormmodel.User) (*domain.User, error) {
	id, err := valueobjects.IDFromUUID(model.ID)
	if err != nil {
		return nil, fmt.Errorf("map user id: %w", err)
	}
	return domain.NewUser(domain.UserParams{
		ID:        id,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	})
}

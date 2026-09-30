package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/user/domain"
	"developers-cv/backend/internal/modules/user/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/platform/database"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("create user: repository dependencies are nil")
	}
	if user == nil || user.ID().IsZero() {
		return fmt.Errorf("create user: invalid user")
	}

	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	model := gormmodel.User{
		ID:        user.ID().UUID(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}
	if err := tx.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

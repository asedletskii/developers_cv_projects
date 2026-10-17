package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/user/domain"
	"developers-cv/backend/internal/modules/user/infrastructure/persistence/mappers"
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
	if user == nil {
		return fmt.Errorf("create user: user is nil")
	}

	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	model, err := mappers.UserToDB(user)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	if err := tx.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("create user: %w", database.MapPostgresError(err))
	}
	return nil
}

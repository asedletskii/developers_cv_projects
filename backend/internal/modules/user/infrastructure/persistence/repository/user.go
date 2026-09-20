package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"developers-cv/backend/internal/modules/user/domain"
	"developers-cv/backend/internal/modules/user/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/modules/user/infrastructure/persistence/mappers"
	"developers-cv/backend/internal/platform/database"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
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

func (r *UserRepository) Exists(ctx context.Context, userID sharedvalueobjects.ID) (bool, error) {
	if r == nil || r.db == nil {
		return false, fmt.Errorf("check user existence: repository dependencies are nil")
	}
	if userID.IsZero() {
		return false, nil
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return false, fmt.Errorf("check user existence: %w", err)
	}
	var count int64
	if err := tx.WithContext(ctx).Model(&gormmodel.User{}).Where("id = ?", userID.UUID()).Count(&count).Error; err != nil {
		return false, fmt.Errorf("check user existence: %w", database.MapPostgresError(err))
	}
	return count > 0, nil
}

func (r *UserRepository) LockForDelete(ctx context.Context, userID sharedvalueobjects.ID) (bool, error) {
	if r == nil || r.db == nil {
		return false, fmt.Errorf("lock user for delete: repository dependencies are nil")
	}
	if userID.IsZero() {
		return false, nil
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return false, fmt.Errorf("lock user for delete: %w", err)
	}
	var model gormmodel.User
	if err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Select("id").
		Where("id = ?", userID.UUID()).
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("lock user for delete: %w", database.MapPostgresError(err))
	}
	return true, nil
}

func (r *UserRepository) Delete(ctx context.Context, userID sharedvalueobjects.ID) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("delete user: repository dependencies are nil")
	}
	if userID.IsZero() {
		return fmt.Errorf("delete user: user id is empty")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	if err := tx.WithContext(ctx).Where("id = ?", userID.UUID()).Delete(&gormmodel.User{}).Error; err != nil {
		return fmt.Errorf("delete user: %w", database.MapPostgresError(err))
	}
	return nil
}

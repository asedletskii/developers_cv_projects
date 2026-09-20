package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/auth/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/platform/database"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type UserDataCleanupRepository struct {
	db *gorm.DB
}

func NewUserDataCleanupRepository(db *gorm.DB) *UserDataCleanupRepository {
	return &UserDataCleanupRepository{db: db}
}

func (r *UserDataCleanupRepository) DeleteByUserID(ctx context.Context, userID sharedvalueobjects.ID) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("delete user auth data: repository dependencies are nil")
	}
	if userID.IsZero() {
		return fmt.Errorf("delete user auth data: user id is empty")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("delete user auth data: %w", err)
	}
	tx = tx.WithContext(ctx)
	if err := tx.Where("user_id = ?", userID.UUID()).Delete(&gormmodel.AuthSession{}).Error; err != nil {
		return fmt.Errorf("delete user auth sessions: %w", mapDBError(err))
	}
	if err := tx.Where("user_id = ?", userID.UUID()).Delete(&gormmodel.UserIdentity{}).Error; err != nil {
		return fmt.Errorf("delete user auth identities: %w", mapDBError(err))
	}
	return nil
}

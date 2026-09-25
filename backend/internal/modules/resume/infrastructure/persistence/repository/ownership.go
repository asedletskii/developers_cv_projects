package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
	"developers-cv/backend/internal/platform/database"
)

func (r *baseRepository) owned(ctx context.Context, resumeID valueobjects.ID, uid useridentity.ID) (*gorm.DB, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("repository dependencies are nil")
	}
	tx := database.DBFromContext(ctx, r.db)
	var n int64
	if e := tx.WithContext(ctx).Model(&gormmodel.Resume{}).Where("id=? AND user_id=?", resumeID.UUID(), uid.UUID()).Count(&n).Error; e != nil {
		return nil, e
	}
	if n == 0 {
		return nil, domain.ErrEntityNotFound
	}
	return tx, nil
}

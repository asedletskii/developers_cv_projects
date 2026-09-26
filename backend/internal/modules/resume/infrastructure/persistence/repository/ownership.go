package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/platform/database"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func (r *baseRepository) owned(ctx context.Context, resumeID valueobjects.ID, uid useridentity.ID) (*gorm.DB, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("repository dependencies are nil")
	}
	tx := database.DBFromContext(ctx, r.db)
	var n int64
	if err := tx.WithContext(ctx).Model(&gormmodel.Resume{}).Where("id=? AND user_id=?", resumeID.UUID(), uid.UUID()).Count(&n).Error; err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, domain.ErrEntityNotFound
	}
	return tx, nil
}

func (r *baseRepository) ownedForWrite(ctx context.Context, resumeID valueobjects.ID, uid useridentity.ID) (*gorm.DB, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("repository dependencies are nil")
	}
	tx := database.DBFromContext(ctx, r.db).WithContext(ctx)
	var resume gormmodel.Resume
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id=? AND user_id=?", resumeID.UUID(), uid.UUID()).First(&resume).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrEntityNotFound
		}
		return nil, err
	}
	return tx, nil
}

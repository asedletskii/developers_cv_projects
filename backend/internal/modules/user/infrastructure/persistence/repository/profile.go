package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"developers-cv/backend/internal/modules/user/domain"
	"developers-cv/backend/internal/modules/user/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/modules/user/infrastructure/persistence/mappers"
	"developers-cv/backend/internal/platform/database"
)

type ProfileRepository struct{ db *gorm.DB }

func NewProfileRepository(db *gorm.DB) *ProfileRepository { return &ProfileRepository{db: db} }

func (r *ProfileRepository) ReplaceProfile(ctx context.Context, profile *domain.Profile) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("replace profile: repository dependencies are nil")
	}
	if profile == nil {
		return fmt.Errorf("replace profile: profile is nil")
	}

	db, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("replace profile: %w", err)
	}

	now := time.Now().UTC()
	model, rows, err := mappers.ProfileToDBModels(profile, now)
	if err != nil {
		return fmt.Errorf("replace profile: %w", err)
	}

	if err := db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.Assignments(map[string]any{
			"photo_path": model.PhotoURL,
			"updated_at": now,
		}),
	}).Create(&model).Error; err != nil {
		return fmt.Errorf("replace profile: upsert profile: %w", mapProfileDBError(err))
	}

	if err := db.WithContext(ctx).Where("user_id = ?", profile.UserID().UUID()).Delete(&gormmodel.UserProfileLocalization{}).Error; err != nil {
		return fmt.Errorf("replace profile: clear localizations: %w", mapProfileDBError(err))
	}
	if len(rows) > 0 {
		if err := db.WithContext(ctx).Create(&rows).Error; err != nil {
			return fmt.Errorf("replace profile: create localizations: %w", mapProfileDBError(err))
		}
	}

	return nil
}

func mapProfileDBError(err error) error {
	mapped := database.MapPostgresError(err)
	if errors.Is(mapped, database.ErrConstraintViolation) || errors.Is(mapped, database.ErrInvalidData) {
		return domain.ErrInvalidProfile
	}
	return mapped
}

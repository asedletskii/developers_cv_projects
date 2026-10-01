package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

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
	if profile == nil || profile.UserID().IsZero() {
		return fmt.Errorf("replace profile: invalid profile")
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
	if err := db.WithContext(ctx).Where("user_id = ?", profile.UserID().UUID()).First(&model).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("replace profile: load profile: %w", err)
		}
		if err := db.WithContext(ctx).Create(&model).Error; err != nil {
			return fmt.Errorf("replace profile: create profile: %w", err)
		}
	} else {
		if err := db.WithContext(ctx).Model(&gormmodel.UserProfile{}).Where("user_id = ?", profile.UserID().UUID()).Updates(map[string]any{"photo_path": profile.PhotoPath(), "updated_at": now}).Error; err != nil {
			return fmt.Errorf("replace profile: update profile: %w", err)
		}
	}

	if err := db.WithContext(ctx).Where("user_id = ?", profile.UserID().UUID()).Delete(&gormmodel.UserProfileLocalization{}).Error; err != nil {
		return fmt.Errorf("replace profile: clear localizations: %w", err)
	}
	if len(rows) > 0 {
		if err := db.WithContext(ctx).Create(&rows).Error; err != nil {
			return fmt.Errorf("replace profile: create localizations: %w", err)
		}
	}

	return nil
}

package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/user/domain"
	"developers-cv/backend/internal/modules/user/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/platform/database"
)

type ProfileRepository struct{ db *gorm.DB }

func NewProfileRepository(db *gorm.DB) *ProfileRepository { return &ProfileRepository{db: db} }

func (r *ProfileRepository) ReplaceProfile(ctx context.Context, profile *domain.Profile) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("replace profile: repository dependencies are nil")
	}
	if profile == nil || profile.UserID.IsZero() {
		return fmt.Errorf("replace profile: invalid profile")
	}

	db, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("replace profile: %w", err)
	}

	now := time.Now().UTC()
	model := gormmodel.UserProfile{UserID: profile.UserID.UUID(), PhotoPath: profile.PhotoPath, CreatedAt: now, UpdatedAt: now}
	if err := db.WithContext(ctx).Where("user_id = ?", profile.UserID.UUID()).First(&model).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return fmt.Errorf("replace profile: load profile: %w", err)
		}
		if err := db.WithContext(ctx).Create(&model).Error; err != nil {
			return fmt.Errorf("replace profile: create profile: %w", err)
		}
	} else {
		if err := db.WithContext(ctx).Model(&gormmodel.UserProfile{}).Where("user_id = ?", profile.UserID.UUID()).Updates(map[string]any{"photo_path": profile.PhotoPath, "updated_at": now}).Error; err != nil {
			return fmt.Errorf("replace profile: update profile: %w", err)
		}
	}

	if err := db.WithContext(ctx).Where("user_id = ?", profile.UserID.UUID()).Delete(&gormmodel.UserProfileLocalization{}).Error; err != nil {
		return fmt.Errorf("replace profile: clear localizations: %w", err)
	}

	rows := make([]gormmodel.UserProfileLocalization, 0, 2)
	for language, localized := range map[string]domain.LocalizedProfile{"ru": profile.RU, "en": profile.EN} {
		if localizedEmpty(localized) {
			continue
		}
		rows = append(rows, gormmodel.UserProfileLocalization{UserID: profile.UserID.UUID(), Language: language, Name: localized.Name, Surname: localized.Surname, Patronymic: localized.Patronymic, Position: localized.Position, Location: localized.Location})
	}
	if len(rows) > 0 {
		if err := db.WithContext(ctx).Create(&rows).Error; err != nil {
			return fmt.Errorf("replace profile: create localizations: %w", err)
		}
	}
	return nil
}

func localizedEmpty(p domain.LocalizedProfile) bool {
	return p.Name == nil && p.Surname == nil && p.Patronymic == nil && p.Position == nil && p.Location == nil
}

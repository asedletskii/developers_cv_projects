package readers

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/user/application/queries"
	"developers-cv/backend/internal/modules/user/domain"
	"developers-cv/backend/internal/modules/user/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/platform/database"
)

type GetProfileReader struct{ db *gorm.DB }

func NewGetProfileReader(db *gorm.DB) *GetProfileReader { return &GetProfileReader{db: db} }

func (r *GetProfileReader) GetProfile(ctx context.Context, query queries.GetProfileQuery) (queries.ProfileView, error) {
	if r == nil || r.db == nil {
		return queries.ProfileView{}, fmt.Errorf("get profile: reader dependencies are nil")
	}
	if query.UserID.IsZero() {
		return queries.ProfileView{}, fmt.Errorf("get profile: invalid user id")
	}

	db := database.DBFromContext(ctx, r.db).WithContext(ctx)
	var model gormmodel.UserProfile
	if err := db.Where("user_id = ?", query.UserID.UUID()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return queries.ProfileView{}, queries.ErrProfileNotFound
		}
		return queries.ProfileView{}, fmt.Errorf("get profile: %w", err)
	}

	var localizations []gormmodel.UserProfileLocalization
	if err := db.Where("user_id = ?", query.UserID.UUID()).Find(&localizations).Error; err != nil {
		return queries.ProfileView{}, fmt.Errorf("get profile: load localizations: %w", err)
	}
	result := queries.ProfileView{ID: model.UserID.String(), PhotoPath: model.PhotoPath}
	for _, localization := range localizations {
		localized := queries.LocalizedProfileView{Name: localization.Name, Surname: localization.Surname, Patronymic: localization.Patronymic, Position: localization.Position, Location: localization.Location}
		switch localization.Language {
		case string(domain.ProfileLanguageRU):
			result.RU = localized
		case string(domain.ProfileLanguageEN):
			result.EN = localized
		}
	}
	return result, nil
}

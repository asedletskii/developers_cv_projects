package mappers

import (
	"fmt"
	"time"

	"developers-cv/backend/internal/modules/user/domain"
	"developers-cv/backend/internal/modules/user/infrastructure/persistence/gormmodel"
)

func ProfileToDBModels(profile *domain.Profile, now time.Time) (gormmodel.UserProfile, []gormmodel.UserProfileLocalization, error) {
	if profile == nil {
		return gormmodel.UserProfile{}, nil, fmt.Errorf("profile is nil")
	}
	model := gormmodel.UserProfile{
		UserID:    profile.UserID().UUID(),
		PhotoURL:  profile.PhotoURL(),
		CreatedAt: now,
		UpdatedAt: now,
	}
	rows := make([]gormmodel.UserProfileLocalization, 0, 2)
	appendLocalization := func(language domain.ProfileLanguage, localized domain.LocalizedProfile) {
		if localized.Name() == nil && localized.Surname() == nil && localized.Patronymic() == nil && localized.Position() == nil && localized.Location() == nil {
			return
		}
		rows = append(rows, gormmodel.UserProfileLocalization{
			UserID: profile.UserID().UUID(), Language: string(language),
			Name: localized.Name(), Surname: localized.Surname(), Patronymic: localized.Patronymic(),
			Position: localized.Position(), Location: localized.Location(),
		})
	}
	appendLocalization(domain.ProfileLanguageRU, profile.RU())
	appendLocalization(domain.ProfileLanguageEN, profile.EN())
	return model, rows, nil
}

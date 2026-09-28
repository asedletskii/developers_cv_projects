package domain

import (
	"fmt"

	"developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type ProfileLanguage string

const (
	ProfileLanguageRU ProfileLanguage = "ru"
	ProfileLanguageEN ProfileLanguage = "en"
)

type LocalizedProfile struct {
	Name       *string
	Surname    *string
	Patronymic *string
	Position   *string
	Location   *string
}

type Profile struct {
	UserID    valueobjects.ID
	RU        LocalizedProfile
	EN        LocalizedProfile
	PhotoPath *string
}

func NewProfile(userID valueobjects.ID, ru, en LocalizedProfile, photoPath *string) (*Profile, error) {
	if userID.IsZero() {
		return nil, fmt.Errorf("%w: profile user id is empty", ErrInvalidUser)
	}
	return &Profile{UserID: userID, RU: ru, EN: en, PhotoPath: photoPath}, nil
}

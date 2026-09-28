package domain

import (
	"fmt"
	"net/url"
	"strings"

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
	if err := ValidateLocalizedProfile(ru, "ru"); err != nil {
		return nil, err
	}
	if err := ValidateLocalizedProfile(en, "en"); err != nil {
		return nil, err
	}
	if err := ValidatePhotoPath(photoPath); err != nil {
		return nil, err
	}
	return &Profile{UserID: userID, RU: ru, EN: en, PhotoPath: photoPath}, nil
}

func ValidateLocalizedProfile(profile LocalizedProfile, language ProfileLanguage) error {
	if profile.Name != nil && strings.TrimSpace(*profile.Name) == "" {
		return fmt.Errorf("%w: %s name is empty", ErrInvalidProfile, language)
	}
	if profile.Surname != nil && strings.TrimSpace(*profile.Surname) == "" {
		return fmt.Errorf("%w: %s surname is empty", ErrInvalidProfile, language)
	}
	return nil
}

func ValidatePhotoPath(photoPath *string) error {
	if photoPath == nil {
		return nil
	}
	value := strings.TrimSpace(*photoPath)
	if value == "" {
		return fmt.Errorf("%w: photo path is empty", ErrInvalidProfile)
	}
	parsed, err := url.ParseRequestURI(value)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return fmt.Errorf("%w: photo path must be an http(s) URL", ErrInvalidProfile)
	}
	return nil
}

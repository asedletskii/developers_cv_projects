package domain

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ProfileLanguage string

const (
	ProfileLanguageRU ProfileLanguage = "ru"
	ProfileLanguageEN ProfileLanguage = "en"
)

func (v ProfileLanguage) IsValid() bool {
	switch v {
	case ProfileLanguageRU, ProfileLanguageEN:
		return true
	default:
		return false
	}
}

type LocalizedProfile struct {
	name       *string
	surname    *string
	patronymic *string
	position   *string
	location   *string
}

func NewLocalizedProfile(name, surname, patronymic, position, location *string) (LocalizedProfile, error) {
	profile := LocalizedProfile{
		name:       trimOptionalString(name),
		surname:    trimOptionalString(surname),
		patronymic: trimOptionalString(patronymic),
		position:   trimOptionalString(position),
		location:   trimOptionalString(location),
	}
	if err := validateLocalizedProfile(profile, ""); err != nil {
		return LocalizedProfile{}, err
	}
	profile.name = normalizeOptionalString(profile.name)
	profile.surname = normalizeOptionalString(profile.surname)
	profile.patronymic = normalizeOptionalString(profile.patronymic)
	profile.position = normalizeOptionalString(profile.position)
	profile.location = normalizeOptionalString(profile.location)
	return profile, nil
}

func (p LocalizedProfile) Name() *string       { return cloneString(p.name) }
func (p LocalizedProfile) Surname() *string    { return cloneString(p.surname) }
func (p LocalizedProfile) Patronymic() *string { return cloneString(p.patronymic) }
func (p LocalizedProfile) Position() *string   { return cloneString(p.position) }
func (p LocalizedProfile) Location() *string   { return cloneString(p.location) }

type Profile struct {
	userID    sharedvalueobjects.ID
	ru        LocalizedProfile
	en        LocalizedProfile
	photoPath *string
}

func NewProfile(userID sharedvalueobjects.ID, ru, en LocalizedProfile, photoPath *string) (*Profile, error) {
	if userID.IsZero() {
		return nil, fmt.Errorf("%w: profile user id is empty", ErrInvalidProfile)
	}
	if err := validateLocalizedProfile(ru, ProfileLanguageRU); err != nil {
		return nil, err
	}
	if err := validateLocalizedProfile(en, ProfileLanguageEN); err != nil {
		return nil, err
	}
	if err := validatePhotoURL(photoPath); err != nil {
		return nil, err
	}
	return &Profile{
		userID:    userID,
		ru:        ru,
		en:        en,
		photoPath: normalizeOptionalString(photoPath),
	}, nil
}

func validateLocalizedProfile(profile LocalizedProfile, language ProfileLanguage) error {
	if language != "" && !language.IsValid() {
		return fmt.Errorf("%w: invalid profile language", ErrInvalidProfile)
	}
	if profile.name != nil && strings.TrimSpace(*profile.name) == "" {
		return fmt.Errorf("%w: %s name is empty", ErrInvalidProfile, language)
	}
	if profile.surname != nil && strings.TrimSpace(*profile.surname) == "" {
		return fmt.Errorf("%w: %s surname is empty", ErrInvalidProfile, language)
	}
	if !withinLength(profile.name, 100) || !withinLength(profile.surname, 100) ||
		!withinLength(profile.patronymic, 100) || !withinLength(profile.position, 200) || !withinLength(profile.location, 200) {
		return fmt.Errorf("%w: localized profile field exceeds configured length", ErrInvalidProfile)
	}
	return nil
}

func validatePhotoURL(photoURL *string) error {
	if photoURL == nil {
		return nil
	}
	value := strings.TrimSpace(*photoURL)
	if value == "" {
		return fmt.Errorf("%w: photo URL is empty", ErrInvalidProfile)
	}
	if len(value) > 2048 {
		return fmt.Errorf("%w: photo URL exceeds 2048 bytes", ErrInvalidProfile)
	}
	parsed, err := url.ParseRequestURI(value)
	if err != nil || parsed.Hostname() == "" || parsed.Scheme != "https" || parsed.User != nil {
		return fmt.Errorf("%w: photo URL must be an HTTPS URL without user information", ErrInvalidProfile)
	}
	return nil
}

func (p *Profile) UserID() sharedvalueobjects.ID {
	if p == nil {
		return sharedvalueobjects.ID{}
	}
	return p.userID
}

func (p *Profile) RU() LocalizedProfile {
	if p == nil {
		return LocalizedProfile{}
	}
	return p.ru
}

func (p *Profile) EN() LocalizedProfile {
	if p == nil {
		return LocalizedProfile{}
	}
	return p.en
}

func (p *Profile) PhotoURL() *string {
	if p == nil {
		return nil
	}
	return cloneString(p.photoPath)
}

func cloneString(value *string) *string {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func trimOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	return &trimmed
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func withinLength(value *string, max int) bool {
	return value == nil || utf8.RuneCountInString(*value) <= max
}

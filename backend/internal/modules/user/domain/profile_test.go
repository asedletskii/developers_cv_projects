package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestNewProfileRejectsWhitespaceOnlyName(t *testing.T) {
	_, err := NewProfile(valueobjects.NewID(), LocalizedProfile{Name: stringPtr("   ")}, LocalizedProfile{}, nil)
	require.ErrorIs(t, err, ErrInvalidProfile)
}

func TestNewProfileRejectsWhitespaceOnlySurname(t *testing.T) {
	_, err := NewProfile(valueobjects.NewID(), LocalizedProfile{Surname: stringPtr("	")}, LocalizedProfile{}, nil)
	require.ErrorIs(t, err, ErrInvalidProfile)
}

func TestNewProfileAllowsPartialLocalization(t *testing.T) {
	name := "Dmitry"
	profile, err := NewProfile(valueobjects.NewID(), LocalizedProfile{Name: &name}, LocalizedProfile{}, nil)
	require.NoError(t, err)
	require.Equal(t, &name, profile.RU.Name)
}

func TestNewProfileRejectsRelativePhotoPath(t *testing.T) {
	photoPath := "/profiles/dmitry.svg"
	_, err := NewProfile(valueobjects.NewID(), LocalizedProfile{}, LocalizedProfile{}, &photoPath)
	require.ErrorIs(t, err, ErrInvalidProfile)
}

func TestValidatePhotoPathRejectsEmptyPhotoPath(t *testing.T) {
	photoPath := "   "
	require.ErrorIs(t, ValidatePhotoPath(&photoPath), ErrInvalidProfile)
}

func TestNewProfileAcceptsHTTPSPhotoPath(t *testing.T) {
	photoPath := "https://example.com/profiles/dmitry.svg"
	profile, err := NewProfile(valueobjects.NewID(), LocalizedProfile{}, LocalizedProfile{}, &photoPath)
	require.NoError(t, err)
	require.Equal(t, &photoPath, profile.PhotoPath)
}

func stringPtr(value string) *string { return &value }

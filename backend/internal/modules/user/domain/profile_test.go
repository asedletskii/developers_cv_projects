package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestNewProfileRejectsWhitespaceOnlyName(t *testing.T) {
	ru := LocalizedProfile{name: stringPtr("   ")}
	_, err := NewProfile(sharedvalueobjects.NewID(), ru, LocalizedProfile{}, nil)
	require.ErrorIs(t, err, ErrInvalidProfile)
}

func TestNewProfileRejectsWhitespaceOnlySurname(t *testing.T) {
	ru := LocalizedProfile{surname: stringPtr("\t")}
	_, err := NewProfile(sharedvalueobjects.NewID(), ru, LocalizedProfile{}, nil)
	require.ErrorIs(t, err, ErrInvalidProfile)
}

func TestNewProfileAllowsPartialLocalization(t *testing.T) {
	name := "Dmitry"
	ru, err := NewLocalizedProfile(&name, nil, nil, nil, nil)
	require.NoError(t, err)
	profile, err := NewProfile(sharedvalueobjects.NewID(), ru, LocalizedProfile{}, nil)
	require.NoError(t, err)
	require.Equal(t, &name, profile.RU().Name())
}

func TestNewProfileRejectsRelativePhotoURL(t *testing.T) {
	photoURL := "/profiles/dmitry.svg"
	_, err := NewProfile(sharedvalueobjects.NewID(), LocalizedProfile{}, LocalizedProfile{}, &photoURL)
	require.ErrorIs(t, err, ErrInvalidProfile)
}

func TestValidatePhotoURLRejectsEmptyPhotoURL(t *testing.T) {
	photoURL := "   "
	require.ErrorIs(t, validatePhotoURL(&photoURL), ErrInvalidProfile)
}

func TestNewProfileAcceptsHTTPSPhotoURL(t *testing.T) {
	photoURL := "https://example.com/profiles/dmitry.svg"
	profile, err := NewProfile(sharedvalueobjects.NewID(), LocalizedProfile{}, LocalizedProfile{}, &photoURL)
	require.NoError(t, err)
	require.Equal(t, &photoURL, profile.PhotoURL())
}

func TestNewLocalizedProfileTrimsOptionalValues(t *testing.T) {
	name := " Dmitry "
	position := " Backend Developer "
	profile, err := NewLocalizedProfile(&name, nil, nil, &position, nil)
	require.NoError(t, err)
	require.Equal(t, "Dmitry", *profile.Name())
	require.Equal(t, "Backend Developer", *profile.Position())
}

func stringPtr(value string) *string { return &value }

func TestNewProfileRejectsZeroUserID(t *testing.T) {
	_, err := NewProfile(sharedvalueobjects.ID{}, LocalizedProfile{}, LocalizedProfile{}, nil)
	require.ErrorIs(t, err, ErrInvalidProfile)
}

func TestNewLocalizedProfileRejectsWhitespaceName(t *testing.T) {
	name := "   "
	_, err := NewLocalizedProfile(&name, nil, nil, nil, nil)
	require.ErrorIs(t, err, ErrInvalidProfile)
}

func TestLocalizedProfileAccessorsReturnCopies(t *testing.T) {
	name := "Dmitry"
	profile, err := NewLocalizedProfile(&name, nil, nil, nil, nil)
	require.NoError(t, err)

	value := profile.Name()
	*value = "Changed"

	require.Equal(t, "Dmitry", *profile.Name())
}

func TestNewProfileRejectsHTTPPhotoURL(t *testing.T) {
	photoURL := "http://example.com/profiles/dmitry.svg"
	_, err := NewProfile(sharedvalueobjects.NewID(), LocalizedProfile{}, LocalizedProfile{}, &photoURL)
	require.ErrorIs(t, err, ErrInvalidProfile)
}

func TestNewProfileRejectsPhotoURLWithUserInformation(t *testing.T) {
	// #nosec G101 -- Это тест, проверяющий, что URL с пользовательской информацией отклоняется.
	photoURL := "https://user:password@example.com/profiles/dmitry.svg"
	_, err := NewProfile(sharedvalueobjects.NewID(), LocalizedProfile{}, LocalizedProfile{}, &photoURL)
	require.ErrorIs(t, err, ErrInvalidProfile)
}

func TestNewLocalizedProfileRejectsOverlongName(t *testing.T) {
	name := strings.Repeat("a", 101)
	_, err := NewLocalizedProfile(&name, nil, nil, nil, nil)
	require.ErrorIs(t, err, ErrInvalidProfile)
}

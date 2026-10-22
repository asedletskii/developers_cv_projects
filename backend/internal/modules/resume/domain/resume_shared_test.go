package domain

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestResumeEnumValidationContracts(t *testing.T) {
	cases := []struct {
		name  string
		valid func() bool
		bad   func() bool
	}{
		{"work format", func() bool { return WorkFormatRemote.isValid() }, func() bool { return WorkFormat("invalid").isValid() }},
		{"resume language", func() bool { return ResumeLanguageEN.isValid() }, func() bool { return ResumeLanguage("invalid").isValid() }},
		{"spoken language", func() bool { return SpokenLanguageEN.isValid() }, func() bool { return SpokenLanguage("invalid").isValid() }},
		{"language grade", func() bool { return LanguageGradeC1.isValid() }, func() bool { return LanguageGrade("invalid").isValid() }},
		{"education level", func() bool { return EducationLevelMaster.isValid() }, func() bool { return EducationLevel(99).isValid() }},
		{"activity type", func() bool { return ActivityTypeOpenSource.isValid() }, func() bool { return ActivityType("invalid").isValid() }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.valid() || tc.bad() {
				t.Fatal("enum validation contract is broken")
			}
		})
	}
}

func TestDateRoundTrip(t *testing.T) {
	value, err := ParseDate("1998-04-12")
	require.NoError(t, err)
	require.Equal(t, "1998-04-12", value.String())
	require.False(t, value.IsZero())
}

func TestYearMonthRoundTrip(t *testing.T) {
	value, err := ParseYearMonth("2026-09")
	require.NoError(t, err)
	require.Equal(t, "2026-09", value.String())
	require.False(t, value.IsZero())
}

func TestYearRoundTrip(t *testing.T) {
	value, err := ParseYear("2026")
	require.NoError(t, err)
	require.Equal(t, "2026", value.String())
	require.False(t, value.IsZero())
}

func TestDateNormalizesToUTCDateOnly(t *testing.T) {
	value, err := NewDate(time.Date(2026, time.September, 22, 23, 59, 0, 0, time.FixedZone("test", 3*60*60)))
	require.NoError(t, err)
	require.Equal(t, "2026-09-22", value.String())
	require.Equal(t, time.UTC, value.Time().Location())
}

func TestExternalURLAllowsHTTPS(t *testing.T) {
	value := "https://example.com/avatar.svg"
	parsed, err := NewExternalURL(value)
	require.NoError(t, err)
	require.Equal(t, value, parsed.String())
}

func TestExternalURLRejectsHTTP(t *testing.T) {
	_, err := NewExternalURL("http://example.com/avatar.svg")
	require.ErrorIs(t, err, ErrInvalidValue)
}

func TestExternalURLRejectsUserInformation(t *testing.T) {
	_, err := NewExternalURL("https://user:password@example.com/avatar.svg")
	require.ErrorIs(t, err, ErrInvalidValue)
}

func TestExternalURLRejectsOverlongValue(t *testing.T) {
	value := "https://example.com/" + strings.Repeat("a", MaxExternalURLLength)
	_, err := NewExternalURL(value)
	require.ErrorIs(t, err, ErrInvalidValue)
}

func TestExternalURLRejectsUnsupportedScheme(t *testing.T) {
	_, err := NewExternalURL("javascript:alert(1)")
	require.Error(t, err)
}

func TestExternalURLRejectsRelativePath(t *testing.T) {
	_, err := NewExternalURL("/assets/avatar.svg")
	require.Error(t, err)
}

func TestParseResumeStatus(t *testing.T) {
	status, err := ParseResumeStatus(string(ResumeStatusOpened))
	require.NoError(t, err)
	require.Equal(t, ResumeStatusOpened, status)

	_, err = ParseResumeStatus("invalid")
	require.ErrorIs(t, err, ErrInvalidValue)
}

func TestResolveSingletonIDUsesCanonicalIdentityPolicy(t *testing.T) {
	currentID := sharedvalueobjects.NewID()
	generatedID, err := resolveSingletonID(sharedvalueobjects.ID{}, sharedvalueobjects.ID{}, "profile")
	require.NoError(t, err)
	require.False(t, generatedID.IsZero())

	preservedID, err := resolveSingletonID(sharedvalueobjects.ID{}, currentID, "profile")
	require.NoError(t, err)
	require.Equal(t, currentID, preservedID)

	updatedID, err := resolveSingletonID(currentID, currentID, "profile")
	require.NoError(t, err)
	require.Equal(t, currentID, updatedID)

	_, err = resolveSingletonID(sharedvalueobjects.NewID(), currentID, "profile")
	require.ErrorIs(t, err, ErrResumeInvariant)
	_, err = resolveSingletonID(currentID, sharedvalueobjects.ID{}, "profile")
	require.ErrorIs(t, err, ErrResumeInvariant)
}

func TestResolveCollectionItemIDUsesCanonicalIdentityPolicy(t *testing.T) {
	currentID := sharedvalueobjects.NewID()
	existing := map[sharedvalueobjects.ID]struct{}{currentID: {}}

	generatedID, err := resolveCollectionItemID(sharedvalueobjects.ID{}, existing, "project")
	require.NoError(t, err)
	require.False(t, generatedID.IsZero())
	require.NotEqual(t, currentID, generatedID)

	preservedID, err := resolveCollectionItemID(currentID, existing, "project")
	require.NoError(t, err)
	require.Equal(t, currentID, preservedID)

	_, err = resolveCollectionItemID(sharedvalueobjects.NewID(), existing, "project")
	require.ErrorIs(t, err, ErrResumeInvariant)
}

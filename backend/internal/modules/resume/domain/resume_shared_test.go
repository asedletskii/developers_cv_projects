package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestResumeEnumValidationContracts(t *testing.T) {
	cases := []struct {
		name  string
		valid func() bool
		bad   func() bool
	}{
		{"work format", func() bool { return WorkFormatRemote.IsValid() }, func() bool { return WorkFormat("invalid").IsValid() }},
		{"resume language", func() bool { return ResumeLanguageEN.IsValid() }, func() bool { return ResumeLanguage("invalid").IsValid() }},
		{"spoken language", func() bool { return SpokenLanguageEN.IsValid() }, func() bool { return SpokenLanguage("invalid").IsValid() }},
		{"language grade", func() bool { return LanguageGradeC1.IsValid() }, func() bool { return LanguageGrade("invalid").IsValid() }},
		{"technology grade", func() bool { return TechnologyGradeSenior.IsValid() }, func() bool { return TechnologyGrade("invalid").IsValid() }},
		{"education level", func() bool { return EducationLevelMaster.IsValid() }, func() bool { return EducationLevel(99).IsValid() }},
		{"contact key", func() bool { return ContactKeyEmail.IsValid() }, func() bool { return ContactKey("invalid").IsValid() }},
		{"activity type", func() bool { return ActivityTypeOpenSource.IsValid() }, func() bool { return ActivityType("invalid").IsValid() }},
		{"hard skill", func() bool { return HardSkillAPI.IsValid() }, func() bool { return HardSkill("invalid").IsValid() }},
		{"soft skill", func() bool { return SoftSkillOwnership.IsValid() }, func() bool { return SoftSkill("invalid").IsValid() }},
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

func TestExternalURLAllowsHTTPAndHTTPS(t *testing.T) {
	for _, value := range []string{
		"https://example.com/avatar.svg",
		"http://example.com/avatar.svg",
	} {
		url, err := NewExternalURL(value)
		require.NoError(t, err)
		require.Equal(t, value, url.String())
	}
}

func TestExternalURLRejectsUnsupportedScheme(t *testing.T) {
	_, err := NewExternalURL("javascript:alert(1)")
	require.Error(t, err)
}

func TestExternalURLRejectsRelativePath(t *testing.T) {
	_, err := NewExternalURL("/assets/avatar.svg")
	require.Error(t, err)
}

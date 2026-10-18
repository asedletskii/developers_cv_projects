package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewResumeOverviewRejectsDuplicateSections(t *testing.T) {
	section, err := NewSectionSettings(SectionSettingsParams{ID: SectionAbout, Enabled: true})
	require.NoError(t, err)
	_, err = newResumeOverview(ResumeOverviewParams{Sections: []SectionSettings{section, section}})
	require.ErrorIs(t, err, ErrInvalidResume)
}

func TestNewResumeOverviewRejectsInvalidSection(t *testing.T) {
	_, err := newResumeOverview(ResumeOverviewParams{Sections: []SectionSettings{{}}})
	require.ErrorIs(t, err, ErrInvalidValue)
}

func TestResumeOverviewDeepCopiesMutableFields(t *testing.T) {
	primary, err := NewColor("#112233")
	require.NoError(t, err)
	theme, err := NewThemeSettings(ThemeSettingsParams{Primary: &primary})
	require.NoError(t, err)
	section, err := NewSectionSettings(SectionSettingsParams{ID: SectionAbout, Enabled: true})
	require.NoError(t, err)
	overview, err := newResumeOverview(ResumeOverviewParams{Theme: &theme, Sections: []SectionSettings{section}})
	require.NoError(t, err)

	view := overview.Theme()
	primaryCopy := view.Primary()
	*primaryCopy = Color{}
	sections := overview.Sections()
	sections[0] = SectionSettings{}

	againTheme := overview.Theme()
	againSections := overview.Sections()
	require.NotNil(t, againTheme)
	require.NotNil(t, againTheme.Primary())
	require.Equal(t, "#112233", againTheme.Primary().String())
	require.Equal(t, SectionAbout, againSections[0].ID())
}

func TestResumeReplaceOverviewReplacesWholeSection(t *testing.T) {
	resume := newTestResume(t)

	first := ResumeOverviewParams{
		Greeting: mustGreeting(t, GreetingParams{First: "First"}),
		Sections: []SectionSettings{},
	}
	require.NoError(t, resume.ReplaceOverview(first))

	second := ResumeOverviewParams{
		Greeting: mustGreeting(t, GreetingParams{First: "Second"}),
		Sections: []SectionSettings{},
	}
	require.NoError(t, resume.ReplaceOverview(second))

	require.Equal(t, "Second", resume.Overview().Greeting().First())
}

func TestResumeReplaceOverviewRejectsInvalidSection(t *testing.T) {
	resume := newTestResume(t)

	err := resume.ReplaceOverview(ResumeOverviewParams{Sections: []SectionSettings{{}}})
	require.ErrorIs(t, err, ErrInvalidValue)
	require.Nil(t, resume.Overview())
}

func TestAvailabilityStatusValidationContract(t *testing.T) {
	if !AvailabilityStatusOpen.isValid() || AvailabilityStatus("invalid").isValid() {
		t.Fatal("availability status validation contract is broken")
	}
}

func TestSectionIDValidationContract(t *testing.T) {
	if !SectionAbout.isValid() || SectionID("invalid").isValid() {
		t.Fatal("section ID validation contract is broken")
	}
}

func TestNewResumeOverviewRejectsEmptyThemeValue(t *testing.T) {
	_, err := newResumeOverview(ResumeOverviewParams{Theme: &ThemeSettings{}})
	require.ErrorIs(t, err, ErrInvalidValue)
}

func TestNewThemeSettingsRejectsZeroValueColors(t *testing.T) {
	color := Color{}
	_, err := NewThemeSettings(ThemeSettingsParams{Primary: &color})
	require.ErrorIs(t, err, ErrInvalidValue)
}

func TestNewResumeOverviewRejectsZeroValueNestedObjects(t *testing.T) {
	availability := Availability{}
	_, err := newResumeOverview(ResumeOverviewParams{Availability: &availability})
	require.ErrorIs(t, err, ErrInvalidValue)
}

func TestParseSectionID(t *testing.T) {
	section, err := ParseSectionID(string(SectionAbout))
	require.NoError(t, err)
	require.Equal(t, SectionAbout, section)

	_, err = ParseSectionID("invalid")
	require.ErrorIs(t, err, ErrInvalidValue)
}

func TestNewGreetingNormalizesAndValidates(t *testing.T) {
	greeting, err := NewGreeting(GreetingParams{First: " Hello ", SecondHighlight: " Developer "})
	require.NoError(t, err)
	require.Equal(t, "Hello", greeting.First())
	require.Equal(t, "Developer", greeting.SecondHighlight())

	_, err = NewGreeting(GreetingParams{First: strings.Repeat("x", MaxShortTextLength+1)})
	require.ErrorIs(t, err, ErrInvalidValue)
}

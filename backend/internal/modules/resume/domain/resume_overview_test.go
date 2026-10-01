package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewResumeOverviewRejectsDuplicateSections(t *testing.T) {
	section, err := NewSectionSettings(SectionSettingsParams{ID: SectionAbout, Enabled: true})
	require.NoError(t, err)
	_, err = NewResumeOverview(ResumeOverviewParams{Sections: []SectionSettings{section, section}})
	require.ErrorIs(t, err, ErrInvalidResume)
}

func TestNewResumeOverviewRejectsInvalidSection(t *testing.T) {
	_, err := NewResumeOverview(ResumeOverviewParams{Sections: []SectionSettings{{}}})
	require.ErrorIs(t, err, ErrInvalidValue)
}

func TestResumeOverviewDeepCopiesMutableFields(t *testing.T) {
	primary, err := NewColor("#112233")
	require.NoError(t, err)
	theme, err := NewThemeSettings(ThemeSettingsParams{Primary: &primary})
	require.NoError(t, err)
	section, err := NewSectionSettings(SectionSettingsParams{ID: SectionAbout, Enabled: true})
	require.NoError(t, err)
	overview, err := NewResumeOverview(ResumeOverviewParams{Theme: &theme, Sections: []SectionSettings{section}})
	require.NoError(t, err)

	view := overview.Theme()
	view.Primary = nil
	sections := overview.Sections()
	sections[0] = SectionSettings{}

	againTheme := overview.Theme()
	againSections := overview.Sections()
	require.NotNil(t, againTheme)
	require.NotNil(t, againTheme.Primary)
	require.Equal(t, SectionAbout, againSections[0].ID())
}

func TestResumeReplaceOverviewReplacesWholeSection(t *testing.T) {
	resume := newTestResume(t)

	first := ResumeOverviewParams{
		Greeting: NewGreeting(GreetingParams{First: "First"}),
		Sections: []SectionSettings{},
	}
	require.NoError(t, resume.ReplaceOverview(first))

	second := ResumeOverviewParams{
		Greeting: NewGreeting(GreetingParams{First: "Second"}),
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
	if !AvailabilityStatusOpen.IsValid() || AvailabilityStatus("invalid").IsValid() {
		t.Fatal("availability status validation contract is broken")
	}
}

func TestSectionIDValidationContract(t *testing.T) {
	if !SectionAbout.IsValid() || SectionID("invalid").IsValid() {
		t.Fatal("section ID validation contract is broken")
	}
}

package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

func TestNewResumeOverviewRejectsDuplicateSections(t *testing.T) {
	section, err := valueobjects.NewSectionSettings(valueobjects.SectionSettingsParams{ID: enums.SectionAbout, Enabled: true})
	require.NoError(t, err)
	_, err = NewResumeOverview(ResumeOverviewParams{Sections: []valueobjects.SectionSettings{section, section}})
	require.ErrorIs(t, err, ErrInvalidResume)
}

func TestNewResumeOverviewRejectsInvalidSection(t *testing.T) {
	_, err := NewResumeOverview(ResumeOverviewParams{Sections: []valueobjects.SectionSettings{{}}})
	require.ErrorIs(t, err, valueobjects.ErrInvalidValue)
}

func TestResumeOverviewDeepCopiesMutableFields(t *testing.T) {
	primary, err := valueobjects.NewColor("#112233")
	require.NoError(t, err)
	theme, err := valueobjects.NewThemeSettings(valueobjects.ThemeSettingsParams{Primary: &primary})
	require.NoError(t, err)
	section, err := valueobjects.NewSectionSettings(valueobjects.SectionSettingsParams{ID: enums.SectionAbout, Enabled: true})
	require.NoError(t, err)
	overview, err := NewResumeOverview(ResumeOverviewParams{Theme: &theme, Sections: []valueobjects.SectionSettings{section}})
	require.NoError(t, err)

	view := overview.Theme()
	view.Primary = nil
	sections := overview.Sections()
	sections[0] = valueobjects.SectionSettings{}

	againTheme := overview.Theme()
	againSections := overview.Sections()
	require.NotNil(t, againTheme)
	require.NotNil(t, againTheme.Primary)
	require.Equal(t, enums.SectionAbout, againSections[0].ID())
}

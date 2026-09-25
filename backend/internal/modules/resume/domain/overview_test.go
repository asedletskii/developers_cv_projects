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
	_, err = NewResumeOverview(ResumeOverviewParams{Presentation: OverviewPresentation{Sections: []valueobjects.SectionSettings{section, section}}})
	require.ErrorIs(t, err, ErrInvalidResume)
}

func TestNewResumeOverviewRejectsInvalidSection(t *testing.T) {
	_, err := NewResumeOverview(ResumeOverviewParams{Presentation: OverviewPresentation{Sections: []valueobjects.SectionSettings{{}}}})
	require.ErrorIs(t, err, valueobjects.ErrInvalidValue)
}

func TestResumeOverviewDeepCopiesMutableFields(t *testing.T) {
	primary, err := valueobjects.NewColor("#112233")
	require.NoError(t, err)
	theme, err := valueobjects.NewThemeSettings(valueobjects.ThemeSettingsParams{Primary: &primary})
	require.NoError(t, err)
	section, err := valueobjects.NewSectionSettings(valueobjects.SectionSettingsParams{ID: enums.SectionAbout, Enabled: true})
	require.NoError(t, err)
	overview, err := NewResumeOverview(ResumeOverviewParams{Presentation: OverviewPresentation{Theme: &theme, Sections: []valueobjects.SectionSettings{section}}})
	require.NoError(t, err)

	view := overview.Presentation()
	view.Theme.Primary = nil
	view.Sections[0] = valueobjects.SectionSettings{}

	again := overview.Presentation()
	require.NotNil(t, again.Theme)
	require.NotNil(t, again.Theme.Primary)
	require.Equal(t, enums.SectionAbout, again.Sections[0].ID())
}

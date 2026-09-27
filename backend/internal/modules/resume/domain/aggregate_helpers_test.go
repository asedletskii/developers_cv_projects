package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

func TestNormalizeOptionalStringTrimsWhitespaceAndTurnsEmptyIntoNil(t *testing.T) {
	value := "  optional value  "
	got := cloneString(&value)
	require.NotNil(t, got)
	require.Equal(t, "optional value", *got)

	empty := "   "
	require.Nil(t, cloneString(&empty))
	require.Nil(t, cloneString(nil))
}

func TestCloneYearMonthReturnsIndependentValue(t *testing.T) {
	value, err := valueobjects.ParseYearMonth("2025-01")
	require.NoError(t, err)

	clone := cloneYearMonth(&value)
	require.NotNil(t, clone)
	require.Equal(t, value, *clone)
	require.NotSame(t, &value, clone)
	require.Nil(t, cloneYearMonth(nil))
}

func TestCloneSectionSettingsReturnsIndependentSlice(t *testing.T) {
	section, err := valueobjects.NewSectionSettings(valueobjects.SectionSettingsParams{
		ID:      enums.SectionAbout,
		Enabled: true,
	})
	require.NoError(t, err)

	values := []valueobjects.SectionSettings{section}
	clone := cloneSectionSettings(values)
	require.Equal(t, values, clone)
	require.NotSame(t, &values[0], &clone[0])

	clone[0] = valueobjects.SectionSettings{}
	require.Equal(t, enums.SectionAbout, values[0].ID())
	require.Nil(t, cloneSectionSettings(nil))
}

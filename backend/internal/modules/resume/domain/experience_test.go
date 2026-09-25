package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

func TestNewExperienceRejectsEndBeforeStart(t *testing.T) {
	id := valueobjects.NewID()
	resumeID := valueobjects.NewID()
	start, err := valueobjects.ParseYearMonth("2024-01")
	require.NoError(t, err)
	end, err := valueobjects.ParseYearMonth("2023-12")
	require.NoError(t, err)

	_, err = NewExperience(ExperienceParams{
		ID: id, ResumeID: resumeID, Company: "Acme", Position: "Developer",
		StartDate: start, EndDate: &end,
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestNewExperienceTrimsValuesAndBuildsChildren(t *testing.T) {
	id := valueobjects.NewID()
	resumeID := valueobjects.NewID()
	start, err := valueobjects.ParseYearMonth("2024-01")
	require.NoError(t, err)
	grade := enums.TechnologyGradeSenior

	item, err := NewExperience(ExperienceParams{
		ID: id, ResumeID: resumeID, Company: " Acme ", Position: " Developer ",
		StartDate:    start,
		Achievements: []string{" shipped v2 ", "reduced latency"},
		Technologies: []ExperienceTechnologyParams{{Name: " Go ", Grade: &grade}},
	})
	require.NoError(t, err)
	require.Equal(t, "Acme", item.Company())
	require.Equal(t, "Developer", item.Position())
	require.Equal(t, []string{"shipped v2", "reduced latency"}, item.Achievements())
	require.Len(t, item.Technologies(), 1)
	require.Equal(t, "Go", item.Technologies()[0].Name())
	require.Equal(t, &grade, item.Technologies()[0].Grade())
}

package domain_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

func TestEducationValidation(t *testing.T) {
	id := valueobjects.NewID()
	resumeID := valueobjects.NewID()
	endDate, err := valueobjects.ParseYear("2024")
	require.NoError(t, err)

	item, err := domain.NewEducation(domain.EducationParams{
		ID: id, ResumeID: resumeID, Location: "Paris", University: "University",
		Grade: enums.EducationLevelMaster, ProgramName: "Computer Science", Faculty: "Engineering", EndDate: &endDate,
	})
	require.NoError(t, err)
	require.Equal(t, id, item.ID())
	require.Equal(t, resumeID, item.ResumeID())
	require.Equal(t, "2024", item.EndDate().String())

	_, err = domain.NewEducation(domain.EducationParams{ID: id, ResumeID: resumeID, Location: "", University: "University", Grade: enums.EducationLevelMaster, ProgramName: "CS", Faculty: "Engineering"})
	require.Error(t, err)
}

package entities

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

func TestNewExperienceRejectsEndBeforeStart(t *testing.T) {
	id, err := valueobjects.IDFromUUID(uuid.New())
	require.NoError(t, err)
	start, err := valueobjects.ParseYearMonth("2026-09")
	require.NoError(t, err)
	end, err := valueobjects.ParseYearMonth("2026-08")
	require.NoError(t, err)

	_, err = NewExperience(ExperienceParams{
		ID:        id,
		Company:   "Acme",
		Position:  "Developer",
		StartDate: start,
		EndDate:   &end,
	})
	require.Error(t, err)
}

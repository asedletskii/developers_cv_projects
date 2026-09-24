package domain_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

func TestActivityValidation(t *testing.T) {
	date, err := valueobjects.ParseYearMonth("2025-01")
	require.NoError(t, err)
	item, err := domain.NewActivity(domain.ActivityParams{ID: valueobjects.NewID(), ResumeID: valueobjects.NewID(), Type: enums.ActivityTypeOpenSource, Title: "Open Source", Date: &date})
	require.NoError(t, err)
	require.Equal(t, enums.ActivityTypeOpenSource, item.Type())
	require.Equal(t, "2025-01", item.Date().String())

	_, err = domain.NewActivity(domain.ActivityParams{ID: valueobjects.NewID(), ResumeID: valueobjects.NewID(), Type: enums.ActivityType("invalid"), Title: "Open Source"})
	require.Error(t, err)
}

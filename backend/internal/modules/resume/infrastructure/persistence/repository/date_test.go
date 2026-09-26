package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

func TestYearMonthDate(t *testing.T) {
	value, err := valueobjects.ParseYearMonth("2026-09")
	require.NoError(t, err)

	actual := yearMonthDate(value)

	require.Equal(t, time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC), actual)
}

func TestYearDate(t *testing.T) {
	value, err := valueobjects.ParseYear("2026")
	require.NoError(t, err)

	actual := yearDate(value)

	require.Equal(t, time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC), actual)
}

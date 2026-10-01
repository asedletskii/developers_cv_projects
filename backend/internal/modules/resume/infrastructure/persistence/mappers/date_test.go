package mappers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain"
)

func TestYearMonthDate(t *testing.T) {
	value, err := domain.ParseYearMonth("2026-09")
	require.NoError(t, err)

	actual := YearMonthToTime(value)

	require.Equal(t, time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC), actual)
}

func TestYearDate(t *testing.T) {
	value, err := domain.ParseYear("2026")
	require.NoError(t, err)

	actual := YearToTime(value)

	require.Equal(t, time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC), actual)
}

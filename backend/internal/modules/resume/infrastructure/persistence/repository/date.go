package repository

import (
	"time"

	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

func yearMonthDate(value valueobjects.YearMonth) time.Time {
	return time.Date(value.Time().Year(), value.Time().Month(), 1, 0, 0, 0, 0, time.UTC)
}

func yearDate(value valueobjects.Year) time.Time {
	return time.Date(value.Int(), time.January, 1, 0, 0, 0, 0, time.UTC)
}

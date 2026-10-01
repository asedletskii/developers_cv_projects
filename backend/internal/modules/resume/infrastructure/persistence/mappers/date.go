package mappers

import (
	"time"

	"developers-cv/backend/internal/modules/resume/domain"
)

func YearMonthToTime(value domain.YearMonth) time.Time {
	return time.Date(value.Time().Year(), value.Time().Month(), 1, 0, 0, 0, 0, time.UTC)
}

func YearToTime(value domain.Year) time.Time {
	return time.Date(value.Int(), time.January, 1, 0, 0, 0, 0, time.UTC)
}

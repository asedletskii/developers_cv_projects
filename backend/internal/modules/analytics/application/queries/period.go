package queries

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrInvalidAnalyticsPeriod = errors.New("invalid analytics period")

type AnalyticsPeriod struct {
	From time.Time
	To   time.Time
	TZ   string
}

func NewAnalyticsPeriod(from, to, tz string, now time.Time) (AnalyticsPeriod, error) {
	if tz == "" {
		tz = "UTC"
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return AnalyticsPeriod{}, fmt.Errorf("%w: unknown timezone", ErrInvalidAnalyticsPeriod)
	}

	localNow := now.In(loc)
	today := time.Date(localNow.Year(), localNow.Month(), localNow.Day(), 0, 0, 0, 0, loc)
	if to == "" {
		to = today.Format("2006-01-02")
	}
	if from == "" {
		from = today.AddDate(0, 0, -29).Format("2006-01-02")
	}

	fromDate, err := time.ParseInLocation("2006-01-02", from, loc)
	if err != nil {
		return AnalyticsPeriod{}, fmt.Errorf("%w: invalid from", ErrInvalidAnalyticsPeriod)
	}
	toDate, err := time.ParseInLocation("2006-01-02", to, loc)
	if err != nil {
		return AnalyticsPeriod{}, fmt.Errorf("%w: invalid to", ErrInvalidAnalyticsPeriod)
	}
	if fromDate.After(toDate) {
		return AnalyticsPeriod{}, fmt.Errorf("%w: invalid range", ErrInvalidAnalyticsPeriod)
	}

	dayCount := calendarDayCount(fromDate, toDate.AddDate(0, 0, 1))
	if dayCount > 366 {
		return AnalyticsPeriod{}, fmt.Errorf("%w: period exceeds 366 calendar days", ErrInvalidAnalyticsPeriod)
	}
	return AnalyticsPeriod{From: fromDate, To: toDate.AddDate(0, 0, 1), TZ: tz}, nil
}

func (p AnalyticsPeriod) Days() []string {
	loc, err := time.LoadLocation(p.TZ)
	if err != nil {
		return nil
	}
	start := p.From.In(loc)
	end := p.To.In(loc)
	result := make([]string, 0, calendarDayCount(start, end))
	for day := start; day.Before(end); day = day.AddDate(0, 0, 1) {
		result = append(result, day.Format("2006-01-02"))
	}
	return result
}

func (p AnalyticsPeriod) Validate() error {
	if p.From.IsZero() || p.To.IsZero() || p.From.After(p.To) || strings.TrimSpace(p.TZ) == "" {
		return ErrInvalidAnalyticsPeriod
	}
	return nil
}

func validateAnalyticsPeriod(period AnalyticsPeriod) error {
	if err := period.Validate(); err != nil {
		return fmt.Errorf("validate analytics period: %w", err)
	}
	return nil
}

func calendarDayCount(start, end time.Time) int {
	count := 0
	for day := start; day.Before(end); day = day.AddDate(0, 0, 1) {
		count++
	}
	return count
}

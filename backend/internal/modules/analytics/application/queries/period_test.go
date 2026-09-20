package queries

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewAnalyticsPeriodUsesInclusiveCalendarDatesInTimezone(t *testing.T) {
	now := time.Date(2026, 3, 29, 12, 0, 0, 0, time.UTC)
	period, err := NewAnalyticsPeriod("2026-03-28", "2026-03-29", "Europe/Berlin", now)
	require.NoError(t, err)
	require.Equal(t, "2026-03-28", period.From.In(mustLocation(t, "Europe/Berlin")).Format("2006-01-02"))
	require.Equal(t, "2026-03-30", period.To.In(mustLocation(t, "Europe/Berlin")).Format("2006-01-02"))
	require.Equal(t, []string{"2026-03-28", "2026-03-29"}, period.Days())
}

func TestNewAnalyticsPeriodAcceptsNovosibirskTimezone(t *testing.T) {
	now := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
	period, err := NewAnalyticsPeriod("2026-09-06", "2026-10-05", "Asia/Novosibirsk", now)
	require.NoError(t, err)
	require.Equal(t, "Asia/Novosibirsk", period.TZ)
	require.Len(t, period.Days(), 30)
}

func TestNewAnalyticsPeriodRejectsInvalidTimezoneAndRange(t *testing.T) {
	now := time.Now().UTC()
	_, err := NewAnalyticsPeriod("2026-01-02", "2026-01-01", "UTC", now)
	require.ErrorIs(t, err, ErrInvalidAnalyticsPeriod)
	_, err = NewAnalyticsPeriod("2026-01-01", "2026-01-01", "Not/AZone", now)
	require.ErrorIs(t, err, ErrInvalidAnalyticsPeriod)
}

func TestVisitCursorRoundTrip(t *testing.T) {
	at := time.Date(2026, 10, 5, 10, 11, 12, 0, time.UTC)
	cursor, err := EncodeVisitCursor(at, "550e8400-e29b-41d4-a716-446655440000")
	require.NoError(t, err)
	gotAt, gotID, err := DecodeVisitCursor(cursor)
	require.NoError(t, err)
	require.Equal(t, at, gotAt)
	require.Equal(t, "550e8400-e29b-41d4-a716-446655440000", gotID)
}

func TestNormalizeVisitLimit(t *testing.T) {
	require.Equal(t, 20, mustLimit(t, 0))
	require.Equal(t, 50, mustLimit(t, 50))
	_, err := NormalizeVisitLimit(51)
	require.Error(t, err)
}

func TestNewAnalyticsPeriodAllowsAtMost366CalendarDays(t *testing.T) {
	now := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
	period, err := NewAnalyticsPeriod("2024-01-01", "2024-12-31", "UTC", now)
	require.NoError(t, err)
	require.Len(t, period.Days(), 366)

	_, err = NewAnalyticsPeriod("2024-01-01", "2025-01-01", "UTC", now)
	require.ErrorIs(t, err, ErrInvalidAnalyticsPeriod)
	_, err = NewAnalyticsPeriod("2024-01-01", "2025-01-02", "UTC", now)
	require.ErrorIs(t, err, ErrInvalidAnalyticsPeriod)
	period, err = NewAnalyticsPeriod("2025-01-01", "2026-01-01", "UTC", now)
	require.NoError(t, err)
	require.Len(t, period.Days(), 366)

	_, err = NewAnalyticsPeriod("2025-01-01", "2026-01-02", "UTC", now)
	require.ErrorIs(t, err, ErrInvalidAnalyticsPeriod)
}

func TestAnalyticsPeriodDaysHandlesDSTWithoutLosingCalendarDays(t *testing.T) {
	now := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
	period, err := NewAnalyticsPeriod("2026-03-27", "2026-03-30", "Europe/Berlin", now)
	require.NoError(t, err)
	require.Equal(t, []string{"2026-03-27", "2026-03-28", "2026-03-29", "2026-03-30"}, period.Days())
}

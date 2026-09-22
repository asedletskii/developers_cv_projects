package valueobjects

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDateRoundTrip(t *testing.T) {
	value, err := ParseDate("1998-04-12")
	require.NoError(t, err)
	require.Equal(t, "1998-04-12", value.String())
	require.False(t, value.IsZero())
}

func TestYearMonthRoundTrip(t *testing.T) {
	value, err := ParseYearMonth("2026-09")
	require.NoError(t, err)
	require.Equal(t, "2026-09", value.String())
	require.False(t, value.IsZero())
}

func TestYearRoundTrip(t *testing.T) {
	value, err := ParseYear("2026")
	require.NoError(t, err)
	require.Equal(t, "2026", value.String())
	require.False(t, value.IsZero())
}

func TestDateNormalizesToUTCDateOnly(t *testing.T) {
	value, err := NewDate(time.Date(2026, time.September, 22, 23, 59, 0, 0, time.FixedZone("test", 3*60*60)))
	require.NoError(t, err)
	require.Equal(t, "2026-09-22", value.String())
	require.Equal(t, time.UTC, value.Time().Location())
}

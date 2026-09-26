package readers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDateFormatting(t *testing.T) {
	value := time.Date(2026, time.September, 17, 13, 45, 0, 0, time.UTC)

	require.Equal(t, "2026-09", yearMonth(value))

	yearMonthValue := value
	require.Equal(t, "2026-09", *optionalYearMonth(&yearMonthValue))

	require.Equal(t, "2026", *optionalYear(&yearMonthValue))
	require.Nil(t, optionalYearMonth(nil))
	require.Nil(t, optionalYear(nil))
}

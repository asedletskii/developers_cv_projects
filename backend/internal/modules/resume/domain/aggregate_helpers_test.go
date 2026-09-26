package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeOptionalStringTrimsWhitespaceAndTurnsEmptyIntoNil(t *testing.T) {
	value := "  optional value  "
	got := cloneString(&value)
	require.NotNil(t, got)
	require.Equal(t, "optional value", *got)

	empty := "   "
	require.Nil(t, cloneString(&empty))
	require.Nil(t, cloneString(nil))
}

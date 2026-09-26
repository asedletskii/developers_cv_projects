package valueobjects

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestIDRoundTrip(t *testing.T) {
	id := NewID()
	parsed, err := ParseID(id.String())
	require.NoError(t, err)
	require.Equal(t, id, parsed)
	require.False(t, parsed.IsZero())
}

func TestIDRejectsZeroUUID(t *testing.T) {
	_, err := IDFromUUID(uuid.Nil)
	require.ErrorIs(t, err, ErrInvalidID)
}

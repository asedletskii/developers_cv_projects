package valueobjects

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewProviderSubjectPreservesOpaqueValue(t *testing.T) {
	value, err := NewProviderSubject("  opaque-subject/123  ")
	require.NoError(t, err)
	require.Equal(t, "  opaque-subject/123  ", value.String())
}

func TestNewProviderSubjectRejectsEmpty(t *testing.T) {
	_, err := NewProviderSubject("   ")
	require.Error(t, err)
}

func TestNewProviderSubjectRejectsTooLong(t *testing.T) {
	_, err := NewProviderSubject(strings.Repeat("a", MaxProviderSubjectLength+1))
	require.Error(t, err)
}

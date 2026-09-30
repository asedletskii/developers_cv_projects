package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewProvider(t *testing.T) {
	provider, err := NewProvider(" telegram ")
	require.NoError(t, err)
	require.Equal(t, ProviderTelegram, provider)
	require.True(t, provider.IsValid())
}

func TestNewProviderRejectsUnknown(t *testing.T) {
	_, err := NewProvider("google")
	require.ErrorIs(t, err, ErrInvalidProvider)
}

func TestNewProviderRejectsEmpty(t *testing.T) {
	_, err := NewProvider("   ")
	require.ErrorIs(t, err, ErrInvalidProvider)
}

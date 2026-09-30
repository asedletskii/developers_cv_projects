package valueobjects

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSessionTokenHashRequiresSHA256Length(t *testing.T) {
	_, err := NewSessionTokenHash(make([]byte, HashSize-1))
	require.Error(t, err)
}

func TestStateHashCopiesInput(t *testing.T) {
	input := bytes.Repeat([]byte{1}, HashSize)
	value, err := NewStateHash(input)
	require.NoError(t, err)
	input[0] = 9
	require.Equal(t, byte(1), value.Bytes()[0])
}

func TestNonceHashRequiresSHA256Length(t *testing.T) {
	_, err := NewNonceHash(make([]byte, HashSize+1))
	require.Error(t, err)
}

func TestBrowserBindingHashRequiresSHA256Length(t *testing.T) {
	_, err := NewBrowserBindingHash(make([]byte, HashSize-1))
	require.Error(t, err)
}

func TestBrowserBindingHashCopiesInput(t *testing.T) {
	input := bytes.Repeat([]byte{4}, HashSize)
	value, err := NewBrowserBindingHash(input)
	require.NoError(t, err)
	input[0] = 9
	require.Equal(t, byte(4), value.Bytes()[0])
}

func TestHashesExposeCopies(t *testing.T) {
	input := bytes.Repeat([]byte{3}, HashSize)
	value, err := NewSessionTokenHash(input)
	require.NoError(t, err)
	output := value.Bytes()
	output[0] = 8
	require.Equal(t, byte(3), value.Bytes()[0])
}

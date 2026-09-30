package token

import (
	"bytes"
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGeneratorProduces32RandomBytes(t *testing.T) {
	generator := NewGenerator()
	values := [][]byte{
		mustGenerate(t, generator.GenerateState),
		mustGenerate(t, generator.GenerateBrowserBinding),
		mustGenerate(t, generator.GenerateNonce),
		mustGenerate(t, generator.GeneratePKCEVerifier),
		mustGenerate(t, generator.GenerateSessionToken),
	}
	for _, value := range values {
		require.Len(t, value, randomTokenSize)
	}
	require.NotEqual(t, values[0], values[1])
}

func TestGeneratorHashIsSHA256AndDoesNotAliasInput(t *testing.T) {
	generator := NewGenerator()
	input := []byte("value")
	hash := generator.Hash(input)
	expected := sha256.Sum256(input)
	require.Equal(t, expected[:], hash)
	require.False(t, bytes.Equal(hash, input))
	input[0] = 'X'
	require.Equal(t, expected[:], hash)
}

func mustGenerate(t *testing.T, generate func() ([]byte, error)) []byte {
	t.Helper()
	value, err := generate()
	require.NoError(t, err)
	return value
}

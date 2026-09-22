package valueobjects

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExternalURLAllowsHTTPAndHTTPS(t *testing.T) {
	for _, value := range []string{"https://example.com/a", "http://example.com/a"} {
		url, err := NewExternalURL(value)
		require.NoError(t, err)
		require.Equal(t, value, url.String())
	}
}

func TestExternalURLRejectsUnsupportedScheme(t *testing.T) {
	_, err := NewExternalURL("javascript:alert(1)")
	require.Error(t, err)
}

func TestAssetPathRequiresInternalAbsolutePath(t *testing.T) {
	path, err := NewAssetPath("/mock/avatar.svg")
	require.NoError(t, err)
	require.Equal(t, "/mock/avatar.svg", path.String())

	_, err = NewAssetPath("https://example.com/avatar.svg")
	require.Error(t, err)
}

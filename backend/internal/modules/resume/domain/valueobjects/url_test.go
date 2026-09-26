package valueobjects

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExternalURLAllowsHTTPAndHTTPS(t *testing.T) {
	for _, value := range []string{
		"https://example.com/avatar.svg",
		"http://example.com/avatar.svg",
	} {
		url, err := NewExternalURL(value)
		require.NoError(t, err)
		require.Equal(t, value, url.String())
	}
}

func TestExternalURLRejectsUnsupportedScheme(t *testing.T) {
	_, err := NewExternalURL("javascript:alert(1)")
	require.Error(t, err)
}

func TestExternalURLRejectsRelativePath(t *testing.T) {
	_, err := NewExternalURL("/assets/avatar.svg")
	require.Error(t, err)
}

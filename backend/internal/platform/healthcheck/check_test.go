package healthcheck

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPortFromEnvUsesDefault(t *testing.T) {
	t.Setenv("HTTP_PORT", "")

	port, err := portFromEnv()

	require.NoError(t, err)
	require.Equal(t, defaultPort, port)
}

func TestPortFromEnvParsesConfiguredPort(t *testing.T) {
	t.Setenv("HTTP_PORT", "9090")

	port, err := portFromEnv()

	require.NoError(t, err)
	require.Equal(t, 9090, port)
}

func TestPortFromEnvRejectsInvalidPort(t *testing.T) {
	t.Setenv("HTTP_PORT", "invalid")

	_, err := portFromEnv()

	require.Error(t, err)
}

func TestCheckReturnsErrorWhenServerIsUnavailable(t *testing.T) {
	t.Setenv("HTTP_PORT", "65535")

	err := Check(context.Background())

	require.Error(t, err)
}

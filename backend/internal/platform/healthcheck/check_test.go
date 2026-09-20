package healthcheck

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
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

func TestCheckSucceedsWhenLivenessEndpointReturnsOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/health", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	setHealthcheckPort(t, server)
	require.NoError(t, Check(context.Background()))
}

func TestCheckRejectsNonOKLivenessEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusServiceUnavailable) }))
	defer server.Close()
	setHealthcheckPort(t, server)
	require.Error(t, Check(context.Background()))
}

func setHealthcheckPort(t *testing.T, server *httptest.Server) {
	t.Helper()
	address, ok := server.Listener.Addr().(*net.TCPAddr)
	require.True(t, ok, "httptest server listener must use TCP")
	t.Setenv("HTTP_PORT", strconv.Itoa(address.Port))
}

func TestCheckReturnsErrorWhenServerIsUnavailable(t *testing.T) {
	t.Setenv("HTTP_PORT", "65535")

	err := Check(context.Background())

	require.Error(t, err)
}

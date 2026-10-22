package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"developers-cv/backend/internal/platform/config"
)

func TestServerRunAndShutdownStopsPublicAndInternalListeners(t *testing.T) {
	publicAddr := availableTCPAddress(t)
	internalAddr := availableTCPAddress(t)
	publicHost, publicPort, err := net.SplitHostPort(publicAddr)
	require.NoError(t, err)
	port, err := strconv.Atoi(publicPort)
	require.NoError(t, err)

	cfg := config.HTTPConfig{
		Host:              publicHost,
		Port:              port,
		MetricsAddr:       internalAddr,
		ReadTimeout:       time.Second,
		ReadHeaderTimeout: time.Second,
		WriteTimeout:      time.Second,
		IdleTimeout:       time.Second,
		ShutdownTimeout:   2 * time.Second,
	}
	publicHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) })
	internalHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) })
	server := NewServer(cfg, publicHandler, zap.NewNop(), internalHandler)

	runErr := make(chan error, 1)
	go func() { runErr <- server.Run() }()

	waitForHTTPStatus(t, "http://"+publicAddr, http.StatusNoContent)
	waitForHTTPStatus(t, "http://"+internalAddr, http.StatusNoContent)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	require.NoError(t, server.Shutdown(ctx))
	select {
	case err := <-runErr:
		require.NoError(t, err)
	case <-time.After(3 * time.Second):
		t.Fatal("Run did not return after Shutdown")
	}

	assertHTTPUnavailable(t, "http://"+publicAddr)
	assertHTTPUnavailable(t, "http://"+internalAddr)
}

func TestServerRunClosesPublicListenerWhenInternalListenerFails(t *testing.T) {
	publicAddr := availableTCPAddress(t)
	internalListener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, internalListener.Close())
	})
	internalAddr := internalListener.Addr().String()
	publicHost, publicPort, err := net.SplitHostPort(publicAddr)
	require.NoError(t, err)
	port, err := strconv.Atoi(publicPort)
	require.NoError(t, err)

	cfg := config.HTTPConfig{
		Host:              publicHost,
		Port:              port,
		MetricsAddr:       internalAddr,
		ReadTimeout:       time.Second,
		ReadHeaderTimeout: time.Second,
		WriteTimeout:      time.Second,
		IdleTimeout:       time.Second,
		ShutdownTimeout:   time.Second,
	}
	server := NewServer(cfg, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}), zap.NewNop(), http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))

	err = server.Run()
	require.Error(t, err)
	require.Contains(t, err.Error(), "listen and serve")

	assertCanListen(t, publicAddr)
}

func availableTCPAddress(t *testing.T) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := listener.Addr().String()
	require.NoError(t, listener.Close())
	return addr
}

func waitForHTTPStatus(t *testing.T, endpoint string, expectedStatus int) {
	t.Helper()
	client := &http.Client{Timeout: 100 * time.Millisecond}
	deadline := time.Now().Add(3 * time.Second)
	var lastErr error
	for time.Now().Before(deadline) {
		resp, err := client.Get(endpoint)
		if err == nil {
			statusCode := resp.StatusCode
			if closeErr := resp.Body.Close(); closeErr != nil {
				lastErr = fmt.Errorf("close response body from %s: %w", endpoint, closeErr)
				time.Sleep(10 * time.Millisecond)
				continue
			}
			if statusCode == expectedStatus {
				return
			}
			lastErr = fmt.Errorf("%s returned HTTP %d, expected %d", endpoint, statusCode, expectedStatus)
		} else {
			lastErr = err
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("server did not become ready at %s: %v", endpoint, lastErr)
}

func assertHTTPUnavailable(t *testing.T, endpoint string) {
	t.Helper()
	client := &http.Client{Timeout: 100 * time.Millisecond}
	resp, err := client.Get(endpoint)
	if resp != nil {
		require.NoError(t, resp.Body.Close())
	}
	require.Error(t, err, "expected %s to stop accepting connections", endpoint)
}

func assertCanListen(t *testing.T, addr string) {
	t.Helper()
	listener, err := net.Listen("tcp", addr)
	require.NoError(t, err, "expected address %s to be released", addr)
	require.NoError(t, listener.Close())
}

func TestEffectiveShutdownTimeoutUsesSecondsForDefault(t *testing.T) {
	require.Equal(t, 10*time.Second, effectiveShutdownTimeout(0))
	require.Equal(t, 10*time.Second, effectiveShutdownTimeout(-time.Second))
	require.Equal(t, 3*time.Second, effectiveShutdownTimeout(3*time.Second))
}

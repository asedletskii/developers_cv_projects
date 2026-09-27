package healthcheck

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	defaultHost = "127.0.0.1"
	defaultPort = 8080
	timeout     = 2 * time.Second
)

func Check(ctx context.Context) error {
	port, err := portFromEnv()
	if err != nil {
		return err
	}

	dialer := net.Dialer{Timeout: timeout}
	address := net.JoinHostPort(defaultHost, strconv.Itoa(port))
	requestCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	connection, err := dialer.DialContext(requestCtx, "tcp", address)
	if err != nil {
		return fmt.Errorf("connect to HTTP server: %w", err)
	}
	if err := connection.Close(); err != nil {
		return fmt.Errorf("close healthcheck connection: %w", err)
	}
	return nil
}

func portFromEnv() (int, error) {
	value := os.Getenv("HTTP_PORT")
	if value == "" {
		return defaultPort, nil
	}

	port, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("parse HTTP_PORT: %w", err)
	}
	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("HTTP_PORT must be between 1 and 65535: %d", port)
	}
	return port, nil
}

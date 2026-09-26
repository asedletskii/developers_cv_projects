package healthcheck

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	defaultHost = "127.0.0.1"
	defaultPort = 8080
	defaultPath = "/ready"
	timeout     = 2 * time.Second
)

func Check(ctx context.Context) error {
	port, err := portFromEnv()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s:%d%s", defaultHost, port, defaultPath)
	requestCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(requestCtx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("create readiness request: %w", err)
	}

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request readiness endpoint: %w", err)
	}

	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil && err == nil {
			err = fmt.Errorf("close readiness response body: %w", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("readiness endpoint returned status %d", resp.StatusCode)
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

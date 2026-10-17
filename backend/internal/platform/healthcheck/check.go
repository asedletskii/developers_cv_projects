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
	timeout     = 2 * time.Second
)

func Check(ctx context.Context) (checkErr error) {
	port, err := portFromEnv()
	if err != nil {
		return err
	}
	address := fmt.Sprintf("http://%s:%d/health", defaultHost, port)
	requestCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	request, err := http.NewRequestWithContext(requestCtx, http.MethodGet, address, nil)
	if err != nil {
		return fmt.Errorf("create liveness request: %w", err)
	}
	client := &http.Client{Timeout: timeout, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("request liveness endpoint: %w", err)
	}
	defer func() {
		if closeErr := response.Body.Close(); closeErr != nil && checkErr == nil {
			checkErr = fmt.Errorf("close liveness response body: %w", closeErr)
		}
	}()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("liveness endpoint returned status %d", response.StatusCode)
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

package security

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"developers-cv/backend/internal/modules/resume_link/application/ports"
)

const (
	DefaultSiteverifyURL = "https://challenges.cloudflare.com/turnstile/v0/siteverify"
)

var (
	ErrVerificationHTTP = errors.New("anti-bot verification service returned an invalid response")
)

type Config struct {
	Secret         string
	SiteverifyURL  string
	Timeout        time.Duration
	ExpectedHost   string
	ExpectedAction string
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type TurnstileVerifier struct {
	secret         string
	siteverifyURL  string
	timeout        time.Duration
	expectedHost   string
	expectedAction string
	client         HTTPClient
}

type siteverifyResponse struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes"`
	Hostname   string   `json:"hostname"`
	Action     string   `json:"action"`
}

func NewTurnstileVerifier(cfg Config, client HTTPClient) (*TurnstileVerifier, error) {
	if strings.TrimSpace(cfg.Secret) == "" {
		return nil, errors.New("turnstile secret must not be empty")
	}
	if client == nil {
		return nil, errors.New("turnstile HTTP client must not be nil")
	}
	if cfg.Timeout <= 0 {
		return nil, errors.New("turnstile timeout must be positive")
	}
	if cfg.SiteverifyURL == "" {
		cfg.SiteverifyURL = DefaultSiteverifyURL
	}
	parsed, err := url.Parse(cfg.SiteverifyURL)
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" || parsed.User != nil {
		return nil, fmt.Errorf("invalid turnstile siteverify URL")
	}

	return &TurnstileVerifier{
		secret:         cfg.Secret,
		siteverifyURL:  cfg.SiteverifyURL,
		timeout:        cfg.Timeout,
		expectedHost:   strings.TrimSpace(cfg.ExpectedHost),
		expectedAction: strings.TrimSpace(cfg.ExpectedAction),
		client:         client,
	}, nil
}

func (v *TurnstileVerifier) Verify(ctx context.Context, token string, remoteIP string) (err error) {
	if strings.TrimSpace(token) == "" || len(token) > ports.MaxTokenLength {
		return ports.ErrTokenVerificationFailed
	}

	verifyCtx, cancel := context.WithTimeout(ctx, v.timeout)
	defer cancel()

	values := url.Values{}
	values.Set("secret", v.secret)
	values.Set("response", token)
	if strings.TrimSpace(remoteIP) != "" {
		values.Set("remoteip", remoteIP)
	}

	req, err := http.NewRequestWithContext(verifyCtx, http.MethodPost, v.siteverifyURL, strings.NewReader(values.Encode()))
	if err != nil {
		return ErrVerificationHTTP
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := v.client.Do(req)
	if err != nil {
		return ErrVerificationHTTP
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
			err = ErrVerificationHTTP
		}
	}()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return ErrVerificationHTTP
	}

	var result siteverifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ErrVerificationHTTP
	}

	if !result.Success {
		return ports.ErrTokenVerificationFailed
	}
	if v.expectedHost != "" && result.Hostname != v.expectedHost {
		return ports.ErrTokenVerificationFailed
	}
	if v.expectedAction != "" && result.Action != v.expectedAction {
		return ports.ErrTokenVerificationFailed
	}

	return nil
}

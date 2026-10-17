package telegram

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"developers-cv/backend/internal/modules/auth/application/ports"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	"developers-cv/backend/internal/platform/telemetry"
)

const (
	discoveryTimeout  = 10 * time.Second
	initialRetryDelay = time.Second
	maximumRetryDelay = 30 * time.Second
)

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	IssuerURL    string
}

type oidcAdapter struct {
	cfg     Config
	client  *http.Client
	metrics *telemetry.Metrics

	mu           sync.Mutex
	provider     *oidcProvider
	initializing chan struct{}
	lastErr      error
	nextRetryAt  time.Time
	retryDelay   time.Duration
}

func New(cfg Config, client *http.Client, metrics *telemetry.Metrics) (ports.TelegramOIDC, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}
	if client == nil {
		return nil, errors.New("telegram oidc http client is nil")
	}

	return &oidcAdapter{
		cfg: cfg, client: client, metrics: metrics, retryDelay: initialRetryDelay,
	}, nil
}

func (a *oidcAdapter) AuthorizationURL(ctx context.Context, params ports.TelegramAuthorizationParams) (string, error) {
	provider, err := a.getProvider(ctx)
	if err != nil {
		return "", err
	}
	return provider.buildAuthorizationURL(params)
}

func (a *oidcAdapter) Exchange(ctx context.Context, code string, redirectURI string, verifier string, expectedNonceHash valueobjects.NonceHash) (ports.VerifiedIdentity, error) {
	provider, err := a.getProvider(ctx)
	if err != nil {
		return ports.VerifiedIdentity{}, err
	}
	return provider.exchangeAndVerifyCode(ctx, code, redirectURI, verifier, expectedNonceHash)
}

func (a *oidcAdapter) getProvider(ctx context.Context) (*oidcProvider, error) {
	for {
		a.mu.Lock()
		if a.provider != nil {
			provider := a.provider
			a.mu.Unlock()
			return provider, nil
		}
		if a.initializing != nil {
			initialized := a.initializing
			a.mu.Unlock()
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-initialized:
				continue
			}
		}
		if time.Now().Before(a.nextRetryAt) {
			retryAt, lastErr := a.nextRetryAt, a.lastErr
			a.mu.Unlock()
			return nil, fmt.Errorf("%w: discovery retry is delayed until %s: %v", ports.ErrProviderUnavailable, retryAt.UTC().Format(time.RFC3339), lastErr)
		}

		initialized := make(chan struct{})
		a.initializing = initialized
		cfg, client, metrics := a.cfg, a.client, a.metrics
		a.mu.Unlock()

		initCtx, cancel := context.WithTimeout(ctx, discoveryTimeout)
		provider, err := discoverProvider(initCtx, cfg, client, metrics)
		cancel()

		a.mu.Lock()
		if err != nil {
			a.lastErr = err
			a.nextRetryAt = time.Now().Add(a.retryDelay)
			a.retryDelay *= 2
			if a.retryDelay > maximumRetryDelay {
				a.retryDelay = maximumRetryDelay
			}
		} else {
			a.provider = provider
			a.lastErr = nil
			a.nextRetryAt = time.Time{}
			a.retryDelay = initialRetryDelay
		}
		a.initializing = nil
		close(initialized)
		a.mu.Unlock()

		if err != nil {
			return nil, fmt.Errorf("%w: discover Telegram OIDC provider: %v", ports.ErrProviderUnavailable, err)
		}
		return provider, nil
	}
}

func validateConfig(cfg Config) error {
	if strings.TrimSpace(cfg.ClientID) == "" {
		return errors.New("telegram oidc client id is empty")
	}
	if strings.TrimSpace(cfg.ClientSecret) == "" {
		return errors.New("telegram oidc client secret is empty")
	}
	if strings.TrimSpace(cfg.RedirectURL) == "" {
		return errors.New("telegram oidc redirect url is empty")
	}
	if strings.TrimSpace(cfg.IssuerURL) == "" {
		return errors.New("telegram oidc issuer url is empty")
	}
	return nil
}

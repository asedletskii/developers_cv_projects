package telegram

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/oauth2"

	"developers-cv/backend/internal/modules/auth/application/ports"
	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	"developers-cv/backend/internal/platform/telemetry"
)

const (
	IssuerURL          = "https://oauth.telegram.org"
	authorizationScope = "profile"
	allowedAlgorithm   = "RS256"
)

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	IssuerURL    string
}

type TelegramOIDC struct {
	oauth2Config oauth2.Config
	verifier     *oidc.IDTokenVerifier
	tracer       trace.Tracer
	metrics      *telemetry.Metrics
}

func New(ctx context.Context, cfg Config, client *http.Client, metrics ...*telemetry.Metrics) (*TelegramOIDC, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}
	if client == nil {
		return nil, errors.New("telegram oidc http client is nil")
	}

	providerCtx := oidc.ClientContext(ctx, client)
	provider, err := oidc.NewProvider(providerCtx, cfg.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("discover telegram oidc provider: %w", err)
	}

	oauth2Config := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  cfg.RedirectURL,
		Scopes:       []string{oidc.ScopeOpenID, authorizationScope},
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID:             cfg.ClientID,
		SupportedSigningAlgs: []string{allowedAlgorithm},
	})

	var metricsValue *telemetry.Metrics
	if len(metrics) > 0 {
		metricsValue = metrics[0]
	}
	return &TelegramOIDC{oauth2Config: oauth2Config, verifier: verifier, tracer: otel.Tracer("developers-cv-api/auth/telegram"), metrics: metricsValue}, nil
}

func (t *TelegramOIDC) AuthorizationURL(params ports.TelegramAuthorizationParams) string {
	if t == nil {
		return ""
	}

	return t.oauth2Config.AuthCodeURL(
		params.State,
		oauth2.SetAuthURLParam("nonce", params.Nonce),
		oauth2.SetAuthURLParam("code_challenge", params.CodeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)
}

func (t *TelegramOIDC) Exchange(ctx context.Context, code string, redirectURI string, verifier string, expectedNonceHash valueobjects.NonceHash) (ports.VerifiedIdentity, error) {
	if t == nil || t.verifier == nil {
		return ports.VerifiedIdentity{}, errors.New("telegram oidc adapter is nil")
	}
	if strings.TrimSpace(code) == "" {
		return ports.VerifiedIdentity{}, errors.New("authorization code is empty")
	}
	if strings.TrimSpace(verifier) == "" {
		return ports.VerifiedIdentity{}, errors.New("pkce verifier is empty")
	}
	if strings.TrimSpace(redirectURI) == "" {
		return ports.VerifiedIdentity{}, errors.New("redirect uri is empty")
	}
	if redirectURI != t.oauth2Config.RedirectURL {
		return ports.VerifiedIdentity{}, errors.New("redirect uri does not match configured telegram redirect uri")
	}

	config := t.oauth2Config
	exchangeCtx, exchangeSpan := t.tracer.Start(ctx, "auth.telegram.exchange")
	exchangeStartedAt := time.Now()
	token, err := config.Exchange(exchangeCtx, code, oauth2.SetAuthURLParam("code_verifier", verifier))
	if t.metrics != nil {
		t.metrics.ObserveAuthOIDCExchange(string(domain.ProviderTelegram), time.Since(exchangeStartedAt))
	}
	exchangeSpan.SetAttributes(attribute.String("auth.provider", string(domain.ProviderTelegram)), attribute.String("auth.result", resultFromError(err)))
	if err != nil {
		exchangeSpan.RecordError(err)
	}
	exchangeSpan.End()
	if err != nil {
		return ports.VerifiedIdentity{}, fmt.Errorf("%w: %v", ports.ErrTokenExchangeFailed, err)
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok || strings.TrimSpace(rawIDToken) == "" {
		return ports.VerifiedIdentity{}, errors.New("telegram token response does not contain id_token")
	}

	verifyCtx, verifySpan := t.tracer.Start(ctx, "auth.telegram.verify")
	verifyStartedAt := time.Now()
	idToken, err := t.verifier.Verify(verifyCtx, rawIDToken)
	if t.metrics != nil {
		t.metrics.ObserveAuthOIDCVerify(string(domain.ProviderTelegram), time.Since(verifyStartedAt))
	}
	verifySpan.SetAttributes(attribute.String("auth.provider", string(domain.ProviderTelegram)), attribute.String("auth.result", resultFromError(err)))
	if err != nil {
		verifySpan.RecordError(err)
	}
	verifySpan.End()
	if err != nil {
		return ports.VerifiedIdentity{}, fmt.Errorf("%w: %v", ports.ErrIDTokenInvalid, err)
	}

	var claims struct {
		Subject string `json:"sub"`
		Nonce   string `json:"nonce"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return ports.VerifiedIdentity{}, fmt.Errorf("decode telegram id token claims: %w", err)
	}
	nonceHash := sha256.Sum256([]byte(claims.Nonce))
	if !bytes.Equal(nonceHash[:], expectedNonceHash.Bytes()) {
		return ports.VerifiedIdentity{}, ports.ErrNonceMismatch
	}

	provider, err := domain.NewProvider(string(domain.ProviderTelegram))
	if err != nil {
		return ports.VerifiedIdentity{}, fmt.Errorf("create telegram provider: %w", err)
	}
	subject, err := valueobjects.NewProviderSubject(claims.Subject)
	if err != nil {
		return ports.VerifiedIdentity{}, fmt.Errorf("create telegram provider subject: %w", err)
	}

	return ports.VerifiedIdentity{Provider: provider, ProviderSubject: subject}, nil
}

func NewHTTPClient(timeout time.Duration) *http.Client {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &http.Client{Timeout: timeout}
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

var _ ports.TelegramOIDC = (*TelegramOIDC)(nil)

func resultFromError(err error) string {
	if err != nil {
		return "failure"
	}
	return "success"
}

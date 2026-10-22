package config

import (
	"errors"
	"net"
	"net/url"
	"strings"
)

func validateProductionSecurity(cfg Config) error {
	if cfg.App.Env != "production" {
		return nil
	}
	if cfg.Auth.Mode != "telegram" {
		return errors.New("AUTH_MODE must be telegram in production")
	}
	if strings.TrimSpace(cfg.App.Version) == "" || strings.EqualFold(strings.TrimSpace(cfg.App.Version), "dev") {
		return errors.New("APP_VERSION must be set to a release version in production")
	}
	if strings.TrimSpace(cfg.HTTP.MetricsAddr) == "" {
		return errors.New("METRICS_ADDR must be a loopback host:port in production")
	}
	metricsHost, metricsPort, err := net.SplitHostPort(cfg.HTTP.MetricsAddr)
	if err != nil || !isValidTCPPort(metricsPort) {
		return errors.New("METRICS_ADDR must be a loopback IP address with a valid TCP port in production")
	}
	metricsIP := net.ParseIP(metricsHost)
	if metricsIP == nil || !metricsIP.IsLoopback() {
		return errors.New("METRICS_ADDR must use a loopback IP address in production")
	}
	redirectURL, err := url.ParseRequestURI(cfg.Auth.Telegram.RedirectURL)
	if err != nil || redirectURL.Scheme != "https" || redirectURL.Host == "" || redirectURL.User != nil || redirectURL.Fragment != "" {
		return errors.New("AUTH_TELEGRAM_REDIRECT_URL must be an https URL in production")
	}
	for _, origin := range cfg.CORS.AllowedOrigins {
		parsed, err := url.ParseRequestURI(origin)
		if err != nil || parsed.Scheme != "https" || parsed.Host == "" ||
			parsed.User != nil || parsed.Path != "" || parsed.RawQuery != "" || parsed.Fragment != "" {
			return errors.New("CORS_ALLOWED_ORIGINS must contain only https origins in production")
		}
	}
	if len(cfg.HTTP.TrustedProxies) == 0 {
		return errors.New("HTTP_TRUSTED_PROXIES must contain the address or subnet of the trusted reverse proxy in production")
	}
	if !cfg.Auth.Cookie.Secure || cfg.Auth.Cookie.Name != "__Host-session" {
		return errors.New("production auth cookie must be secure __Host-session")
	}

	if !cfg.HTTP.RateLimit.Enabled {
		return errors.New("RATE_LIMIT_ENABLED must be true in production")
	}
	if !cfg.HTTP.ContactViewRateLimit.Enabled {
		return errors.New("CONTACT_VIEW_RATE_LIMIT_ENABLED must be true in production")
	}
	if !cfg.HTTP.AnalyticsEventsRateLimit.Enabled {
		return errors.New("ANALYTICS_EVENTS_RATE_LIMIT_ENABLED must be true in production")
	}
	if !cfg.Auth.RateLimit.Enabled {
		return errors.New("AUTH_RATE_LIMIT_ENABLED must be true in production")
	}
	if isPlaceholderSecret(cfg.Analytics.VisitorSecret) {
		return errors.New("ANALYTICS_VISITOR_SECRET must be a deployment-specific secret in production")
	}
	if isPlaceholderSecret(cfg.Postgres.Password) {
		return errors.New("POSTGRES_PASSWORD must be a deployment-specific secret in production")
	}
	if isPlaceholderSecret(cfg.Auth.Telegram.ClientID) {
		return errors.New("AUTH_TELEGRAM_CLIENT_ID must be configured in production")
	}
	if isPlaceholderSecret(cfg.Auth.Telegram.ClientSecret) {
		return errors.New("AUTH_TELEGRAM_CLIENT_SECRET must be a deployment-specific secret in production")
	}
	if !cfg.Turnstile.Enabled {
		return errors.New("TURNSTILE_ENABLED must be true in production")
	}
	if isPlaceholderSecret(cfg.Turnstile.SecretKey) {
		return errors.New("TURNSTILE_SECRET_KEY must be a deployment-specific secret in production")
	}
	if strings.TrimSpace(cfg.Turnstile.ExpectedHost) == "" {
		return errors.New("TURNSTILE_EXPECTED_HOSTNAME is required in production")
	}
	if strings.TrimSpace(cfg.Turnstile.ExpectedAction) == "" {
		return errors.New("TURNSTILE_EXPECTED_ACTION is required in production")
	}
	if cfg.Telemetry.Enabled && cfg.Telemetry.Insecure {
		return errors.New("OTEL_EXPORTER_OTLP_INSECURE must be false in production")
	}
	return nil
}

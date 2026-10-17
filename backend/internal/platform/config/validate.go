package config

import (
	"errors"
	"fmt"
	"math"
	"net"
	"net/url"
	"strconv"
	"strings"
)

func validate(cfg Config) error {
	switch cfg.App.Env {
	case "development", "test", "staging", "production":
	default:
		return fmt.Errorf("invalid APP_ENV %q: expected development, test, staging or production", cfg.App.Env)
	}
	if len([]byte(strings.TrimSpace(cfg.Analytics.VisitorSecret))) < 32 {
		return errors.New("ANALYTICS_VISITOR_SECRET must be at least 32 bytes")
	}
	if cfg.Analytics.SuspectBurst <= 0 || cfg.Analytics.SuspectWindow <= 0 {
		return errors.New("analytics suspect limits must be positive")
	}
	if cfg.Analytics.ClickRetentionDays <= 0 || cfg.Analytics.ClickRetentionDays > 3650 {
		return errors.New("ANALYTICS_CLICK_RETENTION_DAYS must be between 1 and 3650")
	}
	if cfg.Analytics.ViewRetentionDays <= 0 || cfg.Analytics.ViewRetentionDays > 3650 {
		return errors.New("ANALYTICS_VIEW_RETENTION_DAYS must be between 1 and 3650")
	}
	if cfg.Auth.Mode != "mock" && cfg.Auth.Mode != "telegram" {
		return fmt.Errorf("invalid AUTH_MODE: %q", cfg.Auth.Mode)
	}
	if cfg.Auth.Mode == "mock" && cfg.App.Env != "development" && cfg.App.Env != "test" {
		return errors.New("AUTH_MODE=mock is allowed only when APP_ENV is development or test")
	}
	if cfg.Auth.Session.AbsoluteTTL <= 0 || cfg.Auth.Session.IdleTTL <= 0 || cfg.Auth.Session.TouchInterval <= 0 {
		return errors.New("auth session durations must be positive")
	}
	if cfg.Auth.Session.IdleTTL > cfg.Auth.Session.AbsoluteTTL {
		return errors.New("AUTH_SESSION_IDLE_TTL must not exceed AUTH_SESSION_ABSOLUTE_TTL")
	}
	if cfg.Auth.Transaction.TTL <= 0 {
		return errors.New("AUTH_TRANSACTION_TTL must be positive")
	}
	if cfg.Auth.RateLimit.Enabled && (!isPositiveFinite(cfg.Auth.RateLimit.StartRatePerMinute) ||
		cfg.Auth.RateLimit.StartBurst <= 0 ||
		!isPositiveFinite(cfg.Auth.RateLimit.CallbackRatePerMinute) ||
		cfg.Auth.RateLimit.CallbackBurst <= 0) {
		return errors.New("auth rate limits must be positive when enabled")
	}
	if strings.TrimSpace(cfg.Auth.Cookie.Name) == "" {
		return errors.New("AUTH_COOKIE_NAME must not be empty")
	}
	if strings.HasPrefix(cfg.Auth.Cookie.Name, "__Host-") && !cfg.Auth.Cookie.Secure {
		return errors.New("__Host- auth cookie requires AUTH_COOKIE_SECURE=true")
	}
	if cfg.Auth.Mode == "telegram" {
		if strings.TrimSpace(cfg.Auth.Telegram.ClientID) == "" ||
			strings.TrimSpace(cfg.Auth.Telegram.ClientSecret) == "" ||
			strings.TrimSpace(cfg.Auth.Telegram.RedirectURL) == "" ||
			strings.TrimSpace(cfg.Auth.Telegram.IssuerURL) == "" {
			return errors.New("telegram auth configuration is incomplete")
		}
	}
	if cfg.Auth.Telegram.IssuerURL != "" && cfg.Auth.Telegram.IssuerURL != "https://oauth.telegram.org" && cfg.App.Env == "production" {
		return errors.New("AUTH_TELEGRAM_ISSUER_URL must be https://oauth.telegram.org in production")
	}

	if strings.TrimSpace(cfg.Postgres.Host) == "" ||
		strings.TrimSpace(cfg.Postgres.Database) == "" ||
		strings.TrimSpace(cfg.Postgres.User) == "" ||
		strings.TrimSpace(cfg.Postgres.Password) == "" {
		return errors.New("PostgreSQL connection credentials must not be empty")
	}
	if cfg.Postgres.Port < 1 || cfg.Postgres.Port > 65535 {
		return fmt.Errorf("POSTGRES_PORT must be between 1 and 65535: %d", cfg.Postgres.Port)
	}
	if !isValidPostgresSSLMode(cfg.Postgres.SSLMode) {
		return fmt.Errorf("invalid POSTGRES_SSLMODE: %q", cfg.Postgres.SSLMode)
	}
	if cfg.Postgres.MaxConns <= 0 ||
		cfg.Postgres.MaxIdleConns < 0 ||
		cfg.Postgres.MaxIdleConns > cfg.Postgres.MaxConns {
		return errors.New("invalid PostgreSQL pool limits")
	}
	if cfg.Postgres.MaxConnLifetime <= 0 ||
		cfg.Postgres.MaxConnIdleTime <= 0 {
		return errors.New("PostgreSQL connection lifetimes must be positive")
	}
	if cfg.HTTP.Port < 1 || cfg.HTTP.Port > 65535 {
		return fmt.Errorf("HTTP_PORT must be between 1 and 65535: %d", cfg.HTTP.Port)
	}
	if cfg.HTTP.ReadTimeout <= 0 ||
		cfg.HTTP.ReadHeaderTimeout <= 0 ||
		cfg.HTTP.WriteTimeout <= 0 ||
		cfg.HTTP.IdleTimeout <= 0 ||
		cfg.HTTP.ShutdownTimeout <= 0 {
		return errors.New("HTTP timeouts must be positive")
	}
	if cfg.HTTP.TLS.Enabled &&
		(cfg.HTTP.TLS.CertFile == "" || cfg.HTTP.TLS.KeyFile == "") {
		return errors.New("HTTP_TLS_CERT_FILE and HTTP_TLS_KEY_FILE are required when TLS is enabled")
	}
	if cfg.HTTP.RateLimit.Enabled &&
		(!isPositiveFinite(cfg.HTTP.RateLimit.Rate) || cfg.HTTP.RateLimit.Burst <= 0) {
		return errors.New("RATE_LIMIT_RATE and RATE_LIMIT_BURST must be positive")
	}
	if cfg.HTTP.ContactViewRateLimit.Enabled &&
		(!isPositiveFinite(cfg.HTTP.ContactViewRateLimit.Rate) || cfg.HTTP.ContactViewRateLimit.Burst <= 0) {
		return errors.New("CONTACT_VIEW_RATE_LIMIT_PER_MINUTE and CONTACT_VIEW_RATE_LIMIT_BURST must be positive")
	}
	if cfg.HTTP.AnalyticsEventsRateLimit.Enabled &&
		(!isPositiveFinite(cfg.HTTP.AnalyticsEventsRateLimit.Rate) || cfg.HTTP.AnalyticsEventsRateLimit.Burst <= 0) {
		return errors.New("ANALYTICS_EVENTS_RATE_LIMIT_PER_MINUTE and ANALYTICS_EVENTS_RATE_LIMIT_BURST must be positive")
	}
	if cfg.HTTP.MaxBodyBytes <= 0 {
		return fmt.Errorf("HTTP_MAX_BODY_BYTES must be positive: %d", cfg.HTTP.MaxBodyBytes)
	}
	if cfg.Turnstile.Timeout <= 0 {
		return errors.New("TURNSTILE_TIMEOUT must be positive")
	}
	if cfg.CORS.AllowCredentials && len(cfg.CORS.AllowedOrigins) == 0 {
		return errors.New("CORS_ALLOWED_ORIGINS must be configured when credentials are enabled")
	}
	for _, origin := range cfg.CORS.AllowedOrigins {
		if origin == "*" {
			return errors.New("CORS_ALLOWED_ORIGINS must not contain * when using an explicit allowlist")
		}
		parsed, err := url.ParseRequestURI(origin)
		if err != nil ||
			parsed.Scheme != "http" && parsed.Scheme != "https" ||
			parsed.Host == "" ||
			parsed.User != nil ||
			parsed.Path != "" ||
			parsed.RawQuery != "" ||
			parsed.Fragment != "" {
			return fmt.Errorf(
				"invalid CORS origin %q: expected an http(s) origin without path, query or fragment",
				origin,
			)
		}
	}
	for _, proxy := range cfg.HTTP.TrustedProxies {
		if net.ParseIP(proxy) != nil {
			continue
		}
		if _, _, err := net.ParseCIDR(proxy); err != nil {
			return fmt.Errorf("invalid HTTP_TRUSTED_PROXIES entry %q", proxy)
		}
	}
	if cfg.Log.Format != "json" && cfg.Log.Format != "console" {
		return fmt.Errorf("LOG_FORMAT must be json or console: %q", cfg.Log.Format)
	}
	if cfg.Telemetry.SampleRatio < 0 || cfg.Telemetry.SampleRatio > 1 || math.IsNaN(cfg.Telemetry.SampleRatio) {
		return errors.New("OTEL_TRACES_SAMPLE_RATIO must be between 0 and 1")
	}
	if cfg.Postgres.StatementTimeout < 0 || cfg.Postgres.LockTimeout < 0 ||
		cfg.Postgres.IdleInTransactionTimeout < 0 || cfg.Postgres.SlowQueryThreshold < 0 {
		return errors.New("PostgreSQL timeouts must not be negative")
	}
	if cfg.Telemetry.Enabled {
		if strings.TrimSpace(cfg.Telemetry.ServiceName) == "" {
			return errors.New("OTEL_SERVICE_NAME must not be empty when tracing is enabled")
		}
		if strings.TrimSpace(cfg.Telemetry.OTLPEndpoint) == "" {
			return errors.New("OTEL_EXPORTER_OTLP_ENDPOINT is required when tracing is enabled")
		}
		if _, _, err := net.SplitHostPort(cfg.Telemetry.OTLPEndpoint); err != nil {
			return fmt.Errorf(
				"invalid OTEL_EXPORTER_OTLP_ENDPOINT %q: expected host:port",
				cfg.Telemetry.OTLPEndpoint,
			)
		}
	}
	if err := validateProductionSecurity(cfg); err != nil {
		return err
	}
	return nil
}

func isPositiveFinite(value float64) bool {
	return value > 0 && !math.IsNaN(value) && !math.IsInf(value, 0)
}

func isPlaceholderSecret(secret string) bool {
	normalized := strings.ToLower(strings.TrimSpace(secret))
	if normalized == "" {
		return true
	}
	for _, marker := range []string{
		"change-me", "changeme", "replace-me", "replace_this", "development-only",
		"placeholder", "your-secret", "example-secret", "random-32-byte-secret",
	} {
		if strings.Contains(normalized, marker) {
			return true
		}
	}
	return false
}

func isValidTCPPort(value string) bool {
	if value == "" {
		return false
	}
	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}
	port, err := strconv.Atoi(value)
	return err == nil && port >= 1 && port <= 65535
}

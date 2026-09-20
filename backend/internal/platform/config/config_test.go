package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoadRequiresOriginsWhenCredentialsEnabled(t *testing.T) {
	setValidConfigEnvironment(t)
	t.Setenv("CORS_ALLOWED_ORIGINS", "")
	t.Setenv("CORS_ALLOW_CREDENTIALS", "true")

	_, err := Load()

	require.Error(t, err)
	require.Contains(t, err.Error(), "CORS_ALLOWED_ORIGINS")
}

func TestCORSMockHeaderFollowsAuthenticationMode(t *testing.T) {
	setValidConfigEnvironment(t)
	cfg, err := parse()
	require.NoError(t, err)
	require.True(t, cfg.CORS.AllowMockAuthHeader)

	t.Setenv("AUTH_MODE", "telegram")
	t.Setenv("AUTH_TELEGRAM_CLIENT_ID", "test-client-id")
	t.Setenv("AUTH_TELEGRAM_CLIENT_SECRET", "test-client-secret")
	cfg, err = parse()
	require.NoError(t, err)
	require.False(t, cfg.CORS.AllowMockAuthHeader)
}

func TestLoadParsesConfiguration(t *testing.T) {
	setValidConfigEnvironment(t)
	t.Setenv("HTTP_PORT", "9090")
	t.Setenv("HTTP_MAX_BODY_BYTES", "2048")
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:4173")
	t.Setenv("LOG_FORMAT", "console")

	cfg, err := Load()

	require.NoError(t, err)
	require.Equal(t, 9090, cfg.HTTP.Port)
	require.Equal(t, int64(2048), cfg.HTTP.MaxBodyBytes)
	require.Equal(
		t,
		[]string{"http://localhost:5173", "http://localhost:4173"},
		cfg.CORS.AllowedOrigins,
	)
	require.Equal(t, "console", cfg.Log.Format)
}

func TestLoadParsesTurnstileConfiguration(t *testing.T) {
	setValidConfigEnvironment(t)
	t.Setenv("TURNSTILE_ENABLED", "true")
	t.Setenv("TURNSTILE_SECRET_KEY", "secret-key")
	t.Setenv("TURNSTILE_EXPECTED_HOSTNAME", "cv.example.com")
	t.Setenv("TURNSTILE_EXPECTED_ACTION", "contact_view")
	t.Setenv("TURNSTILE_TIMEOUT", "7s")

	cfg, err := parse()

	require.NoError(t, err)
	require.True(t, cfg.Turnstile.Enabled)
	require.Equal(t, "secret-key", cfg.Turnstile.SecretKey)
	require.Equal(t, "cv.example.com", cfg.Turnstile.ExpectedHost)
	require.Equal(t, "contact_view", cfg.Turnstile.ExpectedAction)
	require.Equal(t, 7*time.Second, cfg.Turnstile.Timeout)
}

func TestLoadParsesContactViewRateLimit(t *testing.T) {
	setValidConfigEnvironment(t)
	t.Setenv("CONTACT_VIEW_RATE_LIMIT_ENABLED", "true")
	t.Setenv("CONTACT_VIEW_RATE_LIMIT_PER_MINUTE", "7")
	t.Setenv("CONTACT_VIEW_RATE_LIMIT_BURST", "4")

	cfg, err := parse()

	require.NoError(t, err)
	require.True(t, cfg.HTTP.ContactViewRateLimit.Enabled)
	require.Equal(t, 7.0, cfg.HTTP.ContactViewRateLimit.Rate)
	require.Equal(t, 4, cfg.HTTP.ContactViewRateLimit.Burst)
}

func TestDevelopmentCookieDefaultsAreCompatibleWithLocalHTTP(t *testing.T) {
	setValidConfigEnvironment(t)
	t.Setenv("AUTH_COOKIE_NAME", "")
	t.Setenv("AUTH_COOKIE_SECURE", "")

	cfg, err := parse()

	require.NoError(t, err)
	require.Equal(t, "session", cfg.Auth.Cookie.Name)
	require.False(t, cfg.Auth.Cookie.Secure)
}

func TestLoadParsesAuthMode(t *testing.T) {
	setValidConfigEnvironment(t)
	t.Setenv("AUTH_MODE", "mock")
	t.Setenv("AUTH_COOKIE_NAME", "__Host-session")
	t.Setenv("AUTH_COOKIE_SECURE", "true")
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:5173")

	cfg, err := parse()

	require.NoError(t, err)
	require.Equal(t, "mock", cfg.Auth.Mode)
}

func TestLoadRejectsUnsupportedAuthMode(t *testing.T) {
	setValidConfigEnvironment(t)
	t.Setenv("AUTH_MODE", "jwt")

	_, err := parse()

	require.Error(t, err)
	require.Contains(t, err.Error(), "AUTH_MODE")
}

func TestLoadParsesTrustedProxyConfiguration(t *testing.T) {
	setValidConfigEnvironment(t)
	t.Setenv("HTTP_TRUSTED_PROXIES", "192.0.2.10, 2001:db8::/32")

	cfg, err := parse()
	require.NoError(t, err)
	require.Equal(t, []string{"192.0.2.10", "2001:db8::/32"}, cfg.HTTP.TrustedProxies)
}

func TestLoadParsesDedicatedAuthRateLimits(t *testing.T) {
	setValidConfigEnvironment(t)
	t.Setenv("AUTH_RATE_LIMIT_ENABLED", "true")
	t.Setenv("AUTH_RATE_LIMIT_START_PER_MINUTE", "7")
	t.Setenv("AUTH_RATE_LIMIT_START_BURST", "7")
	t.Setenv("AUTH_RATE_LIMIT_CALLBACK_PER_MINUTE", "25")
	t.Setenv("AUTH_RATE_LIMIT_CALLBACK_BURST", "25")

	cfg, err := parse()

	require.NoError(t, err)
	require.True(t, cfg.Auth.RateLimit.Enabled)
	require.Equal(t, 7.0, cfg.Auth.RateLimit.StartRatePerMinute)
	require.Equal(t, 7, cfg.Auth.RateLimit.StartBurst)
	require.Equal(t, 25.0, cfg.Auth.RateLimit.CallbackRatePerMinute)
	require.Equal(t, 25, cfg.Auth.RateLimit.CallbackBurst)
}

func TestLoadParsesAnalyticsRetentionDays(t *testing.T) {
	setValidConfigEnvironment(t)
	t.Setenv("ANALYTICS_CLICK_RETENTION_DAYS", "45")
	t.Setenv("ANALYTICS_VIEW_RETENTION_DAYS", "180")

	cfg, err := parse()
	require.NoError(t, err)
	require.Equal(t, 45, cfg.Analytics.ClickRetentionDays)
	require.Equal(t, 180, cfg.Analytics.ViewRetentionDays)
}

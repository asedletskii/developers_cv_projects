package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func setValidConfigEnvironment(t *testing.T) {
	t.Helper()

	t.Setenv("AUTH_MODE", "mock")
	t.Setenv("AUTH_COOKIE_NAME", "__Host-session")
	t.Setenv("AUTH_COOKIE_SECURE", "true")
	t.Setenv("AUTH_COOKIE_HTTP_ONLY", "true")
	t.Setenv("AUTH_COOKIE_SAME_SITE", "lax")
	t.Setenv("AUTH_COOKIE_PATH", "/")
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:5173")
	t.Setenv("CORS_ALLOW_CREDENTIALS", "true")
	t.Setenv("POSTGRES_PASSWORD", "test-password")
	t.Setenv("POSTGRES_SSLMODE", "disable")
	t.Setenv("ANALYTICS_VISITOR_SECRET", "test-analytics-visitor-secret-32-bytes-minimum")
}

func TestLoadRequiresOriginsWhenCredentialsEnabled(t *testing.T) {
	setValidConfigEnvironment(t)
	t.Setenv("CORS_ALLOWED_ORIGINS", "")
	t.Setenv("CORS_ALLOW_CREDENTIALS", "true")

	_, err := Load()

	require.Error(t, err)
	require.Contains(t, err.Error(), "CORS_ALLOWED_ORIGINS")
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
	require.True(t, cfg.Analytics.RespectDoNotTrack)
}

func TestLoadParsesTurnstileConfiguration(t *testing.T) {
	setValidConfigEnvironment(t)
	t.Setenv("TURNSTILE_ENABLED", "true")
	t.Setenv("TURNSTILE_SITE_KEY", "site-key")
	t.Setenv("TURNSTILE_SECRET_KEY", "secret-key")
	t.Setenv("TURNSTILE_EXPECTED_HOSTNAME", "cv.example.com")
	t.Setenv("TURNSTILE_EXPECTED_ACTION", "contact_view")
	t.Setenv("TURNSTILE_TIMEOUT", "7s")

	cfg, err := parse()

	require.NoError(t, err)
	require.True(t, cfg.Turnstile.Enabled)
	require.Equal(t, "site-key", cfg.Turnstile.SiteKey)
	require.Equal(t, "secret-key", cfg.Turnstile.SecretKey)
	require.Equal(t, "cv.example.com", cfg.Turnstile.ExpectedHost)
	require.Equal(t, "contact_view", cfg.Turnstile.ExpectedAction)
	require.Equal(t, 7*time.Second, cfg.Turnstile.Timeout)
}

func TestValidateRejectsEnabledTurnstileWithoutSiteKey(t *testing.T) {
	setValidConfigEnvironment(t)
	t.Setenv("TURNSTILE_ENABLED", "true")
	t.Setenv("TURNSTILE_SITE_KEY", "")

	_, err := parse()

	require.Error(t, err)
	require.Contains(t, err.Error(), "TURNSTILE_SITE_KEY")
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

func TestValidateRejectsInvalidContactViewRateLimit(t *testing.T) {
	cfg := productionSecurityTestConfig()
	cfg.HTTP.ContactViewRateLimit.Burst = 0

	err := validate(cfg)

	require.Error(t, err)
	require.Contains(t, err.Error(), "CONTACT_VIEW_RATE_LIMIT")
}

func TestValidateProductionRequiresContactViewRateLimit(t *testing.T) {
	cfg := productionSecurityTestConfig()
	cfg.HTTP.ContactViewRateLimit.Enabled = false

	err := validateProductionSecurity(cfg)

	require.Error(t, err)
	require.Contains(t, err.Error(), "CONTACT_VIEW_RATE_LIMIT_ENABLED")
}

func TestValidateProductionRequiresTurnstile(t *testing.T) {
	cfg := productionSecurityTestConfig()
	cfg.Turnstile.Enabled = false

	err := validateProductionSecurity(cfg)

	require.Error(t, err)
	require.Contains(t, err.Error(), "TURNSTILE_ENABLED")
}

func TestValidateProductionRequiresTurnstileCredentialsAndBindings(t *testing.T) {
	cfg := productionSecurityTestConfig()

	cfg.Turnstile.SecretKey = ""
	err := validateProductionSecurity(cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "TURNSTILE_SECRET_KEY")

	cfg = productionSecurityTestConfig()
	cfg.Turnstile.SiteKey = ""
	err = validateProductionSecurity(cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "TURNSTILE_SITE_KEY")

	cfg = productionSecurityTestConfig()
	cfg.Turnstile.ExpectedHost = ""
	err = validateProductionSecurity(cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "TURNSTILE_EXPECTED_HOSTNAME")

	cfg = productionSecurityTestConfig()
	cfg.Turnstile.ExpectedAction = ""
	err = validateProductionSecurity(cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "TURNSTILE_EXPECTED_ACTION")
}

func TestLoadDotEnvDoesNotOverrideEnvironment(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	err := os.WriteFile(
		path,
		[]byte(
			"TEST_DOTENV_VALUE=from-file\n"+
				"TEST_DOTENV_QUOTED=\"hello world\"\n",
		),
		0o600,
	)
	require.NoError(t, err)

	t.Setenv("TEST_DOTENV_VALUE", "from-environment")

	err = loadDotEnv(path)
	require.NoError(t, err)

	require.Equal(t, "from-environment", os.Getenv("TEST_DOTENV_VALUE"))
	require.Equal(t, "hello world", os.Getenv("TEST_DOTENV_QUOTED"))
}

func TestLoadRejectsInvalidEnvironmentValue(t *testing.T) {
	setValidConfigEnvironment(t)
	t.Setenv("HTTP_PORT", "not-a-port")

	_, err := parse()

	require.Error(t, err)
	require.Contains(t, err.Error(), "HTTP_PORT")
}

func TestLoadRejectsWildcardCORSOrigin(t *testing.T) {
	setValidConfigEnvironment(t)
	t.Setenv("CORS_ALLOWED_ORIGINS", "*")
	t.Setenv("CORS_ALLOW_CREDENTIALS", "true")

	_, err := parse()

	require.Error(t, err)
	require.Contains(t, err.Error(), "must not contain *")
}

func TestPostgresConnectionStringEscapesCredentials(t *testing.T) {
	t.Setenv("TEST_POSTGRES_PASSWORD", "p@ss:w?rd#1")

	cfg := PostgresConfig{
		Host:            "localhost",
		Port:            5432,
		Database:        "developers_cv",
		User:            "developers",
		Password:        os.Getenv("TEST_POSTGRES_PASSWORD"),
		SSLMode:         "require",
		MaxConns:        10,
		MaxIdleConns:    2,
		MaxConnLifetime: time.Hour,
		MaxConnIdleTime: 30 * time.Minute,
	}

	require.Equal(
		t,
		"postgres://developers:p%40ss%3Aw%3Frd%231@localhost:5432/developers_cv?sslmode=require",
		cfg.ConnectionString(),
	)
}

func TestLoadRejectsInvalidPostgresSSLMode(t *testing.T) {
	setValidConfigEnvironment(t)
	t.Setenv("POSTGRES_SSLMODE", "invalid")

	_, err := parse()

	require.Error(t, err)
	require.Contains(t, err.Error(), "POSTGRES_SSLMODE")
}

func TestLoadParsesAuthMode(t *testing.T) {
	setValidConfigEnvironment(t)
	t.Setenv("AUTH_MODE", "mock")
	t.Setenv("AUTH_COOKIE_NAME", "__Host-session")
	t.Setenv("AUTH_COOKIE_SECURE", "true")
	t.Setenv("AUTH_COOKIE_HTTP_ONLY", "true")
	t.Setenv("AUTH_COOKIE_SAME_SITE", "lax")
	t.Setenv("AUTH_COOKIE_PATH", "/")
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

func TestValidateRejectsDisabledPrometheusInProduction(t *testing.T) {
	cfg := productionSecurityTestConfig()
	cfg.Auth.Mode = "telegram"
	cfg.Postgres.SSLMode = "require"
	cfg.Prometheus.Enabled = false

	err := validateProductionSecurity(cfg)

	require.Error(t, err)
	require.Contains(t, err.Error(), "PROMETHEUS_ENABLED")
}

func TestValidateRejectsInsecureTelemetryInProduction(t *testing.T) {
	cfg := productionSecurityTestConfig()
	cfg.Auth.Mode = "telegram"
	cfg.Postgres.SSLMode = "require"
	cfg.Telemetry.Enabled = true
	cfg.Telemetry.OTLPEndpoint = "otel-collector:4317"
	cfg.Telemetry.Insecure = true

	err := validateProductionSecurity(cfg)

	require.Error(t, err)
	require.Contains(t, err.Error(), "OTEL_EXPORTER_OTLP_INSECURE")
}

func TestValidateRejectsInvalidTelemetryEndpoint(t *testing.T) {
	cfg := productionSecurityTestConfig()
	cfg.App.Env = "test"
	cfg.Auth.Mode = "mock"
	cfg.Telemetry.Enabled = true
	cfg.Telemetry.ServiceName = "developers-cv-api"
	cfg.Telemetry.OTLPEndpoint = "http://localhost:4317"

	err := validate(cfg)

	require.Error(t, err)
	require.Contains(t, err.Error(), "OTEL_EXPORTER_OTLP_ENDPOINT")
}

func TestValidateAcceptsTelemetryEndpoint(t *testing.T) {
	cfg := productionSecurityTestConfig()
	cfg.App.Env = "test"
	cfg.Auth.Mode = "mock"
	cfg.Telemetry.Enabled = true
	cfg.Telemetry.ServiceName = "developers-cv-api"
	cfg.Telemetry.OTLPEndpoint = "localhost:4317"

	err := validate(cfg)

	require.NoError(t, err)
}

func TestValidateRejectsDisabledRateLimitInProduction(t *testing.T) {
	cfg := productionSecurityTestConfig()
	cfg.App.Env = "production"
	cfg.Auth.Mode = "telegram"
	cfg.Postgres.SSLMode = "require"
	cfg.HTTP.TLS.Enabled = true
	cfg.HTTP.TLS.CertFile = "/run/secrets/tls.crt"
	cfg.HTTP.TLS.KeyFile = "/run/secrets/tls.key"
	cfg.HTTP.RateLimit.Enabled = false

	err := validateProductionSecurity(cfg)

	require.Error(t, err)
	require.Contains(t, err.Error(), "RATE_LIMIT_ENABLED")
}

func TestValidateRejectsInvalidTrustedProxy(t *testing.T) {
	cfg := productionSecurityTestConfig()
	cfg.HTTP.TrustedProxies = []string{"not-a-proxy"}

	err := validate(cfg)

	require.Error(t, err)
	require.Contains(t, err.Error(), "HTTP_TRUSTED_PROXIES")
}

func productionSecurityTestConfig() Config {
	return Config{
		App:       AppConfig{Env: "production"},
		Analytics: AnalyticsConfig{VisitorSecret: "production-test-analytics-secret-32-bytes-minimum", SuspectBurst: 5, SuspectWindow: time.Minute, RespectDoNotTrack: true},
		Auth: AuthConfig{
			Mode:        "telegram",
			Telegram:    TelegramAuthConfig{ClientID: "client-id", ClientSecret: "client-secret", RedirectURL: "https://example.com/api/v1/auth/telegram/callback", IssuerURL: "https://oauth.telegram.org"},
			Cookie:      CookieAuthConfig{Name: "__Host-session", Secure: true, HTTPOnly: true, SameSite: "lax", Path: "/"},
			Session:     SessionAuthConfig{AbsoluteTTL: 720 * time.Hour, IdleTTL: 168 * time.Hour, TouchInterval: 5 * time.Minute},
			Transaction: TransactionAuthConfig{TTL: 5 * time.Minute},
			RateLimit:   AuthRateLimitConfig{Enabled: true, StartRatePerMinute: 5, StartBurst: 5, CallbackRatePerMinute: 20, CallbackBurst: 20},
		},
		HTTP: HTTPConfig{
			Port: 8080, ReadTimeout: time.Second, ReadHeaderTimeout: time.Second, WriteTimeout: time.Second, IdleTimeout: time.Second, ShutdownTimeout: time.Second, MaxBodyBytes: 1024,
			TLS:                      TLSConfig{Enabled: true, CertFile: "/run/secrets/tls.crt", KeyFile: "/run/secrets/tls.key"},
			RateLimit:                RateLimitConfig{Enabled: true, Rate: 10, Burst: 20},
			ContactViewRateLimit:     RateLimitConfig{Enabled: true, Rate: 3, Burst: 2},
			AnalyticsEventsRateLimit: RateLimitConfig{Enabled: true, Rate: 30, Burst: 10},
		},
		CORS:       CORSConfig{AllowedOrigins: []string{"https://example.com"}, AllowCredentials: true},
		Turnstile:  TurnstileConfig{Enabled: true, SiteKey: "site-key", SecretKey: "secret-key", ExpectedHost: "example.com", ExpectedAction: "contact_view", Timeout: 5 * time.Second},
		Log:        LogConfig{Format: "json"},
		Telemetry:  TelemetryConfig{ServiceName: "developers-cv-api", ServiceVersion: "test", Environment: "production"},
		Prometheus: PrometheusConfig{Enabled: true},
		Postgres:   PostgresConfig{Host: "postgres", Port: 5432, Database: "developers_cv", User: "developers_cv", Password: "secret", SSLMode: "require", MaxConns: 10, MaxIdleConns: 2, MaxConnLifetime: time.Hour, MaxConnIdleTime: time.Minute},
	}
}

func TestValidateProductionRejectsMockAuthentication(t *testing.T) {
	cfg := productionSecurityTestConfig()
	cfg.Auth.Mode = "mock"

	err := validateProductionSecurity(cfg)

	require.Error(t, err)
	require.Contains(t, err.Error(), "AUTH_MODE")
}

func TestValidateProductionAcceptsTelegramAuthentication(t *testing.T) {
	cfg := productionSecurityTestConfig()
	cfg.Auth.Mode = "telegram"
	cfg.Auth.Telegram.ClientID = "client-id"
	cfg.Auth.Telegram.ClientSecret = "client-secret"
	cfg.Auth.Telegram.RedirectURL = "https://example.com/api/v1/auth/telegram/callback"
	cfg.Auth.Telegram.IssuerURL = "https://oauth.telegram.org"
	cfg.Auth.Cookie.Name = "__Host-session"
	cfg.Auth.Cookie.Secure = true
	cfg.Auth.Cookie.HTTPOnly = true
	cfg.Auth.Cookie.Path = "/"
	cfg.Auth.Cookie.SameSite = "lax"

	err := validate(cfg)

	require.NoError(t, err)
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

func TestValidateRejectsDisabledDedicatedAuthRateLimitInProduction(t *testing.T) {
	cfg := productionSecurityTestConfig()
	cfg.Auth.Mode = "telegram"
	cfg.Postgres.SSLMode = "require"
	cfg.HTTP.TLS.Enabled = true
	cfg.HTTP.TLS.CertFile = "/run/secrets/tls.crt"
	cfg.HTTP.TLS.KeyFile = "/run/secrets/tls.key"
	cfg.Auth.RateLimit.Enabled = false

	err := validateProductionSecurity(cfg)

	require.Error(t, err)
	require.Contains(t, err.Error(), "AUTH_RATE_LIMIT_ENABLED")
}

func TestValidateRejectsInvalidDedicatedAuthRateLimit(t *testing.T) {
	cfg := productionSecurityTestConfig()
	cfg.Auth.RateLimit.StartBurst = 0

	err := validate(cfg)

	require.Error(t, err)
	require.Contains(t, err.Error(), "auth rate limits")
}

func TestLoadParsesRespectDoNotTrack(t *testing.T) {
	setValidConfigEnvironment(t)
	t.Setenv("RESPECT_DO_NOT_TRACK", "false")
	cfg, err := parse()
	require.NoError(t, err)
	require.False(t, cfg.Analytics.RespectDoNotTrack)
}

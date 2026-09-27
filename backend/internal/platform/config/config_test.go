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
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:5173")
	t.Setenv("CORS_ALLOW_CREDENTIALS", "true")
	t.Setenv("POSTGRES_PASSWORD", "test-password")
	t.Setenv("POSTGRES_SSLMODE", "disable")
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
	cfg.Auth.Mode = "jwt"
	cfg.Postgres.SSLMode = "require"
	cfg.Prometheus.Enabled = false

	err := validateProductionSecurity(cfg)

	require.Error(t, err)
	require.Contains(t, err.Error(), "PROMETHEUS_ENABLED")
}

func TestValidateRejectsInsecureTelemetryInProduction(t *testing.T) {
	cfg := productionSecurityTestConfig()
	cfg.Auth.Mode = "jwt"
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
	cfg.Auth.Mode = "jwt"
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
		App:  AppConfig{Env: "production"},
		Auth: AuthConfig{Mode: "mock"},
		HTTP: HTTPConfig{
			Port:              8080,
			ReadTimeout:       time.Second,
			ReadHeaderTimeout: time.Second,
			WriteTimeout:      time.Second,
			IdleTimeout:       time.Second,
			ShutdownTimeout:   time.Second,
			MaxBodyBytes:      1024,
			TLS: TLSConfig{
				Enabled:  true,
				CertFile: "/run/secrets/tls.crt",
				KeyFile:  "/run/secrets/tls.key",
			},
			RateLimit: RateLimitConfig{
				Enabled: true,
				Rate:    10,
				Burst:   20,
			},
		},
		CORS: CORSConfig{
			AllowedOrigins:   []string{"https://example.com"},
			AllowCredentials: true,
		},
		Log: LogConfig{Format: "json"},
		Telemetry: TelemetryConfig{
			ServiceName:    "developers-cv-api",
			ServiceVersion: "test",
			Environment:    "production",
		},
		Prometheus: PrometheusConfig{Enabled: true},
		Postgres: PostgresConfig{
			Host:            "postgres",
			Port:            5432,
			Database:        "developers_cv",
			User:            "developers_cv",
			Password:        "secret",
			SSLMode:         "require",
			MaxConns:        10,
			MaxIdleConns:    2,
			MaxConnLifetime: time.Hour,
			MaxConnIdleTime: time.Minute,
		},
	}
}

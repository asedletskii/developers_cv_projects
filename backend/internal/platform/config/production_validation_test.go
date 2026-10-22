package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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

func TestValidateAllowsDisabledPrometheusInProduction(t *testing.T) {
	cfg := productionSecurityTestConfig()
	cfg.Auth.Mode = "telegram"
	cfg.Postgres.SSLMode = "require"
	cfg.Prometheus.Enabled = false

	err := validateProductionSecurity(cfg)

	require.NoError(t, err)
}

func TestValidateProductionRequiresLoopbackMetricsAddress(t *testing.T) {
	tests := []struct {
		name string
		addr string
		want string
	}{
		{name: "empty address", addr: "", want: "METRICS_ADDR"},
		{name: "all interfaces", addr: "0.0.0.0:9090", want: "loopback IP"},
		{name: "public interface", addr: "192.0.2.10:9090", want: "loopback IP"},
		{name: "hostname is not explicit loopback", addr: "localhost:9090", want: "loopback IP"},
		{name: "invalid address", addr: "not-an-address", want: "METRICS_ADDR"},
		{name: "non-numeric port", addr: "127.0.0.1:abc", want: "valid TCP port"},
		{name: "empty port", addr: "127.0.0.1:", want: "valid TCP port"},
		{name: "zero port", addr: "127.0.0.1:0", want: "valid TCP port"},
		{name: "port above range", addr: "127.0.0.1:65536", want: "valid TCP port"},
		{name: "negative port", addr: "127.0.0.1:-1", want: "valid TCP port"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := productionSecurityTestConfig()
			cfg.HTTP.MetricsAddr = tt.addr

			err := validateProductionSecurity(cfg)

			require.Error(t, err)
			require.Contains(t, err.Error(), tt.want)
		})
	}
}

func TestValidateProductionAcceptsIPv4AndIPv6LoopbackMetricsAddresses(t *testing.T) {
	for _, addr := range []string{"127.0.0.1:9090", "[::1]:9090"} {
		t.Run(addr, func(t *testing.T) {
			cfg := productionSecurityTestConfig()
			cfg.HTTP.MetricsAddr = addr

			require.NoError(t, validateProductionSecurity(cfg))
		})
	}
}

func TestValidateProductionAcceptsValidMetricsPortBoundaries(t *testing.T) {
	for _, addr := range []string{"127.0.0.1:1", "127.0.0.1:65535", "127.0.0.1:09090"} {
		t.Run(addr, func(t *testing.T) {
			cfg := productionSecurityTestConfig()
			cfg.HTTP.MetricsAddr = addr

			require.NoError(t, validateProductionSecurity(cfg))
		})
	}
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

func TestValidateProductionRejectsInsecureDeploymentConfiguration(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*Config)
		want   string
	}{
		{
			name:   "development version",
			mutate: func(cfg *Config) { cfg.App.Version = "dev" },
			want:   "APP_VERSION",
		},
		{
			name:   "http telegram redirect",
			mutate: func(cfg *Config) { cfg.Auth.Telegram.RedirectURL = "http://example.com/callback" },
			want:   "AUTH_TELEGRAM_REDIRECT_URL",
		},
		{
			name:   "http cors origin",
			mutate: func(cfg *Config) { cfg.CORS.AllowedOrigins = []string{"http://example.com"} },
			want:   "CORS_ALLOWED_ORIGINS",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := productionSecurityTestConfig()
			tt.mutate(&cfg)
			err := validateProductionSecurity(cfg)
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.want)
		})
	}
}

func TestValidateProductionAllowsPostgresWithoutTLSOnPrivateComposeNetwork(t *testing.T) {
	cfg := productionSecurityTestConfig()
	cfg.Postgres.Host = "postgres"
	cfg.Postgres.SSLMode = "disable"

	require.NoError(t, validateProductionSecurity(cfg))
}

func TestValidateProductionRejectsPlaceholderSecrets(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*Config)
		want   string
	}{
		{
			name:   "database password",
			mutate: func(cfg *Config) { cfg.Postgres.Password = "replace-me-with-a-long-secret" },
			want:   "POSTGRES_PASSWORD",
		},
		{
			name:   "telegram client secret",
			mutate: func(cfg *Config) { cfg.Auth.Telegram.ClientSecret = "replace-me-telegram-client-secret" },
			want:   "AUTH_TELEGRAM_CLIENT_SECRET",
		},
		{
			name:   "turnstile secret",
			mutate: func(cfg *Config) { cfg.Turnstile.SecretKey = "replace-me-turnstile-secret-key" },
			want:   "TURNSTILE_SECRET_KEY",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := productionSecurityTestConfig()
			tt.mutate(&cfg)
			err := validateProductionSecurity(cfg)
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.want)
		})
	}
}

func TestValidateProductionRequiresTrustedProxy(t *testing.T) {
	cfg := productionSecurityTestConfig()
	cfg.HTTP.TrustedProxies = nil

	err := validate(cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "HTTP_TRUSTED_PROXIES")
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

	err := validate(cfg)

	require.NoError(t, err)
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

func TestValidateProductionRejectsPlaceholderVisitorSecret(t *testing.T) {
	cfg := productionSecurityTestConfig()
	cfg.Analytics.VisitorSecret = "change-me-to-a-random-32-byte-secret"

	err := validate(cfg)

	require.Error(t, err)
	require.Contains(t, err.Error(), "ANALYTICS_VISITOR_SECRET")
}

func TestValidateProductionRejectsDevelopmentOnlyVisitorSecret(t *testing.T) {
	cfg := productionSecurityTestConfig()
	cfg.Analytics.VisitorSecret = "development-only-analytics-secret-change-me-32"

	err := validate(cfg)

	require.Error(t, err)
	require.Contains(t, err.Error(), "ANALYTICS_VISITOR_SECRET")
}

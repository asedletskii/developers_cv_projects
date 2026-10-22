package config

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateRejectsInvalidContactViewRateLimit(t *testing.T) {
	cfg := productionSecurityTestConfig()
	cfg.HTTP.ContactViewRateLimit.Burst = 0

	err := validate(cfg)

	require.Error(t, err)
	require.Contains(t, err.Error(), "CONTACT_VIEW_RATE_LIMIT")
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

func TestValidateRejectsNonFiniteRateLimit(t *testing.T) {
	cfg := productionSecurityTestConfig()
	cfg.HTTP.RateLimit.Rate = math.NaN()

	err := validate(cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "RATE_LIMIT_RATE")
}

func TestValidateRejectsInvalidTrustedProxy(t *testing.T) {
	cfg := productionSecurityTestConfig()
	cfg.HTTP.TrustedProxies = []string{"not-a-proxy"}

	err := validate(cfg)

	require.Error(t, err)
	require.Contains(t, err.Error(), "HTTP_TRUSTED_PROXIES")
}

func TestValidateAllowsMockAuthenticationOnlyInDevelopmentAndTest(t *testing.T) {
	for _, env := range []string{"development", "test"} {
		t.Run(env, func(t *testing.T) {
			cfg := productionSecurityTestConfig()
			cfg.App.Env = env
			cfg.Auth.Mode = "mock"
			require.NoError(t, validate(cfg))
		})
	}
	for _, env := range []string{"staging", "production"} {
		t.Run(env, func(t *testing.T) {
			cfg := productionSecurityTestConfig()
			cfg.App.Env = env
			cfg.Auth.Mode = "mock"
			err := validate(cfg)
			require.Error(t, err)
			require.Contains(t, err.Error(), "AUTH_MODE=mock")
		})
	}
}

func TestValidateRejectsInvalidDedicatedAuthRateLimit(t *testing.T) {
	cfg := productionSecurityTestConfig()
	cfg.Auth.RateLimit.StartBurst = 0

	err := validate(cfg)

	require.Error(t, err)
	require.Contains(t, err.Error(), "auth rate limits")
}

func TestLoadRejectsUnknownAppEnvironment(t *testing.T) {
	setValidConfigEnvironment(t)
	t.Setenv("APP_ENV", "prod")

	_, err := parse()

	require.Error(t, err)
	require.Contains(t, err.Error(), "APP_ENV")
}

func TestValidateRejectsInvalidAnalyticsRetentionDays(t *testing.T) {
	cfg := productionSecurityTestConfig()
	cfg.Analytics.ClickRetentionDays = 0
	err := validate(cfg)
	require.ErrorContains(t, err, "ANALYTICS_CLICK_RETENTION_DAYS")

	cfg = productionSecurityTestConfig()
	cfg.Analytics.ViewRetentionDays = 3651
	err = validate(cfg)
	require.ErrorContains(t, err, "ANALYTICS_VIEW_RETENTION_DAYS")
}

func TestValidateRejectsOutOfRangeTraceSampleRatio(t *testing.T) {
	cfg := productionSecurityTestConfig()
	cfg.App.Env = "test"
	cfg.Auth.Mode = "mock"
	cfg.Telemetry.SampleRatio = 1.1

	err := validate(cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "OTEL_TRACES_SAMPLE_RATIO")
}

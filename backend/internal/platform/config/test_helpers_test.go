package config

import (
	"testing"
	"time"
)

func setValidConfigEnvironment(t *testing.T) {
	t.Helper()

	t.Setenv("APP_ENV", "development")
	t.Setenv("AUTH_MODE", "mock")
	t.Setenv("AUTH_COOKIE_NAME", "__Host-session")
	t.Setenv("AUTH_COOKIE_SECURE", "true")
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:5173")
	t.Setenv("CORS_ALLOW_CREDENTIALS", "true")
	t.Setenv("POSTGRES_PASSWORD", "test-password")
	t.Setenv("POSTGRES_SSLMODE", "disable")
	t.Setenv("ANALYTICS_VISITOR_SECRET", "test-analytics-visitor-secret-32-bytes-minimum")
}

func productionSecurityTestConfig() Config {
	return Config{
		App: AppConfig{Env: "production", Version: "0.1.0"},
		Analytics: AnalyticsConfig{
			VisitorSecret:      "production-test-analytics-secret-32-bytes-minimum",
			SuspectBurst:       5,
			SuspectWindow:      time.Minute,
			ClickRetentionDays: 90,
			ViewRetentionDays:  365,
		},
		Auth: AuthConfig{
			Mode: "telegram",
			Telegram: TelegramAuthConfig{
				ClientID:     "client-id",
				ClientSecret: "client-secret",
				RedirectURL:  "https://example.com/api/v1/auth/telegram/callback",
				IssuerURL:    "https://oauth.telegram.org",
			},
			Cookie:      CookieAuthConfig{Name: "__Host-session", Secure: true},
			Session:     SessionAuthConfig{AbsoluteTTL: 720 * time.Hour, IdleTTL: 168 * time.Hour, TouchInterval: 5 * time.Minute},
			Transaction: TransactionAuthConfig{TTL: 5 * time.Minute},
			RateLimit:   AuthRateLimitConfig{Enabled: true, StartRatePerMinute: 5, StartBurst: 5, CallbackRatePerMinute: 20, CallbackBurst: 20},
		},
		HTTP: HTTPConfig{
			TrustedProxies:           []string{"192.0.2.10"},
			MetricsAddr:              "127.0.0.1:9090",
			Port:                     8080,
			ReadTimeout:              time.Second,
			ReadHeaderTimeout:        time.Second,
			WriteTimeout:             time.Second,
			IdleTimeout:              time.Second,
			ShutdownTimeout:          time.Second,
			MaxBodyBytes:             1024,
			TLS:                      TLSConfig{Enabled: true, CertFile: "/run/secrets/tls.crt", KeyFile: "/run/secrets/tls.key"},
			RateLimit:                RateLimitConfig{Enabled: true, Rate: 10, Burst: 20},
			ContactViewRateLimit:     RateLimitConfig{Enabled: true, Rate: 3, Burst: 2},
			AnalyticsEventsRateLimit: RateLimitConfig{Enabled: true, Rate: 30, Burst: 10},
		},
		CORS: CORSConfig{AllowedOrigins: []string{"https://example.com"}, AllowCredentials: true},
		Turnstile: TurnstileConfig{
			Enabled:        true,
			SecretKey:      "secret-key",
			ExpectedHost:   "example.com",
			ExpectedAction: "contact_view",
			Timeout:        5 * time.Second,
		},
		Log:        LogConfig{Format: "json"},
		Telemetry:  TelemetryConfig{ServiceName: "developers-cv-api", ServiceVersion: "test", Environment: "production"},
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

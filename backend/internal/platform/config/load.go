package config

import "fmt"

func Load() (Config, error) {
	if err := loadDotEnv(); err != nil {
		return Config{}, fmt.Errorf("load .env: %w", err)
	}
	return parse()
}

func parse() (Config, error) {
	app := parseAppConfig()
	auth, err := parseAuthConfig(app.Env)
	if err != nil {
		return Config{}, err
	}
	httpConfig, err := parseHTTPConfig()
	if err != nil {
		return Config{}, err
	}
	turnstile, err := parseTurnstileConfig()
	if err != nil {
		return Config{}, err
	}
	logConfig := parseLogConfig()
	cors, err := parseCORSConfig(auth.Mode)
	if err != nil {
		return Config{}, err
	}
	telemetry, err := parseTelemetryConfig(app)
	if err != nil {
		return Config{}, err
	}
	analytics, err := parseAnalyticsConfig()
	if err != nil {
		return Config{}, err
	}
	prometheus, err := parsePrometheusConfig()
	if err != nil {
		return Config{}, err
	}
	postgres, err := parsePostgresConfig()
	if err != nil {
		return Config{}, err
	}

	cfg := Config{App: app, Auth: auth, HTTP: httpConfig, Turnstile: turnstile,
		Log: logConfig, CORS: cors, Telemetry: telemetry, Analytics: analytics,
		Prometheus: prometheus, Postgres: postgres}
	if err := validate(cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

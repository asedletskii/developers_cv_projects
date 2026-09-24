package config

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	App        AppConfig
	Auth       AuthConfig
	HTTP       HTTPConfig
	Log        LogConfig
	CORS       CORSConfig
	Telemetry  TelemetryConfig
	Prometheus PrometheusConfig
	Postgres   PostgresConfig
}

type AuthConfig struct {
	Mode string
}

type AppConfig struct {
	Env     string
	Name    string
	Version string
}

type HTTPConfig struct {
	Host              string
	Port              int
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
	MaxBodyBytes      int64
	TLS               TLSConfig
	TrustedProxies    []string
	RateLimit         RateLimitConfig
}

type TLSConfig struct {
	Enabled  bool
	CertFile string
	KeyFile  string
}

type RateLimitConfig struct {
	Enabled bool
	Rate    float64
	Burst   int
}

type LogConfig struct {
	Level  string
	Format string
}

type CORSConfig struct {
	AllowedOrigins   []string
	AllowCredentials bool
}

type TelemetryConfig struct {
	Enabled      bool
	ServiceName  string
	OTLPEndpoint string
	Insecure     bool
}

type PrometheusConfig struct {
	Enabled bool
}

type PostgresConfig struct {
	Host            string
	Port            int
	Database        string
	User            string
	Password        string
	SSLMode         string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

func Load() (Config, error) {
	if err := loadDotEnv(); err != nil {
		return Config{}, fmt.Errorf("load .env: %w", err)
	}

	cfg, err := parse()
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func parse() (Config, error) {
	appEnv, err := envString("APP_ENV", "development")
	if err != nil {
		return Config{}, err
	}
	appName, err := envString("APP_NAME", "developers-cv-api")
	if err != nil {
		return Config{}, err
	}
	appVersion, err := envString("APP_VERSION", "dev")
	if err != nil {
		return Config{}, err
	}
	authMode, err := envString("AUTH_MODE", "mock")
	if err != nil {
		return Config{}, err
	}

	host, err := envString("HTTP_HOST", "0.0.0.0")
	if err != nil {
		return Config{}, err
	}
	port, err := envInt("HTTP_PORT", 8080)
	if err != nil {
		return Config{}, err
	}
	readTimeout, err := envDuration("HTTP_READ_TIMEOUT", 10*time.Second)
	if err != nil {
		return Config{}, err
	}
	readHeaderTimeout, err := envDuration("HTTP_READ_HEADER_TIMEOUT", 5*time.Second)
	if err != nil {
		return Config{}, err
	}
	writeTimeout, err := envDuration("HTTP_WRITE_TIMEOUT", 15*time.Second)
	if err != nil {
		return Config{}, err
	}
	idleTimeout, err := envDuration("HTTP_IDLE_TIMEOUT", 60*time.Second)
	if err != nil {
		return Config{}, err
	}
	shutdownTimeout, err := envDuration("HTTP_SHUTDOWN_TIMEOUT", 10*time.Second)
	if err != nil {
		return Config{}, err
	}
	maxBodyBytes, err := envInt64("HTTP_MAX_BODY_BYTES", 1<<20)
	if err != nil {
		return Config{}, err
	}
	tlsEnabled, err := envBool("HTTP_TLS_ENABLED", false)
	if err != nil {
		return Config{}, err
	}
	tlsCertFile, err := envString("HTTP_TLS_CERT_FILE", "")
	if err != nil {
		return Config{}, err
	}
	tlsKeyFile, err := envString("HTTP_TLS_KEY_FILE", "")
	if err != nil {
		return Config{}, err
	}
	trustedProxies := splitCSV(os.Getenv("HTTP_TRUSTED_PROXIES"))
	rateLimitEnabled, err := envBool("RATE_LIMIT_ENABLED", true)
	if err != nil {
		return Config{}, err
	}
	rateLimitRate, err := envFloat("RATE_LIMIT_RATE", 50)
	if err != nil {
		return Config{}, err
	}
	rateLimitBurst, err := envInt("RATE_LIMIT_BURST", 100)
	if err != nil {
		return Config{}, err
	}

	logLevel, err := envString("LOG_LEVEL", "info")
	if err != nil {
		return Config{}, err
	}
	logFormat, err := envString("LOG_FORMAT", "json")
	if err != nil {
		return Config{}, err
	}
	corsCredentials, err := envBool("CORS_ALLOW_CREDENTIALS", true)
	if err != nil {
		return Config{}, err
	}
	otelEnabled, err := envBool("OTEL_ENABLED", false)
	if err != nil {
		return Config{}, err
	}
	otelServiceName, err := envString("OTEL_SERVICE_NAME", "developers-cv-api")
	if err != nil {
		return Config{}, err
	}
	otelEndpoint, err := envString("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	if err != nil {
		return Config{}, err
	}
	otelInsecure, err := envBool("OTEL_EXPORTER_OTLP_INSECURE", false)
	if err != nil {
		return Config{}, err
	}
	prometheusEnabled, err := envBool("PROMETHEUS_ENABLED", true)
	if err != nil {
		return Config{}, err
	}

	postgresHost, err := envString("POSTGRES_HOST", "localhost")
	if err != nil {
		return Config{}, err
	}
	postgresPort, err := envInt("POSTGRES_PORT", 5432)
	if err != nil {
		return Config{}, err
	}
	postgresDatabase, err := envString("POSTGRES_DB", "developers_cv")
	if err != nil {
		return Config{}, err
	}
	postgresUser, err := envString("POSTGRES_USER", "developers_cv")
	if err != nil {
		return Config{}, err
	}
	postgresPassword, err := envString("POSTGRES_PASSWORD", "")
	if err != nil {
		return Config{}, err
	}
	postgresSSLMode, err := envString("POSTGRES_SSLMODE", "disable")
	if err != nil {
		return Config{}, err
	}
	postgresMaxConns, err := envInt("POSTGRES_MAX_CONNS", 10)
	if err != nil {
		return Config{}, err
	}
	postgresMinConns, err := envInt("POSTGRES_MIN_CONNS", 2)
	if err != nil {
		return Config{}, err
	}
	postgresMaxLifetime, err := envDuration("POSTGRES_MAX_CONN_LIFETIME", time.Hour)
	if err != nil {
		return Config{}, err
	}
	postgresMaxIdle, err := envDuration("POSTGRES_MAX_CONN_IDLE_TIME", 30*time.Minute)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		App:  AppConfig{Env: appEnv, Name: appName, Version: appVersion},
		Auth: AuthConfig{Mode: authMode},
		HTTP: HTTPConfig{
			Host:              host,
			Port:              port,
			ReadTimeout:       readTimeout,
			ReadHeaderTimeout: readHeaderTimeout,
			WriteTimeout:      writeTimeout,
			IdleTimeout:       idleTimeout,
			ShutdownTimeout:   shutdownTimeout,
			MaxBodyBytes:      maxBodyBytes,
			TLS: TLSConfig{
				Enabled:  tlsEnabled,
				CertFile: tlsCertFile,
				KeyFile:  tlsKeyFile,
			},
			TrustedProxies: trustedProxies,
			RateLimit: RateLimitConfig{
				Enabled: rateLimitEnabled,
				Rate:    rateLimitRate,
				Burst:   rateLimitBurst,
			},
		},
		Log: LogConfig{Level: logLevel, Format: logFormat},
		CORS: CORSConfig{
			AllowedOrigins:   splitCSV(os.Getenv("CORS_ALLOWED_ORIGINS")),
			AllowCredentials: corsCredentials,
		},
		Telemetry: TelemetryConfig{
			Enabled:      otelEnabled,
			ServiceName:  otelServiceName,
			OTLPEndpoint: otelEndpoint,
			Insecure:     otelInsecure,
		},
		Prometheus: PrometheusConfig{Enabled: prometheusEnabled},
		Postgres: PostgresConfig{
			Host:            postgresHost,
			Port:            postgresPort,
			Database:        postgresDatabase,
			User:            postgresUser,
			Password:        postgresPassword,
			SSLMode:         postgresSSLMode,
			MaxConns:        int32(postgresMaxConns),
			MinConns:        int32(postgresMinConns),
			MaxConnLifetime: postgresMaxLifetime,
			MaxConnIdleTime: postgresMaxIdle,
		},
	}

	if err := validate(cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c PostgresConfig) ConnectionString() string {
	values := url.Values{}
	values.Set("sslmode", c.SSLMode)

	return (&url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.User, c.Password),
		Host:     net.JoinHostPort(c.Host, strconv.Itoa(c.Port)),
		Path:     "/" + c.Database,
		RawQuery: values.Encode(),
	}).String()
}

func isValidPostgresSSLMode(value string) bool {
	switch value {
	case "disable", "allow", "prefer", "require", "verify-ca", "verify-full":
		return true
	default:
		return false
	}
}

func validate(cfg Config) error {
	if cfg.Auth.Mode != "mock" {
		return fmt.Errorf("invalid AUTH_MODE: %q", cfg.Auth.Mode)
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
	if cfg.Postgres.MaxConns <= 0 || cfg.Postgres.MinConns < 0 || cfg.Postgres.MinConns > cfg.Postgres.MaxConns {
		return errors.New("invalid PostgreSQL pool limits")
	}
	if cfg.Postgres.MaxConnLifetime <= 0 || cfg.Postgres.MaxConnIdleTime <= 0 {
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
	if cfg.HTTP.TLS.Enabled && (cfg.HTTP.TLS.CertFile == "" || cfg.HTTP.TLS.KeyFile == "") {
		return errors.New("HTTP_TLS_CERT_FILE and HTTP_TLS_KEY_FILE are required when TLS is enabled")
	}
	if cfg.HTTP.RateLimit.Enabled && (cfg.HTTP.RateLimit.Rate <= 0 || cfg.HTTP.RateLimit.Burst <= 0) {
		return errors.New("RATE_LIMIT_RATE and RATE_LIMIT_BURST must be positive")
	}
	if cfg.HTTP.MaxBodyBytes <= 0 {
		return fmt.Errorf("HTTP_MAX_BODY_BYTES must be positive: %d", cfg.HTTP.MaxBodyBytes)
	}
	if cfg.CORS.AllowCredentials && len(cfg.CORS.AllowedOrigins) == 0 {
		return errors.New("CORS_ALLOWED_ORIGINS must be configured when credentials are enabled")
	}
	for _, origin := range cfg.CORS.AllowedOrigins {
		if origin == "*" {
			return errors.New("CORS_ALLOWED_ORIGINS must not contain * when using an explicit allowlist")
		}
		parsed, err := url.ParseRequestURI(origin)
		if err != nil || parsed.Scheme != "http" &&
			parsed.Scheme != "https" || parsed.Host == "" ||
			parsed.Path != "" || parsed.RawQuery != "" ||
			parsed.Fragment != "" {
			return fmt.Errorf("invalid CORS origin %q: expected an http(s) origin without path, query or fragment", origin)
		}
	}
	if cfg.Log.Format != "json" && cfg.Log.Format != "console" {
		return fmt.Errorf("LOG_FORMAT must be json or console: %q", cfg.Log.Format)
	}
	if cfg.Telemetry.Enabled && strings.TrimSpace(cfg.Telemetry.OTLPEndpoint) == "" {
		return errors.New("OTEL_EXPORTER_OTLP_ENDPOINT is required when tracing is enabled")
	}
	return nil
}

func loadDotEnv(paths ...string) error {
	path := ".env"
	if len(paths) > 0 && strings.TrimSpace(paths[0]) != "" {
		path = paths[0]
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return fmt.Errorf("invalid line %q", line)
		}
		key = strings.TrimSpace(key)
		if key == "" {
			return fmt.Errorf("empty variable name in line %q", line)
		}
		value = strings.TrimSpace(value)
		if len(value) >= 2 &&
			((value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'')) {
			value = value[1 : len(value)-1]
			if strings.HasPrefix(line[strings.Index(line, "=")+1:], "\"") {
				if unquoted, err := strconv.Unquote("\"" + value + "\""); err == nil {
					value = unquoted
				}
			}
		}
		if _, exists := os.LookupEnv(key); !exists {
			if err := os.Setenv(key, value); err != nil {
				return fmt.Errorf("set %s: %w", key, err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read file: %w", err)
	}
	return nil
}

func envString(key, fallback string) (string, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}
	return value, nil
}

func envInt(key string, fallback int) (int, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer: %w", key, err)
	}
	return parsed, nil
}

func envInt64(key string, fallback int64) (int64, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer: %w", key, err)
	}
	return parsed, nil
}

func envDuration(key string, fallback time.Duration) (time.Duration, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a duration: %w", key, err)
	}
	return parsed, nil
}

func envFloat(key string, fallback float64) (float64, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("%s must be a number: %w", key, err)
	}
	return parsed, nil
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

func envBool(key string, fallback bool) (bool, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("%s must be a boolean: %w", key, err)
	}
	return parsed, nil
}

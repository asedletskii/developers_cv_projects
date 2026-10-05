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
	Turnstile  TurnstileConfig
	Log        LogConfig
	CORS       CORSConfig
	Telemetry  TelemetryConfig
	Analytics  AnalyticsConfig
	Prometheus PrometheusConfig
	Postgres   PostgresConfig
}

type AuthConfig struct {
	Mode        string
	Telegram    TelegramAuthConfig
	Session     SessionAuthConfig
	Cookie      CookieAuthConfig
	Transaction TransactionAuthConfig
	RateLimit   AuthRateLimitConfig
}

type AuthRateLimitConfig struct {
	Enabled               bool
	StartRatePerMinute    float64
	StartBurst            int
	CallbackRatePerMinute float64
	CallbackBurst         int
}

type TelegramAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	IssuerURL    string
}

type SessionAuthConfig struct {
	AbsoluteTTL   time.Duration
	IdleTTL       time.Duration
	TouchInterval time.Duration
}

type CookieAuthConfig struct {
	Name     string
	Secure   bool
	HTTPOnly bool
	SameSite string
	Path     string
}

type TransactionAuthConfig struct {
	TTL time.Duration
}

type AppConfig struct {
	Env     string
	Name    string
	Version string
}

type HTTPConfig struct {
	Host                     string
	Port                     int
	ReadTimeout              time.Duration
	ReadHeaderTimeout        time.Duration
	WriteTimeout             time.Duration
	IdleTimeout              time.Duration
	ShutdownTimeout          time.Duration
	MaxBodyBytes             int64
	TLS                      TLSConfig
	TrustedProxies           []string
	RateLimit                RateLimitConfig
	ContactViewRateLimit     RateLimitConfig
	AnalyticsEventsRateLimit RateLimitConfig
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

type TurnstileConfig struct {
	Enabled        bool
	SiteKey        string
	SecretKey      string
	ExpectedHost   string
	ExpectedAction string
	Timeout        time.Duration
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
	Enabled        bool
	ServiceName    string
	ServiceVersion string
	Environment    string
	OTLPEndpoint   string
	Insecure       bool
}

type AnalyticsConfig struct {
	VisitorSecret     string
	SuspectBurst      int
	SuspectWindow     time.Duration
	RespectDoNotTrack bool
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
	MaxConns        int
	MaxIdleConns    int
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
	authTelegramClientID, err := envString("AUTH_TELEGRAM_CLIENT_ID", "")
	if err != nil {
		return Config{}, err
	}
	authTelegramClientSecret, err := envString("AUTH_TELEGRAM_CLIENT_SECRET", "")
	if err != nil {
		return Config{}, err
	}
	authTelegramRedirectURL, err := envString("AUTH_TELEGRAM_REDIRECT_URL", "http://localhost:8080/api/v1/auth/telegram/callback")
	if err != nil {
		return Config{}, err
	}
	authTelegramIssuerURL, err := envString("AUTH_TELEGRAM_ISSUER_URL", "https://oauth.telegram.org")
	if err != nil {
		return Config{}, err
	}
	authSessionAbsoluteTTL, err := envDuration("AUTH_SESSION_ABSOLUTE_TTL", 720*time.Hour)
	if err != nil {
		return Config{}, err
	}
	authSessionIdleTTL, err := envDuration("AUTH_SESSION_IDLE_TTL", 168*time.Hour)
	if err != nil {
		return Config{}, err
	}
	authSessionTouchInterval, err := envDuration("AUTH_SESSION_TOUCH_INTERVAL", 5*time.Minute)
	if err != nil {
		return Config{}, err
	}
	authTransactionTTL, err := envDuration("AUTH_TRANSACTION_TTL", 5*time.Minute)
	if err != nil {
		return Config{}, err
	}
	authCookieName, err := envString("AUTH_COOKIE_NAME", "__Host-session")
	if err != nil {
		return Config{}, err
	}
	authCookieSecure, err := envBool("AUTH_COOKIE_SECURE", true)
	if err != nil {
		return Config{}, err
	}
	authCookieHTTPOnly, err := envBool("AUTH_COOKIE_HTTP_ONLY", true)
	if err != nil {
		return Config{}, err
	}
	authCookieSameSite, err := envString("AUTH_COOKIE_SAME_SITE", "lax")
	if err != nil {
		return Config{}, err
	}
	authCookiePath, err := envString("AUTH_COOKIE_PATH", "/")
	if err != nil {
		return Config{}, err
	}
	authRateLimitEnabled, err := envBool("AUTH_RATE_LIMIT_ENABLED", true)
	if err != nil {
		return Config{}, err
	}
	authStartRatePerMinute, err := envFloat("AUTH_RATE_LIMIT_START_PER_MINUTE", 5)
	if err != nil {
		return Config{}, err
	}
	authStartBurst, err := envInt("AUTH_RATE_LIMIT_START_BURST", 5)
	if err != nil {
		return Config{}, err
	}
	authCallbackRatePerMinute, err := envFloat("AUTH_RATE_LIMIT_CALLBACK_PER_MINUTE", 20)
	if err != nil {
		return Config{}, err
	}
	authCallbackBurst, err := envInt("AUTH_RATE_LIMIT_CALLBACK_BURST", 20)
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
	contactViewRateLimitEnabled, err := envBool("CONTACT_VIEW_RATE_LIMIT_ENABLED", true)
	if err != nil {
		return Config{}, err
	}
	contactViewRateLimitPerMinute, err := envFloat("CONTACT_VIEW_RATE_LIMIT_PER_MINUTE", 3)
	if err != nil {
		return Config{}, err
	}
	contactViewRateLimitBurst, err := envInt("CONTACT_VIEW_RATE_LIMIT_BURST", 2)
	if err != nil {
		return Config{}, err
	}
	analyticsEventsRateLimitEnabled, err := envBool("ANALYTICS_EVENTS_RATE_LIMIT_ENABLED", true)
	if err != nil {
		return Config{}, err
	}
	analyticsEventsRateLimitPerMinute, err := envFloat("ANALYTICS_EVENTS_RATE_LIMIT_PER_MINUTE", 30)
	if err != nil {
		return Config{}, err
	}
	analyticsEventsRateLimitBurst, err := envInt("ANALYTICS_EVENTS_RATE_LIMIT_BURST", 10)
	if err != nil {
		return Config{}, err
	}

	turnstileEnabled, err := envBool("TURNSTILE_ENABLED", false)
	if err != nil {
		return Config{}, err
	}
	turnstileSiteKey, err := envString("TURNSTILE_SITE_KEY", "")
	if err != nil {
		return Config{}, err
	}
	turnstileSecretKey, err := envString("TURNSTILE_SECRET_KEY", "")
	if err != nil {
		return Config{}, err
	}
	turnstileExpectedHost, err := envString("TURNSTILE_EXPECTED_HOSTNAME", "")
	if err != nil {
		return Config{}, err
	}
	turnstileExpectedAction, err := envString("TURNSTILE_EXPECTED_ACTION", "contact_view")
	if err != nil {
		return Config{}, err
	}
	turnstileTimeout, err := envDuration("TURNSTILE_TIMEOUT", 5*time.Second)
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
	analyticsVisitorSecret, err := envString("ANALYTICS_VISITOR_SECRET", "development-only-analytics-secret-change-me-32")
	if err != nil {
		return Config{}, err
	}
	analyticsSuspectBurst, err := envInt("ANALYTICS_SUSPECT_BURST", 20)
	if err != nil {
		return Config{}, err
	}
	analyticsSuspectWindow, err := envDuration("ANALYTICS_SUSPECT_WINDOW", time.Minute)
	if err != nil {
		return Config{}, err
	}
	analyticsRespectDoNotTrack, err := envBool("RESPECT_DO_NOT_TRACK", true)
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
	postgresMaxIdleConns, err := envInt("POSTGRES_MAX_IDLE_CONNS", 2)
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
		App: AppConfig{
			Env:     appEnv,
			Name:    appName,
			Version: appVersion,
		},
		Auth: AuthConfig{
			Mode:        authMode,
			Telegram:    TelegramAuthConfig{ClientID: authTelegramClientID, ClientSecret: authTelegramClientSecret, RedirectURL: authTelegramRedirectURL, IssuerURL: authTelegramIssuerURL},
			Session:     SessionAuthConfig{AbsoluteTTL: authSessionAbsoluteTTL, IdleTTL: authSessionIdleTTL, TouchInterval: authSessionTouchInterval},
			Cookie:      CookieAuthConfig{Name: authCookieName, Secure: authCookieSecure, HTTPOnly: authCookieHTTPOnly, SameSite: authCookieSameSite, Path: authCookiePath},
			Transaction: TransactionAuthConfig{TTL: authTransactionTTL},
			RateLimit:   AuthRateLimitConfig{Enabled: authRateLimitEnabled, StartRatePerMinute: authStartRatePerMinute, StartBurst: authStartBurst, CallbackRatePerMinute: authCallbackRatePerMinute, CallbackBurst: authCallbackBurst},
		},
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
			ContactViewRateLimit: RateLimitConfig{
				Enabled: contactViewRateLimitEnabled,
				Rate:    contactViewRateLimitPerMinute,
				Burst:   contactViewRateLimitBurst,
			},
			AnalyticsEventsRateLimit: RateLimitConfig{
				Enabled: analyticsEventsRateLimitEnabled,
				Rate:    analyticsEventsRateLimitPerMinute,
				Burst:   analyticsEventsRateLimitBurst,
			},
		},
		Turnstile: TurnstileConfig{
			Enabled:        turnstileEnabled,
			SiteKey:        turnstileSiteKey,
			SecretKey:      turnstileSecretKey,
			ExpectedHost:   turnstileExpectedHost,
			ExpectedAction: turnstileExpectedAction,
			Timeout:        turnstileTimeout,
		},
		Log: LogConfig{
			Level:  logLevel,
			Format: logFormat,
		},
		CORS: CORSConfig{
			AllowedOrigins:   splitCSV(os.Getenv("CORS_ALLOWED_ORIGINS")),
			AllowCredentials: corsCredentials,
		},
		Telemetry: TelemetryConfig{
			Enabled:        otelEnabled,
			ServiceName:    otelServiceName,
			ServiceVersion: appVersion,
			Environment:    appEnv,
			OTLPEndpoint:   otelEndpoint,
			Insecure:       otelInsecure,
		},
		Prometheus: PrometheusConfig{
			Enabled: prometheusEnabled,
		},
		Analytics: AnalyticsConfig{VisitorSecret: analyticsVisitorSecret, SuspectBurst: analyticsSuspectBurst, SuspectWindow: analyticsSuspectWindow, RespectDoNotTrack: analyticsRespectDoNotTrack},
		Postgres: PostgresConfig{
			Host:            postgresHost,
			Port:            postgresPort,
			Database:        postgresDatabase,
			User:            postgresUser,
			Password:        postgresPassword,
			SSLMode:         postgresSSLMode,
			MaxConns:        postgresMaxConns,
			MaxIdleConns:    postgresMaxIdleConns,
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
	if len([]byte(strings.TrimSpace(cfg.Analytics.VisitorSecret))) < 32 {
		return errors.New("ANALYTICS_VISITOR_SECRET must be at least 32 bytes")
	}
	if cfg.Analytics.SuspectBurst <= 0 || cfg.Analytics.SuspectWindow <= 0 {
		return errors.New("analytics suspect limits must be positive")
	}
	if cfg.Auth.Mode != "mock" && cfg.Auth.Mode != "telegram" {
		return fmt.Errorf("invalid AUTH_MODE: %q", cfg.Auth.Mode)
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
	if cfg.Auth.RateLimit.Enabled && (cfg.Auth.RateLimit.StartRatePerMinute <= 0 || cfg.Auth.RateLimit.StartBurst <= 0 || cfg.Auth.RateLimit.CallbackRatePerMinute <= 0 || cfg.Auth.RateLimit.CallbackBurst <= 0) {
		return errors.New("auth rate limits must be positive when enabled")
	}
	if strings.TrimSpace(cfg.Auth.Cookie.Name) == "" || cfg.Auth.Cookie.Path != "/" || !cfg.Auth.Cookie.HTTPOnly || strings.ToLower(cfg.Auth.Cookie.SameSite) != "lax" {
		return errors.New("invalid auth cookie configuration")
	}
	if strings.HasPrefix(cfg.Auth.Cookie.Name, "__Host-") && (!cfg.Auth.Cookie.Secure || cfg.Auth.Cookie.Path != "/") {
		return errors.New("__Host- auth cookie requires Secure and Path=/")
	}
	if cfg.Auth.Mode == "telegram" {
		if strings.TrimSpace(cfg.Auth.Telegram.ClientID) == "" || strings.TrimSpace(cfg.Auth.Telegram.ClientSecret) == "" || strings.TrimSpace(cfg.Auth.Telegram.RedirectURL) == "" || strings.TrimSpace(cfg.Auth.Telegram.IssuerURL) == "" {
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
		(cfg.HTTP.RateLimit.Rate <= 0 || cfg.HTTP.RateLimit.Burst <= 0) {
		return errors.New("RATE_LIMIT_RATE and RATE_LIMIT_BURST must be positive")
	}
	if cfg.HTTP.ContactViewRateLimit.Enabled &&
		(cfg.HTTP.ContactViewRateLimit.Rate <= 0 || cfg.HTTP.ContactViewRateLimit.Burst <= 0) {
		return errors.New("CONTACT_VIEW_RATE_LIMIT_PER_MINUTE and CONTACT_VIEW_RATE_LIMIT_BURST must be positive")
	}
	if cfg.HTTP.AnalyticsEventsRateLimit.Enabled &&
		(cfg.HTTP.AnalyticsEventsRateLimit.Rate <= 0 || cfg.HTTP.AnalyticsEventsRateLimit.Burst <= 0) {
		return errors.New("ANALYTICS_EVENTS_RATE_LIMIT_PER_MINUTE and ANALYTICS_EVENTS_RATE_LIMIT_BURST must be positive")
	}
	if cfg.Analytics.SuspectBurst <= 0 {
		return errors.New("ANALYTICS_SUSPECT_BURST must be positive")
	}
	if cfg.Analytics.SuspectWindow <= 0 {
		return errors.New("ANALYTICS_SUSPECT_WINDOW must be positive")
	}
	if cfg.HTTP.MaxBodyBytes <= 0 {
		return fmt.Errorf("HTTP_MAX_BODY_BYTES must be positive: %d", cfg.HTTP.MaxBodyBytes)
	}
	if cfg.Turnstile.Timeout <= 0 {
		return errors.New("TURNSTILE_TIMEOUT must be positive")
	}
	if cfg.Turnstile.Enabled && strings.TrimSpace(cfg.Turnstile.SiteKey) == "" {
		return errors.New("TURNSTILE_SITE_KEY is required when Turnstile is enabled")
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

func validateProductionSecurity(cfg Config) error {
	if cfg.App.Env != "production" {
		return nil
	}
	if cfg.Auth.Mode != "telegram" {
		return errors.New("AUTH_MODE must be telegram in production")
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
	if !cfg.Prometheus.Enabled {
		return errors.New("PROMETHEUS_ENABLED must be true in production")
	}
	if strings.HasPrefix(cfg.Analytics.VisitorSecret, "development-only-") {
		return errors.New("ANALYTICS_VISITOR_SECRET must be changed in production")
	}
	if !cfg.Turnstile.Enabled {
		return errors.New("TURNSTILE_ENABLED must be true in production")
	}
	if strings.TrimSpace(cfg.Turnstile.SiteKey) == "" {
		return errors.New("TURNSTILE_SITE_KEY is required in production")
	}
	if strings.TrimSpace(cfg.Turnstile.SecretKey) == "" {
		return errors.New("TURNSTILE_SECRET_KEY is required in production")
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

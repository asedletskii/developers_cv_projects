package config

import "time"

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
	Name   string
	Secure bool
}

type TransactionAuthConfig struct {
	TTL time.Duration
}

type AppConfig struct {
	Env     string
	Version string
}

type HTTPConfig struct {
	Host                     string
	MetricsAddr              string
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
	AllowedOrigins      []string
	AllowCredentials    bool
	AllowMockAuthHeader bool
}

type TelemetryConfig struct {
	Enabled        bool
	ServiceName    string
	ServiceVersion string
	Environment    string
	OTLPEndpoint   string
	Insecure       bool
	SampleRatio    float64
}

type AnalyticsConfig struct {
	VisitorSecret      string
	SuspectBurst       int
	SuspectWindow      time.Duration
	ClickRetentionDays int
	ViewRetentionDays  int
}

type PrometheusConfig struct {
	Enabled bool
}

type PostgresConfig struct {
	Host                     string
	Port                     int
	Database                 string
	User                     string
	Password                 string
	SSLMode                  string
	StatementTimeout         time.Duration
	LockTimeout              time.Duration
	IdleInTransactionTimeout time.Duration
	SlowQueryThreshold       time.Duration
	MaxConns                 int
	MaxIdleConns             int
	MaxConnLifetime          time.Duration
	MaxConnIdleTime          time.Duration
}

package config

import (
	"os"
	"time"
)

func parseHTTPConfig() (HTTPConfig, error) {
	port, err := envInt("HTTP_PORT", 8080)
	if err != nil {
		return HTTPConfig{}, err
	}
	readTimeout, err := envDuration("HTTP_READ_TIMEOUT", 10*time.Second)
	if err != nil {
		return HTTPConfig{}, err
	}
	readHeaderTimeout, err := envDuration("HTTP_READ_HEADER_TIMEOUT", 5*time.Second)
	if err != nil {
		return HTTPConfig{}, err
	}
	writeTimeout, err := envDuration("HTTP_WRITE_TIMEOUT", 15*time.Second)
	if err != nil {
		return HTTPConfig{}, err
	}
	idleTimeout, err := envDuration("HTTP_IDLE_TIMEOUT", 60*time.Second)
	if err != nil {
		return HTTPConfig{}, err
	}
	shutdownTimeout, err := envDuration("HTTP_SHUTDOWN_TIMEOUT", 10*time.Second)
	if err != nil {
		return HTTPConfig{}, err
	}
	maxBodyBytes, err := envInt64("HTTP_MAX_BODY_BYTES", 1<<20)
	if err != nil {
		return HTTPConfig{}, err
	}
	tlsEnabled, err := envBool("HTTP_TLS_ENABLED", false)
	if err != nil {
		return HTTPConfig{}, err
	}
	limitEnabled, err := envBool("RATE_LIMIT_ENABLED", true)
	if err != nil {
		return HTTPConfig{}, err
	}
	limitRate, err := envFloat("RATE_LIMIT_RATE", 50)
	if err != nil {
		return HTTPConfig{}, err
	}
	limitBurst, err := envInt("RATE_LIMIT_BURST", 100)
	if err != nil {
		return HTTPConfig{}, err
	}
	contactEnabled, err := envBool("CONTACT_VIEW_RATE_LIMIT_ENABLED", true)
	if err != nil {
		return HTTPConfig{}, err
	}
	contactRate, err := envFloat("CONTACT_VIEW_RATE_LIMIT_PER_MINUTE", 3)
	if err != nil {
		return HTTPConfig{}, err
	}
	contactBurst, err := envInt("CONTACT_VIEW_RATE_LIMIT_BURST", 2)
	if err != nil {
		return HTTPConfig{}, err
	}
	analyticsEnabled, err := envBool("ANALYTICS_EVENTS_RATE_LIMIT_ENABLED", true)
	if err != nil {
		return HTTPConfig{}, err
	}
	analyticsRate, err := envFloat("ANALYTICS_EVENTS_RATE_LIMIT_PER_MINUTE", 30)
	if err != nil {
		return HTTPConfig{}, err
	}
	analyticsBurst, err := envInt("ANALYTICS_EVENTS_RATE_LIMIT_BURST", 10)
	if err != nil {
		return HTTPConfig{}, err
	}
	return HTTPConfig{
		Host: envString("HTTP_HOST", "0.0.0.0"), MetricsAddr: envString("METRICS_ADDR", "127.0.0.1:9090"), Port: port,
		ReadTimeout: readTimeout, ReadHeaderTimeout: readHeaderTimeout, WriteTimeout: writeTimeout, IdleTimeout: idleTimeout,
		ShutdownTimeout: shutdownTimeout, MaxBodyBytes: maxBodyBytes,
		TLS: TLSConfig{
			Enabled:  tlsEnabled,
			CertFile: envString("HTTP_TLS_CERT_FILE", ""),
			KeyFile:  envString("HTTP_TLS_KEY_FILE", ""),
		},
		TrustedProxies:           splitCSV(os.Getenv("HTTP_TRUSTED_PROXIES")),
		RateLimit:                RateLimitConfig{Enabled: limitEnabled, Rate: limitRate, Burst: limitBurst},
		ContactViewRateLimit:     RateLimitConfig{Enabled: contactEnabled, Rate: contactRate, Burst: contactBurst},
		AnalyticsEventsRateLimit: RateLimitConfig{Enabled: analyticsEnabled, Rate: analyticsRate, Burst: analyticsBurst},
	}, nil
}

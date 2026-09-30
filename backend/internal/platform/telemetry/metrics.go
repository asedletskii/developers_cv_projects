package telemetry

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	requests                 *prometheus.CounterVec
	duration                 *prometheus.HistogramVec
	inFlight                 prometheus.Gauge
	authLoginAttempts        *prometheus.CounterVec
	authLoginSuccess         *prometheus.CounterVec
	authLoginFailures        *prometheus.CounterVec
	authLogout               *prometheus.CounterVec
	authSessionRejections    *prometheus.CounterVec
	authTransactionCleanup   *prometheus.CounterVec
	authCleanupFailures      *prometheus.CounterVec
	authLoginDuration        *prometheus.HistogramVec
	authOIDCExchangeDuration *prometheus.HistogramVec
	authOIDCVerifyDuration   *prometheus.HistogramVec
	registry                 *prometheus.Registry
}

func NewMetrics(enabled bool) (*Metrics, error) {
	if !enabled {
		return nil, nil
	}

	requests := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "developers_cv",
		Subsystem: "http",
		Name:      "requests_total",
		Help:      "Total number of HTTP requests.",
	}, []string{"method", "route", "status"})
	duration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "developers_cv",
		Subsystem: "http",
		Name:      "request_duration_seconds",
		Help:      "HTTP request duration in seconds.",
	}, []string{"method", "route"})
	inFlight := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "developers_cv",
		Subsystem: "http",
		Name:      "requests_in_flight",
		Help:      "Current number of in-flight HTTP requests.",
	})

	authLoginAttempts := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "developers_cv", Name: "auth_login_attempts_total", Help: "Total number of authentication login attempts.",
	}, []string{"provider"})
	authLoginSuccess := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "developers_cv", Name: "auth_login_success_total", Help: "Total number of successful authentication logins.",
	}, []string{"provider"})
	authLoginFailures := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "developers_cv", Name: "auth_login_failures_total", Help: "Total number of failed authentication logins.",
	}, []string{"provider", "reason"})
	authLogout := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "developers_cv", Name: "auth_logout_total", Help: "Total number of authentication logout operations.",
	}, []string{"provider", "reason"})
	authSessionRejections := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "developers_cv", Name: "auth_session_rejections_total", Help: "Total number of rejected authentication sessions.",
	}, []string{"provider", "reason"})
	authTransactionCleanup := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "developers_cv", Name: "auth_transaction_cleanup_total", Help: "Total number of expired authentication transactions removed.",
	}, []string{"provider", "reason"})
	authCleanupFailures := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "developers_cv", Name: "auth_cleanup_failures_total", Help: "Total number of authentication cleanup failures.",
	}, []string{"resource"})
	authLoginDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "developers_cv", Name: "auth_login_duration_seconds", Help: "Authentication login duration in seconds.",
	}, []string{"provider"})
	authOIDCExchangeDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "developers_cv", Name: "auth_oidc_exchange_duration_seconds", Help: "Telegram OIDC token exchange duration in seconds.",
	}, []string{"provider"})
	authOIDCVerifyDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "developers_cv", Name: "auth_oidc_verify_duration_seconds", Help: "Telegram OIDC ID token verification duration in seconds.",
	}, []string{"provider"})

	registry := prometheus.NewRegistry()
	collectors := []prometheus.Collector{
		requests, duration, inFlight,
		authLoginAttempts, authLoginSuccess, authLoginFailures, authLogout,
		authSessionRejections, authTransactionCleanup, authCleanupFailures, authLoginDuration,
		authOIDCExchangeDuration, authOIDCVerifyDuration,
	}
	for _, collector := range collectors {
		if err := registry.Register(collector); err != nil {
			return nil, fmt.Errorf("register telemetry metric: %w", err)
		}
	}
	return &Metrics{
		requests: requests, duration: duration, inFlight: inFlight,
		authLoginAttempts: authLoginAttempts, authLoginSuccess: authLoginSuccess,
		authLoginFailures: authLoginFailures, authLogout: authLogout,
		authSessionRejections: authSessionRejections, authTransactionCleanup: authTransactionCleanup,
		authCleanupFailures: authCleanupFailures, authLoginDuration: authLoginDuration, authOIDCExchangeDuration: authOIDCExchangeDuration,
		authOIDCVerifyDuration: authOIDCVerifyDuration, registry: registry,
	}, nil
}

func (m *Metrics) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		m.inFlight.Inc()
		defer m.inFlight.Dec()
		c.Next()

		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}
		m.requests.WithLabelValues(c.Request.Method, route, strconv.Itoa(c.Writer.Status())).Inc()
		m.duration.WithLabelValues(c.Request.Method, route).Observe(time.Since(start).Seconds())
	}
}

func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}

func (m *Metrics) RecordAuthLoginAttempt(provider string) {
	if m == nil || m.authLoginAttempts == nil {
		return
	}
	m.authLoginAttempts.WithLabelValues(provider).Inc()
}

func (m *Metrics) RecordAuthLoginSuccess(provider string, duration time.Duration) {
	if m == nil || m.authLoginSuccess == nil {
		return
	}
	m.authLoginSuccess.WithLabelValues(provider).Inc()
	m.authLoginDuration.WithLabelValues(provider).Observe(duration.Seconds())
}

func (m *Metrics) RecordAuthLoginFailure(provider, reason string, duration time.Duration) {
	if m == nil || m.authLoginFailures == nil {
		return
	}
	m.authLoginFailures.WithLabelValues(provider, reason).Inc()
	m.authLoginDuration.WithLabelValues(provider).Observe(duration.Seconds())
}

func (m *Metrics) RecordAuthLogout(provider, reason string) {
	if m == nil || m.authLogout == nil {
		return
	}
	m.authLogout.WithLabelValues(provider, reason).Inc()
}

func (m *Metrics) RecordAuthSessionRejection(provider, reason string) {
	if m == nil || m.authSessionRejections == nil {
		return
	}
	m.authSessionRejections.WithLabelValues(provider, reason).Inc()
}

func (m *Metrics) RecordAuthTransactionCleanup(provider, reason string, count int) {
	if m == nil || m.authTransactionCleanup == nil || count <= 0 {
		return
	}
	m.authTransactionCleanup.WithLabelValues(provider, reason).Add(float64(count))
}

func (m *Metrics) ObserveAuthOIDCExchange(provider string, duration time.Duration) {
	if m == nil || m.authOIDCExchangeDuration == nil {
		return
	}
	m.authOIDCExchangeDuration.WithLabelValues(provider).Observe(duration.Seconds())
}

func (m *Metrics) ObserveAuthOIDCVerify(provider string, duration time.Duration) {
	if m == nil || m.authOIDCVerifyDuration == nil {
		return
	}
	m.authOIDCVerifyDuration.WithLabelValues(provider).Observe(duration.Seconds())
}

func (m *Metrics) RecordAuthCleanupFailure(resource string) {
	if m == nil || m.authCleanupFailures == nil {
		return
	}
	m.authCleanupFailures.WithLabelValues(resource).Inc()
}

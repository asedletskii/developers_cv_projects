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
	requests                   *prometheus.CounterVec
	duration                   *prometheus.HistogramVec
	inFlight                   prometheus.Gauge
	authLoginAttempts          *prometheus.CounterVec
	authLoginSuccess           *prometheus.CounterVec
	authLoginFailures          *prometheus.CounterVec
	authLogout                 *prometheus.CounterVec
	authSessionRejections      *prometheus.CounterVec
	authTransactionCleanup     *prometheus.CounterVec
	authSessionCleanup         *prometheus.CounterVec
	authCleanupFailures        *prometheus.CounterVec
	authLoginDuration          *prometheus.HistogramVec
	authOIDCExchangeDuration   *prometheus.HistogramVec
	authOIDCVerifyDuration     *prometheus.HistogramVec
	analyticsDroppedInput      *prometheus.CounterVec
	analyticsBatchOutcomes     *prometheus.CounterVec
	analyticsRetentionRuns     *prometheus.CounterVec
	analyticsRetentionDeleted  *prometheus.CounterVec
	analyticsRetentionDuration *prometheus.HistogramVec
	schemaOrphanRows           *prometheus.GaugeVec
	registry                   *prometheus.Registry
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
	authSessionCleanup := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "developers_cv", Name: "auth_session_cleanup_total", Help: "Total number of expired authentication sessions removed.",
	}, []string{"reason"})
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
	analyticsDroppedInput := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "developers_cv",
		Name:      "analytics_dropped_input_total",
		Help:      "Total number of analytics inputs dropped during validation.",
	}, []string{"kind", "reason"})
	analyticsBatchOutcomes := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "developers_cv", Subsystem: "analytics", Name: "batch_outcomes_total",
		Help: "Total analytics ingestion batch outcomes by bounded outcome and reason.",
	}, []string{"outcome", "reason"})
	analyticsRetentionRuns := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "developers_cv", Subsystem: "analytics_retention", Name: "runs_total",
		Help: "Total analytics retention cleanup runs by outcome.",
	}, []string{"outcome"})
	analyticsRetentionDeleted := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "developers_cv", Subsystem: "analytics_retention", Name: "deleted_total",
		Help: "Total expired analytics records deleted by resource type.",
	}, []string{"resource"})
	analyticsRetentionDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "developers_cv", Subsystem: "analytics_retention", Name: "duration_seconds",
		Help:    "Duration of analytics retention cleanup runs in seconds.",
		Buckets: prometheus.DefBuckets,
	}, []string{"outcome"})
	schemaOrphanRows := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "developers_cv", Name: "schema_orphan_rows",
		Help: "Current number of cross-module orphan rows by relation.",
	}, []string{"relation"})

	registry := prometheus.NewRegistry()
	collectors := []prometheus.Collector{
		requests, duration, inFlight,
		authLoginAttempts, authLoginSuccess, authLoginFailures, authLogout,
		authSessionRejections, authTransactionCleanup, authSessionCleanup, authCleanupFailures, authLoginDuration,
		authOIDCExchangeDuration, authOIDCVerifyDuration, analyticsDroppedInput, analyticsBatchOutcomes,
		analyticsRetentionRuns, analyticsRetentionDeleted, analyticsRetentionDuration,
		schemaOrphanRows,
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
		authSessionCleanup: authSessionCleanup, authCleanupFailures: authCleanupFailures, authLoginDuration: authLoginDuration, authOIDCExchangeDuration: authOIDCExchangeDuration,
		authOIDCVerifyDuration: authOIDCVerifyDuration, analyticsDroppedInput: analyticsDroppedInput,
		analyticsBatchOutcomes: analyticsBatchOutcomes, analyticsRetentionRuns: analyticsRetentionRuns, analyticsRetentionDeleted: analyticsRetentionDeleted,
		analyticsRetentionDuration: analyticsRetentionDuration, schemaOrphanRows: schemaOrphanRows, registry: registry,
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

func (m *Metrics) RecordAuthSessionCleanup(reason string, count int) {
	if m == nil || m.authSessionCleanup == nil || count <= 0 {
		return
	}
	m.authSessionCleanup.WithLabelValues(reason).Add(float64(count))
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

func (m *Metrics) RecordDroppedAnalyticsInput(kind, reason string) {
	if m == nil || m.analyticsDroppedInput == nil {
		return
	}
	kind = boundedAnalyticsKind(kind)
	reason = boundedAnalyticsDropReason(reason)
	m.analyticsDroppedInput.WithLabelValues(kind, reason).Inc()
}

func (m *Metrics) RecordAnalyticsBatchOutcome(outcome, reason string) {
	if m == nil || m.analyticsBatchOutcomes == nil {
		return
	}
	switch outcome {
	case "accepted", "ignored", "rejected", "failure":
	default:
		outcome = "failure"
		reason = "unknown"
	}
	reason = boundedAnalyticsBatchReason(reason)
	m.analyticsBatchOutcomes.WithLabelValues(outcome, reason).Inc()
}

func boundedAnalyticsKind(kind string) string {
	switch kind {
	case "event", "click":
		return kind
	default:
		return "other"
	}
}

func boundedAnalyticsDropReason(reason string) string {
	switch reason {
	case "invalid_json_shape", "invalid_fields", "negative_timestamp", "invalid_client_id",
		"missing_client_id", "duplicate_client_id", "technology_not_in_resume",
		"invalid_payload_or_type", "section_not_in_layout", "invalid_click",
		"contact_view_routed_separately",
		"per_view_limit", "print_deduplicated":
		return reason
	default:
		return "other"
	}
}

func boundedAnalyticsBatchReason(reason string) string {
	switch reason {
	case "success", "bot_traffic", "stale_sequence", "same_sequence", "invalid_batch",
		"batch_limit", "view_rejected", "invalid_state", "load_failed", "save_failed",
		"transaction_failed", "repository_mismatch", "unknown":
		return reason
	default:
		return "unknown"
	}
}

func (m *Metrics) RecordAnalyticsRetentionCleanup(outcome string, clicks, views int, duration time.Duration) {
	if m == nil || m.analyticsRetentionRuns == nil {
		return
	}
	switch outcome {
	case "success", "failure":
	default:
		outcome = "failure"
	}
	m.analyticsRetentionRuns.WithLabelValues(outcome).Inc()
	m.analyticsRetentionDuration.WithLabelValues(outcome).Observe(duration.Seconds())
	if clicks > 0 {
		m.analyticsRetentionDeleted.WithLabelValues("clicks").Add(float64(clicks))
	}
	if views > 0 {
		m.analyticsRetentionDeleted.WithLabelValues("views").Add(float64(views))
	}
}

func (m *Metrics) RecordSchemaOrphanRows(relation string, count int64) {
	if m == nil || m.schemaOrphanRows == nil {
		return
	}
	switch relation {
	case "analytics_without_link", "links_without_resume", "resumes_without_user",
		"auth_identities_without_user", "auth_sessions_without_user":
	default:
		relation = "other"
	}
	m.schemaOrphanRows.WithLabelValues(relation).Set(float64(count))
}

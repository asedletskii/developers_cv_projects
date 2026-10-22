package telemetry

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestNewMetricsDisabled(t *testing.T) {
	metrics, err := NewMetrics(false)

	require.NoError(t, err)
	require.Nil(t, metrics)
}

func TestMetricsMiddlewareRecordsRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	metrics, err := NewMetrics(true)
	require.NoError(t, err)

	engine := gin.New()
	engine.Use(metrics.Middleware())
	engine.GET("/resumes/:resumeId", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodGet, "/resumes/123", nil)
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	require.Equal(t, http.StatusOK, response.Code)

	metricsResponse := httptest.NewRecorder()
	metrics.Handler().ServeHTTP(metricsResponse, httptest.NewRequest(http.MethodGet, "/metrics", nil))

	require.Equal(t, http.StatusOK, metricsResponse.Code)
	body := metricsResponse.Body.String()
	require.Contains(t, body, `developers_cv_http_requests_total`)
	require.Contains(t, body, `route="/resumes/:resumeId"`)
	require.Contains(t, body, `status="200"`)
	require.Contains(t, body, `developers_cv_http_request_duration_seconds`)
	require.Contains(t, body, `developers_cv_http_requests_in_flight`)
	require.NotContains(t, body, `route="/resumes/123"`)

	metrics.RecordAuthLoginAttempt("telegram")
	metrics.RecordAuthLoginSuccess("telegram", time.Second)
	metrics.RecordAuthLoginFailure("telegram", "invalid_request", time.Second)
	metrics.RecordAuthLogout("telegram", "success")
	metrics.RecordAuthSessionRejection("telegram", "session_missing")
	metrics.RecordAuthTransactionCleanup("telegram", "expired", 3)
	metrics.RecordAuthSessionCleanup("expired", 5)
	metrics.RecordAuthCleanupFailure("sessions")
	metrics.ObserveAuthOIDCExchange("telegram", time.Second)
	metrics.ObserveAuthOIDCVerify("telegram", time.Second)
	metrics.RecordDroppedAnalyticsInput("event", "invalid_payload_or_type")
	metrics.RecordDroppedAnalyticsInput("unbounded-kind", "user-provided-reason")
	metrics.RecordAnalyticsBatchOutcome("accepted", "success")
	metrics.RecordAnalyticsBatchOutcome("failure", "save_failed")
	metrics.RecordAnalyticsBatchOutcome("unbounded-user-value", "user-provided-reason")
	metrics.RecordAnalyticsRetentionCleanup("success", 4, 7, 250*time.Millisecond)
	metrics.RecordAnalyticsRetentionCleanup("failure", 0, 0, time.Second)

	metricsResponse = httptest.NewRecorder()
	metrics.Handler().ServeHTTP(metricsResponse, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	body = metricsResponse.Body.String()
	require.Contains(t, body, `developers_cv_auth_login_attempts_total`)
	require.Contains(t, body, `developers_cv_auth_login_success_total`)
	require.Contains(t, body, `developers_cv_auth_login_failures_total`)
	require.Contains(t, body, `developers_cv_auth_logout_total`)
	require.Contains(t, body, `developers_cv_auth_session_rejections_total`)
	require.Contains(t, body, `developers_cv_auth_transaction_cleanup_total`)
	require.Contains(t, body, `developers_cv_auth_session_cleanup_total`)
	require.Contains(t, body, `developers_cv_auth_cleanup_failures_total`)
	require.Contains(t, body, `developers_cv_auth_login_duration_seconds`)
	require.Contains(t, body, `developers_cv_auth_oidc_exchange_duration_seconds`)
	require.Contains(t, body, `developers_cv_auth_oidc_verify_duration_seconds`)
	require.Contains(t, body, `developers_cv_analytics_dropped_input_total`)
	require.Contains(t, body, `kind="event"`)
	require.Contains(t, body, `reason="invalid_payload_or_type"`)
	require.Contains(t, body, `kind="other",reason="other"`)
	require.NotContains(t, body, `user-provided-reason`)
	require.Contains(t, body, `developers_cv_analytics_batch_outcomes_total`)
	require.Contains(t, body, `outcome="accepted",reason="success"`)
	require.Contains(t, body, `outcome="failure",reason="save_failed"`)
	require.Contains(t, body, `outcome="failure",reason="unknown"`)
	require.NotContains(t, body, `user-provided-reason`)
	require.Contains(t, body, `developers_cv_analytics_retention_runs_total`)
	require.Contains(t, body, `outcome="success"`)
	require.Contains(t, body, `outcome="failure"`)
	require.Contains(t, body, `developers_cv_analytics_retention_deleted_total`)
	require.Contains(t, body, `resource="clicks"`)
	require.Contains(t, body, `resource="views"`)
	require.Contains(t, body, `developers_cv_analytics_retention_duration_seconds`)
}

func TestMetricsMiddlewareUsesUnmatchedRouteWithoutHighCardinalityPath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	metrics, err := NewMetrics(true)
	require.NoError(t, err)

	engine := gin.New()
	engine.Use(metrics.Middleware())

	request := httptest.NewRequest(http.MethodGet, "/unknown/123", nil)
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	metricsResponse := httptest.NewRecorder()
	metrics.Handler().ServeHTTP(metricsResponse, httptest.NewRequest(http.MethodGet, "/metrics", nil))

	require.NotContains(t, strings.TrimSpace(metricsResponse.Body.String()), `route="/unknown/123"`)
}

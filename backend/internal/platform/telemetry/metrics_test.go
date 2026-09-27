package telemetry

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

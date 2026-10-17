package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	apperrors "developers-cv/backend/internal/platform/http/errors"
)

func TestLoggerUsesRouteTemplateInsteadOfRawPath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	core, logs := observer.New(zapcore.InfoLevel)
	logger := zap.New(core)

	engine := gin.New()
	engine.Use(RequestID(), Logger(logger))
	engine.GET("/resumes/:resumeId", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodGet, "/resumes/550e8400-e29b-41d4-a716-446655440001", nil)
	request.Header.Set("X-Request-ID", "request-123")
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, 1, logs.Len())
	entry := logs.All()[0]
	fields := entry.ContextMap()
	require.Equal(t, "/resumes/:resumeId", fields["route"])
	require.NotContains(t, entry.Message, "550e8400-e29b-41d4-a716-446655440001")
}

func TestLoggerRecordsInternalErrorCauseAndRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	core, logs := observer.New(zapcore.DebugLevel)
	logger := zap.New(core)

	engine := gin.New()
	engine.Use(RequestID(), Logger(logger))
	engine.GET("/failure", func(c *gin.Context) {
		apperrors.WriteError(c, apperrors.Internal(errors.New("database exploded")))
	})

	request := httptest.NewRequest(http.MethodGet, "/failure", nil)
	request.Header.Set("X-Request-ID", "request-error-123")
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	require.Equal(t, http.StatusInternalServerError, response.Code)
	entries := logs.FilterMessage("http request failed").All()
	require.Len(t, entries, 1)
	fields := entries[0].ContextMap()
	require.EqualValues(t, http.StatusInternalServerError, fields["status"])
	require.Equal(t, "request-error-123", fields["request_id"])
	require.Contains(t, fields["error"], "database exploded")
	require.NotContains(t, response.Body.String(), "database exploded")
}

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"

	"developers-cv/backend/internal/platform/http/requestid"
)

func TestRecoveryUsesCanonicalErrorContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	engine := gin.New()
	engine.Use(RequestID(), Recovery(zaptest.NewLogger(t)))
	engine.GET("/panic", func(c *gin.Context) {
		panic("secret panic details")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.Equal(t, w.Header().Get(requestid.Header), requestIDFromBody(t, w.Body.Bytes()))
	require.Contains(t, w.Body.String(), `"code":"INTERNAL_ERROR"`)
	require.Contains(t, w.Body.String(), `"message":"Internal server error"`)
	require.NotContains(t, w.Body.String(), "secret panic details")
}

func TestRecoveryAllowsAccessLoggerToObserveRecoveredPanic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	core, observed := observer.New(zapcore.DebugLevel)
	log := zap.New(core)
	engine := gin.New()
	engine.Use(RequestID(), Logger(log), Recovery(log))
	engine.GET("/panic", func(c *gin.Context) { panic("private panic detail") })

	request := httptest.NewRequest(http.MethodGet, "/panic", nil)
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	require.Equal(t, http.StatusInternalServerError, response.Code)
	entries := observed.FilterMessage("http request failed").All()
	require.Len(t, entries, 1)
	require.EqualValues(t, http.StatusInternalServerError, entries[0].ContextMap()["status"])
	fields := entries[0].ContextMap()
	require.NotContains(t, fields["error"], "private panic detail")
}

func TestRecoveryProtectsMiddlewareRegisteredAfterIt(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := zaptest.NewLogger(t)
	engine := gin.New()
	engine.Use(
		RequestID(),
		Logger(log),
		Recovery(log),
		func(c *gin.Context) { panic("middleware panic detail") },
	)
	engine.GET("/test", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	response := httptest.NewRecorder()
	engine.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/test", nil))

	require.Equal(t, http.StatusInternalServerError, response.Code)
	require.Contains(t, response.Body.String(), "INTERNAL_ERROR")
	require.NotContains(t, response.Body.String(), "middleware panic detail")
}

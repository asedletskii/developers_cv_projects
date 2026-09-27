package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestSecurityHeadersWithoutTLS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	engine := gin.New()
	engine.Use(SecurityHeaders(false))
	engine.GET("/test", func(c *gin.Context) {
		c.Status(204)
	})

	w := httptest.NewRecorder()
	engine.ServeHTTP(w, httptest.NewRequest("GET", "/test", nil))

	require.Equal(t, 204, w.Code)
	require.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	require.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	require.Equal(t, "strict-origin-when-cross-origin", w.Header().Get("Referrer-Policy"))
	require.Equal(t, "camera=(), microphone=(), geolocation=()", w.Header().Get("Permissions-Policy"))
	require.Equal(t, "default-src 'none'; frame-ancestors 'none'", w.Header().Get("Content-Security-Policy"))
	require.Equal(t, "no-store", w.Header().Get("Cache-Control"))
	require.Empty(t, w.Header().Get("Strict-Transport-Security"))
}

func TestSecurityHeadersWithTLS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	engine := gin.New()
	engine.Use(SecurityHeaders(true))
	engine.GET("/test", func(c *gin.Context) {
		c.Status(204)
	})

	w := httptest.NewRecorder()
	engine.ServeHTTP(w, httptest.NewRequest("GET", "/test", nil))

	require.Equal(t, "max-age=31536000; includeSubDomains", w.Header().Get("Strict-Transport-Security"))
}

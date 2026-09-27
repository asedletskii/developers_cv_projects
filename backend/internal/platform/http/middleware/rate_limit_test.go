package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/platform/http/requestid"
)

func TestRateLimiterUsesCanonicalErrorContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	engine := gin.New()
	engine.Use(RequestID(), NewRateLimiter(0, 1).Middleware())
	engine.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	first := httptest.NewRecorder()
	engine.ServeHTTP(first, httptest.NewRequest(http.MethodGet, "/test", nil))
	require.Equal(t, http.StatusNoContent, first.Code)

	second := httptest.NewRecorder()
	engine.ServeHTTP(second, httptest.NewRequest(http.MethodGet, "/test", nil))

	require.Equal(t, http.StatusTooManyRequests, second.Code)
	require.Equal(t, "1", second.Header().Get("Retry-After"))
	require.Equal(t, second.Header().Get(requestid.Header), requestIDFromBody(t, second.Body.Bytes()))
	require.Contains(t, second.Body.String(), `"code":"TOO_MANY_REQUESTS"`)
	require.Contains(t, second.Body.String(), `"message":"Too many requests"`)
}

package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/platform/http/requestid"
)

func TestRequestIDPreservesValidHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	engine := gin.New()
	engine.Use(RequestID())
	engine.GET("/test", func(c *gin.Context) {
		c.Status(204)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set(requestid.Header, "request-from-client")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, 204, w.Code)
	require.Equal(t, "request-from-client", w.Header().Get(requestid.Header))
}

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/platform/http/requestid"

	"go.uber.org/zap/zaptest"
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

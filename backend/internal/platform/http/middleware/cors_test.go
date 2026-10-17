package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/platform/config"
)

func TestCORSDoesNotAllowPATCH(t *testing.T) {
	gin.SetMode(gin.TestMode)

	engine := gin.New()
	engine.Use(CORS(config.CORSConfig{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowCredentials: true,
	}))
	engine.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", http.MethodPatch)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.NotContains(t, w.Header().Get("Access-Control-Allow-Methods"), http.MethodPatch)
}

func TestCORSAllowsMockAuthHeaderOnlyWhenMockAuthEnabled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, tc := range []struct {
		name           string
		allowMock      bool
		wantMockHeader bool
	}{
		{name: "real authentication", allowMock: false, wantMockHeader: false},
		{name: "mock authentication", allowMock: true, wantMockHeader: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			engine := gin.New()
			engine.Use(CORS(config.CORSConfig{
				AllowedOrigins:      []string{"http://localhost:5173"},
				AllowMockAuthHeader: tc.allowMock,
			}))
			engine.GET("/test", func(c *gin.Context) { c.Status(http.StatusNoContent) })
			req := httptest.NewRequest(http.MethodOptions, "/test", nil)
			req.Header.Set("Origin", "http://localhost:5173")
			req.Header.Set("Access-Control-Request-Method", http.MethodGet)
			req.Header.Set("Access-Control-Request-Headers", "X-Mock-User-ID")
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			allowed := strings.Contains(strings.ToLower(w.Header().Get("Access-Control-Allow-Headers")), strings.ToLower("X-Mock-User-ID"))
			require.Equal(t, tc.wantMockHeader, allowed)
		})
	}
}

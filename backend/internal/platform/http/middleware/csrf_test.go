package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newCSRFTestEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(CSRFProtection([]string{"https://app.example.com"}))
	engine.POST("/protected", func(c *gin.Context) { c.Status(http.StatusNoContent) })
	engine.GET("/protected", func(c *gin.Context) { c.Status(http.StatusNoContent) })
	return engine
}

func TestCSRFProtectionAllowsSafeMethods(t *testing.T) {
	engine := newCSRFTestEngine()

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
}

func TestCSRFProtectionAllowsSameOriginFetchMetadata(t *testing.T) {
	engine := newCSRFTestEngine()

	req := httptest.NewRequest(http.MethodPost, "/protected", nil)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
}

func TestCSRFProtectionBlocksCrossSiteFetchMetadata(t *testing.T) {
	for _, site := range []string{"cross-site", "none"} {
		t.Run(site, func(t *testing.T) {
			engine := newCSRFTestEngine()

			req := httptest.NewRequest(http.MethodPost, "/protected", nil)
			req.Header.Set("X-Requested-With", "XMLHttpRequest")
			req.Header.Set("Sec-Fetch-Site", site)
			req.Header.Set("Origin", "https://app.example.com")
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)

			require.Equal(t, http.StatusForbidden, w.Code)
		})
	}
}

func TestCSRFProtectionAllowsAllowlistedOrigin(t *testing.T) {
	engine := newCSRFTestEngine()

	req := httptest.NewRequest(http.MethodPost, "/protected", nil)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Origin", "https://app.example.com")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
}

func TestCSRFProtectionAllowsAllowlistedReferrer(t *testing.T) {
	engine := newCSRFTestEngine()

	req := httptest.NewRequest(http.MethodPost, "/protected", nil)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Referer", "https://app.example.com/lk")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
}

func TestCSRFProtectionBlocksForeignOrigin(t *testing.T) {
	engine := newCSRFTestEngine()

	req := httptest.NewRequest(http.MethodPost, "/protected", nil)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Origin", "https://evil.example.com")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusForbidden, w.Code)
}

func TestCSRFProtectionBlocksSameSiteWithoutAllowlistedOrigin(t *testing.T) {
	engine := newCSRFTestEngine()

	req := httptest.NewRequest(http.MethodPost, "/protected", nil)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Origin", "https://evil.example.com")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusForbidden, w.Code)
}

func TestCSRFProtectionBlocksMissingRequestedWith(t *testing.T) {
	engine := newCSRFTestEngine()

	req := httptest.NewRequest(http.MethodPost, "/protected", nil)
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Origin", "https://app.example.com")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusForbidden, w.Code)
	require.Contains(t, w.Body.String(), `"code":"FORBIDDEN"`)
}

func TestCSRFProtectionBlocksInvalidRequestedWith(t *testing.T) {
	engine := newCSRFTestEngine()

	req := httptest.NewRequest(http.MethodPost, "/protected", nil)
	req.Header.Set("X-Requested-With", "fetch")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Origin", "https://app.example.com")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusForbidden, w.Code)
}

func TestCSRFProtectionAllowsRequestedWithWithoutOriginSignals(t *testing.T) {
	engine := newCSRFTestEngine()

	req := httptest.NewRequest(http.MethodPost, "/protected", nil)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
}

func TestCSRFProtectionBlocksMissingAllSignals(t *testing.T) {
	engine := newCSRFTestEngine()

	req := httptest.NewRequest(http.MethodPost, "/protected", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusForbidden, w.Code)
	require.Contains(t, w.Body.String(), `"code":"FORBIDDEN"`)
}

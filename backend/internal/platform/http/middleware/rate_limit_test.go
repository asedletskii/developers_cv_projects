package middleware

import (
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"

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

func TestLimiterKeyNormalizesIPv6To64(t *testing.T) {
	require.Equal(t, "2001:db8:abcd:12::/64", limiterKey("2001:db8:abcd:12:1:2:3:4"))
	require.Equal(t, "2001:db8:abcd:12::/64", limiterKey("2001:db8:abcd:12:ffff:ffff:ffff:ffff"))
	require.Equal(t, "192.0.2.1", limiterKey("192.0.2.1"))
	require.Equal(t, "192.0.2.1", limiterKey("::ffff:192.0.2.1"))
	require.NotEqual(t, limiterKey("2001:db8:abcd:12::1"), limiterKey("2001:db8:abcd:13::1"))
	require.Equal(t, "fe80::/64", limiterKey("fe80::1%eth0"))
	require.Equal(t, limiterKey("fe80::1%eth0"), limiterKey("fe80::2%eth1"))
}

func TestRetryAfterUsesConfiguredRate(t *testing.T) {
	limiter := rate.NewLimiter(rate.Limit(3.0/60.0), 1)
	now := time.Now()
	require.True(t, limiter.AllowN(now, 1))
	require.False(t, limiter.AllowN(now, 1))
	require.Equal(t, "20", retryAfterSeconds(limiter, now, rate.Limit(3.0/60.0)))
}

func TestRetryAfterSaturatesForExtremelySmallRates(t *testing.T) {
	limiter := rate.NewLimiter(rate.Limit(math.SmallestNonzeroFloat64), 1)
	now := time.Now()
	require.True(t, limiter.AllowN(now, 1))
	require.False(t, limiter.AllowN(now, 1))
	require.Equal(t, strconv.Itoa(int(^uint(0)>>1)), retryAfterSeconds(limiter, now, rate.Limit(math.SmallestNonzeroFloat64)))
}

func TestRateLimiterBoundsVisitorStateAndUsesOverflowLimiter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	limiter := NewRateLimiter(0, 1)
	now := time.Now()
	limiter.lastCleanup = now
	for i := 0; i < maxVisitors; i++ {
		key := fmt.Sprintf("visitor-%d", i)
		limiter.visitors[key] = visitor{limiter: rate.NewLimiter(0, 1), lastSeen: now}
	}

	engine := gin.New()
	engine.Use(limiter.Middleware())
	engine.GET("/test", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	firstRequest := httptest.NewRequest(http.MethodGet, "/test", nil)
	firstRequest.RemoteAddr = "203.0.113.10:12345"
	first := httptest.NewRecorder()
	engine.ServeHTTP(first, firstRequest)
	require.Equal(t, http.StatusNoContent, first.Code)

	secondRequest := httptest.NewRequest(http.MethodGet, "/test", nil)
	secondRequest.RemoteAddr = "203.0.113.11:12345"
	second := httptest.NewRecorder()
	engine.ServeHTTP(second, secondRequest)
	require.Equal(t, http.StatusTooManyRequests, second.Code)

	limiter.mu.Lock()
	require.Len(t, limiter.visitors, maxVisitors)
	limiter.mu.Unlock()
}

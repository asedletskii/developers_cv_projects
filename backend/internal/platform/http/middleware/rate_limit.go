package middleware

import (
	"math"
	"net/http"
	"net/netip"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	apperrors "developers-cv/backend/internal/platform/http/errors"

	"golang.org/x/time/rate"
)

const (
	maxVisitors     = 10_000
	visitorTTL      = 10 * time.Minute
	cleanupInterval = 1 * time.Minute
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	mu          sync.Mutex
	visitors    map[string]visitor
	rate        rate.Limit
	burst       int
	lastCleanup time.Time
	overflow    *rate.Limiter
}

func NewRateLimiter(r float64, burst int) *RateLimiter {
	return &RateLimiter{
		visitors:    make(map[string]visitor),
		rate:        rate.Limit(r),
		burst:       burst,
		lastCleanup: time.Now(),
		overflow:    rate.NewLimiter(rate.Limit(r), burst),
	}
}

func (l *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if isRateLimitExempt(c.Request.URL.Path) {
			c.Next()
			return
		}

		now := time.Now()
		key := limiterKey(c.ClientIP())

		l.mu.Lock()
		l.cleanupLocked(now)

		current, ok := l.visitors[key]
		if !ok {
			if len(l.visitors) < maxVisitors {
				current = visitor{limiter: rate.NewLimiter(l.rate, l.burst)}
				ok = true
			} else {
				current = visitor{limiter: l.overflow}
			}
		}
		current.lastSeen = now
		if ok {
			l.visitors[key] = current
		}
		allowed := current.limiter.AllowN(now, 1)
		retryAfter := ""
		if !allowed {
			retryAfter = retryAfterSeconds(current.limiter, now, l.rate)
		}
		l.mu.Unlock()

		if !allowed {
			c.Header("Retry-After", retryAfter)
			apperrors.WriteError(c, apperrors.New(apperrors.CodeTooMany, http.StatusTooManyRequests, "Too many requests"))
			return
		}
		c.Next()
	}
}

func limiterKey(ip string) string {
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return ip
	}
	addr = addr.Unmap()
	if addr.Is6() {
		addr = addr.WithZone("")
		prefix, err := addr.Prefix(64)
		if err == nil {
			return prefix.Masked().String()
		}
	}
	return addr.String()
}

func retryAfterSeconds(limiter *rate.Limiter, now time.Time, configuredRate rate.Limit) string {
	ratePerSecond := float64(configuredRate)
	if ratePerSecond <= 0 || math.IsNaN(ratePerSecond) || math.IsInf(ratePerSecond, 0) {
		return "1"
	}
	missingTokens := 1 - limiter.TokensAt(now)
	if missingTokens <= 0 {
		return "1"
	}
	seconds := math.Ceil(missingTokens / ratePerSecond)
	if seconds < 1 {
		seconds = 1
	}
	maxInt := int(^uint(0) >> 1)
	if math.IsInf(seconds, 1) || math.IsNaN(seconds) || seconds >= float64(maxInt) {
		return strconv.Itoa(maxInt)
	}
	return strconv.Itoa(int(seconds))
}

func isRateLimitExempt(path string) bool {
	return path == "/health" || path == "/ready" || path == "/metrics"
}

func (l *RateLimiter) cleanupLocked(now time.Time) {
	if now.Sub(l.lastCleanup) < cleanupInterval {
		return
	}
	cutoff := now.Add(-visitorTTL)
	for ip, current := range l.visitors {
		if current.lastSeen.Before(cutoff) {
			delete(l.visitors, ip)
		}
	}
	l.lastCleanup = now
}

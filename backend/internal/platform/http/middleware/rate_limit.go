package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
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
}

func NewRateLimiter(r float64, burst int) *RateLimiter {
	return &RateLimiter{
		visitors:    make(map[string]visitor),
		rate:        rate.Limit(r),
		burst:       burst,
		lastCleanup: time.Now(),
	}
}

func (l *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if isRateLimitExempt(c.Request.URL.Path) {
			c.Next()
			return
		}

		now := time.Now()
		ip := c.ClientIP()

		l.mu.Lock()
		l.cleanupLocked(now)

		current, ok := l.visitors[ip]
		if !ok {
			if len(l.visitors) >= maxVisitors {
				l.evictOneLocked()
			}
			current = visitor{limiter: rate.NewLimiter(l.rate, l.burst)}
		}
		current.lastSeen = now
		l.visitors[ip] = current
		allowed := current.limiter.AllowN(now, 1)
		l.mu.Unlock()

		if !allowed {
			c.Header("Retry-After", "1")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": gin.H{
				"code":      "TOO_MANY_REQUESTS",
				"message":   "Too many requests",
				"requestId": GetRequestID(c),
			}})
			return
		}
		c.Next()
	}
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

func (l *RateLimiter) evictOneLocked() {
	for ip := range l.visitors {
		delete(l.visitors, ip)
		return
	}
}

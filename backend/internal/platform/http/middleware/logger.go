package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"developers-cv/backend/internal/platform/http/requestid"
)

func Logger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}
		fields := []zap.Field{
			zap.String("request_id", requestid.Get(c)),
			zap.String("method", c.Request.Method),
			zap.String("route", route),
			zap.Int("status", c.Writer.Status()),
			zap.String("client_ip", c.ClientIP()),
			zap.Int("bytes", c.Writer.Size()),
			zap.Duration("duration", time.Since(start)),
		}
		spanContext := trace.SpanContextFromContext(c.Request.Context())
		if spanContext.IsValid() {
			fields = append(fields, zap.String("trace_id", spanContext.TraceID().String()), zap.String("span_id", spanContext.SpanID().String()))
		}
		for _, requestErr := range c.Errors {
			fields = append(fields, zap.NamedError("error", requestErr.Err))
		}
		if c.Request.URL.Path == "/health" || c.Request.URL.Path == "/ready" || c.Request.URL.Path == "/metrics" {
			return
		}
		switch status := c.Writer.Status(); {
		case status >= http.StatusInternalServerError:
			log.Error("http request failed", fields...)
		case status >= http.StatusBadRequest:
			log.Warn("http request rejected", fields...)
		default:
			log.Info("http request completed", fields...)
		}
	}
}

package requestid

import (
	"context"

	"github.com/gin-gonic/gin"
)

const (
	Header = "X-Request-ID"
	Key    = "request_id"
)

type contextKey struct{}

func Get(c *gin.Context) string {
	if c == nil {
		return ""
	}
	value, exists := c.Get(Key)
	if !exists {
		return ""
	}
	requestID, ok := value.(string)
	if !ok {
		return ""
	}
	return requestID
}

func WithContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, contextKey{}, requestID)
}

func FromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	requestID, ok := ctx.Value(contextKey{}).(string)
	if !ok {
		return ""
	}
	return requestID
}

package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"developers-cv/backend/internal/platform/http/requestid"
)

func isValidRequestID(value string) bool {
	if len(value) < 1 || len(value) > 128 {
		return false
	}
	for _, r := range value {
		if r < 0x21 || r > 0x7e {
			return false
		}
		switch {
		case r >= 'a' && r <= 'z',
			r >= 'A' && r <= 'Z',
			r >= '0' && r <= '9',
			r == '-', r == '_', r == '.', r == ':':
		default:
			return false
		}
	}
	return true
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := strings.TrimSpace(c.GetHeader(requestid.Header))
		if !isValidRequestID(requestID) {
			requestID = uuid.NewString()
		}

		c.Set(requestid.Key, requestID)
		c.Request = c.Request.WithContext(requestid.WithContext(c.Request.Context(), requestID))
		c.Header(requestid.Header, requestID)
		c.Next()
	}
}

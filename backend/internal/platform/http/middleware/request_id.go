package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"developers-cv/backend/internal/platform/http/requestid"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := strings.TrimSpace(c.GetHeader(requestid.Header))
		if requestID == "" || len(requestID) > 128 {
			requestID = uuid.NewString()
		}

		c.Set(requestid.Key, requestID)
		c.Header(requestid.Header, requestID)
		c.Next()
	}
}

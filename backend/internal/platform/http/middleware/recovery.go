package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Error(
					"panic recovered",
					zap.String("panic_type", fmt.Sprintf("%T", recovered)),
					zap.Stack("stack"),
					zap.String("request_id", GetRequestID(c)),
				)
				c.AbortWithStatusJSON(500, gin.H{"error": gin.H{
					"code":      "INTERNAL_ERROR",
					"message":   "Internal server error",
					"requestId": GetRequestID(c),
				}})
			}
		}()
		c.Next()
	}
}

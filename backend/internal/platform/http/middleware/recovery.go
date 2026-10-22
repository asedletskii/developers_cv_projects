package middleware

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"

	apperrors "developers-cv/backend/internal/platform/http/errors"

	"go.uber.org/zap"

	"developers-cv/backend/internal/platform/http/requestid"
)

func Recovery(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Error(
					"panic recovered",
					zap.String("panic_type", fmt.Sprintf("%T", recovered)),
					zap.Stack("stack"),
					zap.String("request_id", requestid.Get(c)),
				)
				apperrors.WriteError(c, apperrors.Internal(errors.New("panic recovered")))
			}
		}()
		c.Next()
	}
}

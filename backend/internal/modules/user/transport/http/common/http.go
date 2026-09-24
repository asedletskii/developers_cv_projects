package common

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "developers-cv/backend/internal/platform/http/errors"
	"developers-cv/backend/internal/platform/http/middleware"
)

func WriteError(c *gin.Context, err error) {
	appErr := apperrors.From(err)
	requestID := middleware.GetRequestID(c)
	c.AbortWithStatusJSON(appErr.Status, gin.H{
		"error": gin.H{
			"code":      appErr.Code,
			"message":   appErr.Message,
			"requestId": requestID,
		},
	})
}

func Unauthorized(c *gin.Context) {
	WriteError(c, apperrors.New(apperrors.CodeUnauthorized, http.StatusUnauthorized, "Authentication required"))
}

func Internal(c *gin.Context, err error) {
	if err == nil {
		err = errors.New("internal server error")
	}
	WriteError(c, apperrors.Internal(err))
}

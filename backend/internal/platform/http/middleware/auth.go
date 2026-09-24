package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var ErrUnauthenticated = errors.New("unauthenticated")

type AuthenticateFunc func(context.Context, *http.Request) (uuid.UUID, error)

const userIDKey = "user_id"

func Authentication(authenticate AuthenticateFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if authenticate == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": gin.H{
				"code":      "INTERNAL_ERROR",
				"message":   "Internal server error",
				"requestId": GetRequestID(c),
			}})
			return
		}

		userID, err := authenticate(c.Request.Context(), c.Request)
		if err != nil {
			if errors.Is(err, ErrUnauthenticated) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": gin.H{
					"code":      "UNAUTHORIZED",
					"message":   "Authentication required",
					"requestId": GetRequestID(c),
				}})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": gin.H{
				"code":      "INTERNAL_ERROR",
				"message":   "Internal server error",
				"requestId": GetRequestID(c),
			}})
			return
		}
		if userID == uuid.Nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": gin.H{
				"code":      "INTERNAL_ERROR",
				"message":   "Internal server error",
				"requestId": GetRequestID(c),
			}})
			return
		}

		c.Set(userIDKey, userID)
		c.Next()
	}
}

func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	value, exists := c.Get(userIDKey)
	if !exists {
		return uuid.Nil, false
	}
	userID, ok := value.(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return uuid.Nil, false
	}
	return userID, true
}

package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "developers-cv/backend/internal/platform/http/errors"

	"github.com/google/uuid"
)

var ErrUnauthenticated = errors.New("unauthenticated")

type AuthenticateFunc func(context.Context, *http.Request) (uuid.UUID, error)

const userIDKey = "user_id"

func Authentication(authenticate AuthenticateFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if authenticate == nil {
			apperrors.WriteError(c, apperrors.Internal(errors.New("authentication function is nil")))
			return
		}

		userID, err := authenticate(c.Request.Context(), c.Request)
		if err != nil {
			if errors.Is(err, ErrUnauthenticated) {
				apperrors.Unauthorized(c)
				return
			}
			apperrors.WriteError(c, apperrors.Internal(err))
			return
		}
		if userID == uuid.Nil {
			apperrors.WriteError(c, apperrors.Internal(errors.New("authentication returned empty user id")))
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

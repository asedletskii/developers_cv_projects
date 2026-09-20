package common

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "developers-cv/backend/internal/platform/http/errors"
	"developers-cv/backend/internal/platform/http/middleware"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func RequireUserID(c *gin.Context) (sharedvalueobjects.ID, bool) {
	userUUID, ok := middleware.GetUserID(c)
	if !ok {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeUnauthorized, http.StatusUnauthorized, "Authentication required"))
		return sharedvalueobjects.ID{}, false
	}
	id, err := sharedvalueobjects.IDFromUUID(userUUID)
	if err != nil {
		apperrors.WriteError(c, apperrors.Internal(err))
		return sharedvalueobjects.ID{}, false
	}
	return id, true
}

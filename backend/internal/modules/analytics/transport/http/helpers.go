package http

import (
	stderrors "errors"
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "developers-cv/backend/internal/platform/http/errors"
	jsonhttp "developers-cv/backend/internal/platform/http/json"
	"developers-cv/backend/internal/platform/http/middleware"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func parseAnalyticsID(c *gin.Context, parameter, message string) (sharedvalueobjects.ID, bool) {
	id, err := sharedvalueobjects.ParseID(c.Param(parameter))
	if err != nil {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, message))
		return sharedvalueobjects.ID{}, false
	}
	return id, true
}

func requireAnalyticsUserID(c *gin.Context) (sharedvalueobjects.ID, bool) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		apperrors.Unauthorized(c)
		return sharedvalueobjects.ID{}, false
	}
	id, err := sharedvalueobjects.IDFromUUID(userID)
	if err != nil {
		apperrors.WriteError(c, apperrors.Internal(err))
		return sharedvalueobjects.ID{}, false
	}
	return id, true
}

func writeAnalyticsBindError(c *gin.Context, err error) {
	if stderrors.Is(err, jsonhttp.ErrRequestBodyTooLarge) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeRequestTooLarge, http.StatusRequestEntityTooLarge, "Request body too large"))
		return
	}
	apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid request body"))
}

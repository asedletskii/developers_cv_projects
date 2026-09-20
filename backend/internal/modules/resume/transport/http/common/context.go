package common

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "developers-cv/backend/internal/platform/http/errors"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func RequireResumeContext(c *gin.Context) (sharedvalueobjects.ID, sharedvalueobjects.ID, bool) {
	userID, ok := RequireUserID(c)
	if !ok {
		return sharedvalueobjects.ID{}, sharedvalueobjects.ID{}, false
	}
	resumeID, ok := ParseResumeID(c)
	if !ok {
		return sharedvalueobjects.ID{}, sharedvalueobjects.ID{}, false
	}
	return userID, resumeID, true
}

func RequirePositiveVersion(c *gin.Context, version int64) bool {
	if version > 0 {
		return true
	}
	apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Version must be positive"))
	return false
}

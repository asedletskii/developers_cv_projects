package common

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "developers-cv/backend/internal/platform/http/errors"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func RequireResumeContext(c *gin.Context) (identity.ID, identity.ID, bool) {
	userID, ok := RequireUserID(c)
	if !ok {
		return identity.ID{}, identity.ID{}, false
	}
	resumeID, ok := ParseResumeID(c)
	if !ok {
		return identity.ID{}, identity.ID{}, false
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

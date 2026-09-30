package common

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	apperrors "developers-cv/backend/internal/platform/http/errors"
	jsonhttp "developers-cv/backend/internal/platform/http/json"
	"developers-cv/backend/internal/platform/http/middleware"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

var ErrRequestBodyTooLarge = jsonhttp.ErrRequestBodyTooLarge

func RequireUserID(c *gin.Context) (useridentity.ID, bool) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeUnauthorized, http.StatusUnauthorized, "Authentication required"))
		return useridentity.ID{}, false
	}
	identity, err := useridentity.IDFromUUID(userID)
	if err != nil {
		apperrors.WriteError(c, apperrors.Internal(err))
		return useridentity.ID{}, false
	}
	return identity, true
}

func ParseResumeID(c *gin.Context) (valueobjects.ID, bool) {
	parsed, err := uuid.Parse(c.Param("resumeId"))
	if err != nil {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid resume id"))
		return valueobjects.ID{}, false
	}
	id, err := valueobjects.IDFromUUID(parsed)
	if err != nil {
		apperrors.WriteError(c, apperrors.Internal(err))
		return valueobjects.ID{}, false
	}
	return id, true
}

func WriteCommandError(c *gin.Context, err error) {
	if errors.Is(err, domain.ErrEntityNotFound) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Resume resource not found"))
		return
	}
	if errors.Is(err, domain.ErrDuplicateEntity) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeConflict, http.StatusConflict, "Resume resource already exists"))
		return
	}
	if IsValidationError(err) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid request"))
		return
	}
	apperrors.WriteError(c, apperrors.Internal(err))
}

func WriteInvalidBody(c *gin.Context) {
	apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid request body"))
}

func WriteBindError(c *gin.Context, err error) {
	if errors.Is(err, ErrRequestBodyTooLarge) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeRequestTooLarge, http.StatusRequestEntityTooLarge, "Request body too large"))
		return
	}
	WriteInvalidBody(c)
}

func WriteQueryError(c *gin.Context, err error) {
	if errors.Is(err, queries.ErrResumeNotFound) {
		WriteResumeNotFound(c)
		return
	}
	apperrors.WriteError(c, err)
}

func WriteResumeNotFound(c *gin.Context) {
	apperrors.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Resume not found"))
}

func WriteContactNotFound(c *gin.Context) {
	apperrors.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Contact not found"))
}

func WriteContactNotHidden(c *gin.Context) {
	apperrors.WriteError(c, apperrors.New(apperrors.CodeConflict, http.StatusConflict, "Contact is not hidden"))
}

func BindJSON(c *gin.Context, target any) error {
	return jsonhttp.BindJSON(c, target)
}

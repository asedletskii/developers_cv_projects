package common

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
	"developers-cv/backend/internal/modules/resume/domain"
	apperrors "developers-cv/backend/internal/platform/http/errors"
	jsonhttp "developers-cv/backend/internal/platform/http/json"
	"developers-cv/backend/internal/platform/http/middleware"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

var ErrRequestBodyTooLarge = jsonhttp.ErrRequestBodyTooLarge

func BindJSON(c *gin.Context, target any) error {
	return jsonhttp.BindJSON(c, target)
}

func RequireUserID(c *gin.Context) (sharedvalueobjects.ID, bool) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeUnauthorized, http.StatusUnauthorized, "Authentication required"))
		return sharedvalueobjects.ID{}, false
	}
	parsedID, err := sharedvalueobjects.IDFromUUID(userID)
	if err != nil {
		apperrors.WriteError(c, apperrors.Internal(err))
		return sharedvalueobjects.ID{}, false
	}
	return parsedID, true
}

func ParseResumeID(c *gin.Context) (sharedvalueobjects.ID, bool) {
	parsed, err := uuid.Parse(c.Param("resumeId"))
	if err != nil {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid resume id"))
		return sharedvalueobjects.ID{}, false
	}
	id, err := sharedvalueobjects.IDFromUUID(parsed)
	if err != nil {
		apperrors.WriteError(c, apperrors.Internal(err))
		return sharedvalueobjects.ID{}, false
	}
	return id, true
}

func WriteCommandError(c *gin.Context, err error) {
	if errors.Is(err, domain.ErrConcurrencyConflict) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodePreconditionFailed, http.StatusPreconditionFailed, "Resume was modified"))
		return
	}
	if errors.Is(err, domain.ErrInvalidResumeStatusTransition) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeConflict, http.StatusConflict, "Invalid resume status transition"))
		return
	}
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
	if errors.Is(err, resumecontracts.ErrResumeNotFound) {
		WriteResumeNotFound(c)
		return
	}
	apperrors.WriteError(c, err)
}

func WriteResumeNotFound(c *gin.Context) {
	apperrors.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Resume not found"))
}

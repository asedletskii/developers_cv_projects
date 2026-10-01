package common

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/user/application/queries"
	"developers-cv/backend/internal/modules/user/domain"
	apperrors "developers-cv/backend/internal/platform/http/errors"
	jsonhttp "developers-cv/backend/internal/platform/http/json"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

var ErrRequestBodyTooLarge = jsonhttp.ErrRequestBodyTooLarge

func BindJSON(c *gin.Context, target any) error {
	return jsonhttp.BindJSON(c, target)
}

func WriteBindError(c *gin.Context, err error) {
	if errors.Is(err, ErrRequestBodyTooLarge) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeRequestTooLarge, http.StatusRequestEntityTooLarge, "Request body too large"))
		return
	}
	WriteInvalidBody(c)
}

func WriteInvalidBody(c *gin.Context) {
	apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid request body"))
}

func WriteCommandError(c *gin.Context, err error) {
	if IsValidationError(err) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid request"))
		return
	}
	apperrors.WriteError(c, apperrors.Internal(err))
}

func WriteQueryError(c *gin.Context, err error) {
	if errors.Is(err, queries.ErrProfileNotFound) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Profile not found"))
		return
	}
	if errors.Is(err, queries.ErrUserNotFound) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "User not found"))
		return
	}
	apperrors.WriteError(c, err)
}

func IsValidationError(err error) bool {
	return errors.Is(err, domain.ErrInvalidUser) ||
		errors.Is(err, domain.ErrInvalidProfile) ||
		errors.Is(err, sharedvalueobjects.ErrInvalidID)
}

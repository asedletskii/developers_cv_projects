package common

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
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

func RequireUserID(c *gin.Context) (useridentity.ID, bool) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		WriteError(c, apperrors.New(apperrors.CodeUnauthorized, http.StatusUnauthorized, "Authentication required"))
		return useridentity.ID{}, false
	}
	identity, err := useridentity.IDFromUUID(userID)
	if err != nil {
		WriteError(c, apperrors.Internal(err))
		return useridentity.ID{}, false
	}
	return identity, true
}

func ParseResumeID(c *gin.Context) (valueobjects.ID, bool) {
	parsed, err := uuid.Parse(c.Param("resumeId"))
	if err != nil {
		WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid resume id"))
		return valueobjects.ID{}, false
	}
	id, err := valueobjects.IDFromUUID(parsed)
	if err != nil {
		WriteError(c, apperrors.Internal(err))
		return valueobjects.ID{}, false
	}
	return id, true
}

func ParseLanguage(c *gin.Context) (string, bool) {
	language := c.DefaultQuery("lang", "ru")
	if language != "ru" && language != "en" {
		WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid language"))
		return "", false
	}
	return language, true
}

func WriteCommandError(c *gin.Context, err error) {
	if errors.Is(err, domain.ErrEntityNotFound) {
		WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Resume resource not found"))
		return
	}
	if errors.Is(err, domain.ErrDuplicateEntity) {
		WriteError(c, apperrors.New(apperrors.CodeConflict, http.StatusConflict, "Resume resource already exists"))
		return
	}
	if IsValidationError(err) {
		WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, err.Error()))
		return
	}
	WriteError(c, apperrors.Internal(err))
}

func WriteInvalidBody(c *gin.Context) {
	WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid request body"))
}

func WriteQueryError(c *gin.Context, err error) {
	if errors.Is(err, queries.ErrResumeNotFound) {
		WriteResumeNotFound(c)
		return
	}
	WriteError(c, err)
}

func WriteResumeNotFound(c *gin.Context) {
	WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Resume not found"))
}

func BindJSON(c *gin.Context, target any) error {
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return errors.New("request body must contain a single JSON value")
		}
		return err
	}
	return nil
}

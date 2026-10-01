package common

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/domain"
	apperrors "developers-cv/backend/internal/platform/http/errors"
	jsonhttp "developers-cv/backend/internal/platform/http/json"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

var ErrRequestBodyTooLarge = jsonhttp.ErrRequestBodyTooLarge

func BindJSON(c *gin.Context, target any) error { return jsonhttp.BindJSON(c, target) }

func ParseResumeID(c *gin.Context) (sharedvalueobjects.ID, bool) {
	return parseResumeID(c.Query("resumeId"), c)
}

func ParseResumeIDValue(c *gin.Context, value string) (sharedvalueobjects.ID, bool) {
	return parseResumeID(value, c)
}

func parseResumeID(value string, c *gin.Context) (sharedvalueobjects.ID, bool) {
	parsed, err := uuid.Parse(value)
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

func ParseLinkID(c *gin.Context) (sharedvalueobjects.ID, bool) {
	id, err := sharedvalueobjects.ParseID(c.Param("linkId"))
	if err != nil {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid link id"))
		return sharedvalueobjects.ID{}, false
	}
	return id, true
}

func RequirePositiveVersion(c *gin.Context, version int64) bool {
	if version > 0 {
		return true
	}
	apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Version must be positive"))
	return false
}

func WriteBindError(c *gin.Context, err error) {
	if errors.Is(err, ErrRequestBodyTooLarge) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeRequestTooLarge, http.StatusRequestEntityTooLarge, "Request body too large"))
		return
	}
	apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid request body"))
}

func WriteCommandError(c *gin.Context, err error) {
	if errors.Is(err, resumequeries.ErrResumeNotFound) || errors.Is(err, resumedomain.ErrEntityNotFound) {
		WriteResumeLinkNotFound(c)
		return
	}
	if errors.Is(err, domain.ErrConcurrencyConflict) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodePreconditionFailed, http.StatusPreconditionFailed, "Resume link was modified"))
		return
	}
	if errors.Is(err, domain.ErrDuplicateResumeLink) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeConflict, http.StatusConflict, "Resume link already exists"))
		return
	}
	if errors.Is(err, domain.ErrResumeLinkLimitExceeded) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeConflict, http.StatusConflict, "Resume link limit exceeded"))
		return
	}
	if errors.Is(err, domain.ErrEntityNotFound) {
		WriteResumeLinkNotFound(c)
		return
	}
	if errors.Is(err, domain.ErrInvalidResumeLink) || errors.Is(err, domain.ErrInvalidResumeLinkName) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid resume link"))
		return
	}
	apperrors.WriteError(c, err)
}

func WriteCreateError(c *gin.Context, err error) {
	if errors.Is(err, resumequeries.ErrResumeNotFound) || errors.Is(err, resumedomain.ErrEntityNotFound) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Resume not found"))
		return
	}
	if errors.Is(err, domain.ErrDuplicateResumeLink) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeConflict, http.StatusConflict, "Resume link already exists"))
		return
	}
	if errors.Is(err, domain.ErrResumeLinkLimitExceeded) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeConflict, http.StatusConflict, "Resume link limit exceeded"))
		return
	}
	if errors.Is(err, domain.ErrInvalidResumeLink) || errors.Is(err, domain.ErrInvalidResumeLinkName) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid resume link"))
		return
	}
	apperrors.WriteError(c, err)
}

func WriteResumeLinkNotFound(c *gin.Context) {
	apperrors.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Resume link not found"))
}

func WriteListError(c *gin.Context, err error) {
	if errors.Is(err, resumequeries.ErrResumeNotFound) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Resume not found"))
		return
	}
	if errors.Is(err, queries.ErrResumeLinkNotFound) || errors.Is(err, domain.ErrEntityNotFound) {
		WriteResumeLinkNotFound(c)
		return
	}
	apperrors.WriteError(c, err)
}

func WriteStatisticsError(c *gin.Context, err error) {
	if errors.Is(err, resumequeries.ErrResumeNotFound) || errors.Is(err, queries.ErrResumeLinkNotFound) || errors.Is(err, domain.ErrEntityNotFound) {

		WriteResumeLinkNotFound(c)
		return
	}
	apperrors.WriteError(c, err)
}

func WritePublicError(c *gin.Context, err error) {
	if errors.Is(err, domain.ErrResumeLinkInactive) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeGone, http.StatusGone, "Resume link is no longer active"))
		return
	}
	if errors.Is(err, queries.ErrResumeLinkNotFound) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Resume not found"))
		return
	}
	if errors.Is(err, resumequeries.ErrResumeNotFound) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Resume not found"))
		return
	}
	apperrors.WriteError(c, err)
}

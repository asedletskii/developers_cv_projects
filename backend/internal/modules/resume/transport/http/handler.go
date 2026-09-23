package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	apperrors "developers-cv/backend/internal/platform/http/errors"
	"developers-cv/backend/internal/platform/http/middleware"
)

type Handler struct {
	getResume    *queries.GetResumeHandler
	createResume *commands.CreateResumeHandler
}

func NewHandler(
	getResume *queries.GetResumeHandler,
	createResume ...*commands.CreateResumeHandler,
) *Handler {
	var create *commands.CreateResumeHandler
	if len(createResume) > 0 {
		create = createResume[0]
	}
	return &Handler{getResume: getResume, createResume: create}
}

func (h *Handler) GetResume(c *gin.Context) {
	id, err := uuid.Parse(c.Param("cvId"))
	if err != nil {
		writeError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid CV id"))
		return
	}

	lang := queries.Language(c.DefaultQuery("lang", string(queries.LanguageRU)))
	if lang != queries.LanguageRU && lang != queries.LanguageEN {
		writeError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid language"))
		return
	}

	resume, err := h.getResume.Handle(c.Request.Context(), queries.GetResumeQuery{ID: id, Lang: lang})
	if err != nil {
		if errors.Is(err, queries.ErrResumeNotFound) {
			writeError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "CV not found"))
			return
		}
		writeError(c, apperrors.Internal(err))
		return
	}

	c.JSON(http.StatusOK, resume)
}

func writeError(c *gin.Context, err error) {
	appErr := apperrors.From(err)
	requestID := middleware.GetRequestID(c)
	c.AbortWithStatusJSON(
		appErr.Status,
		gin.H{
			"error": gin.H{
				"code":      appErr.Code,
				"message":   appErr.Message,
				"requestId": requestID,
			},
		},
	)
}

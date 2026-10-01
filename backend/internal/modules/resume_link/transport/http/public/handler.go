package public

import (
	"context"
	"errors"

	"net/http"

	"github.com/gin-gonic/gin"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/domain"
	apperrors "developers-cv/backend/internal/platform/http/errors"
	"developers-cv/backend/internal/shared/domain/valueobjects"
)

type Handler struct {
	query publicResumeLinkQuery
}

type publicResumeLinkQuery interface {
	Handle(context.Context, queries.PublicResumeLinkQuery) (resumequeries.ResumeView, error)
}

func NewHandler(query publicResumeLinkQuery) *Handler {
	return &Handler{query: query}
}

func (h *Handler) Handle(c *gin.Context) {
	if h == nil || h.query == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("public resume link query handler is not configured")))
		return
	}

	linkID, err := valueobjects.ParseID(c.Param("linkId"))
	if err != nil {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid link id"))
		return
	}

	view, err := h.query.Handle(c.Request.Context(), queries.PublicResumeLinkQuery{LinkID: linkID})
	if err != nil {
		if errors.Is(err, domain.ErrResumeLinkInactive) {
			apperrors.WriteError(c, apperrors.New(apperrors.CodeGone, http.StatusGone, "Resume link is no longer active"))
			return
		}
		if errors.Is(err, domain.ErrEntityNotFound) || errors.Is(err, resumequeries.ErrResumeNotFound) {
			apperrors.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Resume not found"))
			return
		}
		apperrors.WriteError(c, apperrors.Internal(err))
		return
	}

	c.JSON(http.StatusOK, view)
}

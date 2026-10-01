package public

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/transport/http/common"
	apperrors "developers-cv/backend/internal/platform/http/errors"
)

type Handler struct{ query publicResumeLinkQuery }

type publicResumeLinkQuery interface {
	Handle(context.Context, queries.PublicResumeLinkQuery) (resumequeries.ResumeView, error)
}

func NewHandler(query publicResumeLinkQuery) *Handler { return &Handler{query: query} }

func (h *Handler) Handle(c *gin.Context) {
	if h == nil || h.query == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("public resume link query handler is not configured")))
		return
	}
	linkID, ok := common.ParseLinkID(c)
	if !ok {
		return
	}
	view, err := h.query.Handle(c.Request.Context(), queries.PublicResumeLinkQuery{LinkID: linkID})
	if err != nil {
		common.WritePublicError(c, err)
		return
	}
	c.JSON(http.StatusOK, view)
}

package public

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
)

type Handler struct {
	query *queries.PublicResumeQueryHandler
}

func NewHandler(query *queries.PublicResumeQueryHandler) *Handler {
	return &Handler{query: query}
}

func (h *Handler) Handle(c *gin.Context) {
	if h == nil || h.query == nil {
		common.WriteError(c, errors.New("public resume handler dependencies are nil"))
		return
	}
	resumeID, ok := common.ParseResumeID(c)
	if !ok {
		return
	}
	view, err := h.query.Handle(c.Request.Context(), queries.PublicResumeQuery{ID: resumeID.UUID()})
	if err != nil {
		common.WriteQueryError(c, err)
		return
	}
	c.JSON(http.StatusOK, view)
}

package list

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
)

type Handler struct {
	query *queries.ListResumesHandler
}

func NewHandler(query *queries.ListResumesHandler) *Handler {
	return &Handler{query: query}
}

func (h *Handler) Handle(c *gin.Context) {
	if h == nil || h.query == nil {
		common.WriteError(c, errors.New("list resumes handler dependencies are nil"))
		return
	}
	userID, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	items, err := h.query.Handle(c.Request.Context(), queries.ListResumesQuery{UserID: userID})
	if err != nil {
		common.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

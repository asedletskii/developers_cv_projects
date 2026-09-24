package get

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
)

type Handler struct {
	query *queries.GetResumeHandler
}

func NewHandler(query *queries.GetResumeHandler) *Handler {
	return &Handler{query: query}
}

func (h *Handler) Handle(c *gin.Context, public bool) {
	if h == nil || h.query == nil {
		common.WriteError(c, errors.New("get resume handler dependencies are nil"))
		return
	}

	resumeID, ok := common.ParseResumeID(c)
	if !ok {
		return
	}
	language, ok := common.ParseLanguage(c)
	if !ok {
		return
	}

	query := queries.GetResumeQuery{
		ID:   resumeID.UUID(),
		Lang: queries.Language(language),
	}
	if !public {
		userID, ok := common.RequireUserID(c)
		if !ok {
			return
		}
		query.UserID = &userID
	}

	resume, err := h.query.Handle(c.Request.Context(), query)
	if err != nil {
		common.WriteQueryError(c, err)
		return
	}
	c.JSON(http.StatusOK, resume)
}

func (h *Handler) Private(c *gin.Context) {
	h.Handle(c, false)
}

func (h *Handler) Public(c *gin.Context) {
	h.Handle(c, true)
}

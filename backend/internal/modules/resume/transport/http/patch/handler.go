package patch

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
)

type Request struct {
	Title *string `json:"title"`
}

type Handler struct {
	patch *commands.PatchResumeHandler
}

func NewHandler(patchHandler *commands.PatchResumeHandler) *Handler {
	return &Handler{patch: patchHandler}
}

func (h *Handler) Handle(c *gin.Context) {
	if h == nil || h.patch == nil {
		common.WriteError(c, errors.New("patch resume handler dependencies are nil"))
		return
	}
	userID, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	resumeID, ok := common.ParseResumeID(c)
	if !ok {
		return
	}

	var body Request
	if err := common.BindJSON(c, &body); err != nil || body.Title == nil {
		common.WriteInvalidBody(c)
		return
	}
	result, err := h.patch.Handle(c.Request.Context(), commands.PatchResumeCommand{
		ResumeID: resumeID,
		UserID:   userID,
		Title:    *body.Title,
	})
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.JSON(http.StatusOK, queries.ResumeListItem{ID: result.ResumeID.String(), Title: result.Title})
}

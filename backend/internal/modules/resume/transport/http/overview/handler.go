package overview

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
)

type Handler struct {
	replace *commands.ReplaceOverviewHandler
}

func NewHandler(replace *commands.ReplaceOverviewHandler) *Handler { return &Handler{replace: replace} }

func (h *Handler) Replace(c *gin.Context) {
	uid, id, ok := common.RequireResumeContext(c)
	if !ok {
		return
	}
	var body overviewRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteBindError(c, err)
		return
	}
	if !common.RequirePositiveVersion(c, body.Version) {
		return
	}
	cmd, err := body.ToCommand(uid, id)
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	if err = h.replace.Handle(c.Request.Context(), cmd); err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

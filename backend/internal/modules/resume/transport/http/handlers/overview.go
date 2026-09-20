package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	"developers-cv/backend/internal/modules/resume/transport/http/requests"
)

type ResumeOverviewHandler struct {
	replace *commands.ReplaceOverviewHandler
}

func NewResumeOverviewHandler(replace *commands.ReplaceOverviewHandler) *ResumeOverviewHandler {
	return &ResumeOverviewHandler{replace: replace}
}

func (h *ResumeOverviewHandler) Replace(c *gin.Context) {
	uid, id, ok := common.RequireResumeContext(c)
	if !ok {
		return
	}
	var body requests.ResumeOverviewRequest
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

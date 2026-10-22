package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	"developers-cv/backend/internal/modules/resume/transport/http/requests"
)

type SkillsHandler struct {
	replace *commands.ReplaceSkillsHandler
}

func NewSkillsHandler(replace *commands.ReplaceSkillsHandler) *SkillsHandler {
	return &SkillsHandler{replace: replace}
}

func (h *SkillsHandler) Replace(c *gin.Context) {
	uid, rid, ok := common.RequireResumeContext(c)
	if !ok {
		return
	}
	var body requests.SkillsRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteBindError(c, err)
		return
	}
	if !common.RequirePositiveVersion(c, body.Version) {
		return
	}
	input, err := body.ToInput()
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	if err := h.replace.Handle(c.Request.Context(), commands.ReplaceSkillsCommand{
		ResumeID: rid,
		UserID:   uid,
		Version:  body.Version,
		Input:    input,
	}); err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

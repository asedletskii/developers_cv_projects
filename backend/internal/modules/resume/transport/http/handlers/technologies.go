package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	"developers-cv/backend/internal/modules/resume/transport/http/requests"
)

type ResumeTechnologiesHandler struct {
	replace *commands.ReplaceTechnologiesHandler
}

func NewResumeTechnologiesHandler(replace *commands.ReplaceTechnologiesHandler) *ResumeTechnologiesHandler {
	return &ResumeTechnologiesHandler{replace: replace}
}

func (h *ResumeTechnologiesHandler) Replace(c *gin.Context) {
	uid, rid, ok := common.RequireResumeContext(c)
	if !ok {
		return
	}
	var body requests.ResumeTechnologiesRequest
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
	if err := h.replace.Handle(c.Request.Context(), commands.ReplaceTechnologiesCommand{
		ResumeID: rid, UserID: uid, Version: body.Version, Input: input,
	}); err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	"developers-cv/backend/internal/modules/resume/transport/http/requests"
)

type ProfileHandler struct {
	replace *commands.ReplaceProfileHandler
}

func NewProfileHandler(replace *commands.ReplaceProfileHandler) *ProfileHandler {
	return &ProfileHandler{replace: replace}
}

func (h *ProfileHandler) Replace(c *gin.Context) {
	uid, rid, ok := common.RequireResumeContext(c)
	if !ok {
		return
	}
	var body requests.ProfileRequest
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
	err = h.replace.Handle(c.Request.Context(), commands.ReplaceProfileCommand{ResumeID: rid, UserID: uid, Version: body.Version, Input: input})
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

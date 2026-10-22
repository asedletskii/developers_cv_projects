package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	"developers-cv/backend/internal/modules/resume/transport/http/requests"
)

type ContactsHandler struct {
	replace *commands.ReplaceContactsHandler
}

func NewContactsHandler(replace *commands.ReplaceContactsHandler) *ContactsHandler {
	return &ContactsHandler{replace: replace}
}

func (h *ContactsHandler) Replace(c *gin.Context) {
	uid, resumeID, ok := common.RequireResumeContext(c)
	if !ok {
		return
	}

	var body requests.ContactsRequest
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

	err = h.replace.Handle(c.Request.Context(), commands.ReplaceContactsCommand{
		ResumeID: resumeID,
		UserID:   uid,
		Version:  body.Version,
		Input:    input,
	})
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

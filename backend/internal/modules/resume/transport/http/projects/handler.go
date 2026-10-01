package projects

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
)

type Handler struct {
	replace *commands.ReplaceProjectsHandler
}

func NewHandler(replace *commands.ReplaceProjectsHandler) *Handler {
	return &Handler{replace: replace}
}

func (h *Handler) Replace(c *gin.Context) {
	uid, rid, ok := common.RequireResumeContext(c)
	if !ok {
		return
	}
	var body replaceProjectsRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteBindError(c, err)
		return
	}
	if !common.RequirePositiveVersion(c, body.Version) {
		return
	}

	items := make([]commands.ProjectInput, len(body.Items))
	for i, item := range body.Items {
		input, err := projectInput(item)
		if err != nil {
			common.WriteCommandError(c, err)
			return
		}
		items[i] = input
	}

	err := h.replace.Handle(c.Request.Context(), commands.ReplaceProjectsCommand{
		ResumeID: rid,
		UserID:   uid,
		Version:  body.Version,
		Input:    commands.ProjectsInput{Items: items},
	})
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

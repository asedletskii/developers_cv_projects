package projects

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type Handler struct {
	replace *commands.ReplaceProjectsHandler
}

func NewHandler(replace *commands.ReplaceProjectsHandler) *Handler {
	return &Handler{replace: replace}
}

func (h *Handler) Replace(c *gin.Context) {
	uid, rid, ok := ids(c)
	if !ok {
		return
	}
	var body replaceProjectsRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteBindError(c, err)
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
		if errors.Is(err, domain.ErrConcurrencyConflict) {
			c.Status(http.StatusPreconditionFailed)
			return
		}
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func ids(c *gin.Context) (useridentity.ID, valueobjects.ID, bool) {
	uid, ok := common.RequireUserID(c)
	if !ok {
		return uid, valueobjects.ID{}, false
	}
	rid, ok := common.ParseResumeID(c)
	return uid, rid, ok
}

package technologies

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type Handler struct {
	replace *commands.ReplaceTechnologiesHandler
}

func NewHandler(replace *commands.ReplaceTechnologiesHandler) *Handler {
	return &Handler{replace: replace}
}

type technologiesRequest struct {
	Version int64               `json:"version"`
	Items   []technologyRequest `json:"items"`
}

type technologyRequest struct {
	ID       string  `json:"id,omitempty"`
	Name     string  `json:"name"`
	Grade    *string `json:"grade,omitempty"`
	Position int     `json:"position"`
}

func (r technologiesRequest) input() (commands.TechnologiesInput, error) {
	items := make([]commands.TechnologyInput, 0, len(r.Items))
	for _, item := range r.Items {
		var id identity.ID
		if item.ID != "" {
			parsed, err := identity.ParseID(item.ID)
			if err != nil {
				return commands.TechnologiesInput{}, err
			}
			id = parsed
		}
		var grade *domain.TechnologyGrade
		if item.Grade != nil {
			value := domain.TechnologyGrade(*item.Grade)
			grade = &value
		}
		items = append(items, commands.TechnologyInput{
			ID: id, Name: item.Name, Grade: grade, Position: item.Position,
		})
	}
	return commands.TechnologiesInput{Entries: items}, nil
}

func (h *Handler) Replace(c *gin.Context) {
	uid, rid, ok := common.RequireResumeContext(c)
	if !ok {
		return
	}
	var body technologiesRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteBindError(c, err)
		return
	}
	if !common.RequirePositiveVersion(c, body.Version) {
		return
	}
	input, err := body.input()
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

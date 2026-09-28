package technologies

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
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
		var id valueobjects.ID
		if item.ID != "" {
			parsed, err := valueobjects.ParseID(item.ID)
			if err != nil {
				return commands.TechnologiesInput{}, err
			}
			id = parsed
		}
		var grade *enums.TechnologyGrade
		if item.Grade != nil {
			value := enums.TechnologyGrade(*item.Grade)
			grade = &value
		}
		items = append(items, commands.TechnologyInput{
			ID: id, Name: item.Name, Grade: grade, Position: item.Position,
		})
	}
	return commands.TechnologiesInput{Entries: items}, nil
}

func (h *Handler) Replace(c *gin.Context) {
	uid, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	rid, ok := common.ParseResumeID(c)
	if !ok {
		return
	}
	var body technologiesRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteBindError(c, err)
		return
	}
	if body.Version < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "version must be positive"})
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
		if errors.Is(err, domain.ErrConcurrencyConflict) {
			c.Status(http.StatusPreconditionFailed)
			return
		}
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

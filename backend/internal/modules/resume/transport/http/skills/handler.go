package skills

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
	replace *commands.ReplaceSkillsHandler
}

func NewHandler(replace *commands.ReplaceSkillsHandler) *Handler {
	return &Handler{replace: replace}
}

type skillsRequest struct {
	Version    int64    `json:"version"`
	ID         string   `json:"id,omitempty"`
	HardSkills []string `json:"hardSkills"`
	SoftSkills []string `json:"softSkills"`
}

func (r skillsRequest) input() (commands.SkillsInput, error) {
	var id valueobjects.ID
	if r.ID != "" {
		parsed, err := valueobjects.ParseID(r.ID)
		if err != nil {
			return commands.SkillsInput{}, err
		}
		id = parsed
	}
	hard := make([]enums.HardSkill, len(r.HardSkills))
	for i, value := range r.HardSkills {
		hard[i] = enums.HardSkill(value)
	}
	soft := make([]enums.SoftSkill, len(r.SoftSkills))
	for i, value := range r.SoftSkills {
		soft[i] = enums.SoftSkill(value)
	}
	return commands.SkillsInput{ID: id, Hard: hard, Soft: soft}, nil
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
	var body skillsRequest
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
	if err := h.replace.Handle(c.Request.Context(), commands.ReplaceSkillsCommand{
		ResumeID: rid,
		UserID:   uid,
		Version:  body.Version,
		Input:    input,
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

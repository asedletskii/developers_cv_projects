package skills

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
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
	var id identity.ID
	if r.ID != "" {
		parsed, err := identity.ParseID(r.ID)
		if err != nil {
			return commands.SkillsInput{}, err
		}
		id = parsed
	}
	hard := make([]domain.HardSkill, len(r.HardSkills))
	for i, value := range r.HardSkills {
		hard[i] = domain.HardSkill(value)
	}
	soft := make([]domain.SoftSkill, len(r.SoftSkills))
	for i, value := range r.SoftSkills {
		soft[i] = domain.SoftSkill(value)
	}
	return commands.SkillsInput{ID: id, Hard: hard, Soft: soft}, nil
}

func (h *Handler) Replace(c *gin.Context) {
	uid, rid, ok := common.RequireResumeContext(c)
	if !ok {
		return
	}
	var body skillsRequest
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

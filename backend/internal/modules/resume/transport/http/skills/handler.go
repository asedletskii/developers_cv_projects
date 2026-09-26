package skills

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type Handler struct {
	create *commands.CreateSkillsHandler
	get    *queries.GetSkillsQueryHandler
	update *commands.UpdateSkillsHandler
}

func NewHandler(create *commands.CreateSkillsHandler, get *queries.GetSkillsQueryHandler, update *commands.UpdateSkillsHandler) *Handler {
	return &Handler{create: create, get: get, update: update}
}

type skillsRequest struct {
	HardSkills []string `json:"hardSkills"`
	SoftSkills []string `json:"softSkills"`
}

func (h *Handler) Create(c *gin.Context) {
	uid, rid, ok := ids(c)
	if !ok {
		return
	}
	var b skillsRequest
	if common.BindJSON(c, &b) != nil {
		common.WriteInvalidBody(c)
		return
	}
	id := valueobjects.NewID()
	_, e := h.create.Handle(c.Request.Context(), commands.CreateSkillsCommand{ID: id, ResumeID: rid, UserID: uid, Hard: toHardSkills(b.HardSkills), Soft: toSoftSkills(b.SoftSkills)})
	if e != nil {
		common.WriteCommandError(c, e)
		return
	}
	c.JSON(http.StatusCreated, common.CreatedIDResponse{ID: id.String()})
}
func (h *Handler) Get(c *gin.Context) {
	uid, rid, ok := ids(c)
	if !ok {
		return
	}
	a, e := h.get.Handle(c.Request.Context(), queries.GetSkillsQuery{ResumeID: rid, UserID: uid})
	if e != nil {
		common.WriteQueryError(c, e)
		return
	}
	c.JSON(http.StatusOK, a)
}
func (h *Handler) Update(c *gin.Context) {
	uid, rid, ok := ids(c)
	if !ok {
		return
	}
	var b skillsRequest
	if common.BindJSON(c, &b) != nil {
		common.WriteInvalidBody(c)
		return
	}
	if e := h.update.Handle(c.Request.Context(), commands.UpdateSkillsCommand{ResumeID: rid, UserID: uid, Hard: toHardSkills(b.HardSkills), Soft: toSoftSkills(b.SoftSkills)}); e != nil {
		common.WriteCommandError(c, e)
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

func toHardSkills(values []string) []enums.HardSkill {
	out := make([]enums.HardSkill, len(values))
	for i, value := range values {
		out[i] = enums.HardSkill(value)
	}
	return out
}

func toSoftSkills(values []string) []enums.SoftSkill {
	out := make([]enums.SoftSkill, len(values))
	for i, value := range values {
		out[i] = enums.SoftSkill(value)
	}
	return out
}

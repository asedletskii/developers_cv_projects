package skills

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
	apperrors "developers-cv/backend/internal/platform/http/errors"
)

type Handler struct {
	create *commands.CreateSkillsHandler
	get    *commands.GetSkillsHandler
	update *commands.UpdateSkillsHandler
}

func NewHandler(create *commands.CreateSkillsHandler, get *commands.GetSkillsHandler, update *commands.UpdateSkillsHandler) *Handler {
	return &Handler{create: create, get: get, update: update}
}

type SkillsRequest struct {
	HardSkills []string `json:"hardSkills"`
	SoftSkills []string `json:"softSkills"`
}
type SkillsResponse struct {
	ID         string   `json:"id"`
	ResumeID   string   `json:"resumeId"`
	HardSkills []string `json:"hardSkills"`
	SoftSkills []string `json:"softSkills"`
}

func (h *Handler) Create(c *gin.Context) {
	uid, rid, ok := ids(c)
	if !ok {
		return
	}
	var b SkillsRequest
	if common.BindJSON(c, &b) != nil {
		common.WriteInvalidBody(c)
		return
	}
	id, e := h.create.Handle(c.Request.Context(), commands.CreateSkillsCommand{ID: valueobjects.NewID(), ResumeID: rid, UserID: uid, Hard: b.HardSkills, Soft: b.SoftSkills})
	if e != nil {
		write(c, e)
		return
	}
	a, e := h.get.Handle(c.Request.Context(), commands.GetSkillsQuery{ResumeID: rid, UserID: uid})
	if e != nil {
		write(c, e)
		return
	}
	c.JSON(http.StatusCreated, skillsView(id, rid, a))
}
func (h *Handler) Get(c *gin.Context) {
	uid, rid, ok := ids(c)
	if !ok {
		return
	}
	a, e := h.get.Handle(c.Request.Context(), commands.GetSkillsQuery{ResumeID: rid, UserID: uid})
	if e != nil {
		write(c, e)
		return
	}
	c.JSON(http.StatusOK, skillsView(a.ID(), rid, a))
}
func (h *Handler) Update(c *gin.Context) {
	uid, rid, ok := ids(c)
	if !ok {
		return
	}
	var b SkillsRequest
	if common.BindJSON(c, &b) != nil {
		common.WriteInvalidBody(c)
		return
	}
	if e := h.update.Handle(c.Request.Context(), commands.UpdateSkillsCommand{ResumeID: rid, UserID: uid, Hard: b.HardSkills, Soft: b.SoftSkills}); e != nil {
		write(c, e)
		return
	}
	a, e := h.get.Handle(c.Request.Context(), commands.GetSkillsQuery{ResumeID: rid, UserID: uid})
	if e != nil {
		write(c, e)
		return
	}
	c.JSON(http.StatusOK, skillsView(a.ID(), rid, a))
}
func skillsView(id, rid valueobjects.ID, a *domain.Skills) SkillsResponse {
	h := a.HardSkills()
	s := a.SoftSkills()
	hh := make([]string, len(h))
	ss := make([]string, len(s))
	for i, v := range h {
		hh[i] = string(v)
	}
	for i, v := range s {
		ss[i] = string(v)
	}
	return SkillsResponse{ID: id.String(), ResumeID: rid.String(), HardSkills: hh, SoftSkills: ss}
}

func ids(c *gin.Context) (useridentity.ID, valueobjects.ID, bool) {
	uid, ok := common.RequireUserID(c)
	if !ok {
		return useridentity.ID{}, valueobjects.ID{}, false
	}
	rid, ok := common.ParseResumeID(c)
	return uid, rid, ok
}

func write(c *gin.Context, e error) {
	if errors.Is(e, domain.ErrEntityNotFound) {
		common.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Resume resource not found"))
		return
	}
	if errors.Is(e, domain.ErrDuplicateEntity) {
		common.WriteError(c, apperrors.New(apperrors.CodeConflict, http.StatusConflict, "Resume resource already exists"))
		return
	}
	common.WriteCommandError(c, e)
}

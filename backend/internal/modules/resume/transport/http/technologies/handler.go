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
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
	apperrors "developers-cv/backend/internal/platform/http/errors"
)

type Handler struct {
	create *commands.CreateTechnologyHandler
	list   *commands.ListTechnologiesHandler
	get    *commands.GetTechnologyHandler
	update *commands.UpdateTechnologyHandler
	delete *commands.DeleteTechnologyHandler
}

func NewHandler(create *commands.CreateTechnologyHandler, list *commands.ListTechnologiesHandler, get *commands.GetTechnologyHandler, update *commands.UpdateTechnologyHandler, delete *commands.DeleteTechnologyHandler) *Handler {
	return &Handler{create: create, list: list, get: get, update: update, delete: delete}
}

type TechnologyRequest struct {
	Name     string  `json:"name"`
	Grade    *string `json:"grade,omitempty"`
	Position int     `json:"position"`
}
type TechnologyResponse struct {
	ID       string  `json:"id"`
	ResumeID string  `json:"resumeId"`
	Name     string  `json:"name"`
	Grade    *string `json:"grade,omitempty"`
	Position int     `json:"position"`
}

func (h *Handler) Create(c *gin.Context) {
	uid, rid, ok := ids(c)
	if !ok {
		return
	}
	var b TechnologyRequest
	if common.BindJSON(c, &b) != nil {
		common.WriteInvalidBody(c)
		return
	}
	id, e := h.create.Handle(c.Request.Context(), commands.CreateTechnologyCommand{ID: valueobjects.NewID(), ResumeID: rid, UserID: uid, Input: techInput(b)})
	if e != nil {
		write(c, e)
		return
	}
	x, e := h.get.Handle(c.Request.Context(), commands.GetTechnologyQuery{ResumeID: rid, TechnologyID: id, UserID: uid})
	if e != nil {
		write(c, e)
		return
	}
	c.JSON(http.StatusCreated, techView(rid, x))
}
func (h *Handler) List(c *gin.Context) {
	uid, rid, ok := ids(c)
	if !ok {
		return
	}
	xs, e := h.list.Handle(c.Request.Context(), commands.ListTechnologiesQuery{ResumeID: rid, UserID: uid})
	if e != nil {
		write(c, e)
		return
	}
	out := make([]TechnologyResponse, len(xs))
	for i, x := range xs {
		out[i] = techView(rid, x)
	}
	c.JSON(http.StatusOK, out)
}
func (h *Handler) Get(c *gin.Context) {
	uid, rid, id, ok := itemIDs(c, "technologyId")
	if !ok {
		return
	}
	x, e := h.get.Handle(c.Request.Context(), commands.GetTechnologyQuery{ResumeID: rid, TechnologyID: id, UserID: uid})
	if e != nil {
		write(c, e)
		return
	}
	c.JSON(http.StatusOK, techView(rid, x))
}
func (h *Handler) Update(c *gin.Context) {
	uid, rid, id, ok := itemIDs(c, "technologyId")
	if !ok {
		return
	}
	var b TechnologyRequest
	if common.BindJSON(c, &b) != nil {
		common.WriteInvalidBody(c)
		return
	}
	e := h.update.Handle(c.Request.Context(), commands.UpdateTechnologyCommand{ResumeID: rid, TechnologyID: id, UserID: uid, Input: techInput(b)})
	if e != nil {
		write(c, e)
		return
	}
	x, e := h.get.Handle(c.Request.Context(), commands.GetTechnologyQuery{ResumeID: rid, TechnologyID: id, UserID: uid})
	if e != nil {
		write(c, e)
		return
	}
	c.JSON(http.StatusOK, techView(rid, x))
}
func (h *Handler) Delete(c *gin.Context) {
	uid, rid, id, ok := itemIDs(c, "technologyId")
	if !ok {
		return
	}
	if e := h.delete.Handle(c.Request.Context(), commands.DeleteTechnologyCommand{ResumeID: rid, TechnologyID: id, UserID: uid}); e != nil {
		write(c, e)
		return
	}
	c.Status(http.StatusNoContent)
}
func techInput(b TechnologyRequest) commands.TechnologyInput {
	in := commands.TechnologyInput{ID: valueobjects.NewID(), Name: b.Name, Position: b.Position}
	if b.Grade != nil {
		x := enums.TechnologyGrade(*b.Grade)
		in.Grade = &x
	}
	return in
}
func techView(rid valueobjects.ID, x domain.ResumeTechnologyEntry) TechnologyResponse {
	var g *string
	if v := x.Grade(); v != nil {
		s := string(*v)
		g = &s
	}
	return TechnologyResponse{ID: x.ID().String(), ResumeID: rid.String(), Name: x.Name(), Grade: g, Position: x.Position()}
}

func ids(c *gin.Context) (useridentity.ID, valueobjects.ID, bool) {
	uid, ok := common.RequireUserID(c)
	if !ok {
		return useridentity.ID{}, valueobjects.ID{}, false
	}
	rid, ok := common.ParseResumeID(c)
	return uid, rid, ok
}

func itemIDs(c *gin.Context, name string) (useridentity.ID, valueobjects.ID, valueobjects.ID, bool) {
	uid, rid, ok := ids(c)
	if !ok {
		return uid, rid, valueobjects.ID{}, false
	}
	id, e := valueobjects.ParseID(c.Param(name))
	if e != nil {
		common.WriteInvalidBody(c)
		return uid, rid, valueobjects.ID{}, false
	}
	return uid, rid, id, true
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

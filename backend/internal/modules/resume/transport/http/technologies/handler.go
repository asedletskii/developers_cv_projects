package technologies

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type Handler struct {
	create *commands.CreateTechnologyHandler
	get    *queries.GetTechnologyQueryHandler
	list   *queries.ListTechnologiesQueryHandler
	update *commands.UpdateTechnologyHandler
	delete *commands.DeleteTechnologyHandler
}

func NewHandler(create *commands.CreateTechnologyHandler, list *queries.ListTechnologiesQueryHandler, get *queries.GetTechnologyQueryHandler, update *commands.UpdateTechnologyHandler, delete *commands.DeleteTechnologyHandler) *Handler {
	return &Handler{create: create, list: list, get: get, update: update, delete: delete}
}

type TechnologyRequest struct {
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
		common.WriteCommandError(c, e)
		return
	}
	x, e := h.get.Handle(c.Request.Context(), queries.GetTechnologyQuery{ResumeID: rid, TechnologyID: id, UserID: uid})
	if e != nil {
		common.WriteQueryError(c, e)
		return
	}
	c.JSON(http.StatusCreated, x)
}
func (h *Handler) List(c *gin.Context) {
	uid, rid, ok := ids(c)
	if !ok {
		return
	}
	xs, e := h.list.Handle(c.Request.Context(), queries.ListTechnologiesQuery{ResumeID: rid, UserID: uid})
	if e != nil {
		common.WriteQueryError(c, e)
		return
	}
	c.JSON(http.StatusOK, xs)
}
func (h *Handler) Get(c *gin.Context) {
	uid, rid, id, ok := itemIDs(c, "technologyId")
	if !ok {
		return
	}
	x, e := h.get.Handle(c.Request.Context(), queries.GetTechnologyQuery{ResumeID: rid, TechnologyID: id, UserID: uid})
	if e != nil {
		common.WriteQueryError(c, e)
		return
	}
	c.JSON(http.StatusOK, x)
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
		common.WriteCommandError(c, e)
		return
	}
	x, e := h.get.Handle(c.Request.Context(), queries.GetTechnologyQuery{ResumeID: rid, TechnologyID: id, UserID: uid})
	if e != nil {
		common.WriteQueryError(c, e)
		return
	}
	c.JSON(http.StatusOK, x)
}
func (h *Handler) Delete(c *gin.Context) {
	uid, rid, id, ok := itemIDs(c, "technologyId")
	if !ok {
		return
	}
	if e := h.delete.Handle(c.Request.Context(), commands.DeleteTechnologyCommand{ResumeID: rid, TechnologyID: id, UserID: uid}); e != nil {
		common.WriteCommandError(c, e)
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

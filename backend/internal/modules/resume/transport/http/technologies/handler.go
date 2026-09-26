package technologies

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
	create  *commands.CreateTechnologyHandler
	get     *queries.GetTechnologyQueryHandler
	list    *queries.ListTechnologiesQueryHandler
	update  *commands.UpdateTechnologyHandler
	reorder *commands.ReorderTechnologiesHandler
	delete  *commands.DeleteTechnologyHandler
}

func NewHandler(create *commands.CreateTechnologyHandler, list *queries.ListTechnologiesQueryHandler, get *queries.GetTechnologyQueryHandler, update *commands.UpdateTechnologyHandler, reorder *commands.ReorderTechnologiesHandler, delete *commands.DeleteTechnologyHandler) *Handler {
	return &Handler{create: create, list: list, get: get, update: update, reorder: reorder, delete: delete}
}

type technologyRequest struct {
	Name  string  `json:"name"`
	Grade *string `json:"grade,omitempty"`
}

type reorderTechnologiesRequest struct {
	IDs []string `json:"ids"`
}

func (h *Handler) Create(c *gin.Context) {
	uid, rid, ok := ids(c)
	if !ok {
		return
	}
	var body technologyRequest
	if common.BindJSON(c, &body) != nil {
		common.WriteInvalidBody(c)
		return
	}
	id, err := h.create.Handle(c.Request.Context(), commands.CreateTechnologyCommand{ID: valueobjects.NewID(), ResumeID: rid, UserID: uid, Input: techInput(body)})
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.JSON(http.StatusCreated, common.CreatedIDResponse{ID: id.String()})
}

func (h *Handler) List(c *gin.Context) {
	uid, rid, ok := ids(c)
	if !ok {
		return
	}
	xs, err := h.list.Handle(c.Request.Context(), queries.ListTechnologiesQuery{ResumeID: rid, UserID: uid})
	if err != nil {
		common.WriteQueryError(c, err)
		return
	}
	c.JSON(http.StatusOK, xs)
}

func (h *Handler) Get(c *gin.Context) {
	uid, rid, id, ok := itemIDs(c, "technologyId")
	if !ok {
		return
	}
	x, err := h.get.Handle(c.Request.Context(), queries.GetTechnologyQuery{ResumeID: rid, TechnologyID: id, UserID: uid})
	if err != nil {
		common.WriteQueryError(c, err)
		return
	}
	c.JSON(http.StatusOK, x)
}

func (h *Handler) Update(c *gin.Context) {
	uid, rid, id, ok := itemIDs(c, "technologyId")
	if !ok {
		return
	}
	var body technologyRequest
	if common.BindJSON(c, &body) != nil {
		common.WriteInvalidBody(c)
		return
	}
	if err := h.update.Handle(c.Request.Context(), commands.UpdateTechnologyCommand{ResumeID: rid, TechnologyID: id, UserID: uid, Input: techInput(body)}); err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) Reorder(c *gin.Context) {
	uid, rid, ok := ids(c)
	if !ok {
		return
	}
	var body reorderTechnologiesRequest
	if common.BindJSON(c, &body) != nil {
		common.WriteInvalidBody(c)
		return
	}
	ids := make([]valueobjects.ID, len(body.IDs))
	for i, raw := range body.IDs {
		id, err := valueobjects.ParseID(raw)
		if err != nil {
			common.WriteInvalidBody(c)
			return
		}
		ids[i] = id
	}
	if err := h.reorder.Handle(c.Request.Context(), commands.ReorderTechnologiesCommand{ResumeID: rid, UserID: uid, TechnologyIDs: ids}); err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) Delete(c *gin.Context) {
	uid, rid, id, ok := itemIDs(c, "technologyId")
	if !ok {
		return
	}
	if err := h.delete.Handle(c.Request.Context(), commands.DeleteTechnologyCommand{ResumeID: rid, TechnologyID: id, UserID: uid}); err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func techInput(body technologyRequest) commands.TechnologyInput {
	input := commands.TechnologyInput{ID: valueobjects.NewID(), Name: body.Name}
	if body.Grade != nil {
		grade := enums.TechnologyGrade(*body.Grade)
		input.Grade = &grade
	}
	return input
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
	id, err := valueobjects.ParseID(c.Param(name))
	if err != nil {
		common.WriteInvalidBody(c)
		return uid, rid, valueobjects.ID{}, false
	}
	return uid, rid, id, true
}

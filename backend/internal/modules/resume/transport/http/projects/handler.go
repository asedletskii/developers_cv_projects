package projects

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type Handler struct {
	create *commands.CreateProjectHandler
	list   *queries.ListProjectsQueryHandler
	get    *queries.GetProjectQueryHandler
	update *commands.UpdateProjectHandler
	delete *commands.DeleteProjectHandler
}

func NewHandler(create *commands.CreateProjectHandler, list *queries.ListProjectsQueryHandler, get *queries.GetProjectQueryHandler, update *commands.UpdateProjectHandler, delete *commands.DeleteProjectHandler) *Handler {
	return &Handler{create: create, list: list, get: get, update: update, delete: delete}
}

func (h *Handler) Create(c *gin.Context) {
	uid, rid, ok := ids(c)
	if !ok {
		return
	}
	var body projectRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteBindError(c, err)
		return
	}
	input, err := projectInput(body)
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	id, err := h.create.Handle(c.Request.Context(), commands.CreateProjectCommand{ID: valueobjects.NewID(), ResumeID: rid, UserID: uid, Input: input})
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
	items, err := h.list.Handle(c.Request.Context(), queries.ListProjectsQuery{ResumeID: rid, UserID: uid})
	if err != nil {
		common.WriteQueryError(c, err)
		return
	}
	out := make([]ProjectResponse, len(items))
	for i, item := range items {
		out[i] = projectQueryView(item)
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) Get(c *gin.Context) {
	uid, rid, id, ok := itemIDs(c)
	if !ok {
		return
	}
	item, err := h.get.Handle(c.Request.Context(), queries.GetProjectQuery{ResumeID: rid, ProjectID: id, UserID: uid})
	if err != nil {
		common.WriteQueryError(c, err)
		return
	}
	c.JSON(http.StatusOK, projectQueryView(item))
}

func (h *Handler) Update(c *gin.Context) {
	uid, rid, id, ok := itemIDs(c)
	if !ok {
		return
	}
	var body projectRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteBindError(c, err)
		return
	}
	input, err := projectInput(body)
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	if err = h.update.Handle(c.Request.Context(), commands.UpdateProjectCommand{ResumeID: rid, ProjectID: id, UserID: uid, Input: input}); err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) Delete(c *gin.Context) {
	uid, rid, id, ok := itemIDs(c)
	if !ok {
		return
	}
	if err := h.delete.Handle(c.Request.Context(), commands.DeleteProjectCommand{ResumeID: rid, ProjectID: id, UserID: uid}); err != nil {
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

func itemIDs(c *gin.Context) (useridentity.ID, valueobjects.ID, valueobjects.ID, bool) {
	uid, rid, ok := ids(c)
	if !ok {
		return uid, rid, valueobjects.ID{}, false
	}
	id, err := valueobjects.ParseID(c.Param("projectId"))
	if err != nil {
		common.WriteInvalidBody(c)
		return uid, rid, valueobjects.ID{}, false
	}
	return uid, rid, id, true
}

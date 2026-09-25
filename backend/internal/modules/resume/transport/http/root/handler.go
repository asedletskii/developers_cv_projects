package root

import (
	"errors"
	"net/http"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/transport/http/common"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	create *commands.CreateResumeRootHandler
	get    *queries.GetResumeRootQueryHandler
	update *commands.UpdateResumeRootHandler
	delete *commands.DeleteResumeRootHandler
}

func NewHandler(c *commands.CreateResumeRootHandler, g *queries.GetResumeRootQueryHandler, u *commands.UpdateResumeRootHandler, d *commands.DeleteResumeRootHandler) *Handler {
	return &Handler{create: c, get: g, update: u, delete: d}
}
func (h *Handler) Create(c *gin.Context) {
	if h == nil || h.create == nil {
		common.WriteError(c, errors.New("create resume root handler dependencies are nil"))
		return
	}
	uid, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	var body ResumeRootRequest
	if common.BindJSON(c, &body) != nil {
		common.WriteInvalidBody(c)
		return
	}
	cmd, err := body.ToCreateCommand(uid.UUID())
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	id, err := h.create.Handle(c.Request.Context(), cmd)
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.JSON(http.StatusCreated, common.CreatedIDResponse{ID: id.String()})
}
func (h *Handler) Get(c *gin.Context) {
	uid, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	id, ok := common.ParseResumeID(c)
	if !ok {
		return
	}
	if h == nil || h.get == nil {
		common.WriteError(c, errors.New("get resume root handler dependencies are nil"))
		return
	}
	resume, err := h.get.Handle(c.Request.Context(), queries.GetResumeRootQuery{UserID: uid, ResumeID: id})
	if err != nil {
		common.WriteQueryError(c, err)
		return
	}
	c.JSON(http.StatusOK, resume)
}
func (h *Handler) Update(c *gin.Context) {
	if h == nil || h.update == nil {
		common.WriteError(c, errors.New("update resume root handler dependencies are nil"))
		return
	}
	uid, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	id, ok := common.ParseResumeID(c)
	if !ok {
		return
	}
	var body ResumeRootUpdateRequest
	if common.BindJSON(c, &body) != nil {
		common.WriteInvalidBody(c)
		return
	}
	cmd, err := body.ToUpdateCommand(uid.UUID(), id)
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	if err = h.update.Handle(c.Request.Context(), cmd); err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
func (h *Handler) Delete(c *gin.Context) {
	if h == nil || h.delete == nil {
		common.WriteError(c, errors.New("delete resume root handler dependencies are nil"))
		return
	}
	uid, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	id, ok := common.ParseResumeID(c)
	if !ok {
		return
	}
	if err := h.delete.Handle(c.Request.Context(), commands.DeleteResumeRootCommand{UserID: uid, ResumeID: id}); err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

package http

import (
	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type resumeHandler struct {
	create *commands.CreateResumeHandler
	get    *queries.GetResumeQueryHandler
	update *commands.UpdateResumeHandler
	delete *commands.DeleteResumeHandler
}

func newResumeHandler(c *commands.CreateResumeHandler, g *queries.GetResumeQueryHandler, u *commands.UpdateResumeHandler, d *commands.DeleteResumeHandler) *resumeHandler {
	return &resumeHandler{create: c, get: g, update: u, delete: d}
}

type createResumeRequest struct {
	Title    string `json:"title"`
	Language string `json:"language"`
}

func (h *resumeHandler) Create(c *gin.Context) {
	if h == nil || h.create == nil {
		common.WriteError(c, errors.New("create resume handler dependencies are nil"))
		return
	}
	uid, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	var b createResumeRequest
	if common.BindJSON(c, &b) != nil {
		common.WriteInvalidBody(c)
		return
	}
	cmd := commands.CreateResumeCommand{UserID: uid, Title: b.Title, Language: enums.ResumeLanguage(b.Language)}
	id, e := h.create.Handle(c.Request.Context(), cmd)
	if e != nil {
		common.WriteCommandError(c, e)
		return
	}
	c.JSON(http.StatusCreated, common.CreatedIDResponse{ID: id.String()})
}
func (h *resumeHandler) Get(c *gin.Context) {
	if h == nil || h.get == nil {
		common.WriteError(c, errors.New("get resume handler dependencies are nil"))
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
	v, err := h.get.Handle(c.Request.Context(), queries.GetResumeQuery{ResumeID: id, UserID: uid})
	if err != nil {
		common.WriteQueryError(c, err)
		return
	}
	c.JSON(http.StatusOK, v)
}

type UpdateResumeRequest struct {
	Title string `json:"title"`
}

func (h *resumeHandler) Update(c *gin.Context) {
	uid, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	id, ok := common.ParseResumeID(c)
	if !ok {
		return
	}
	if h == nil || h.update == nil {
		common.WriteError(c, errors.New("update resume handler dependencies are nil"))
		return
	}
	var body UpdateResumeRequest
	if common.BindJSON(c, &body) != nil {
		common.WriteInvalidBody(c)
		return
	}
	if err := h.update.Handle(c.Request.Context(), commands.UpdateResumeCommand{ResumeID: id, UserID: uid, Title: body.Title}); err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *resumeHandler) Delete(c *gin.Context) {
	if h == nil || h.delete == nil {
		common.WriteError(c, errors.New("delete resume handler dependencies are nil"))
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
	if e := h.delete.Handle(c.Request.Context(), commands.DeleteResumeCommand{ResumeID: id, UserID: uid}); e != nil {
		common.WriteCommandError(c, e)
		return
	}
	c.Status(http.StatusNoContent)
}

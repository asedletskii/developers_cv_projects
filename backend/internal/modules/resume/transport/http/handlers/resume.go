package handlers

import (
	"net/http"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	"developers-cv/backend/internal/modules/resume/transport/http/requests"

	"github.com/gin-gonic/gin"
)

type ResumeHandler struct {
	create       *commands.CreateResumeHandler
	get          *queries.GetResumeQueryHandler
	replaceTitle *commands.ReplaceResumeTitleHandler
	delete       *commands.DeleteResumeHandler
	changeStatus *commands.ChangeResumeStatusHandler
	list         *queries.ListResumesQueryHandler
}

func NewResumeHandler(c *commands.CreateResumeHandler, g *queries.GetResumeQueryHandler, r *commands.ReplaceResumeTitleHandler, d *commands.DeleteResumeHandler, s *commands.ChangeResumeStatusHandler, l *queries.ListResumesQueryHandler) *ResumeHandler {
	return &ResumeHandler{create: c, get: g, replaceTitle: r, delete: d, changeStatus: s, list: l}
}

func (h *ResumeHandler) List(c *gin.Context) {
	userID, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	items, err := h.list.Handle(c.Request.Context(), queries.ListResumesQuery{UserID: userID})
	if err != nil {
		common.WriteQueryError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *ResumeHandler) Create(c *gin.Context) {
	uid, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	var b requests.CreateResumeRequest
	if err := common.BindJSON(c, &b); err != nil {
		common.WriteBindError(c, err)
		return
	}
	cmd := commands.CreateResumeCommand{UserID: uid, Title: b.Title, Language: domain.ResumeLanguage(b.Language)}
	id, e := h.create.Handle(c.Request.Context(), cmd)
	if e != nil {
		common.WriteCommandError(c, e)
		return
	}
	c.JSON(http.StatusCreated, common.CreatedIDResponse{ID: id.String()})
}
func (h *ResumeHandler) Get(c *gin.Context) {
	uid, id, ok := common.RequireResumeContext(c)
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

func (h *ResumeHandler) ReplaceTitle(c *gin.Context) {
	uid, id, ok := common.RequireResumeContext(c)
	if !ok {
		return
	}
	var body requests.ReplaceResumeTitleRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteBindError(c, err)
		return
	}
	if !common.RequirePositiveVersion(c, body.Version) {
		return
	}
	if err := h.replaceTitle.Handle(c.Request.Context(), commands.ReplaceResumeTitleCommand{ResumeID: id, UserID: uid, Version: body.Version, Title: body.Title}); err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ResumeHandler) ChangeStatus(c *gin.Context) {
	uid, id, ok := common.RequireResumeContext(c)
	if !ok {
		return
	}
	var body requests.ChangeResumeStatusRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteBindError(c, err)
		return
	}
	if !common.RequirePositiveVersion(c, body.Version) {
		return
	}
	err := h.changeStatus.Handle(c.Request.Context(), commands.ChangeResumeStatusCommand{
		ResumeID: id, UserID: uid, Status: domain.ResumeStatus(body.Status), Version: body.Version,
	})
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ResumeHandler) Delete(c *gin.Context) {
	uid, id, ok := common.RequireResumeContext(c)
	if !ok {
		return
	}
	var body requests.DeleteResumeRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteBindError(c, err)
		return
	}
	if !common.RequirePositiveVersion(c, body.Version) {
		return
	}
	if e := h.delete.Handle(c.Request.Context(), commands.DeleteResumeCommand{ResumeID: id, UserID: uid, Version: body.Version}); e != nil {
		common.WriteCommandError(c, e)
		return
	}
	c.Status(http.StatusNoContent)
}

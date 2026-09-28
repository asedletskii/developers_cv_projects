package http

import (
	"errors"
	"net/http"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	apperrors "developers-cv/backend/internal/platform/http/errors"

	"github.com/gin-gonic/gin"
)

type resumeHandler struct {
	create       *commands.CreateResumeHandler
	get          *queries.GetResumeQueryHandler
	replaceTitle *commands.ReplaceResumeTitleHandler
	delete       *commands.DeleteResumeHandler
	list         *queries.ListResumesQueryHandler
}

func newResumeHandler(c *commands.CreateResumeHandler, g *queries.GetResumeQueryHandler, r *commands.ReplaceResumeTitleHandler, d *commands.DeleteResumeHandler, l *queries.ListResumesQueryHandler) *resumeHandler {
	return &resumeHandler{create: c, get: g, replaceTitle: r, delete: d, list: l}
}

func (h *resumeHandler) List(c *gin.Context) {
	userID, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	items, err := h.list.Handle(c.Request.Context(), queries.ListResumesQuery{UserID: userID})
	if err != nil {
		apperrors.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

type createResumeRequest struct {
	Title    string `json:"title"`
	Language string `json:"language"`
}

func (h *resumeHandler) Create(c *gin.Context) {
	uid, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	var b createResumeRequest
	if err := common.BindJSON(c, &b); err != nil {
		common.WriteBindError(c, err)
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

type replaceResumeTitleRequest struct {
	Version int64  `json:"version"`
	Title   string `json:"title"`
}

func (h *resumeHandler) ReplaceTitle(c *gin.Context) {
	uid, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	id, ok := common.ParseResumeID(c)
	if !ok {
		return
	}
	var body replaceResumeTitleRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteBindError(c, err)
		return
	}
	if err := h.replaceTitle.Handle(c.Request.Context(), commands.ReplaceResumeTitleCommand{ResumeID: id, UserID: uid, Version: body.Version, Title: body.Title}); err != nil {
		if errors.Is(err, domain.ErrConcurrencyConflict) {
			c.Status(http.StatusPreconditionFailed)
			return
		}
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *resumeHandler) Delete(c *gin.Context) {
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

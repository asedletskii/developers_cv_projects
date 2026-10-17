package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume_link/application/commands"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/modules/resume_link/transport/http/common"
	"developers-cv/backend/internal/modules/resume_link/transport/http/requests"
	apperrors "developers-cv/backend/internal/platform/http/errors"
)

type ResumeLinkManagementHandler struct {
	create    *commands.CreateResumeLinkHandler
	list      *queries.ListResumeLinksQueryHandler
	delete    *commands.DeleteResumeLinkHandler
	lifecycle *commands.ReplaceStatusResumeLinkHandler
}

func NewResumeLinkManagementHandler(
	create *commands.CreateResumeLinkHandler,
	list *queries.ListResumeLinksQueryHandler,
	delete *commands.DeleteResumeLinkHandler,
	lifecycle *commands.ReplaceStatusResumeLinkHandler,
) *ResumeLinkManagementHandler {
	return &ResumeLinkManagementHandler{create: create, list: list, delete: delete, lifecycle: lifecycle}
}

func (h *ResumeLinkManagementHandler) Create(c *gin.Context) {
	if h == nil || h.create == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("create resume link handler is not configured")))
		return
	}
	userID, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	var body requests.CreateResumeLinkRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteBindError(c, err)
		return
	}
	resumeID, ok := common.ParseResumeIDValue(c, body.ResumeID)
	if !ok {
		return
	}
	id, err := h.create.Handle(c.Request.Context(), commands.CreateResumeLinkCommand{UserID: userID, ResumeID: resumeID, Name: body.Name})
	if err != nil {
		common.WriteCreateError(c, err)
		return
	}
	c.JSON(http.StatusCreated, common.CreatedIDResponse{ID: id.String()})
}

func (h *ResumeLinkManagementHandler) List(c *gin.Context) {
	if h == nil || h.list == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("list resume links handler is not configured")))
		return
	}
	userID, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	resumeID, ok := common.ParseResumeIDQuery(c)
	if !ok {
		return
	}
	items, err := h.list.Handle(c.Request.Context(), queries.ListResumeLinksQuery{UserID: userID, ResumeID: resumeID})
	if err != nil {
		common.WriteListError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *ResumeLinkManagementHandler) ReplaceStatus(c *gin.Context) {
	if h == nil || h.lifecycle == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("replace resume link status handler is not configured")))
		return
	}
	userID, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	linkID, ok := common.ParseLinkID(c)
	if !ok {
		return
	}
	var body requests.ReplaceResumeLinkStatusRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteBindError(c, err)
		return
	}
	if !common.RequirePositiveVersion(c, body.Version) {
		return
	}
	status := domain.ResumeLinkStatus(body.Status)
	if !status.IsValid() {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid status"))
		return
	}
	active := status == domain.ResumeLinkStatusActive
	if err := h.lifecycle.Handle(
		c.Request.Context(),
		commands.ReplaceStatusResumeLinkCommand{
			UserID: userID, LinkID: linkID, Active: active, Version: body.Version,
		},
	); err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ResumeLinkManagementHandler) Delete(c *gin.Context) {
	if h == nil || h.delete == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("delete resume link handler is not configured")))
		return
	}
	userID, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	linkID, ok := common.ParseLinkID(c)
	if !ok {
		return
	}
	if err := h.delete.Handle(c.Request.Context(), commands.DeleteResumeLinkCommand{UserID: userID, LinkID: linkID}); err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

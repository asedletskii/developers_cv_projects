package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume_link/application/commands"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/modules/resume_link/transport/http/common"
	"developers-cv/backend/internal/modules/resume_link/transport/http/public"
	apperrors "developers-cv/backend/internal/platform/http/errors"
)

type Handler struct {
	create      *commands.CreateResumeLinkHandler
	list        *queries.ListResumeLinksQueryHandler
	delete      *commands.DeleteResumeLinkHandler
	public      *public.Handler
	contactView *public.ContactViewHandler
	statistics  *queries.GetResumeLinkStatisticsQueryHandler
	lifecycle   *commands.ReplaceStatusResumeLinkHandler
}

type HandlerDependencies struct {
	Create      *commands.CreateResumeLinkHandler
	List        *queries.ListResumeLinksQueryHandler
	Delete      *commands.DeleteResumeLinkHandler
	Public      *public.Handler
	ContactView *public.ContactViewHandler
	Statistics  *queries.GetResumeLinkStatisticsQueryHandler
	Lifecycle   *commands.ReplaceStatusResumeLinkHandler
}

func NewHandler(deps HandlerDependencies) *Handler {
	return &Handler{
		create: deps.Create, list: deps.List, delete: deps.Delete, public: deps.Public,
		contactView: deps.ContactView, statistics: deps.Statistics, lifecycle: deps.Lifecycle,
	}
}

type createRequest struct {
	ResumeID string `json:"resumeId"`
	Name     string `json:"name"`
}

type replaceStatusRequest struct {
	Status  domain.ResumeLinkStatus `json:"status"`
	Version int64                   `json:"version"`
}

func (h *Handler) Create(c *gin.Context) {
	if h == nil || h.create == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("create resume link handler is not configured")))
		return
	}
	userID, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	var body createRequest
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

func (h *Handler) List(c *gin.Context) {
	if h == nil || h.list == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("list resume links handler is not configured")))
		return
	}
	userID, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	resumeID, ok := common.ParseResumeID(c)
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

func (h *Handler) Statistics(c *gin.Context) {
	if h == nil || h.statistics == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("resume link statistics handler is not configured")))
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
	statistics, err := h.statistics.Handle(c.Request.Context(), queries.GetResumeLinkStatisticsQuery{UserID: userID, LinkID: linkID})
	if err != nil {
		common.WriteStatisticsError(c, err)
		return
	}
	c.JSON(http.StatusOK, statistics)
}

func (h *Handler) ReplaceStatus(c *gin.Context) {
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
	var body replaceStatusRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteBindError(c, err)
		return
	}
	if !common.RequirePositiveVersion(c, body.Version) {
		return
	}
	if !body.Status.IsValid() {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid status"))
		return
	}
	active := body.Status == domain.ResumeLinkStatusActive
	if err := h.lifecycle.Handle(c.Request.Context(), commands.ReplaceStatusResumeLinkCommand{UserID: userID, LinkID: linkID, Active: active, Version: body.Version}); err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) ViewContact(c *gin.Context) {
	if h == nil || h.contactView == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("contact view handler is not configured")))
		return
	}
	h.contactView.Handle(c)
}

func (h *Handler) Delete(c *gin.Context) {
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

func RegisterRoutes(handler *Handler, api *gin.RouterGroup, lk *gin.RouterGroup, contactViewRateLimiter gin.HandlerFunc) {
	api.GET("/cv/:linkId", handler.public.Handle)
	if contactViewRateLimiter != nil {
		api.POST("/cv/:linkId/contact-views", contactViewRateLimiter, handler.ViewContact)
	} else {
		api.POST("/cv/:linkId/contact-views", handler.ViewContact)
	}
	lk.GET("/links", handler.List)
	lk.POST("/link", handler.Create)
	lk.DELETE("/link/:linkId", handler.Delete)
	lk.GET("/link/:linkId/statistics", handler.Statistics)
	lk.PUT("/link/:linkId/status", handler.ReplaceStatus)
}

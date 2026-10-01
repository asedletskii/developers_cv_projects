package http

import (
	"errors"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/application/commands"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/modules/resume_link/transport/http/public"
	apperrors "developers-cv/backend/internal/platform/http/errors"
	jsonhttp "developers-cv/backend/internal/platform/http/json"
	"developers-cv/backend/internal/platform/http/middleware"
	"developers-cv/backend/internal/shared/domain/valueobjects"
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
	return &Handler{create: deps.Create, list: deps.List, delete: deps.Delete, public: deps.Public, contactView: deps.ContactView, statistics: deps.Statistics, lifecycle: deps.Lifecycle}
}

type createRequest struct {
	ResumeID string `json:"resumeId"`
	Name     string `json:"name"`
}

type createdResponse struct {
	ID string `json:"id"`
}

func (h *Handler) Create(c *gin.Context) {
	if h == nil || h.create == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("create resume link handler is not configured")))
		return
	}

	userUUID, ok := middleware.GetUserID(c)
	if !ok {
		apperrors.Unauthorized(c)
		return
	}
	userID, err := valueobjects.IDFromUUID(userUUID)
	if err != nil {
		apperrors.WriteError(c, apperrors.Internal(err))
		return
	}

	var body createRequest
	if err := jsonhttp.BindJSON(c, &body); err != nil {
		writeBindError(c, err)
		return
	}
	resumeUUID, err := uuid.Parse(body.ResumeID)
	if err != nil {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid resume id"))
		return
	}
	resumeID, err := valueobjects.IDFromUUID(resumeUUID)
	if err != nil {
		apperrors.WriteError(c, apperrors.Internal(err))
		return
	}

	id, err := h.create.Handle(c.Request.Context(), commands.CreateResumeLinkCommand{
		UserID: userID, ResumeID: resumeID, Name: body.Name,
	})
	if err != nil {
		writeCreateError(c, err)
		return
	}

	c.JSON(http.StatusCreated, createdResponse{ID: id.String()})
}

func (h *Handler) List(c *gin.Context) {
	if h == nil || h.list == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("list resume links handler is not configured")))
		return
	}
	userUUID, ok := middleware.GetUserID(c)
	if !ok {
		apperrors.Unauthorized(c)
		return
	}
	userID, err := valueobjects.IDFromUUID(userUUID)
	if err != nil {
		apperrors.WriteError(c, apperrors.Internal(err))
		return
	}
	resumeUUID, err := uuid.Parse(c.Query("resumeId"))
	if err != nil {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid resume id"))
		return
	}
	resumeID, err := valueobjects.IDFromUUID(resumeUUID)
	if err != nil {
		apperrors.WriteError(c, apperrors.Internal(err))
		return
	}

	items, err := h.list.Handle(c.Request.Context(), queries.ListResumeLinksQuery{
		UserID: userID, ResumeID: resumeID,
	})
	if err != nil {
		if errors.Is(err, resumequeries.ErrResumeNotFound) {
			apperrors.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Resume not found"))
			return
		}
		apperrors.WriteError(c, apperrors.Internal(err))
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) Statistics(c *gin.Context) {
	if h == nil || h.statistics == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("resume link statistics handler is not configured")))
		return
	}

	userUUID, ok := middleware.GetUserID(c)
	if !ok {
		apperrors.Unauthorized(c)
		return
	}
	userID, err := valueobjects.IDFromUUID(userUUID)
	if err != nil {
		apperrors.WriteError(c, apperrors.Internal(err))
		return
	}
	linkID, err := valueobjects.ParseID(c.Param("linkId"))
	if err != nil {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid link id"))
		return
	}

	statistics, err := h.statistics.Handle(c.Request.Context(), queries.GetResumeLinkStatisticsQuery{
		UserID: userID,
		LinkID: linkID,
	})
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) || errors.Is(err, resumequeries.ErrResumeNotFound) {
			apperrors.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Resume link not found"))
			return
		}
		apperrors.WriteError(c, apperrors.Internal(err))
		return
	}

	c.JSON(http.StatusOK, statistics)
}

type replaceStatusRequest struct {
	Status  string `json:"status"`
	Version int64  `json:"version"`
}

func (h *Handler) ReplaceStatus(c *gin.Context) {
	if h == nil || h.lifecycle == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("replace resume link status handler is not configured")))
		return
	}
	userUUID, ok := middleware.GetUserID(c)
	if !ok {
		apperrors.Unauthorized(c)
		return
	}
	userID, err := valueobjects.IDFromUUID(userUUID)
	if err != nil {
		apperrors.WriteError(c, apperrors.Internal(err))
		return
	}
	linkID, err := valueobjects.ParseID(c.Param("linkId"))
	if err != nil {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid link id"))
		return
	}
	var body replaceStatusRequest
	if err := jsonhttp.BindJSON(c, &body); err != nil {
		writeBindError(c, err)
		return
	}
	if body.Version < 1 {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Version must be positive"))
		return
	}
	var active bool
	switch body.Status {
	case domain.ResumeLinkStatusActive:
		active = true
	case domain.ResumeLinkStatusInactive:
		active = false
	default:
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid status"))
		return
	}
	err = h.lifecycle.Handle(c.Request.Context(), commands.ReplaceStatusResumeLinkCommand{
		UserID: userID, LinkID: linkID, Active: active, Version: body.Version,
	})
	if err != nil {
		if errors.Is(err, domain.ErrConcurrencyConflict) {
			apperrors.WriteError(c, apperrors.New(apperrors.CodePreconditionFailed, http.StatusPreconditionFailed, "Resume link was modified"))
			return
		}
		if errors.Is(err, domain.ErrEntityNotFound) || errors.Is(err, resumequeries.ErrResumeNotFound) {
			apperrors.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Resume link not found"))
			return
		}
		apperrors.WriteError(c, apperrors.Internal(err))
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
	userUUID, ok := middleware.GetUserID(c)
	if !ok {
		apperrors.Unauthorized(c)
		return
	}
	userID, err := valueobjects.IDFromUUID(userUUID)
	if err != nil {
		apperrors.WriteError(c, apperrors.Internal(err))
		return
	}
	linkID, err := valueobjects.ParseID(c.Param("linkId"))
	if err != nil {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid link id"))
		return
	}

	err = h.delete.Handle(c.Request.Context(), commands.DeleteResumeLinkCommand{UserID: userID, LinkID: linkID})
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			apperrors.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Resume link not found"))
			return
		}
		if errors.Is(err, resumequeries.ErrResumeNotFound) {
			apperrors.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Resume link not found"))
			return
		}
		apperrors.WriteError(c, apperrors.Internal(err))
		return
	}
	c.Status(http.StatusNoContent)
}

func RegisterRoutes(handler *Handler, api *gin.RouterGroup, lk *gin.RouterGroup) {
	api.GET("/cv/:linkId", handler.public.Handle)
	api.POST("/cv/:linkId/contact-views", handler.ViewContact)
	lk.GET("/links", handler.List)
	lk.POST("/link", handler.Create)
	lk.DELETE("/link/:linkId", handler.Delete)
	lk.GET("/link/:linkId/statistics", handler.Statistics)
	lk.PUT("/link/:linkId/status", handler.ReplaceStatus)
}

func writeBindError(c *gin.Context, err error) {
	if err == nil {
		return
	}
	apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid request body"))
}

func writeCreateError(c *gin.Context, err error) {
	if errors.Is(err, resumequeries.ErrResumeNotFound) || errors.Is(err, domain.ErrEntityNotFound) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Resume not found"))
		return
	}
	if errors.Is(err, domain.ErrDuplicateResumeLink) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeConflict, http.StatusConflict, "Resume link already exists"))
		return
	}
	if errors.Is(err, domain.ErrResumeLinkLimitExceeded) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeConflict, http.StatusConflict, "Resume link limit exceeded"))
		return
	}
	if errors.Is(err, domain.ErrInvalidResumeLinkName) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid link name"))
		return
	}
	if errors.Is(err, domain.ErrInvalidResumeLink) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid resume link"))
		return
	}
	apperrors.WriteError(c, apperrors.Internal(err))
}

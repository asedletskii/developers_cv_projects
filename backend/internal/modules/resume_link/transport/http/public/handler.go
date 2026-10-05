package public

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/application/commands"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/transport/http/common"
	apperrors "developers-cv/backend/internal/platform/http/errors"
	"developers-cv/backend/internal/platform/http/middleware"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type Handler struct{ query publicResumeLinkQuery }

type publicResumeLinkQuery interface {
	HandleWithVisit(context.Context, queries.PublicResumeLinkQuery) (resumequeries.PublicResumeView, ports.VisitorIdentityResult, error)
}

func NewHandler(query publicResumeLinkQuery) *Handler { return &Handler{query: query} }

func (h *Handler) Handle(c *gin.Context) {
	if h == nil || h.query == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("public resume link query handler is not configured")))
		return
	}
	linkID, ok := common.ParseLinkID(c)
	if !ok {
		return
	}
	referrerHost := ""
	if raw := c.Request.Referer(); raw != "" {
		if parsed, err := url.Parse(raw); err == nil {
			referrerHost = parsed.Hostname()
		}
	}
	input := &commands.CreatePublicVisitInput{
		LinkID: linkID,
		Visitor: ports.PublicVisitInput{
			IP:        c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			Referrer:  referrerHost,
		},
	}
	if userID, ok := middleware.GetUserID(c); ok {
		ownerID, err := sharedvalueobjects.IDFromUUID(userID)
		if err == nil {
			input.Visitor.OwnerID = ownerID
		}
	}
	cookieName := "cv_vid"
	if c.Request.TLS != nil || c.Request.Header.Get("X-Forwarded-Proto") == "https" {
		cookieName = "__Secure-cv_vid"
	}
	if cookie, err := c.Request.Cookie(cookieName); err == nil {
		if id, err := uuid.Parse(cookie.Value); err == nil && id != uuid.Nil {
			if visitorID, err := sharedvalueobjects.IDFromUUID(id); err == nil {
				input.Visitor.VisitorID = visitorID
			}
			input.Visitor.HasCookie = true
		}
	}
	view, identity, err := h.query.HandleWithVisit(c.Request.Context(), queries.PublicResumeLinkQuery{LinkID: linkID, Visit: input})
	if err != nil {
		common.WritePublicError(c, err)
		return
	}
	if identity.SetCookie && identity.Cookie != nil {
		http.SetCookie(c.Writer, identity.Cookie)
	}
	if !identity.ViewID.IsZero() {
		c.Header("X-View-Id", identity.ViewID.String())
	}
	c.Header("Cache-Control", "no-store")
	c.JSON(http.StatusOK, view)
}

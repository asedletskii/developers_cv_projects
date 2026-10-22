package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	analyticscontracts "developers-cv/backend/internal/modules/analytics/contracts"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/transport/http/common"
	"developers-cv/backend/internal/modules/resume_link/transport/http/responses"
	apperrors "developers-cv/backend/internal/platform/http/errors"
	"developers-cv/backend/internal/platform/http/middleware"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type PublicResumeLinkHandler struct{ query PublicResumeLinkQuery }

type PublicResumeLinkQuery interface {
	HandleWithVisit(context.Context, queries.PublicResumeLinkQuery) (queries.PublicResumeResult, analyticscontracts.VisitorIdentityResult, error)
}

const referrerHostHeader = "X-Referrer-Host"

func parseReferrerHost(raw string) string {
	host := strings.ToLower(strings.TrimSpace(raw))
	if host == "" || len(host) > 253 || strings.ContainsAny(host, "/ \t:?#@") {
		return ""
	}
	return host
}

func NewPublicResumeLinkHandler(query PublicResumeLinkQuery) *PublicResumeLinkHandler {
	return &PublicResumeLinkHandler{query: query}
}

func (h *PublicResumeLinkHandler) Handle(c *gin.Context) {
	if h == nil || h.query == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("public resume link query handler is not configured")))
		return
	}
	linkID, ok := common.ParseLinkID(c)
	if !ok {
		return
	}
	referrerHost := parseReferrerHost(c.GetHeader(referrerHostHeader))
	input := &analyticscontracts.CreatePublicVisitInput{
		LinkID: linkID,
		Visitor: analyticscontracts.PublicVisitInput{
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
	for _, cookieName := range []string{analyticscontracts.ProductionVisitorCookieName, analyticscontracts.LocalVisitorCookieName} {
		cookie, err := c.Request.Cookie(cookieName)
		if err != nil {
			continue
		}
		id, err := uuid.Parse(cookie.Value)
		if err != nil || id == uuid.Nil {
			continue
		}
		visitorID, err := sharedvalueobjects.IDFromUUID(id)
		if err != nil {
			continue
		}
		input.Visitor.VisitorID = visitorID
		input.Visitor.HasCookie = true
		break
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
	c.JSON(http.StatusOK, responses.NewPublicResumeResponse(view))
}

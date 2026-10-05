package public

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/transport/http/common"
	apperrors "developers-cv/backend/internal/platform/http/errors"
)

type contactViewQuery interface {
	Handle(context.Context, queries.ContactViewThroughLinkQuery) (resumequeries.ContactView, error)
}

type ContactViewHandler struct {
	query    contactViewQuery
	verifier ports.TokenVerifier
}

func NewContactViewHandler(query contactViewQuery, verifier ports.TokenVerifier) *ContactViewHandler {
	return &ContactViewHandler{query: query, verifier: verifier}
}

type contactViewRequest struct {
	Contact        string `json:"contact"`
	TurnstileToken string `json:"turnstileToken"`
}

func (h *ContactViewHandler) Handle(c *gin.Context) {
	if h == nil || h.query == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("contact view handler is not configured")))
		return
	}
	linkID, ok := common.ParseLinkID(c)
	if !ok {
		return
	}
	var body contactViewRequest
	if err := common.BindJSON(c, &body); err != nil {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid request body"))
		return
	}
	contact := resumedomain.ContactKey(body.Contact)
	if !contact.IsValid() {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid request body"))
		return
	}
	if h.verifier != nil {
		if strings.TrimSpace(body.TurnstileToken) == "" || len(body.TurnstileToken) > ports.MaxTokenLength {
			apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid request body"))
			return
		}
		if err := h.verifier.Verify(c.Request.Context(), body.TurnstileToken, c.ClientIP()); err != nil {
			apperrors.WriteError(c, apperrors.New(apperrors.CodeForbidden, http.StatusForbidden, "Contact verification failed"))
			return
		}
	}
	view, err := h.query.Handle(c.Request.Context(), queries.ContactViewThroughLinkQuery{LinkID: linkID, Contact: contact})
	if err != nil {
		if errors.Is(err, resumequeries.ErrContactNotFound) {
			apperrors.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Contact not found"))
			return
		}
		if errors.Is(err, resumequeries.ErrContactNotHidden) {
			apperrors.WriteError(c, apperrors.New(apperrors.CodeConflict, http.StatusConflict, "Contact is not hidden"))
			return
		}
		common.WritePublicError(c, err)
		return
	}
	c.Header("Cache-Control", "no-store")
	c.JSON(http.StatusOK, view)
}

package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/transport/http/common"
	"developers-cv/backend/internal/modules/resume_link/transport/http/requests"
	apperrors "developers-cv/backend/internal/platform/http/errors"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ContactViewQuery interface {
	Handle(context.Context, queries.ContactViewThroughLinkQuery) (queries.ContactViewThroughLinkResult, error)
}

type PublicContactViewHandler struct {
	query    ContactViewQuery
	verifier ports.TokenVerifier
}

func NewPublicContactViewHandler(query ContactViewQuery, verifier ports.TokenVerifier) *PublicContactViewHandler {
	return &PublicContactViewHandler{query: query, verifier: verifier}
}

func (h *PublicContactViewHandler) Handle(c *gin.Context) {
	if h == nil || h.query == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("contact view handler is not configured")))
		return
	}
	linkID, ok := common.ParseLinkID(c)
	if !ok {
		return
	}
	var body requests.PublicContactViewRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteBindError(c, err)
		return
	}
	initiator := queries.ContactViewInitiator(body.Initiator)
	if initiator == "" && body.Contact != nil && strings.TrimSpace(*body.Contact) != "" {
		initiator = queries.ContactViewInitiatorContact
	}
	var contact resumedomain.ContactKey
	switch initiator {
	case queries.ContactViewInitiatorContact:
		if body.Contact == nil {
			apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid request body"))
			return
		}
		parsed, err := resumedomain.ParseContactKey(*body.Contact)
		if err != nil {
			apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid request body"))
			return
		}
		contact = parsed
	case queries.ContactViewInitiatorDownload:
		if body.Contact != nil {
			apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid request body"))
			return
		}
	default:
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
	var viewID sharedvalueobjects.ID
	if strings.TrimSpace(body.ViewID) != "" {
		if parsed, err := uuid.Parse(strings.TrimSpace(body.ViewID)); err == nil && parsed != uuid.Nil {
			parsedID, idErr := sharedvalueobjects.IDFromUUID(parsed)
			if idErr == nil {
				viewID = parsedID
			}
		}
	}
	view, err := h.query.Handle(c.Request.Context(), queries.ContactViewThroughLinkQuery{
		LinkID:    linkID,
		Initiator: initiator,
		Contact:   contact,
		ViewID:    viewID,
	})
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

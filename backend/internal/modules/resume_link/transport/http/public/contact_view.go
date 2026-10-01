package public

import (
	"context"
	"errors"

	"net/http"

	"github.com/gin-gonic/gin"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	linkqueries "developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/domain"
	apperrors "developers-cv/backend/internal/platform/http/errors"
	jsonhttp "developers-cv/backend/internal/platform/http/json"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type contactViewQuery interface {
	Handle(context.Context, linkqueries.ContactViewThroughLinkQuery) (resumequeries.ContactView, error)
}

type ContactViewHandler struct {
	query contactViewQuery
}

func NewContactViewHandler(query contactViewQuery) *ContactViewHandler {
	return &ContactViewHandler{query: query}
}

type contactViewRequest struct {
	Contact string `json:"contact"`
}

func (h *ContactViewHandler) Handle(c *gin.Context) {
	if h == nil || h.query == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("contact view handler is not configured")))
		return
	}

	linkID, err := identity.ParseID(c.Param("linkId"))
	if err != nil {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid link id"))
		return
	}

	var body contactViewRequest
	if err := jsonhttp.BindJSON(c, &body); err != nil {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid request body"))
		return
	}

	contact := resumedomain.ContactKey(body.Contact)
	if !contact.IsValid() {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid request body"))
		return
	}

	view, err := h.query.Handle(c.Request.Context(), linkqueries.ContactViewThroughLinkQuery{
		LinkID:  linkID,
		Contact: contact,
	})
	if err != nil {
		if errors.Is(err, domain.ErrResumeLinkInactive) {
			apperrors.WriteError(c, apperrors.New(apperrors.CodeGone, http.StatusGone, "Resume link is no longer active"))
			return
		}
		if errors.Is(err, domain.ErrEntityNotFound) || errors.Is(err, resumequeries.ErrResumeNotFound) {
			apperrors.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Resume not found"))
			return
		}
		if errors.Is(err, resumequeries.ErrContactNotFound) {
			apperrors.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Contact not found"))
			return
		}
		if errors.Is(err, resumequeries.ErrContactNotHidden) {
			apperrors.WriteError(c, apperrors.New(apperrors.CodeConflict, http.StatusConflict, "Contact is not hidden"))
			return
		}
		apperrors.WriteError(c, apperrors.Internal(err))
		return
	}

	c.Header("Cache-Control", "no-store")
	c.JSON(http.StatusOK, view)
}

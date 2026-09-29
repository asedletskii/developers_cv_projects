package public

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	linkqueries "developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/modules/resume_link/domain/valueobjects"
	apperrors "developers-cv/backend/internal/platform/http/errors"
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

	linkID, err := valueobjects.ParseID(c.Param("linkId"))
	if err != nil {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid link id"))
		return
	}

	var body contactViewRequest
	if err := bindJSON(c, &body); err != nil {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid request body"))
		return
	}

	contact := enums.ContactKey(body.Contact)
	if !contact.IsValid() {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid request body"))
		return
	}

	view, err := h.query.Handle(c.Request.Context(), linkqueries.ContactViewThroughLinkQuery{
		LinkID:  linkID,
		Contact: contact,
	})
	if err != nil {
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

func bindJSON(c *gin.Context, target any) error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return io.ErrUnexpectedEOF
		}
		return err
	}
	return nil
}

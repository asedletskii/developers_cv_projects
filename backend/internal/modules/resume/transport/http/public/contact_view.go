package public

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
)

type contactViewQuery interface {
	Handle(context.Context, queries.ContactViewQuery) (queries.ContactView, error)
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
	resumeID, ok := common.ParseResumeID(c)
	if !ok {
		return
	}

	var body contactViewRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteBindError(c, err)
		return
	}

	contact := enums.ContactKey(body.Contact)
	if !contact.IsValid() {
		common.WriteInvalidBody(c)
		return
	}

	view, err := h.query.Handle(c.Request.Context(), queries.ContactViewQuery{
		ResumeID: resumeID,
		Contact:  contact,
	})
	if err != nil {
		if errors.Is(err, queries.ErrContactNotFound) {
			common.WriteContactNotFound(c)
			return
		}
		if errors.Is(err, queries.ErrContactNotHidden) {
			common.WriteContactNotHidden(c)
			return
		}
		common.WriteQueryError(c, err)
		return
	}

	c.Header("Cache-Control", "no-store")
	c.JSON(http.StatusOK, view)
}

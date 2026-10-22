package queries

import (
	"context"
	"fmt"

	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
	contactkey "developers-cv/backend/internal/shared/domain/contact"
)

type ContactViewReader interface {
	GetContactForView(context.Context, resumecontracts.ContactViewQuery) (resumecontracts.ContactView, error)
}

type ContactViewQueryHandler struct {
	reader ContactViewReader
}

func NewContactViewQueryHandler(reader ContactViewReader) *ContactViewQueryHandler {
	return &ContactViewQueryHandler{reader: reader}
}

func (h *ContactViewQueryHandler) Handle(ctx context.Context, query resumecontracts.ContactViewQuery) (resumecontracts.ContactView, error) {
	if h == nil || h.reader == nil {
		return resumecontracts.ContactView{}, fmt.Errorf("contact view query handler dependencies are nil")
	}
	if query.ResumeID.IsZero() {
		return resumecontracts.ContactView{}, fmt.Errorf("invalid contact view query")
	}
	if _, err := contactkey.Parse(string(query.Contact)); err != nil {
		return resumecontracts.ContactView{}, fmt.Errorf("invalid contact view query: %w", err)
	}
	return h.reader.GetContactForView(ctx, query)
}

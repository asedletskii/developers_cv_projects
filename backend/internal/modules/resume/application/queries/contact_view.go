package queries

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

type ContactViewQuery struct {
	ResumeID valueobjects.ID
	Contact  enums.ContactKey
}

type ContactView struct {
	Contact enums.ContactKey `json:"contact"`
	Value   string           `json:"value"`
}

type ContactViewReader interface {
	GetContactForView(context.Context, ContactViewQuery) (ContactView, error)
}

type ContactViewQueryHandler struct {
	reader ContactViewReader
}

func NewContactViewQueryHandler(reader ContactViewReader) *ContactViewQueryHandler {
	return &ContactViewQueryHandler{reader: reader}
}

func (h *ContactViewQueryHandler) Handle(ctx context.Context, query ContactViewQuery) (ContactView, error) {
	if h == nil || h.reader == nil {
		return ContactView{}, fmt.Errorf("contact view query handler dependencies are nil")
	}
	if query.ResumeID.IsZero() || !query.Contact.IsValid() {
		return ContactView{}, fmt.Errorf("invalid contact view query")
	}
	return h.reader.GetContactForView(ctx, query)
}

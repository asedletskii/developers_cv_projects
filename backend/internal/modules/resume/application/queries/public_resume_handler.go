package queries

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func NewPublicResumeQueryHandler(reader PublicResumeReader) *PublicResumeQueryHandler {
	return &PublicResumeQueryHandler{reader: reader}
}

func (h *PublicResumeQueryHandler) Handle(ctx context.Context, query PublicResumeQuery) (PublicResumeView, error) {
	if h == nil || h.reader == nil {
		return PublicResumeView{}, fmt.Errorf("public resume query handler dependencies are nil")
	}
	if query.ID == uuid.Nil {
		return PublicResumeView{}, fmt.Errorf("invalid resume id")
	}
	return h.reader.GetPublicResume(ctx, query)
}

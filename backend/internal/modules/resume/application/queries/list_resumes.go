package queries

import (
	"context"
	"fmt"

	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type GetListResumesQuery struct {
	UserID useridentity.ID
}

type ListResumesReader interface {
	ListResumes(ctx context.Context, query GetListResumesQuery) ([]ResumeListItem, error)
}

type ListResumesHandler struct {
	reader ListResumesReader
}

func NewListResumesHandler(reader ListResumesReader) *ListResumesHandler {
	return &ListResumesHandler{reader: reader}
}

func (h *ListResumesHandler) Handle(ctx context.Context, query GetListResumesQuery) ([]ResumeListItem, error) {
	if query.UserID.IsZero() {
		return nil, fmt.Errorf("invalid user id")
	}
	return h.reader.ListResumes(ctx, query)
}

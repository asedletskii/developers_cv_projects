package queries

import (
	"context"
	"fmt"

	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type ListResumesQuery struct {
	UserID useridentity.ID
}

type ListResumesReader interface {
	ListResumes(ctx context.Context, query ListResumesQuery) ([]ResumeListItem, error)
}

type ListResumesQueryHandler struct {
	reader ListResumesReader
}

func NewListResumesQueryHandler(reader ListResumesReader) *ListResumesQueryHandler {
	return &ListResumesQueryHandler{reader: reader}
}

func (h *ListResumesQueryHandler) Handle(ctx context.Context, query ListResumesQuery) ([]ResumeListItem, error) {
	if h == nil || h.reader == nil {
		return nil, fmt.Errorf("list resumes handler dependencies are nil")
	}
	if query.UserID.IsZero() {
		return nil, fmt.Errorf("invalid user id")
	}
	return h.reader.ListResumes(ctx, query)
}

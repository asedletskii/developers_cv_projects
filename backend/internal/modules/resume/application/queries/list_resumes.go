package queries

import (
	"context"
	"fmt"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ListResumesQuery struct {
	UserID sharedvalueobjects.ID
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

type ResumeListItem struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Language string `json:"language"`
	Status   string `json:"status"`
	Version  int64  `json:"version"`
}

package queries

import (
	"context"
	"fmt"
	"time"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumevalueobjects "developers-cv/backend/internal/modules/resume/domain/valueobjects"
	uservalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ListResumeLinksQuery struct {
	UserID   uservalueobjects.ID
	ResumeID resumevalueobjects.ID
}

type ResumeLinkListItem struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Version   int64     `json:"version"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ResumeLinkReader interface {
	ListByResumeID(context.Context, ListResumeLinksQuery) ([]ResumeLinkListItem, error)
}

type ResumeOwnershipReader interface {
	GetOwned(context.Context, resumequeries.GetResumeQuery) (resumequeries.ResumeView, error)
}

type ListResumeLinksQueryHandler struct {
	reader  ResumeLinkReader
	resumes ResumeOwnershipReader
}

func NewListResumeLinksQueryHandler(reader ResumeLinkReader, resumes ResumeOwnershipReader) *ListResumeLinksQueryHandler {
	return &ListResumeLinksQueryHandler{reader: reader, resumes: resumes}
}

func (h *ListResumeLinksQueryHandler) Handle(ctx context.Context, query ListResumeLinksQuery) ([]ResumeLinkListItem, error) {
	if h == nil || h.reader == nil || h.resumes == nil {
		return nil, fmt.Errorf("list resume links query handler dependencies are nil")
	}
	if query.UserID.IsZero() || query.ResumeID.IsZero() {
		return nil, fmt.Errorf("invalid list resume links query")
	}
	if _, err := h.resumes.GetOwned(ctx, resumequeries.GetResumeQuery{ResumeID: query.ResumeID, UserID: query.UserID}); err != nil {
		return nil, err
	}
	return h.reader.ListByResumeID(ctx, query)
}

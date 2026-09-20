package queries

import (
	"context"
	"fmt"
	"time"

	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
	application "developers-cv/backend/internal/modules/resume_link/application"
	"developers-cv/backend/internal/modules/resume_link/application/ownership"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ListResumeLinksQuery struct {
	UserID   sharedvalueobjects.ID
	ResumeID sharedvalueobjects.ID
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

type ListResumeLinksQueryHandler struct {
	reader  ResumeLinkReader
	resumes resumecontracts.ResumeOwnershipChecker
}

func NewListResumeLinksQueryHandler(reader ResumeLinkReader, resumes resumecontracts.ResumeOwnershipChecker) *ListResumeLinksQueryHandler {
	return &ListResumeLinksQueryHandler{reader: reader, resumes: resumes}
}

func (h *ListResumeLinksQueryHandler) Handle(ctx context.Context, query ListResumeLinksQuery) ([]ResumeLinkListItem, error) {
	if h == nil || h.reader == nil || h.resumes == nil {
		return nil, fmt.Errorf("list resume links query handler dependencies are nil")
	}
	if query.UserID.IsZero() || query.ResumeID.IsZero() {
		return nil, fmt.Errorf("invalid list resume links query")
	}
	if err := ownership.EnsureResumeOwned(ctx, h.resumes, resumecontracts.ResumeOwnershipQuery{ResumeID: query.ResumeID, UserID: query.UserID}); err != nil {
		return nil, err
	}
	items, err := h.reader.ListByResumeID(ctx, query)
	if err != nil {
		return nil, application.NormalizeLinkNotFoundError(err)
	}
	return items, nil
}

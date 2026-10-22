package queries

import (
	"context"
	"fmt"

	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
)

type OwnedResumeSectionsReader interface {
	GetOwnedSections(context.Context, resumecontracts.OwnedResumeSectionsQuery) ([]ResumeSectionSettings, error)
}

type OwnedResumeSectionsQueryHandler struct {
	reader OwnedResumeSectionsReader
}

func NewOwnedResumeSectionsQueryHandler(reader OwnedResumeSectionsReader) *OwnedResumeSectionsQueryHandler {
	return &OwnedResumeSectionsQueryHandler{reader: reader}
}

func (h *OwnedResumeSectionsQueryHandler) Handle(ctx context.Context, query resumecontracts.OwnedResumeSectionsQuery) ([]resumecontracts.ResumeSectionSettings, error) {
	if h == nil || h.reader == nil {
		return nil, fmt.Errorf("owned resume sections query handler dependencies are nil")
	}
	if query.ResumeID.IsZero() || query.UserID.IsZero() {
		return nil, fmt.Errorf("invalid owned resume sections query")
	}
	sections, err := h.reader.GetOwnedSections(ctx, query)
	if err != nil {
		return nil, err
	}
	result := make([]resumecontracts.ResumeSectionSettings, 0, len(sections))
	for _, section := range sections {
		result = append(result, resumecontracts.ResumeSectionSettings(section))
	}
	return result, nil
}

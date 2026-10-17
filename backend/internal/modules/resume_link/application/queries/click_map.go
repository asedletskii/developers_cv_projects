package queries

import (
	"context"
	"fmt"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

const ClickGridSize = 20

type LinkClicksQuery struct {
	UserID       sharedvalueobjects.ID
	LinkID       sharedvalueobjects.ID
	Section      string
	Breakpoint   string
	IncludeStale bool
	CurrentHash  string
}

type ClickCell struct {
	X     int64 `json:"x"`
	Y     int64 `json:"y"`
	Count int64 `json:"count"`
}

type ClickGrid struct {
	Section    string      `json:"section"`
	Grid       int         `json:"grid"`
	Cells      []ClickCell `json:"cells"`
	Total      int64       `json:"total"`
	StaleTotal int64       `json:"staleTotal"`
}

type AnalyticsClickReader interface {
	GetClicks(context.Context, LinkClicksQuery) (ClickGrid, error)
}

type AnalyticsClickQueryHandler struct {
	links        ResumeLinkPublicReader
	resumes      ports.ResumeOwnershipReader
	publicResume ResumePublicQuery
	layout       ports.LayoutSnapshotBuilder
	reader       AnalyticsClickReader
}

func NewAnalyticsClickQueryHandler(
	links ResumeLinkPublicReader,
	resumes ports.ResumeOwnershipReader,
	publicResume ResumePublicQuery,
	layout ports.LayoutSnapshotBuilder,
	reader AnalyticsClickReader,
) *AnalyticsClickQueryHandler {
	return &AnalyticsClickQueryHandler{links: links, resumes: resumes, publicResume: publicResume, layout: layout, reader: reader}
}

func (h *AnalyticsClickQueryHandler) Handle(ctx context.Context, q LinkClicksQuery) (ClickGrid, error) {
	if h == nil || h.links == nil || h.resumes == nil || h.publicResume == nil || h.layout == nil || h.reader == nil {
		return ClickGrid{}, fmt.Errorf("analytics click query handler dependencies are nil")
	}
	if q.UserID.IsZero() || q.LinkID.IsZero() || q.Section == "" {
		return ClickGrid{}, fmt.Errorf("%w: invalid click query", domain.ErrInvalidAnalyticsClick)
	}
	if !validClickBreakpoint(q.Breakpoint) {
		return ClickGrid{}, fmt.Errorf("%w: invalid breakpoint", domain.ErrInvalidAnalyticsClick)
	}

	_, resumeID, err := authorizeAnalyticsLink(ctx, h.links, h.resumes, q.UserID, q.LinkID)
	if err != nil {
		return ClickGrid{}, err
	}
	publicResume, err := h.publicResume.Handle(ctx, resumequeries.PublicResumeQuery{ID: resumeID})
	if err != nil {
		return ClickGrid{}, err
	}
	snapshot, err := h.layout.Build(publicResume)
	if err != nil {
		return ClickGrid{}, fmt.Errorf("build current analytics layout: %w", err)
	}
	section := resumedomain.SectionID(q.Section)
	currentHash, ok := snapshot.HashFor(section)
	if !ok {
		return ClickGrid{}, fmt.Errorf("%w: section %q is not visible", domain.ErrInvalidAnalyticsClick, q.Section)
	}
	q.CurrentHash = currentHash
	return h.reader.GetClicks(ctx, q)
}

func validClickBreakpoint(value string) bool {
	switch value {
	case "xs", "md", "lg":
		return true
	default:
		return false
	}
}

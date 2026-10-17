package queries

import (
	"context"
	"fmt"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type LinkSectionsQuery struct {
	UserID       sharedvalueobjects.ID
	LinkID       sharedvalueobjects.ID
	Period       AnalyticsPeriod
	SectionOrder []string
}

type SectionsStatistics struct {
	Reads    int64               `json:"reads"`
	Sections []SectionStatistics `json:"sections"`
}

type SectionStatistics struct {
	Section string `json:"section"`
	Reached int64  `json:"reached"`
	Exits   int64  `json:"exits"`
}

type SectionsReader interface {
	GetSections(context.Context, LinkSectionsQuery) (SectionsStatistics, error)
}

type ResumeOverviewReader interface {
	GetOwnedOverview(context.Context, resumequeries.GetResumeQuery) (resumequeries.ResumeOverviewView, error)
}

type SectionsHandler struct {
	links          ResumeLinkPublicReader
	resumeOverview ResumeOverviewReader
	reader         SectionsReader
}

func NewSectionsHandler(
	links ResumeLinkPublicReader,
	resumeOverview ResumeOverviewReader,
	reader SectionsReader,
) *SectionsHandler {
	return &SectionsHandler{links: links, resumeOverview: resumeOverview, reader: reader}
}

func (h *SectionsHandler) Handle(ctx context.Context, q LinkSectionsQuery) (SectionsStatistics, error) {
	if h == nil || h.reader == nil || h.resumeOverview == nil {
		return SectionsStatistics{}, fmt.Errorf("sections handler dependencies are nil")
	}
	if err := validateAnalyticsPeriod(q.Period); err != nil {
		return SectionsStatistics{}, err
	}
	_, resumeID, err := resolveAnalyticsLink(ctx, h.links, q.UserID, q.LinkID)
	if err != nil {
		return SectionsStatistics{}, err
	}
	overview, err := h.resumeOverview.GetOwnedOverview(ctx, resumequeries.GetResumeQuery{ResumeID: resumeID, UserID: q.UserID})
	if err != nil {
		return SectionsStatistics{}, err
	}
	if len(q.SectionOrder) == 0 {
		q.SectionOrder = currentSectionOrder(overview)
	}
	return h.reader.GetSections(ctx, q)
}

func currentSectionOrder(overview resumequeries.ResumeOverviewView) []string {
	out := []string{"about"}
	seen := map[string]bool{"about": true}
	for _, section := range overview.Sections {
		if section.Enabled && section.ID != "techStack" && !seen[section.ID] {
			out = append(out, section.ID)
			seen[section.ID] = true
		}
	}
	return out
}

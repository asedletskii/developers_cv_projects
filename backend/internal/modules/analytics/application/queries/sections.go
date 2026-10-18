package queries

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/analytics/application/ports"
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

type SectionsHandler struct {
	access ports.AnalyticsAccessReader
	reader SectionsReader
}

func NewSectionsHandler(
	access ports.AnalyticsAccessReader,
	reader SectionsReader,
) *SectionsHandler {
	return &SectionsHandler{access: access, reader: reader}
}

func (h *SectionsHandler) Handle(ctx context.Context, q LinkSectionsQuery) (SectionsStatistics, error) {
	if h == nil || h.reader == nil || h.access == nil {
		return SectionsStatistics{}, fmt.Errorf("sections handler dependencies are nil")
	}
	if err := validateAnalyticsPeriod(q.Period); err != nil {
		return SectionsStatistics{}, err
	}
	resumeID, err := h.access.ResolveLink(ctx, q.LinkID)
	if err != nil {
		return SectionsStatistics{}, err
	}
	sectionOrder, err := h.access.ResumeSectionOrder(ctx, q.UserID, resumeID)
	if err != nil {
		return SectionsStatistics{}, err
	}
	if len(q.SectionOrder) == 0 {
		q.SectionOrder = sectionOrder
	}
	return h.reader.GetSections(ctx, q)
}

package queries

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume_link/application/ports"
)

type TechnologyInterest struct {
	Name       string `json:"name"`
	Selections int64  `json:"selections"`
	Visitors   int64  `json:"visitors"`
}

type TechnologiesStatistics struct {
	Technologies []TechnologyInterest `json:"technologies"`
}

type TechnologiesReader interface {
	GetTechnologies(context.Context, LinkAnalyticsQuery) (TechnologiesStatistics, error)
}

type TechnologiesHandler struct {
	links   ResumeLinkPublicReader
	resumes ports.ResumeOwnershipReader
	reader  TechnologiesReader
}

func NewTechnologiesHandler(
	links ResumeLinkPublicReader,
	resumes ports.ResumeOwnershipReader,
	reader TechnologiesReader,
) *TechnologiesHandler {
	return &TechnologiesHandler{links: links, resumes: resumes, reader: reader}
}

func (h *TechnologiesHandler) Handle(ctx context.Context, q LinkAnalyticsQuery) (TechnologiesStatistics, error) {
	if h == nil || h.reader == nil {
		return TechnologiesStatistics{}, fmt.Errorf("technologies handler dependencies are nil")
	}
	if err := validateAnalyticsPeriod(q.Period); err != nil {
		return TechnologiesStatistics{}, err
	}
	if _, _, err := authorizeAnalyticsLink(ctx, h.links, h.resumes, q.UserID, q.LinkID); err != nil {
		return TechnologiesStatistics{}, err
	}
	return h.reader.GetTechnologies(ctx, q)
}

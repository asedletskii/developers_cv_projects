package queries

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/analytics/application/ports"
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
	access ports.AnalyticsAccessReader
	reader TechnologiesReader
}

func NewTechnologiesHandler(
	access ports.AnalyticsAccessReader,
	reader TechnologiesReader,
) *TechnologiesHandler {
	return &TechnologiesHandler{access: access, reader: reader}
}

func (h *TechnologiesHandler) Handle(ctx context.Context, q LinkAnalyticsQuery) (TechnologiesStatistics, error) {
	if h == nil || h.reader == nil {
		return TechnologiesStatistics{}, fmt.Errorf("technologies handler dependencies are nil")
	}
	if err := validateAnalyticsPeriod(q.Period); err != nil {
		return TechnologiesStatistics{}, err
	}
	if err := authorizeAnalyticsLink(ctx, h.access, q.UserID, q.LinkID); err != nil {
		return TechnologiesStatistics{}, err
	}
	return h.reader.GetTechnologies(ctx, q)
}

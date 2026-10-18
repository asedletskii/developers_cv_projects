package queries

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/analytics/application/ports"
)

type DailyStatisticsReader interface {
	GetDaily(context.Context, LinkAnalyticsQuery) (LinkDaily, error)
}

type DailyStatisticsHandler struct {
	access ports.AnalyticsAccessReader
	reader DailyStatisticsReader
}

func NewDailyStatisticsHandler(
	access ports.AnalyticsAccessReader,
	reader DailyStatisticsReader,
) *DailyStatisticsHandler {
	return &DailyStatisticsHandler{access: access, reader: reader}
}

func (h *DailyStatisticsHandler) Handle(ctx context.Context, q LinkAnalyticsQuery) (LinkDaily, error) {
	if h == nil || h.reader == nil {
		return LinkDaily{}, fmt.Errorf("daily statistics handler dependencies are nil")
	}
	if err := validateAnalyticsPeriod(q.Period); err != nil {
		return LinkDaily{}, err
	}
	if err := authorizeAnalyticsLink(ctx, h.access, q.UserID, q.LinkID); err != nil {
		return LinkDaily{}, err
	}
	return h.reader.GetDaily(ctx, q)
}

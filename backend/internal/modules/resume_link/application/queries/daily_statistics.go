package queries

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume_link/application/ports"
)

type DailyStatisticsReader interface {
	GetDaily(context.Context, LinkAnalyticsQuery) (LinkDaily, error)
}

type DailyStatisticsHandler struct {
	links   ResumeLinkPublicReader
	resumes ports.ResumeOwnershipReader
	reader  DailyStatisticsReader
}

func NewDailyStatisticsHandler(
	links ResumeLinkPublicReader,
	resumes ports.ResumeOwnershipReader,
	reader DailyStatisticsReader,
) *DailyStatisticsHandler {
	return &DailyStatisticsHandler{links: links, resumes: resumes, reader: reader}
}

func (h *DailyStatisticsHandler) Handle(ctx context.Context, q LinkAnalyticsQuery) (LinkDaily, error) {
	if h == nil || h.reader == nil {
		return LinkDaily{}, fmt.Errorf("daily statistics handler dependencies are nil")
	}
	if err := validateAnalyticsPeriod(q.Period); err != nil {
		return LinkDaily{}, err
	}
	if _, _, err := authorizeAnalyticsLink(ctx, h.links, h.resumes, q.UserID, q.LinkID); err != nil {
		return LinkDaily{}, err
	}
	return h.reader.GetDaily(ctx, q)
}

package queries

import (
	"context"
	"fmt"

	analyticscontracts "developers-cv/backend/internal/modules/analytics/contracts"
)

type ResumeStatisticsReader interface {
	GetResumeStatistics(context.Context, ResumeAnalyticsQuery) (ResumeStatistics, error)
}

type ResumeStatisticsHandler struct {
	access analyticscontracts.AnalyticsAccessReader
	reader ResumeStatisticsReader
}

func NewResumeStatisticsHandler(access analyticscontracts.AnalyticsAccessReader, reader ResumeStatisticsReader) *ResumeStatisticsHandler {
	return &ResumeStatisticsHandler{access: access, reader: reader}
}

func (h *ResumeStatisticsHandler) Handle(ctx context.Context, q ResumeAnalyticsQuery) (ResumeStatistics, error) {
	if h == nil || h.reader == nil {
		return ResumeStatistics{}, fmt.Errorf("resume statistics handler dependencies are nil")
	}
	if err := validateAnalyticsPeriod(q.Period); err != nil {
		return ResumeStatistics{}, err
	}
	if err := authorizeAnalyticsResume(ctx, h.access, q.UserID, q.ResumeID); err != nil {
		return ResumeStatistics{}, err
	}
	result, err := h.reader.GetResumeStatistics(ctx, q)
	if err != nil {
		return ResumeStatistics{}, err
	}
	result.NormalizeResponse()
	return result, nil
}

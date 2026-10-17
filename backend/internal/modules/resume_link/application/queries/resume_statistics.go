package queries

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume_link/application/ports"
)

type ResumeStatisticsReader interface {
	GetResumeStatistics(context.Context, ResumeAnalyticsQuery) (ResumeStatistics, error)
}

type ResumeStatisticsHandler struct {
	resumes ports.ResumeOwnershipReader
	reader  ResumeStatisticsReader
}

func NewResumeStatisticsHandler(resumes ports.ResumeOwnershipReader, reader ResumeStatisticsReader) *ResumeStatisticsHandler {
	return &ResumeStatisticsHandler{resumes: resumes, reader: reader}
}

func (h *ResumeStatisticsHandler) Handle(ctx context.Context, q ResumeAnalyticsQuery) (ResumeStatistics, error) {
	if h == nil || h.reader == nil {
		return ResumeStatistics{}, fmt.Errorf("resume statistics handler dependencies are nil")
	}
	if err := validateAnalyticsPeriod(q.Period); err != nil {
		return ResumeStatistics{}, err
	}
	if err := authorizeAnalyticsResume(ctx, h.resumes, q.UserID, q.ResumeID); err != nil {
		return ResumeStatistics{}, err
	}
	result, err := h.reader.GetResumeStatistics(ctx, q)
	if err != nil {
		return ResumeStatistics{}, err
	}
	result.NormalizeResponse()
	return result, nil
}

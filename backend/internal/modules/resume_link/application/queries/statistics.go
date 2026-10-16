package queries

import "developers-cv/backend/internal/modules/resume_link/application/ports"

type GetResumeLinkStatisticsQuery = LinkAnalyticsQuery
type ResumeLinkStatistics = LinkStatistics
type GetResumeLinkStatisticsQueryHandler = AnalyticsReadHandler

func NewGetResumeLinkStatisticsQueryHandler(links ResumeLinkPublicReader, resumes ports.ResumeOwnershipReader, reader AnalyticsReadReader) *AnalyticsReadHandler {
	return NewAnalyticsReadHandler(links, resumes, reader)
}

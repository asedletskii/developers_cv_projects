package queries

import (
	"context"
	"fmt"
	"time"

	"developers-cv/backend/internal/modules/resume_link/application/ports"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type LinkAnalyticsQuery struct {
	UserID sharedvalueobjects.ID
	LinkID sharedvalueobjects.ID
	Period AnalyticsPeriod
}

type ResumeAnalyticsQuery struct {
	UserID   sharedvalueobjects.ID
	ResumeID sharedvalueobjects.ID
	Period   AnalyticsPeriod
}

type LinkStatistics struct {
	Views               int64              `json:"views"`
	Reads               int64              `json:"reads"`
	UniqueVisitors      int64              `json:"uniqueVisitors"`
	ReturningVisitors   int64              `json:"returningVisitors"`
	ContactViews        int64              `json:"contactViews"`
	Downloads           int64              `json:"downloads"`
	FirstViewAt         *time.Time         `json:"firstViewAt"`
	LastViewAt          *time.Time         `json:"lastViewAt"`
	MedianActiveSeconds int64              `json:"medianActiveSeconds"`
	Devices             map[string]int64   `json:"devices"`
	Referrers           []ReferrerStat     `json:"referrers"`
	Filtered            FilteredStatistics `json:"filtered"`
}

type FilteredStatistics struct {
	Bot     int64 `json:"bot"`
	Owner   int64 `json:"owner"`
	Suspect int64 `json:"suspect"`
}

type ReferrerStat struct {
	Host  *string `json:"host"`
	Views int64   `json:"views"`
}

type DailyPoint struct {
	Date           string `json:"date"`
	Views          int64  `json:"views"`
	Reads          int64  `json:"reads"`
	UniqueVisitors int64  `json:"uniqueVisitors"`
}

type LinkDaily struct {
	Days []DailyPoint `json:"days"`
}

type ResumeStatistics struct {
	LinkStatistics
	Links []ResumeLinkStatistic `json:"links"`
}

type ResumeLinkStatistic struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Views      int64      `json:"views"`
	Reads      int64      `json:"reads"`
	LastViewAt *time.Time `json:"lastViewAt"`
}

func (s *LinkStatistics) NormalizeResponse() {
	if s == nil {
		return
	}
	if s.Devices == nil {
		s.Devices = make(map[string]int64, 4)
	}
	for _, device := range []string{"desktop", "mobile", "tablet", "unknown"} {
		if _, ok := s.Devices[device]; !ok {
			s.Devices[device] = 0
		}
	}
	if s.Referrers == nil {
		s.Referrers = make([]ReferrerStat, 0)
	}
}

func (s *ResumeStatistics) NormalizeResponse() {
	if s == nil {
		return
	}
	s.LinkStatistics.NormalizeResponse()
	if s.Links == nil {
		s.Links = make([]ResumeLinkStatistic, 0)
	}
}

type LinkStatisticsReader interface {
	GetLinkStatistics(context.Context, LinkAnalyticsQuery) (LinkStatistics, error)
}

type LinkStatisticsHandler struct {
	links   ResumeLinkPublicReader
	resumes ports.ResumeOwnershipReader
	reader  LinkStatisticsReader
}

func NewLinkStatisticsHandler(
	links ResumeLinkPublicReader,
	resumes ports.ResumeOwnershipReader,
	reader LinkStatisticsReader,
) *LinkStatisticsHandler {
	return &LinkStatisticsHandler{links: links, resumes: resumes, reader: reader}
}

func (h *LinkStatisticsHandler) Handle(ctx context.Context, q LinkAnalyticsQuery) (LinkStatistics, error) {
	if h == nil || h.reader == nil {
		return LinkStatistics{}, fmt.Errorf("link statistics handler dependencies are nil")
	}
	if err := validateAnalyticsPeriod(q.Period); err != nil {
		return LinkStatistics{}, err
	}
	if _, _, err := authorizeAnalyticsLink(ctx, h.links, h.resumes, q.UserID, q.LinkID); err != nil {
		return LinkStatistics{}, err
	}
	result, err := h.reader.GetLinkStatistics(ctx, q)
	if err != nil {
		return LinkStatistics{}, err
	}
	result.NormalizeResponse()
	return result, nil
}

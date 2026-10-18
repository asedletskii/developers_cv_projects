package queries

import (
	"context"
	"testing"
	"time"

	"developers-cv/backend/internal/modules/analytics/application/ports"
	"developers-cv/backend/internal/modules/analytics/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type analyticsAccessReaderStub struct {
	resumeID     sharedvalueobjects.ID
	err          error
	sectionOrder []string
	layout       ports.ResumeAnalyticsLayout
}

func (s *analyticsAccessReaderStub) ResolveLink(context.Context, sharedvalueobjects.ID) (sharedvalueobjects.ID, error) {
	return s.resumeID, s.err
}

func (s *analyticsAccessReaderStub) EnsureResumeOwned(context.Context, sharedvalueobjects.ID, sharedvalueobjects.ID) error {
	return s.err
}

func (s *analyticsAccessReaderStub) ResumeSectionOrder(context.Context, sharedvalueobjects.ID, sharedvalueobjects.ID) ([]string, error) {
	return s.sectionOrder, s.err
}

func (s *analyticsAccessReaderStub) ResumeAnalyticsLayout(context.Context, sharedvalueobjects.ID, sharedvalueobjects.ID) (ports.ResumeAnalyticsLayout, error) {
	return s.layout, s.err
}

type statisticsReaderStub struct {
	result LinkStatistics
	query  LinkAnalyticsQuery
}

func (s *statisticsReaderStub) GetLinkStatistics(_ context.Context, query LinkAnalyticsQuery) (LinkStatistics, error) {
	s.query = query
	return s.result, nil
}

type analyticsUseCaseReaderStub struct {
	calls map[string]int
}

func newAnalyticsUseCaseReaderStub() *analyticsUseCaseReaderStub {
	return &analyticsUseCaseReaderStub{calls: make(map[string]int)}
}

func (s *analyticsUseCaseReaderStub) GetLinkStatistics(context.Context, LinkAnalyticsQuery) (LinkStatistics, error) {
	s.calls["statistics"]++
	return LinkStatistics{Views: 7}, nil
}

func (s *analyticsUseCaseReaderStub) GetDaily(context.Context, LinkAnalyticsQuery) (LinkDaily, error) {
	s.calls["daily"]++
	return LinkDaily{Days: []DailyPoint{}}, nil
}

func (s *analyticsUseCaseReaderStub) GetVisits(context.Context, VisitsReadQuery) (LinkVisitsPage, error) {
	s.calls["visits"]++
	return LinkVisitsPage{Items: []LinkVisit{}}, nil
}

func (s *analyticsUseCaseReaderStub) GetSections(_ context.Context, query LinkSectionsQuery) (SectionsStatistics, error) {
	s.calls["sections"]++
	return SectionsStatistics{Sections: []SectionStatistics{{Section: query.SectionOrder[0]}}}, nil
}

func (s *analyticsUseCaseReaderStub) GetTechnologies(context.Context, LinkAnalyticsQuery) (TechnologiesStatistics, error) {
	s.calls["technologies"]++
	return TechnologiesStatistics{Technologies: []TechnologyInterest{}}, nil
}

func (s *analyticsUseCaseReaderStub) GetResumeStatistics(context.Context, ResumeAnalyticsQuery) (ResumeStatistics, error) {
	s.calls["resume_statistics"]++
	return ResumeStatistics{Links: []ResumeLinkStatistic{}}, nil
}

type clickMapReaderStub struct{ query LinkClicksQuery }

func (s *clickMapReaderStub) GetClicks(_ context.Context, query LinkClicksQuery) (ClickGrid, error) {
	s.query = query
	return ClickGrid{Section: query.Section, Grid: ClickGridSize, Cells: []ClickCell{{X: 1, Y: 2, Count: 3}}, Total: 3}, nil
}

type layoutSnapshotBuilderStub struct {
	snapshot domain.LayoutSnapshot
	err      error
}

func (s layoutSnapshotBuilderStub) Build(ports.ResumeAnalyticsLayout) (domain.LayoutSnapshot, error) {
	return s.snapshot, s.err
}

func validAnalyticsTestPeriod(t *testing.T) AnalyticsPeriod {
	t.Helper()
	period, err := NewAnalyticsPeriod("2026-10-01", "2026-10-31", "UTC", time.Date(2026, 10, 17, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("create analytics period: %v", err)
	}
	return period
}

func mustLocation(t *testing.T, name string) *time.Location {
	t.Helper()
	location, err := time.LoadLocation(name)
	if err != nil {
		t.Fatalf("load time zone %q: %v", name, err)
	}
	return location
}

func mustLimit(t *testing.T, value int) int {
	t.Helper()
	limit, err := NormalizeVisitLimit(value)
	if err != nil {
		t.Fatalf("normalize visit limit: %v", err)
	}
	return limit
}

package queries

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	analyticsports "developers-cv/backend/internal/modules/analytics/application/ports"
	analyticsqueries "developers-cv/backend/internal/modules/analytics/application/queries"
	analyticsdomain "developers-cv/backend/internal/modules/analytics/domain"
	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ownershipReaderStub struct {
	err            error
	viewID         string
	query          resumequeries.GetResumeQuery
	calls          int
	getCalls       int
	ownershipCalls int
}

func (s *ownershipReaderStub) GetOwned(_ context.Context, query resumequeries.GetResumeQuery) (resumequeries.ResumeView, error) {
	s.calls++
	s.getCalls++
	s.query = query
	if s.err != nil {
		return resumequeries.ResumeView{}, s.err
	}
	viewID := s.viewID
	if viewID == "" {
		viewID = query.ResumeID.String()
	}
	return resumequeries.ResumeView{ID: viewID}, nil
}
func (s *ownershipReaderStub) GetOwnedOverview(_ context.Context, query resumequeries.GetResumeQuery) (resumequeries.ResumeOverviewView, error) {
	s.calls++
	s.getCalls++
	s.query = query
	if s.err != nil {
		return resumequeries.ResumeOverviewView{}, s.err
	}
	return resumequeries.ResumeOverviewView{}, nil
}
func (s *ownershipReaderStub) OwnsResume(_ context.Context, query resumequeries.GetResumeQuery) (bool, error) {
	s.calls++
	s.ownershipCalls++
	s.query = query
	if s.err != nil {
		return false, s.err
	}
	return s.viewID == "" || s.viewID == query.ResumeID.String(), nil
}

type resumeLinkReaderStub struct {
	view  ResumeLinkPublicView
	err   error
	query PublicResumeLinkQuery
}

func (s *resumeLinkReaderStub) GetByID(_ context.Context, query PublicResumeLinkQuery) (ResumeLinkPublicView, error) {
	s.query = query
	return s.view, s.err
}

type listReaderStub struct {
	items []ResumeLinkListItem
	err   error
	query ListResumeLinksQuery
}

func (s *listReaderStub) ListByResumeID(_ context.Context, query ListResumeLinksQuery) ([]ResumeLinkListItem, error) {
	s.query = query
	return s.items, s.err
}

type publicResumeQueryStub struct {
	view  resumequeries.PublicResumeView
	err   error
	query resumequeries.PublicResumeQuery
}

func (s *publicResumeQueryStub) Handle(_ context.Context, query resumequeries.PublicResumeQuery) (resumequeries.PublicResumeView, error) {
	s.query = query
	return s.view, s.err
}

type publicVisitCreatorStub struct {
	identity analyticsports.VisitorIdentityResult
	err      error
	input    analyticsports.CreatePublicVisitInput
	calls    int
}

func (s *publicVisitCreatorStub) Handle(_ context.Context, input analyticsports.CreatePublicVisitInput) (*analyticsdomain.View, analyticsports.VisitorIdentityResult, error) {
	s.calls++
	s.input = input
	return nil, s.identity, s.err
}

type contactViewResumeQueryStub struct {
	view  resumequeries.ContactView
	err   error
	query resumequeries.ContactViewQuery
}

func (s *contactViewResumeQueryStub) Handle(_ context.Context, query resumequeries.ContactViewQuery) (resumequeries.ContactView, error) {
	s.query = query
	return s.view, s.err
}

type contactViewResumeDownloadQueryStub struct {
	view  resumequeries.DownloadContactsView
	err   error
	query resumequeries.DownloadContactsQuery
}

func (s *contactViewResumeDownloadQueryStub) Handle(_ context.Context, query resumequeries.DownloadContactsQuery) (resumequeries.DownloadContactsView, error) {
	s.query = query
	return s.view, s.err
}

type contactViewAnalyticsRecorderStub struct {
	err       error
	calls     int
	initiator string
	viewID    sharedvalueobjects.ID
	linkID    sharedvalueobjects.ID
}

func (s *contactViewAnalyticsRecorderStub) RecordContactView(_ context.Context, viewID, linkID sharedvalueobjects.ID, initiator string) error {
	s.calls++
	s.initiator = "contact"
	s.viewID = viewID
	s.linkID = linkID
	return s.err
}
func (s *contactViewAnalyticsRecorderStub) RecordDownload(_ context.Context, viewID, linkID sharedvalueobjects.ID) error {
	s.calls++
	s.initiator = "download"
	s.viewID = viewID
	s.linkID = linkID
	return s.err
}

type statisticsReaderStub struct {
	result analyticsqueries.LinkStatistics
	err    error
	query  analyticsqueries.LinkAnalyticsQuery
}

func (s *statisticsReaderStub) GetLinkStatistics(_ context.Context, query analyticsqueries.LinkAnalyticsQuery) (analyticsqueries.LinkStatistics, error) {
	s.query = query
	return s.result, s.err
}
func (s *statisticsReaderStub) GetDaily(context.Context, analyticsqueries.LinkAnalyticsQuery) (analyticsqueries.LinkDaily, error) {
	return analyticsqueries.LinkDaily{}, s.err
}
func (s *statisticsReaderStub) GetVisits(context.Context, analyticsqueries.VisitsReadQuery) (analyticsqueries.LinkVisitsPage, error) {
	return analyticsqueries.LinkVisitsPage{}, s.err
}
func (s *statisticsReaderStub) GetSections(context.Context, analyticsqueries.LinkSectionsQuery) (analyticsqueries.SectionsStatistics, error) {
	return analyticsqueries.SectionsStatistics{}, s.err
}
func (s *statisticsReaderStub) GetTechnologies(context.Context, analyticsqueries.LinkAnalyticsQuery) (analyticsqueries.TechnologiesStatistics, error) {
	return analyticsqueries.TechnologiesStatistics{}, s.err
}
func (s *statisticsReaderStub) GetResumeStatistics(context.Context, analyticsqueries.ResumeAnalyticsQuery) (analyticsqueries.ResumeStatistics, error) {
	return analyticsqueries.ResumeStatistics{}, s.err
}

type analyticsUseCaseReaderStub struct {
	calls map[string]int
}

func newAnalyticsUseCaseReaderStub() *analyticsUseCaseReaderStub {
	return &analyticsUseCaseReaderStub{calls: make(map[string]int)}
}
func (s *analyticsUseCaseReaderStub) GetLinkStatistics(context.Context, analyticsqueries.LinkAnalyticsQuery) (analyticsqueries.LinkStatistics, error) {
	s.calls["statistics"]++
	return analyticsqueries.LinkStatistics{Views: 7}, nil
}
func (s *analyticsUseCaseReaderStub) GetDaily(context.Context, analyticsqueries.LinkAnalyticsQuery) (analyticsqueries.LinkDaily, error) {
	s.calls["daily"]++
	return analyticsqueries.LinkDaily{Days: []analyticsqueries.DailyPoint{}}, nil
}
func (s *analyticsUseCaseReaderStub) GetVisits(context.Context, analyticsqueries.VisitsReadQuery) (analyticsqueries.LinkVisitsPage, error) {
	s.calls["visits"]++
	return analyticsqueries.LinkVisitsPage{Items: []analyticsqueries.LinkVisit{}}, nil
}
func (s *analyticsUseCaseReaderStub) GetSections(_ context.Context, query analyticsqueries.LinkSectionsQuery) (analyticsqueries.SectionsStatistics, error) {
	s.calls["sections"]++
	return analyticsqueries.SectionsStatistics{Sections: []analyticsqueries.SectionStatistics{{Section: query.SectionOrder[0]}}}, nil
}
func (s *analyticsUseCaseReaderStub) GetTechnologies(context.Context, analyticsqueries.LinkAnalyticsQuery) (analyticsqueries.TechnologiesStatistics, error) {
	s.calls["technologies"]++
	return analyticsqueries.TechnologiesStatistics{Technologies: []analyticsqueries.TechnologyInterest{}}, nil
}
func (s *analyticsUseCaseReaderStub) GetResumeStatistics(context.Context, analyticsqueries.ResumeAnalyticsQuery) (analyticsqueries.ResumeStatistics, error) {
	s.calls["resume_statistics"]++
	return analyticsqueries.ResumeStatistics{Links: []analyticsqueries.ResumeLinkStatistic{}}, nil
}

type layoutSnapshotBuilderStub struct {
	snapshot analyticsdomain.LayoutSnapshot
	err      error
}

func (s layoutSnapshotBuilderStub) Build(analyticsports.ResumeAnalyticsLayout) (analyticsdomain.LayoutSnapshot, error) {
	return s.snapshot, s.err
}

type clickMapReaderStub struct {
	query analyticsqueries.LinkClicksQuery
}

func (s *clickMapReaderStub) GetClicks(_ context.Context, query analyticsqueries.LinkClicksQuery) (analyticsqueries.ClickGrid, error) {
	s.query = query
	return analyticsqueries.ClickGrid{Section: query.Section, Grid: analyticsqueries.ClickGridSize, Cells: []analyticsqueries.ClickCell{{X: 1, Y: 2, Count: 3}}, Total: 3}, nil
}

func newPublicAnalyticsTestHandler(
	owner ports.ResumeOwnershipReader,
	visit analyticsports.PublicVisitCreator,
	reporter PublicAnalyticsErrorReporter,
) (*PublicResumeLinkQueryHandler, PublicResumeLinkQuery, *publicResumeQueryStub) {
	linkID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	links := &resumeLinkReaderStub{view: ResumeLinkPublicView{
		ID: linkID.String(), ResumeID: resumeID.String(), Status: string(domain.ResumeLinkStatusActive),
	}}
	resume := &publicResumeQueryStub{view: resumequeries.PublicResumeView{Title: "Public CV"}}
	handler := NewPublicResumeLinkQueryHandlerWithVisit(NewActivePublicResumeLinkResolver(links), resume, owner, visit, reporter)
	query := PublicResumeLinkQuery{
		LinkID: linkID,
		Visit: &analyticsports.CreatePublicVisitInput{
			LinkID:  linkID,
			Visitor: analyticsports.PublicVisitInput{OwnerID: sharedvalueobjects.NewID()},
		},
	}
	return handler, query, resume
}

func validAnalyticsTestPeriod(t *testing.T) analyticsqueries.AnalyticsPeriod {
	t.Helper()
	period, err := analyticsqueries.NewAnalyticsPeriod("2026-10-01", "2026-10-01", "UTC", time.Date(2026, 10, 2, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	return period
}

func mustLimit(t *testing.T, value int) int {
	t.Helper()
	result, err := analyticsqueries.NormalizeVisitLimit(value)
	require.NoError(t, err)
	return result
}

func mustLocation(t *testing.T, name string) *time.Location {
	t.Helper()
	location, err := time.LoadLocation(name)
	require.NoError(t, err)
	return location
}

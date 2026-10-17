package queries

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

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
	identity ports.VisitorIdentityResult
	err      error
	input    ports.CreatePublicVisitInput
	calls    int
}

func (s *publicVisitCreatorStub) Handle(_ context.Context, input ports.CreatePublicVisitInput) (*domain.View, ports.VisitorIdentityResult, error) {
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
	result LinkStatistics
	err    error
	query  LinkAnalyticsQuery
}

func (s *statisticsReaderStub) GetLinkStatistics(_ context.Context, query LinkAnalyticsQuery) (LinkStatistics, error) {
	s.query = query
	return s.result, s.err
}
func (s *statisticsReaderStub) GetDaily(context.Context, LinkAnalyticsQuery) (LinkDaily, error) {
	return LinkDaily{}, s.err
}
func (s *statisticsReaderStub) GetVisits(context.Context, VisitsReadQuery) (LinkVisitsPage, error) {
	return LinkVisitsPage{}, s.err
}
func (s *statisticsReaderStub) GetSections(context.Context, LinkSectionsQuery) (SectionsStatistics, error) {
	return SectionsStatistics{}, s.err
}
func (s *statisticsReaderStub) GetTechnologies(context.Context, LinkAnalyticsQuery) (TechnologiesStatistics, error) {
	return TechnologiesStatistics{}, s.err
}
func (s *statisticsReaderStub) GetResumeStatistics(context.Context, ResumeAnalyticsQuery) (ResumeStatistics, error) {
	return ResumeStatistics{}, s.err
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

type layoutSnapshotBuilderStub struct {
	snapshot domain.LayoutSnapshot
	err      error
}

func (s layoutSnapshotBuilderStub) Build(resumequeries.PublicResumeView) (domain.LayoutSnapshot, error) {
	return s.snapshot, s.err
}

type clickMapReaderStub struct{ query LinkClicksQuery }

func (s *clickMapReaderStub) GetClicks(_ context.Context, query LinkClicksQuery) (ClickGrid, error) {
	s.query = query
	return ClickGrid{Section: query.Section, Grid: ClickGridSize, Cells: []ClickCell{{X: 1, Y: 2, Count: 3}}, Total: 3}, nil
}

func newPublicAnalyticsTestHandler(
	owner ports.ResumeOwnershipReader,
	visit ports.PublicVisitCreator,
	reporter PublicAnalyticsErrorReporter,
) (*PublicResumeLinkQueryHandler, PublicResumeLinkQuery, *publicResumeQueryStub) {
	linkID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	links := &resumeLinkReaderStub{view: ResumeLinkPublicView{
		ID: linkID.String(), ResumeID: resumeID.String(), Status: string(domain.ResumeLinkStatusActive),
	}}
	resume := &publicResumeQueryStub{view: resumequeries.PublicResumeView{Title: "Public CV"}}
	handler := NewPublicResumeLinkQueryHandlerWithVisit(links, resume, owner, visit, reporter)
	query := PublicResumeLinkQuery{
		LinkID: linkID,
		Visit: &ports.CreatePublicVisitInput{
			LinkID:  linkID,
			Visitor: ports.PublicVisitInput{OwnerID: sharedvalueobjects.NewID()},
		},
	}
	return handler, query, resume
}

func validAnalyticsTestPeriod(t *testing.T) AnalyticsPeriod {
	t.Helper()
	period, err := NewAnalyticsPeriod("2026-10-01", "2026-10-01", "UTC", time.Date(2026, 10, 2, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	return period
}

func mustLimit(t *testing.T, value int) int {
	t.Helper()
	result, err := NormalizeVisitLimit(value)
	require.NoError(t, err)
	return result
}

func mustLocation(t *testing.T, name string) *time.Location {
	t.Helper()
	location, err := time.LoadLocation(name)
	require.NoError(t, err)
	return location
}

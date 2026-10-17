package http

import (
	"context"
	"errors"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/platform/http/middleware"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type statisticsHTTPLinkReader struct {
	view queries.ResumeLinkPublicView
	err  error
}

func (s statisticsHTTPLinkReader) GetByID(context.Context, queries.PublicResumeLinkQuery) (queries.ResumeLinkPublicView, error) {
	return s.view, s.err
}

type statisticsHTTPOwnershipReader struct {
	err error
}

func (s statisticsHTTPOwnershipReader) GetOwned(context.Context, resumequeries.GetResumeQuery) (resumequeries.ResumeView, error) {
	if s.err != nil {
		return resumequeries.ResumeView{}, s.err
	}
	return resumequeries.ResumeView{ID: "resume"}, nil
}

func (s statisticsHTTPOwnershipReader) GetOwnedOverview(_ context.Context, _ resumequeries.GetResumeQuery) (resumequeries.ResumeOverviewView, error) {
	if s.err != nil {
		return resumequeries.ResumeOverviewView{}, s.err
	}
	return resumequeries.ResumeOverviewView{}, nil
}

func (s statisticsHTTPOwnershipReader) OwnsResume(context.Context, resumequeries.GetResumeQuery) (bool, error) {
	return s.err == nil, s.err
}

type statisticsHTTPReader struct {
	result queries.LinkStatistics
	err    error
}

func (s statisticsHTTPReader) GetLinkStatistics(context.Context, queries.LinkAnalyticsQuery) (queries.LinkStatistics, error) {
	return s.result, s.err
}

func (s statisticsHTTPReader) GetDaily(context.Context, queries.LinkAnalyticsQuery) (queries.LinkDaily, error) {
	return queries.LinkDaily{}, s.err
}
func (s statisticsHTTPReader) GetVisits(context.Context, queries.VisitsReadQuery) (queries.LinkVisitsPage, error) {
	return queries.LinkVisitsPage{}, s.err
}
func (s statisticsHTTPReader) GetSections(context.Context, queries.LinkSectionsQuery) (queries.SectionsStatistics, error) {
	return queries.SectionsStatistics{}, s.err
}
func (s statisticsHTTPReader) GetTechnologies(context.Context, queries.LinkAnalyticsQuery) (queries.TechnologiesStatistics, error) {
	return queries.TechnologiesStatistics{}, s.err
}
func (s statisticsHTTPReader) GetResumeStatistics(context.Context, queries.ResumeAnalyticsQuery) (queries.ResumeStatistics, error) {
	return queries.ResumeStatistics{}, s.err
}

func newStatisticsHTTPHandler(
	linkReader queries.ResumeLinkPublicReader,
	ownership ports.ResumeOwnershipReader,
	resumeOverview queries.ResumeOverviewReader,
	reader statisticsHTTPReader,
) *Handler {
	return NewHandler(HandlerDependencies{
		AnalyticsStatistics:       queries.NewLinkStatisticsHandler(linkReader, ownership, reader),
		AnalyticsDaily:            queries.NewDailyStatisticsHandler(linkReader, ownership, reader),
		AnalyticsVisits:           queries.NewVisitsHandler(linkReader, ownership, reader),
		AnalyticsSections:         queries.NewSectionsHandler(linkReader, resumeOverview, reader),
		AnalyticsTechnologies:     queries.NewTechnologiesHandler(linkReader, ownership, reader),
		AnalyticsResumeStatistics: queries.NewResumeStatisticsHandler(ownership, reader),
	})
}

func newStatisticsHTTPServer(handler *Handler, userID uuid.UUID) *gin.Engine {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	RegisterRoutes(handler, engine.Group("/api/v1"), engine.Group("/api/v1"), engine.Group("/api/v1/lk"), nil, nil)
	return engine
}

func TestStatisticsHandlerReturnsStatistics(t *testing.T) {
	userID := uuid.New()
	linkID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	engine := newStatisticsHTTPServer(newStatisticsHTTPHandler(
		statisticsHTTPLinkReader{view: queries.ResumeLinkPublicView{ResumeID: resumeID.String()}},
		statisticsHTTPOwnershipReader{},
		statisticsHTTPOwnershipReader{},
		statisticsHTTPReader{result: queries.LinkStatistics{Views: 125, ContactViews: 17}},
	), userID)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/link/"+linkID.String()+"/statistics", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"views":125,"reads":0,"uniqueVisitors":0,"returningVisitors":0,"contactViews":17,"downloads":0,"firstViewAt":null,"lastViewAt":null,"medianActiveSeconds":0,"devices":{"desktop":0,"mobile":0,"tablet":0,"unknown":0},"referrers":[],"filtered":{"bot":0,"owner":0,"suspect":0}}`, w.Body.String())
}

func TestStatisticsHandlerReturnsZeroState(t *testing.T) {
	engine := newStatisticsHTTPServer(newStatisticsHTTPHandler(
		statisticsHTTPLinkReader{view: queries.ResumeLinkPublicView{ResumeID: sharedvalueobjects.NewID().String()}},
		statisticsHTTPOwnershipReader{},
		statisticsHTTPOwnershipReader{},
		statisticsHTTPReader{result: queries.LinkStatistics{}},
	), uuid.New())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/link/"+sharedvalueobjects.NewID().String()+"/statistics", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"views":0,"reads":0,"uniqueVisitors":0,"returningVisitors":0,"contactViews":0,"downloads":0,"firstViewAt":null,"lastViewAt":null,"medianActiveSeconds":0,"devices":{"desktop":0,"mobile":0,"tablet":0,"unknown":0},"referrers":[],"filtered":{"bot":0,"owner":0,"suspect":0}}`, w.Body.String())
}

func TestStatisticsHandlerRejectsForeignLink(t *testing.T) {
	engine := newStatisticsHTTPServer(newStatisticsHTTPHandler(
		statisticsHTTPLinkReader{view: queries.ResumeLinkPublicView{ResumeID: sharedvalueobjects.NewID().String()}},
		statisticsHTTPOwnershipReader{err: resumequeries.ErrResumeNotFound},
		statisticsHTTPOwnershipReader{err: resumequeries.ErrResumeNotFound},
		statisticsHTTPReader{},
	), uuid.New())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/link/"+sharedvalueobjects.NewID().String()+"/statistics", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestStatisticsHandlerReturnsNotFoundForMissingLink(t *testing.T) {
	engine := newStatisticsHTTPServer(newStatisticsHTTPHandler(
		statisticsHTTPLinkReader{err: queries.ErrResumeLinkNotFound},
		statisticsHTTPOwnershipReader{},
		statisticsHTTPOwnershipReader{},
		statisticsHTTPReader{},
	), uuid.New())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/link/"+sharedvalueobjects.NewID().String()+"/statistics", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestStatisticsHandlerReturnsBadRequestForInvalidUUID(t *testing.T) {
	engine := newStatisticsHTTPServer(newStatisticsHTTPHandler(
		statisticsHTTPLinkReader{},
		statisticsHTTPOwnershipReader{},
		statisticsHTTPOwnershipReader{},
		statisticsHTTPReader{},
	), uuid.New())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/link/not-a-uuid/statistics", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStatisticsHandlerReturnsInternalErrorForReaderFailure(t *testing.T) {
	engine := newStatisticsHTTPServer(newStatisticsHTTPHandler(
		statisticsHTTPLinkReader{view: queries.ResumeLinkPublicView{ResumeID: sharedvalueobjects.NewID().String()}},
		statisticsHTTPOwnershipReader{},
		statisticsHTTPOwnershipReader{},
		statisticsHTTPReader{err: errors.New("db down")},
	), uuid.New())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/link/"+sharedvalueobjects.NewID().String()+"/statistics", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAnalyticsReadRejectsInvalidPeriod(t *testing.T) {
	engine := newStatisticsHTTPServer(newStatisticsHTTPHandler(
		statisticsHTTPLinkReader{view: queries.ResumeLinkPublicView{ResumeID: sharedvalueobjects.NewID().String()}},
		statisticsHTTPOwnershipReader{},
		statisticsHTTPOwnershipReader{},
		statisticsHTTPReader{},
	), uuid.New())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/link/"+sharedvalueobjects.NewID().String()+"/statistics/daily?from=2026-02-02&to=2026-02-01&tz=UTC", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAnalyticsReadDailyRouteReturnsZeroDays(t *testing.T) {
	engine := newStatisticsHTTPServer(newStatisticsHTTPHandler(
		statisticsHTTPLinkReader{view: queries.ResumeLinkPublicView{ResumeID: sharedvalueobjects.NewID().String()}},
		statisticsHTTPOwnershipReader{},
		statisticsHTTPOwnershipReader{},
		statisticsHTTPReader{},
	), uuid.New())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/link/"+sharedvalueobjects.NewID().String()+"/statistics/daily?from=2026-01-01&to=2026-01-02&tz=UTC", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

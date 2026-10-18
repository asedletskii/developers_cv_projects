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

	analyticsqueries "developers-cv/backend/internal/modules/analytics/application/queries"
	analyticshttp "developers-cv/backend/internal/modules/analytics/transport/http"
	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/platform/http/middleware"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type statisticsHTTPReader struct {
	result analyticsqueries.LinkStatistics
	err    error
}

func (s statisticsHTTPReader) GetLinkStatistics(context.Context, analyticsqueries.LinkAnalyticsQuery) (analyticsqueries.LinkStatistics, error) {
	return s.result, s.err
}

func (s statisticsHTTPReader) GetDaily(context.Context, analyticsqueries.LinkAnalyticsQuery) (analyticsqueries.LinkDaily, error) {
	return analyticsqueries.LinkDaily{}, s.err
}
func (s statisticsHTTPReader) GetVisits(context.Context, analyticsqueries.VisitsReadQuery) (analyticsqueries.LinkVisitsPage, error) {
	return analyticsqueries.LinkVisitsPage{}, s.err
}
func (s statisticsHTTPReader) GetSections(context.Context, analyticsqueries.LinkSectionsQuery) (analyticsqueries.SectionsStatistics, error) {
	return analyticsqueries.SectionsStatistics{}, s.err
}
func (s statisticsHTTPReader) GetTechnologies(context.Context, analyticsqueries.LinkAnalyticsQuery) (analyticsqueries.TechnologiesStatistics, error) {
	return analyticsqueries.TechnologiesStatistics{}, s.err
}
func (s statisticsHTTPReader) GetResumeStatistics(context.Context, analyticsqueries.ResumeAnalyticsQuery) (analyticsqueries.ResumeStatistics, error) {
	return analyticsqueries.ResumeStatistics{}, s.err
}

func newStatisticsHTTPHandler(
	linkReader queries.ResumeLinkPublicReader,
	ownership ports.ResumeOwnershipReader,
	resumeOverview queries.AnalyticsResumeOverviewReader,
	reader statisticsHTTPReader,
) *analyticshttp.Handler {
	access := queries.NewAnalyticsAccessReader(linkReader, ownership, resumeOverview, nil)
	return analyticshttp.NewHandler(analyticshttp.HandlerDependencies{
		Statistics:       analyticsqueries.NewLinkStatisticsHandler(access, reader),
		Daily:            analyticsqueries.NewDailyStatisticsHandler(access, reader),
		Visits:           analyticsqueries.NewVisitsHandler(access, reader),
		Sections:         analyticsqueries.NewSectionsHandler(access, reader),
		Technologies:     analyticsqueries.NewTechnologiesHandler(access, reader),
		ResumeStatistics: analyticsqueries.NewResumeStatisticsHandler(access, reader),
	})
}

func newStatisticsHTTPServer(handler *analyticshttp.Handler, userID uuid.UUID) *gin.Engine {
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) {
		return userID, nil
	}))
	api := engine.Group("/api/v1")
	analyticshttp.RegisterRoutes(handler, api, api.Group("/lk"), nil)
	return engine
}

func TestStatisticsHandlerReturnsStatistics(t *testing.T) {
	userID := uuid.New()
	linkID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	engine := newStatisticsHTTPServer(newStatisticsHTTPHandler(
		resumeLinkPublicReaderStub{view: queries.ResumeLinkPublicView{ResumeID: resumeID.String()}},
		resumeOwnershipReaderStub{},
		resumeOwnershipReaderStub{},
		statisticsHTTPReader{result: analyticsqueries.LinkStatistics{Views: 125, ContactViews: 17}},
	), userID)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/link/"+linkID.String()+"/statistics", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"views":125,"reads":0,"uniqueVisitors":0,"returningVisitors":0,"contactViews":17,"downloads":0,"firstViewAt":null,"lastViewAt":null,"medianActiveSeconds":0,"devices":{"desktop":0,"mobile":0,"tablet":0,"unknown":0},"referrers":[],"filtered":{"bot":0,"owner":0,"suspect":0}}`, w.Body.String())
}

func TestStatisticsHandlerReturnsZeroState(t *testing.T) {
	engine := newStatisticsHTTPServer(newStatisticsHTTPHandler(
		resumeLinkPublicReaderStub{view: queries.ResumeLinkPublicView{ResumeID: sharedvalueobjects.NewID().String()}},
		resumeOwnershipReaderStub{},
		resumeOwnershipReaderStub{},
		statisticsHTTPReader{result: analyticsqueries.LinkStatistics{}},
	), uuid.New())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/link/"+sharedvalueobjects.NewID().String()+"/statistics", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"views":0,"reads":0,"uniqueVisitors":0,"returningVisitors":0,"contactViews":0,"downloads":0,"firstViewAt":null,"lastViewAt":null,"medianActiveSeconds":0,"devices":{"desktop":0,"mobile":0,"tablet":0,"unknown":0},"referrers":[],"filtered":{"bot":0,"owner":0,"suspect":0}}`, w.Body.String())
}

func TestStatisticsHandlerRejectsForeignLink(t *testing.T) {
	engine := newStatisticsHTTPServer(newStatisticsHTTPHandler(
		resumeLinkPublicReaderStub{view: queries.ResumeLinkPublicView{ResumeID: sharedvalueobjects.NewID().String()}},
		resumeOwnershipReaderStub{err: resumequeries.ErrResumeNotFound},
		resumeOwnershipReaderStub{err: resumequeries.ErrResumeNotFound},
		statisticsHTTPReader{},
	), uuid.New())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/link/"+sharedvalueobjects.NewID().String()+"/statistics", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestStatisticsHandlerReturnsNotFoundForMissingLink(t *testing.T) {
	engine := newStatisticsHTTPServer(newStatisticsHTTPHandler(
		resumeLinkPublicReaderStub{err: queries.ErrResumeLinkNotFound},
		resumeOwnershipReaderStub{},
		resumeOwnershipReaderStub{},
		statisticsHTTPReader{},
	), uuid.New())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/link/"+sharedvalueobjects.NewID().String()+"/statistics", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestStatisticsHandlerReturnsBadRequestForInvalidUUID(t *testing.T) {
	engine := newStatisticsHTTPServer(newStatisticsHTTPHandler(
		resumeLinkPublicReaderStub{},
		resumeOwnershipReaderStub{},
		resumeOwnershipReaderStub{},
		statisticsHTTPReader{},
	), uuid.New())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/link/not-a-uuid/statistics", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStatisticsHandlerReturnsInternalErrorForReaderFailure(t *testing.T) {
	engine := newStatisticsHTTPServer(newStatisticsHTTPHandler(
		resumeLinkPublicReaderStub{view: queries.ResumeLinkPublicView{ResumeID: sharedvalueobjects.NewID().String()}},
		resumeOwnershipReaderStub{},
		resumeOwnershipReaderStub{},
		statisticsHTTPReader{err: errors.New("db down")},
	), uuid.New())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/link/"+sharedvalueobjects.NewID().String()+"/statistics", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAnalyticsReadRejectsInvalidPeriod(t *testing.T) {
	engine := newStatisticsHTTPServer(newStatisticsHTTPHandler(
		resumeLinkPublicReaderStub{view: queries.ResumeLinkPublicView{ResumeID: sharedvalueobjects.NewID().String()}},
		resumeOwnershipReaderStub{},
		resumeOwnershipReaderStub{},
		statisticsHTTPReader{},
	), uuid.New())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/link/"+sharedvalueobjects.NewID().String()+"/statistics/daily?from=2026-02-02&to=2026-02-01&tz=UTC", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAnalyticsReadDailyRouteReturnsZeroDays(t *testing.T) {
	engine := newStatisticsHTTPServer(newStatisticsHTTPHandler(
		resumeLinkPublicReaderStub{view: queries.ResumeLinkPublicView{ResumeID: sharedvalueobjects.NewID().String()}},
		resumeOwnershipReaderStub{},
		resumeOwnershipReaderStub{},
		statisticsHTTPReader{},
	), uuid.New())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/link/"+sharedvalueobjects.NewID().String()+"/statistics/daily?from=2026-01-01&to=2026-01-02&tz=UTC", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

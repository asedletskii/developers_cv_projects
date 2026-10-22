package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	analyticsqueries "developers-cv/backend/internal/modules/analytics/application/queries"
	analyticscontracts "developers-cv/backend/internal/modules/analytics/contracts"
	analyticsdomain "developers-cv/backend/internal/modules/analytics/domain"
	analyticshttp "developers-cv/backend/internal/modules/analytics/transport/http"
	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/infrastructure/analyticsaccess"
	"developers-cv/backend/internal/platform/http/middleware"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type clickHTTPPublicStub struct{}

func (clickHTTPPublicStub) Handle(context.Context, resumecontracts.PublicResumeQuery) (resumecontracts.PublicResumeView, error) {
	return resumecontracts.PublicResumeView{}, nil
}

type clickHTTPLayoutStub struct {
	snapshot analyticsdomain.LayoutSnapshot
}

func (s clickHTTPLayoutStub) Build(analyticscontracts.ResumeAnalyticsLayout) (analyticsdomain.LayoutSnapshot, error) {
	return s.snapshot, nil
}

type clickHTTPReaderStub struct{}

func (clickHTTPReaderStub) GetClicks(context.Context, analyticsqueries.LinkClicksQuery) (analyticsqueries.ClickGrid, error) {
	return analyticsqueries.ClickGrid{Section: "about", Grid: 20, Cells: []analyticsqueries.ClickCell{{X: 0, Y: 1, Count: 2}}, Total: 2}, nil
}

func newClickHTTPHandler(t *testing.T) *analyticshttp.Handler {
	t.Helper()
	layout, err := analyticsdomain.NewLayoutSnapshot(
		[]analyticsdomain.AnalyticsSectionID{analyticsdomain.AnalyticsSectionAbout},
		map[analyticsdomain.AnalyticsSectionID]string{analyticsdomain.AnalyticsSectionAbout: "0123456789abcdef"},
	)
	require.NoError(t, err)
	access := analyticsaccess.NewReader(
		resumeLinkPublicReaderStub{view: queries.ResumeLinkPublicView{ResumeID: sharedvalueobjects.NewID().String()}},
		resumeOwnershipReaderStub{},
		nil,
		nil,
		clickHTTPPublicStub{},
	)
	clicks := analyticsqueries.NewAnalyticsClickQueryHandler(
		access,
		clickHTTPLayoutStub{snapshot: layout},
		clickHTTPReaderStub{},
	)
	return analyticshttp.NewHandler(analyticshttp.HandlerDependencies{Clicks: clicks})
}

func TestClicksHandlerReturnsClickGrid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	linkID := sharedvalueobjects.NewID()
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) { return userID, nil }))
	api := engine.Group("/api/v1")
	analyticshttp.RegisterRoutes(newClickHTTPHandler(t), api, api.Group("/lk"), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/link/"+linkID.String()+"/clicks?section=about&breakpoint=md", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"section":"about","grid":20,"cells":[{"x":0,"y":1,"count":2}],"total":2,"staleTotal":0}`, w.Body.String())
}

func TestClicksHandlerRejectsInvalidIncludeStale(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) { return uuid.New(), nil }))
	api := engine.Group("/api/v1")
	analyticshttp.RegisterRoutes(newClickHTTPHandler(t), api, api.Group("/lk"), nil)

	path := "/api/v1/lk/link/" + sharedvalueobjects.NewID().String() + "/clicks"
	req := httptest.NewRequest(http.MethodGet, path+"?section=about&breakpoint=md&includeStale=nope", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

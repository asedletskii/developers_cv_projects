package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/platform/http/middleware"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type clickHTTPLinkStub struct{ resumeID string }

func (s clickHTTPLinkStub) GetByID(context.Context, queries.PublicResumeLinkQuery) (queries.ResumeLinkPublicView, error) {
	return queries.ResumeLinkPublicView{ResumeID: s.resumeID}, nil
}

type clickHTTPOwnershipStub struct{}

func (clickHTTPOwnershipStub) GetOwned(context.Context, resumequeries.GetResumeQuery) (resumequeries.ResumeView, error) {
	return resumequeries.ResumeView{ID: "resume"}, nil
}

type clickHTTPPublicStub struct{}

func (clickHTTPPublicStub) Handle(context.Context, resumequeries.PublicResumeQuery) (resumequeries.PublicResumeView, error) {
	return resumequeries.PublicResumeView{}, nil
}

type clickHTTPLayoutStub struct{ snapshot domain.LayoutSnapshot }

func (s clickHTTPLayoutStub) Build(resumequeries.PublicResumeView) (domain.LayoutSnapshot, error) {
	return s.snapshot, nil
}

var _ ports.LayoutSnapshotBuilder = clickHTTPLayoutStub{}

type clickHTTPReaderStub struct{}

func (clickHTTPReaderStub) GetClicks(context.Context, queries.LinkClicksQuery) (queries.ClickGrid, error) {
	return queries.ClickGrid{Section: "about", Grid: 20, Cells: []queries.ClickCell{{X: 0, Y: 1, Count: 2}}, Total: 2}, nil
}

func newClickHTTPHandler() *Handler {
	layout, err := domain.NewLayoutSnapshot([]resumedomain.SectionID{resumedomain.SectionAbout}, map[resumedomain.SectionID]string{resumedomain.SectionAbout: "0123456789abcdef"})
	require.NoError(nil, err)
	clicks := queries.NewAnalyticsClickQueryHandler(clickHTTPLinkStub{resumeID: sharedvalueobjects.NewID().String()}, clickHTTPOwnershipStub{}, clickHTTPPublicStub{}, clickHTTPLayoutStub{snapshot: layout}, clickHTTPReaderStub{})
	return NewHandler(HandlerDependencies{Clicks: clicks})
}

func TestClicksHandlerReturnsClickGrid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	linkID := sharedvalueobjects.NewID()
	engine := gin.New()
	engine.Use(middleware.Authentication(func(context.Context, *http.Request) (uuid.UUID, error) { return userID, nil }))
	RegisterRoutes(newClickHTTPHandler(), engine.Group("/api/v1"), engine.Group("/api/v1"), engine.Group("/api/v1/lk"), nil, nil)

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
	RegisterRoutes(newClickHTTPHandler(), engine.Group("/api/v1"), engine.Group("/api/v1"), engine.Group("/api/v1/lk"), nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lk/link/"+sharedvalueobjects.NewID().String()+"/clicks?section=about&breakpoint=md&includeStale=nope", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume_link/application/commands"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type analyticsHandlerRepoStub struct{}

func (analyticsHandlerRepoStub) ApplyBatch(context.Context, sharedvalueobjects.ID, sharedvalueobjects.ID, domain.ViewStateUpdate, []*domain.AnalyticsEvent, []*domain.AnalyticsClick) error {
	return nil
}

type analyticsHandlerViewStub struct{}

func (analyticsHandlerViewStub) CreateView(context.Context, *domain.View) error { return nil }
func (analyticsHandlerViewStub) UpdateViewState(context.Context, sharedvalueobjects.ID, domain.ViewStateUpdate) error {
	return nil
}
func (analyticsHandlerViewStub) GetView(_ context.Context, id sharedvalueobjects.ID) (*domain.View, error) {
	layout, err := domain.NewLayoutSnapshot([]resumedomain.SectionID{resumedomain.SectionAbout}, map[resumedomain.SectionID]string{resumedomain.SectionAbout: "0123456789abcdef"})
	if err != nil {
		return nil, err
	}
	return domain.NewView(domain.ViewParams{ID: id, LinkID: sharedvalueobjects.NewID(), VisitorID: sharedvalueobjects.NewID(), IPHash: make([]byte, 32), Device: domain.DeviceDesktop, Browser: domain.BrowserChrome, OS: domain.OSLinux, Traffic: domain.TrafficHuman, Layout: layout, StartedAt: time.Now().UTC()})
}

type analyticsHandlerUOWStub struct{}

func (analyticsHandlerUOWStub) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

var _ ports.AnalyticsBatchRepository = analyticsHandlerRepoStub{}

func TestAnalyticsEventsHandlerAcceptsTextPlain(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewAnalyticsEventsHandler(commands.NewIngestAnalyticsBatchHandler(analyticsHandlerRepoStub{}, analyticsHandlerViewStub{}, analyticsHandlerUOWStub{}))
	engine := gin.New()
	engine.POST("/api/v1/cv/:linkId/events", handler.Handle)

	linkID := uuid.New()
	viewID := uuid.New()
	body := `{"viewId":"` + viewID.String() + `","seq":1,"final":false,"state":{"activeMs":6000,"maxSection":"about"},"events":[{"id":"` + uuid.New().String() + `","type":"download"}],"clicks":[]}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cv/"+linkID.String()+"/events", strings.NewReader(body))
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
}

func TestAnalyticsEventsHandlerRejectsInvalidViewID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewAnalyticsEventsHandler(commands.NewIngestAnalyticsBatchHandler(analyticsHandlerRepoStub{}, analyticsHandlerViewStub{}, analyticsHandlerUOWStub{}))
	engine := gin.New()
	engine.POST("/api/v1/cv/:linkId/events", handler.Handle)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/cv/"+uuid.New().String()+"/events", strings.NewReader(`{"viewId":"bad","seq":1,"state":{"activeMs":0},"events":[],"clicks":[]}`))
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAnalyticsEventsHandlerMapsInactiveLinkToGone(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewAnalyticsEventsHandler(commands.NewIngestAnalyticsBatchHandler(analyticsHandlerInactiveRepoStub{}, analyticsHandlerViewStub{}, analyticsHandlerUOWStub{}))
	engine := gin.New()
	engine.POST("/api/v1/cv/:linkId/events", handler.Handle)
	viewID := uuid.New()
	body := `{"viewId":"` + viewID.String() + `","seq":1,"final":false,"state":{"activeMs":0},"events":[],"clicks":[]}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cv/"+uuid.New().String()+"/events", strings.NewReader(body))
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusGone, w.Code)
}

type analyticsHandlerInactiveRepoStub struct{}

func (analyticsHandlerInactiveRepoStub) ApplyBatch(context.Context, sharedvalueobjects.ID, sharedvalueobjects.ID, domain.ViewStateUpdate, []*domain.AnalyticsEvent, []*domain.AnalyticsClick) error {
	return domain.ErrResumeLinkInactive
}

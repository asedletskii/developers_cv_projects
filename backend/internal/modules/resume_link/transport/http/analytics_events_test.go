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
	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/modules/resume_link/transport/http/handlers"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type analyticsHandlerRepoStub struct{}

func (analyticsHandlerRepoStub) LoadAnalyticsBatchForUpdate(_ context.Context, viewID, linkID sharedvalueobjects.ID, _, _, _ bool) (domain.AnalyticsBatchSnapshot, error) {
	layout, err := domain.NewLayoutSnapshot(
		[]resumedomain.SectionID{resumedomain.SectionAbout},
		map[resumedomain.SectionID]string{resumedomain.SectionAbout: "0123456789abcdef"},
	)
	if err != nil {
		return domain.AnalyticsBatchSnapshot{}, err
	}
	view, err := domain.RestoreView(domain.ViewSnapshot{
		ID: viewID, LinkID: linkID, VisitorID: sharedvalueobjects.NewID(), IPHash: make([]byte, 32),
		Device: domain.DeviceDesktop, Browser: domain.BrowserChrome, OS: domain.OSLinux,
		Traffic: domain.TrafficHuman, Layout: layout, StartedAt: time.Now().UTC(), CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		return domain.AnalyticsBatchSnapshot{}, err
	}
	return domain.AnalyticsBatchSnapshot{View: view, LinkStatus: domain.ResumeLinkStatusActive, ResumeStatus: resumedomain.ResumeStatusOpened}, nil
}

func (analyticsHandlerRepoStub) SaveAnalyticsBatch(context.Context, *domain.View, []*domain.AnalyticsEvent, []*domain.AnalyticsClick) error {
	return nil
}

type analyticsHandlerUOWStub struct{}

func (analyticsHandlerUOWStub) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

type analyticsHandlerObserverStub struct {
	dropped  map[string]int
	outcomes map[string]int
}

func (s *analyticsHandlerObserverStub) RecordAnalyticsBatchOutcome(outcome, reason string) {
	if s.outcomes == nil {
		s.outcomes = make(map[string]int)
	}
	s.outcomes[outcome+":"+reason]++
}

func (s *analyticsHandlerObserverStub) RecordDroppedAnalyticsInput(kind, reason string) {
	if s.dropped == nil {
		s.dropped = make(map[string]int)
	}
	s.dropped[kind+":"+reason]++
}

func TestAnalyticsEventsHandlerObservesMalformedElementsAndKeepsPartialAcceptance(t *testing.T) {
	gin.SetMode(gin.TestMode)
	observer := &analyticsHandlerObserverStub{}
	ingest := commands.NewIngestAnalyticsBatchHandler(analyticsHandlerRepoStub{}, analyticsHandlerUOWStub{}, observer)
	handler := handlers.NewResumeLinkAnalyticsEventsHandler(ingest)
	engine := gin.New()
	engine.POST("/api/v1/cv/:linkId/events", handler.Handle)

	validEventID := uuid.New()
	body := `{"viewId":"` + uuid.New().String() + `","seq":1,"final":false,"state":{"activeMs":0,"maxSection":"about"},` +
		`"events":[null,{"id":"bad","type":"download","at":1},{"id":"` + uuid.New().String() + `","type":"download","at":-1},{"id":"` + validEventID.String() + `","type":"download","at":1}],` +
		`"clicks":[null,{"section":"about","x":10,"y":20,"vw":800,"target":"button","at":-1}]}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cv/"+uuid.New().String()+"/events", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.Equal(t, 1, observer.dropped["event:invalid_json_shape"])
	require.Equal(t, 1, observer.dropped["event:invalid_client_id"])
	require.Equal(t, 1, observer.dropped["event:negative_timestamp"])
	require.Equal(t, 1, observer.dropped["click:invalid_json_shape"])
	require.Equal(t, 1, observer.dropped["click:negative_timestamp"])
	require.Equal(t, 1, observer.outcomes["accepted:success"])
}

func TestAnalyticsEventsHandlerAcceptsTextPlain(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := handlers.NewResumeLinkAnalyticsEventsHandler(commands.NewIngestAnalyticsBatchHandler(analyticsHandlerRepoStub{}, analyticsHandlerUOWStub{}, nil))
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
	observer := &analyticsHandlerObserverStub{}
	handler := handlers.NewResumeLinkAnalyticsEventsHandler(commands.NewIngestAnalyticsBatchHandler(analyticsHandlerRepoStub{}, analyticsHandlerUOWStub{}, observer))
	engine := gin.New()
	engine.POST("/api/v1/cv/:linkId/events", handler.Handle)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/cv/"+uuid.New().String()+"/events", strings.NewReader(`{"viewId":"bad","seq":1,"state":{"activeMs":0},"events":[],"clicks":[]}`))
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Equal(t, 1, observer.outcomes["rejected:invalid_batch"])
}

func TestAnalyticsEventsHandlerMapsInactiveLinkToGone(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := handlers.NewResumeLinkAnalyticsEventsHandler(commands.NewIngestAnalyticsBatchHandler(analyticsHandlerInactiveRepoStub{}, analyticsHandlerUOWStub{}, nil))
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

func (analyticsHandlerInactiveRepoStub) LoadAnalyticsBatchForUpdate(context.Context, sharedvalueobjects.ID, sharedvalueobjects.ID, bool, bool, bool) (domain.AnalyticsBatchSnapshot, error) {
	return domain.AnalyticsBatchSnapshot{}, domain.ErrResumeLinkInactive
}

func (analyticsHandlerInactiveRepoStub) SaveAnalyticsBatch(context.Context, *domain.View, []*domain.AnalyticsEvent, []*domain.AnalyticsClick) error {
	return nil
}

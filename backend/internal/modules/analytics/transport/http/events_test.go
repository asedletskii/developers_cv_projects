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

	"developers-cv/backend/internal/modules/analytics/application/commands"
	"developers-cv/backend/internal/modules/analytics/application/ports"
	analyticscontracts "developers-cv/backend/internal/modules/analytics/contracts"
	"developers-cv/backend/internal/modules/analytics/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type analyticsHandlerSubjectsStub struct{}

func (analyticsHandlerSubjectsStub) LoadLinkLifecycle(_ context.Context, linkID sharedvalueobjects.ID) (analyticscontracts.LinkLifecycle, error) {
	return analyticscontracts.LinkLifecycle{
		LinkID: linkID, ResumeID: sharedvalueobjects.NewID(),
		LinkActive: true, ResumePublished: true,
	}, nil
}

func (analyticsHandlerSubjectsStub) ListResumeLinks(context.Context, sharedvalueobjects.ID) ([]analyticscontracts.ResumeLinkRef, error) {
	return nil, nil
}

type analyticsHandlerRepoStub struct{}

func (analyticsHandlerRepoStub) LoadAnalyticsBatchForUpdate(_ context.Context, viewID, linkID sharedvalueobjects.ID, _, _, _ bool) (ports.AnalyticsBatchSnapshot, error) {
	layout, err := domain.NewLayoutSnapshot(
		[]domain.AnalyticsSectionID{domain.AnalyticsSectionAbout},
		map[domain.AnalyticsSectionID]string{domain.AnalyticsSectionAbout: "0123456789abcdef"},
	)
	if err != nil {
		return ports.AnalyticsBatchSnapshot{}, err
	}
	view, err := domain.RestoreView(domain.ViewSnapshot{
		ID: viewID, LinkID: linkID, ResumeID: sharedvalueobjects.NewID(), VisitorID: sharedvalueobjects.NewID(), IPHash: make([]byte, 32),
		Device: domain.DeviceDesktop, Browser: domain.BrowserChrome, OS: domain.OSLinux,
		Traffic: domain.TrafficHuman, Layout: layout, StartedAt: time.Now().UTC(), CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		return ports.AnalyticsBatchSnapshot{}, err
	}
	return ports.AnalyticsBatchSnapshot{
		View: view,
		LinkContext: ports.LinkAnalyticsContext{
			LinkID: linkID, ResumeID: sharedvalueobjects.NewID(), LinkActive: true, ResumePublished: true,
		},
	}, nil
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
	ingest := commands.NewIngestAnalyticsBatchHandler(analyticsHandlerRepoStub{}, analyticsHandlerUOWStub{}, analyticsHandlerSubjectsStub{}, observer)
	engine := gin.New()
	registerAnalyticsTestRoutes(engine, ingest)

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
	engine := gin.New()
	registerAnalyticsTestRoutes(engine, commands.NewIngestAnalyticsBatchHandler(analyticsHandlerRepoStub{}, analyticsHandlerUOWStub{}, analyticsHandlerSubjectsStub{}, nil))

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
	engine := gin.New()
	registerAnalyticsTestRoutes(engine, commands.NewIngestAnalyticsBatchHandler(analyticsHandlerRepoStub{}, analyticsHandlerUOWStub{}, analyticsHandlerSubjectsStub{}, observer))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/cv/"+uuid.New().String()+"/events", strings.NewReader(`{"viewId":"bad","seq":1,"state":{"activeMs":0},"events":[],"clicks":[]}`))
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Equal(t, 1, observer.outcomes["rejected:invalid_batch"])
}

func TestAnalyticsEventsHandlerMapsInactiveLinkToGone(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	registerAnalyticsTestRoutes(engine, commands.NewIngestAnalyticsBatchHandler(analyticsHandlerInactiveRepoStub{}, analyticsHandlerUOWStub{}, analyticsHandlerSubjectsStub{}, nil))
	viewID := uuid.New()
	body := `{"viewId":"` + viewID.String() + `","seq":1,"final":false,"state":{"activeMs":0},"events":[],"clicks":[]}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cv/"+uuid.New().String()+"/events", strings.NewReader(body))
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusGone, w.Code)
}

func registerAnalyticsTestRoutes(engine *gin.Engine, ingest *commands.IngestAnalyticsBatchHandler) {
	api := engine.Group("/api/v1")
	RegisterRoutes(NewHandler(HandlerDependencies{Ingest: ingest}), api, api, nil)
}

type analyticsHandlerInactiveRepoStub struct{}

func (analyticsHandlerInactiveRepoStub) LoadAnalyticsBatchForUpdate(context.Context, sharedvalueobjects.ID, sharedvalueobjects.ID, bool, bool, bool) (ports.AnalyticsBatchSnapshot, error) {
	return ports.AnalyticsBatchSnapshot{}, domain.ErrAnalyticsLinkInactive
}

func (analyticsHandlerInactiveRepoStub) SaveAnalyticsBatch(context.Context, *domain.View, []*domain.AnalyticsEvent, []*domain.AnalyticsClick) error {
	return nil
}

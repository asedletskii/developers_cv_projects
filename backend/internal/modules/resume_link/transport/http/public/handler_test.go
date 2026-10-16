package public

import (
	"context"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type queryStub struct {
	view     resumequeries.PublicResumeView
	err      error
	query    queries.PublicResumeLinkQuery
	identity ports.VisitorIdentityResult
}

func (s *queryStub) Handle(_ context.Context, query queries.PublicResumeLinkQuery) (resumequeries.PublicResumeView, error) {
	s.query = query
	return s.view, s.err
}
func (s *queryStub) HandleWithVisit(_ context.Context, query queries.PublicResumeLinkQuery) (resumequeries.PublicResumeView, ports.VisitorIdentityResult, error) {
	s.query = query
	return s.view, s.identity, s.err
}

func TestHandlerReturnsPublicResume(t *testing.T) {
	gin.SetMode(gin.TestMode)
	linkID := sharedvalueobjects.NewID()
	query := &queryStub{view: resumequeries.PublicResumeView{Title: "CV"}}

	engine := gin.New()
	engine.GET("/api/v1/cv/:linkId", NewHandler(query).Handle)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cv/"+linkID.String(), nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "no-store", w.Header().Get("Cache-Control"))
	require.NotContains(t, w.Body.String(), `"resumeId"`)
	require.NotContains(t, w.Body.String(), `"id":"resume-id"`)
	require.Contains(t, w.Body.String(), `"title":"CV"`)
	require.Equal(t, linkID, query.query.LinkID)
}

func TestHandlerRejectsInvalidLinkID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	query := &queryStub{}
	engine := gin.New()
	engine.GET("/api/v1/cv/:linkId", NewHandler(query).Handle)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cv/not-a-uuid", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandlerMapsInactiveLinkToGone(t *testing.T) {
	gin.SetMode(gin.TestMode)
	query := &queryStub{err: domain.ErrResumeLinkInactive}
	engine := gin.New()
	engine.GET("/api/v1/cv/:linkId", NewHandler(query).Handle)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cv/"+sharedvalueobjects.NewID().String(), nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusGone, w.Code)
	require.JSONEq(t, `{"error":{"code":"GONE","message":"Resume link is no longer active","requestId":""}}`, w.Body.String())
}

func TestHandlerEmitsViewIDAndVisitorCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	linkID := sharedvalueobjects.NewID()
	viewID := sharedvalueobjects.NewID()
	cookieID := sharedvalueobjects.NewID()
	query := &queryStub{view: resumequeries.PublicResumeView{Title: "CV"}, identity: ports.VisitorIdentityResult{ViewID: viewID, SetCookie: true, Cookie: &http.Cookie{Name: "cv_vid", Value: cookieID.String(), Path: "/api/v1/cv", MaxAge: 31536000, HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode}}}
	engine := gin.New()
	engine.GET("/api/v1/cv/:linkId", NewHandler(query).Handle)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/cv/"+linkID.String(), nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, viewID.String(), w.Header().Get("X-View-Id"))
	require.Contains(t, w.Header().Get("Set-Cookie"), "cv_vid="+cookieID.String())
}

func TestHandlerDoesNotEmitAnalyticsHeadersWhenVisitIsUnavailable(t *testing.T) {
	gin.SetMode(gin.TestMode)
	query := &queryStub{view: resumequeries.PublicResumeView{Title: "CV"}}
	engine := gin.New()
	engine.GET("/api/v1/cv/:linkId", NewHandler(query).Handle)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/cv/"+sharedvalueobjects.NewID().String(), nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	require.Empty(t, w.Header().Get("X-View-Id"))
	require.Empty(t, w.Header().Get("Set-Cookie"))
	require.Equal(t, "no-store", w.Header().Get("Cache-Control"))
}

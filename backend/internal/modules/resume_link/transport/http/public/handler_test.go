package public

import (
	"context"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type queryStub struct {
	view  resumequeries.PublicResumeView
	err   error
	query queries.PublicResumeLinkQuery
}

func (s *queryStub) Handle(_ context.Context, query queries.PublicResumeLinkQuery) (resumequeries.PublicResumeView, error) {
	s.query = query
	return s.view, s.err
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

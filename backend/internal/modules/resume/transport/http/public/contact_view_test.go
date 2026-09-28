package public

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain/enums"
)

type contactViewQueryStub struct {
	view  queries.ContactView
	err   error
	query queries.ContactViewQuery
}

func (s *contactViewQueryStub) Handle(_ context.Context, query queries.ContactViewQuery) (queries.ContactView, error) {
	s.query = query
	if s.err != nil {
		return queries.ContactView{}, s.err
	}
	return s.view, nil
}

func TestContactViewHandlerReturnsContactValue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	query := &contactViewQueryStub{view: queries.ContactView{Contact: enums.ContactKeyEmail, Value: "person@example.com"}}
	handler := NewContactViewHandler(query)
	engine := gin.New()
	engine.POST("/api/v1/cv/:resumeId/contact-views", handler.Handle)

	id := uuid.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cv/"+id.String()+"/contact-views", jsonBody(t, `{"contact":"email"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, id, query.query.ResumeID.UUID())
	require.Equal(t, enums.ContactKeyEmail, query.query.Contact)
	require.Equal(t, "no-store", w.Header().Get("Cache-Control"))
	require.JSONEq(t, `{"contact":"email","value":"person@example.com"}`, w.Body.String())
}

func TestContactViewHandlerRejectsInvalidContact(t *testing.T) {
	gin.SetMode(gin.TestMode)
	query := &contactViewQueryStub{}
	handler := NewContactViewHandler(query)
	engine := gin.New()
	engine.POST("/api/v1/cv/:resumeId/contact-views", handler.Handle)

	id := uuid.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cv/"+id.String()+"/contact-views", jsonBody(t, `{"contact":"fax"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Empty(t, query.query)
}

func TestContactViewHandlerMapsNotHidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	query := &contactViewQueryStub{err: queries.ErrContactNotHidden}
	handler := NewContactViewHandler(query)
	engine := gin.New()
	engine.POST("/api/v1/cv/:resumeId/contact-views", handler.Handle)

	id := uuid.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cv/"+id.String()+"/contact-views", jsonBody(t, `{"contact":"email"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusConflict, w.Code)
}

func jsonBody(t *testing.T, value string) *strings.Reader {
	t.Helper()
	return strings.NewReader(value)
}

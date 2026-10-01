package public

import (
	"context"
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type contactViewQueryStub struct {
	view  resumequeries.ContactView
	err   error
	query queries.ContactViewThroughLinkQuery
}

func (s *contactViewQueryStub) Handle(_ context.Context, query queries.ContactViewThroughLinkQuery) (resumequeries.ContactView, error) {
	s.query = query
	if s.err != nil {
		return resumequeries.ContactView{}, s.err
	}
	return s.view, nil
}

func TestContactViewHandlerReturnsContactValue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	linkID := sharedvalueobjects.NewID()
	query := &contactViewQueryStub{view: resumequeries.ContactView{Contact: resumedomain.ContactKeyEmail, Value: "person@example.com"}}
	engine := gin.New()
	engine.POST("/api/v1/cv/:linkId/contact-views", NewContactViewHandler(query).Handle)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/cv/"+linkID.String()+"/contact-views", strings.NewReader(`{"contact":"email"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, linkID, query.query.LinkID)
	require.Equal(t, resumedomain.ContactKeyEmail, query.query.Contact)
	require.Equal(t, "no-store", w.Header().Get("Cache-Control"))
	require.JSONEq(t, `{"contact":"email","value":"person@example.com"}`, w.Body.String())
}

func TestContactViewHandlerRejectsInvalidLinkID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	query := &contactViewQueryStub{}
	engine := gin.New()
	engine.POST("/api/v1/cv/:linkId/contact-views", NewContactViewHandler(query).Handle)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/cv/not-a-uuid/contact-views", strings.NewReader(`{"contact":"email"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.True(t, query.query.LinkID.IsZero())
}

func TestContactViewHandlerMapsMissingOrInactiveLink(t *testing.T) {
	cases := []struct {
		name string
		err  error
		code int
	}{
		{name: "not found", err: queries.ErrResumeLinkNotFound, code: http.StatusNotFound},
		{name: "inactive", err: domain.ErrResumeLinkInactive, code: http.StatusGone},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			query := &contactViewQueryStub{err: tt.err}
			engine := gin.New()
			engine.POST("/api/v1/cv/:linkId/contact-views", NewContactViewHandler(query).Handle)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/cv/"+sharedvalueobjects.NewID().String()+"/contact-views", strings.NewReader(`{"contact":"email"}`))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)

			require.Equal(t, tt.code, w.Code)
		})
	}
}

func TestContactViewHandlerMapsResumeContactErrors(t *testing.T) {
	cases := []struct {
		name string
		err  error
		code int
	}{
		{name: "missing contact", err: resumequeries.ErrContactNotFound, code: http.StatusNotFound},
		{name: "not hidden", err: resumequeries.ErrContactNotHidden, code: http.StatusConflict},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			query := &contactViewQueryStub{err: tt.err}
			engine := gin.New()
			engine.POST("/api/v1/cv/:linkId/contact-views", NewContactViewHandler(query).Handle)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/cv/"+sharedvalueobjects.NewID().String()+"/contact-views", strings.NewReader(`{"contact":"email"}`))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)

			require.Equal(t, tt.code, w.Code)
		})
	}
}

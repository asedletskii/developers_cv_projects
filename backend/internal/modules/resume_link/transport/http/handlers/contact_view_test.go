package handlers

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
	resumelinkports "developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/platform/http/middleware"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type tokenVerifierStub struct {
	err   error
	token string
	ip    string
}

func (s *tokenVerifierStub) Verify(_ context.Context, token string, remoteIP string) error {
	s.token = token
	s.ip = remoteIP
	return s.err
}

type contactViewQueryStub struct {
	view  queries.ContactViewThroughLinkResult
	err   error
	query queries.ContactViewThroughLinkQuery
}

func (s *contactViewQueryStub) Handle(_ context.Context, query queries.ContactViewThroughLinkQuery) (queries.ContactViewThroughLinkResult, error) {
	s.query = query
	if s.err != nil {
		return queries.ContactViewThroughLinkResult{}, s.err
	}
	return s.view, nil
}

func TestContactViewHandlerReturnsContactValue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	linkID := sharedvalueobjects.NewID()
	contact := resumedomain.ContactKeyEmail
	value := "person@example.com"
	query := &contactViewQueryStub{view: queries.ContactViewThroughLinkResult{Contact: &contact, Value: &value}}
	verifier := &tokenVerifierStub{}
	engine := gin.New()
	engine.POST("/api/v1/cv/:linkId/contact-views", NewPublicContactViewHandler(query, verifier).Handle)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/cv/"+linkID.String()+"/contact-views", strings.NewReader(`{"contact":"email","turnstileToken":"test-token"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, linkID, query.query.LinkID)
	require.Equal(t, resumedomain.ContactKeyEmail, query.query.Contact)
	require.Equal(t, queries.ContactViewInitiatorContact, query.query.Initiator)
	require.Equal(t, "no-store", w.Header().Get("Cache-Control"))
	require.JSONEq(t, `{"contact":"email","value":"person@example.com"}`, w.Body.String())
	require.NotContains(t, w.Body.String(), "test-token")
}

func TestContactViewHandlerPassesResolvedClientIPToTurnstile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	contact := resumedomain.ContactKeyEmail
	value := "person@example.com"
	query := &contactViewQueryStub{view: queries.ContactViewThroughLinkResult{Contact: &contact, Value: &value}}
	verifier := &tokenVerifierStub{}
	engine := gin.New()
	require.NoError(t, engine.SetTrustedProxies([]string{"192.0.2.0/24"}))
	engine.POST("/api/v1/cv/:linkId/contact-views", NewPublicContactViewHandler(query, verifier).Handle)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/cv/"+sharedvalueobjects.NewID().String()+"/contact-views", strings.NewReader(`{"contact":"email","turnstileToken":"test-token"}`))
	req.RemoteAddr = "192.0.2.10:12345"
	req.Header.Set("X-Forwarded-For", "203.0.113.7")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "203.0.113.7", verifier.ip)
}

func TestContactViewHandlerRejectsInvalidLinkID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	query := &contactViewQueryStub{}
	verifier := &tokenVerifierStub{}
	engine := gin.New()
	engine.POST("/api/v1/cv/:linkId/contact-views", NewPublicContactViewHandler(query, verifier).Handle)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/cv/not-a-uuid/contact-views", strings.NewReader(`{"contact":"email","turnstileToken":"test-token"}`))
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
			verifier := &tokenVerifierStub{}
			engine := gin.New()
			engine.POST("/api/v1/cv/:linkId/contact-views", NewPublicContactViewHandler(query, verifier).Handle)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/cv/"+sharedvalueobjects.NewID().String()+"/contact-views", strings.NewReader(`{"contact":"email","turnstileToken":"test-token"}`))
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
			verifier := &tokenVerifierStub{}
			engine := gin.New()
			engine.POST("/api/v1/cv/:linkId/contact-views", NewPublicContactViewHandler(query, verifier).Handle)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/cv/"+sharedvalueobjects.NewID().String()+"/contact-views", strings.NewReader(`{"contact":"email","turnstileToken":"test-token"}`))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)

			require.Equal(t, tt.code, w.Code)
		})
	}
}

func TestContactViewHandlerRejectsInvalidTokenBeforeVerifier(t *testing.T) {
	gin.SetMode(gin.TestMode)
	query := &contactViewQueryStub{}
	verifier := &tokenVerifierStub{}
	engine := gin.New()
	engine.POST("/api/v1/cv/:linkId/contact-views", NewPublicContactViewHandler(query, verifier).Handle)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/cv/"+sharedvalueobjects.NewID().String()+"/contact-views", strings.NewReader(`{"contact":"email","turnstileToken":""}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Empty(t, verifier.token)
	require.True(t, query.query.LinkID.IsZero())
}

func TestContactViewHandlerMapsTokenVerificationFailureToForbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	contact := resumedomain.ContactKeyEmail
	value := "person@example.com"
	query := &contactViewQueryStub{view: queries.ContactViewThroughLinkResult{Contact: &contact, Value: &value}}
	verifier := &tokenVerifierStub{err: resumelinkports.ErrTokenVerificationFailed}
	engine := gin.New()
	engine.POST("/api/v1/cv/:linkId/contact-views", NewPublicContactViewHandler(query, verifier).Handle)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/cv/"+sharedvalueobjects.NewID().String()+"/contact-views", strings.NewReader(`{"contact":"email","turnstileToken":"test-token"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusForbidden, w.Code)
	require.Contains(t, w.Body.String(), "Contact verification failed")
	require.NotContains(t, w.Body.String(), "test-token")
	require.NotContains(t, w.Body.String(), "invalid-input-secret")
	require.True(t, query.query.LinkID.IsZero())
}

func TestContactViewHandlerReturnsAllContactsForDownload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	linkID := sharedvalueobjects.NewID()
	contacts := []queries.ContactViewThroughLinkContact{
		{Contact: resumedomain.ContactKeyPhone, Value: "+33123456789"},
		{Contact: resumedomain.ContactKeyEmail, Value: "person@example.com"},
	}
	query := &contactViewQueryStub{view: queries.ContactViewThroughLinkResult{Contacts: &contacts}}
	verifier := &tokenVerifierStub{}
	engine := gin.New()
	engine.POST("/api/v1/cv/:linkId/contact-views", NewPublicContactViewHandler(query, verifier).Handle)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/cv/"+linkID.String()+"/contact-views",
		strings.NewReader(`{"initiator":"download","turnstileToken":"test-token"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, queries.ContactViewInitiatorDownload, query.query.Initiator)
	require.Empty(t, query.query.Contact)
	require.Equal(t, "test-token", verifier.token)
	require.JSONEq(t, `{"contacts":[{"contact":"phone","value":"+33123456789"},{"contact":"email","value":"person@example.com"}]}`, w.Body.String())
}

func TestContactViewHandlerReturnsEmptyContactsForDownload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	linkID := sharedvalueobjects.NewID()
	contacts := make([]queries.ContactViewThroughLinkContact, 0)
	query := &contactViewQueryStub{view: queries.ContactViewThroughLinkResult{Contacts: &contacts}}
	engine := gin.New()
	engine.POST("/api/v1/cv/:linkId/contact-views", NewPublicContactViewHandler(query, &tokenVerifierStub{}).Handle)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/cv/"+linkID.String()+"/contact-views",
		strings.NewReader(`{"initiator":"download","turnstileToken":"test-token"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"contacts":[]}`, w.Body.String())
}

func TestContactViewHandlerRejectsInvalidInitiatorAndContactCombinations(t *testing.T) {
	for _, body := range []string{
		`{"initiator":"other","contact":"email","turnstileToken":"test-token"}`,
		`{"initiator":"contact","turnstileToken":"test-token"}`,
		`{"initiator":"download","contact":"email","turnstileToken":"test-token"}`,
		`{"initiator":"download","contact":"","turnstileToken":"test-token"}`,
		`{"turnstileToken":"test-token"}`,
	} {
		t.Run(body, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			query := &contactViewQueryStub{}
			verifier := &tokenVerifierStub{}
			engine := gin.New()
			engine.POST("/api/v1/cv/:linkId/contact-views", NewPublicContactViewHandler(query, verifier).Handle)

			req := httptest.NewRequest(
				http.MethodPost,
				"/api/v1/cv/"+sharedvalueobjects.NewID().String()+"/contact-views",
				strings.NewReader(body),
			)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code)
			require.True(t, query.query.LinkID.IsZero())
			require.Empty(t, verifier.token)
		})
	}
}

func TestContactViewHandlerReturnsPayloadTooLargeForOversizedBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	query := &contactViewQueryStub{}
	verifier := &tokenVerifierStub{}
	engine := gin.New()
	engine.Use(middleware.BodyLimit(32))
	engine.POST("/api/v1/cv/:linkId/contact-views", NewPublicContactViewHandler(query, verifier).Handle)

	body := `{"contact":"email","turnstileToken":"` + strings.Repeat("x", 64) + `"}`
	request := httptest.NewRequest(http.MethodPost, "/api/v1/cv/"+sharedvalueobjects.NewID().String()+"/contact-views", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	require.Equal(t, http.StatusRequestEntityTooLarge, response.Code)
	require.Contains(t, response.Body.String(), "REQUEST_TOO_LARGE")
}

package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	contactkey "developers-cv/backend/internal/shared/domain/contact"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type goldenContactViewProvider struct {
	view resumecontracts.ContactView
	err  error
}

func (s goldenContactViewProvider) Handle(context.Context, resumecontracts.ContactViewQuery) (resumecontracts.ContactView, error) {
	return s.view, s.err
}

type goldenDownloadContactsProvider struct {
	view resumecontracts.DownloadContactsView
	err  error
}

func (s goldenDownloadContactsProvider) Handle(context.Context, resumecontracts.DownloadContactsQuery) (resumecontracts.DownloadContactsView, error) {
	return s.view, s.err
}

func TestContactViewHandlerGoldenReveal(t *testing.T) {
	gin.SetMode(gin.TestMode)
	linkID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	provider := goldenContactViewProvider{view: resumecontracts.ContactView{
		Contact: contactkey.Email,
		Value:   "person@example.com",
	}}
	query := queries.NewContactViewThroughLinkQueryHandler(goldenLinkResolver{resumeID: resumeID}, provider)

	engine := gin.New()
	engine.POST("/api/v1/cv/:linkId/contact-views", NewPublicContactViewHandler(query, &tokenVerifierStub{}).Handle)
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/cv/"+linkID.String()+"/contact-views",
		strings.NewReader(`{"contact":"email","turnstileToken":"test-token"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	golden, err := os.ReadFile("testdata/contact_view.golden.json")
	require.NoError(t, err)
	require.JSONEq(t, string(golden), w.Body.String())

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "no-store", w.Header().Get("Cache-Control"))
	require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	require.NotContains(t, w.Body.String(), "test-token")
}

func TestContactViewHandlerGoldenDownload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	linkID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	phone := "+33123456789"
	telegram := "ada"
	email := "person@example.com"
	provider := goldenDownloadContactsProvider{view: resumecontracts.DownloadContactsView{
		Phone:    &phone,
		Telegram: &telegram,
		Email:    &email,
	}}
	query := queries.NewContactViewThroughLinkQueryHandlerWithAnalytics(
		goldenLinkResolver{resumeID: resumeID},
		nil,
		provider,
		nil,
		nil,
	)

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

	golden, err := os.ReadFile("testdata/contact_download.golden.json")
	require.NoError(t, err)
	require.JSONEq(t, string(golden), w.Body.String())

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "no-store", w.Header().Get("Cache-Control"))
	require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	require.NotContains(t, w.Body.String(), "test-token")
}

package http

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/auth/application"
	"developers-cv/backend/internal/modules/auth/application/commands"
	"developers-cv/backend/internal/modules/auth/application/ports"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	"developers-cv/backend/internal/modules/auth/infrastructure/token"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type fakeStartTelegramLogin struct {
	location string
	err      error
}

func (f *fakeStartTelegramLogin) Handle(context.Context) (string, error) {
	return f.location, f.err
}

type fakeCompleteTelegramLogin struct {
	session commands.AuthenticatedSession
	err     error
}

func (f *fakeCompleteTelegramLogin) Handle(context.Context, string, string) (commands.AuthenticatedSession, error) {
	return f.session, f.err
}

type fakeLogoutHandler struct {
	hash  valueobjects.SessionTokenHash
	now   time.Time
	err   error
	calls int
}

func (f *fakeLogoutHandler) Handle(_ context.Context, hash valueobjects.SessionTokenHash, now time.Time) error {
	f.hash = hash
	f.now = now
	f.calls++
	return f.err
}

func newTestHandler(start startTelegramLogin, complete completeTelegramLogin, logout logoutHandler) *Handler {
	return NewHandler(HandlerDependencies{
		StartTelegramLogin: start,
		CompleteTelegram:   complete,
		Logout:             logout,
		TokenGenerator:     token.NewGenerator(),
		Cookie:             DefaultCookieConfig(30 * 24 * 60 * 60),
		SuccessLocation:    "/cv/lk",
		FailureLocation:    "/cv/auth/error",
	})
}

func performRequest(handler gin.HandlerFunc, method, target string, configure func(*http.Request)) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	request := httptest.NewRequest(method, target, nil)
	router.Handle(method, request.URL.Path, handler)
	recorder := httptest.NewRecorder()
	if configure != nil {
		configure(request)
	}
	router.ServeHTTP(recorder, request)
	return recorder
}

func TestHandlerStartTelegramRedirects(t *testing.T) {
	handler := newTestHandler(&fakeStartTelegramLogin{location: "https://oauth.example/auth?state=opaque"}, nil, nil)

	response := performRequest(handler.StartTelegram, http.MethodGet, "/api/v1/auth/telegram", nil)

	require.Equal(t, http.StatusFound, response.Code)
	require.Equal(t, "https://oauth.example/auth?state=opaque", response.Header().Get("Location"))
}

func TestHandlerStartTelegramMapsApplicationError(t *testing.T) {
	handler := newTestHandler(&fakeStartTelegramLogin{err: application.ErrInvalidAuthFlow}, nil, nil)

	response := performRequest(handler.StartTelegram, http.MethodGet, "/api/v1/auth/telegram", nil)

	require.Equal(t, http.StatusBadRequest, response.Code)
	require.Contains(t, response.Body.String(), `"code":"BAD_REQUEST"`)
	require.Contains(t, response.Body.String(), `"message":"Authentication failed"`)
}

func TestHandlerCompleteTelegramSetsSecureSessionCookieAndRedirects(t *testing.T) {
	userID := sharedvalueobjects.NewID()
	handler := newTestHandler(nil, &fakeCompleteTelegramLogin{session: commands.AuthenticatedSession{
		UserID:          userID,
		RawSessionToken: "test-session-value",
		ExpiresAt:       time.Now().UTC().Add(time.Hour),
	}}, nil)

	response := performRequest(handler.CompleteTelegram, http.MethodGet, "/api/v1/auth/telegram/callback?code=authorization-code&state=state", nil)

	require.Equal(t, http.StatusFound, response.Code)
	require.Equal(t, "/cv/lk", response.Header().Get("Location"))
	setCookie := response.Header().Get("Set-Cookie")
	require.Contains(t, setCookie, "__Host-session=test-session-value")
	require.Contains(t, setCookie, "Path=/")
	require.Contains(t, setCookie, "Max-Age=2592000")
	require.Contains(t, setCookie, "HttpOnly")
	require.Contains(t, setCookie, "Secure")
	require.Contains(t, setCookie, "SameSite=Lax")
	require.NotContains(t, response.Body.String(), "authorization-code")
	require.NotContains(t, response.Body.String(), userID.String())
}

func TestHandlerCompleteTelegramRejectsMissingCallbackValues(t *testing.T) {
	handler := newTestHandler(nil, &fakeCompleteTelegramLogin{}, nil)

	response := performRequest(handler.CompleteTelegram, http.MethodGet, "/api/v1/auth/telegram/callback?code=authorization-code", nil)

	require.Equal(t, http.StatusBadRequest, response.Code)
	require.Contains(t, response.Body.String(), `"message":"Authentication failed"`)
}

func TestHandlerCompleteTelegramRejectsInvalidState(t *testing.T) {
	handler := newTestHandler(nil, &fakeCompleteTelegramLogin{err: application.ErrStateInvalid}, nil)

	response := performRequest(handler.CompleteTelegram, http.MethodGet, "/api/v1/auth/telegram/callback?code=code&state=state", nil)

	require.Equal(t, http.StatusBadRequest, response.Code)
	require.Contains(t, response.Body.String(), `"message":"Authentication failed"`)
	require.NotContains(t, response.Body.String(), application.ErrStateInvalid.Error())
}

func TestHandlerCompleteTelegramProviderDenialRedirects(t *testing.T) {
	handler := newTestHandler(nil, &fakeCompleteTelegramLogin{}, nil)

	response := performRequest(handler.CompleteTelegram, http.MethodGet, "/api/v1/auth/telegram/callback?error=access_denied", nil)

	require.Equal(t, http.StatusFound, response.Code)
	require.Equal(t, "/cv/auth/error?reason=provider_denied", response.Header().Get("Location"))
}

func TestHandlerCompleteTelegramMapsProviderFailureWithoutDetails(t *testing.T) {
	handler := newTestHandler(nil, &fakeCompleteTelegramLogin{err: errors.New("telegram secret token=should-not-leak")}, nil)

	response := performRequest(handler.CompleteTelegram, http.MethodGet, "/api/v1/auth/telegram/callback?code=code&state=state", nil)

	require.Equal(t, http.StatusInternalServerError, response.Code)
	require.NotContains(t, response.Body.String(), "telegram secret")
	require.NotContains(t, response.Body.String(), "should-not-leak")
}

func TestHandlerCompleteTelegramMapsIdentityVerificationFailure(t *testing.T) {
	handler := newTestHandler(nil, &fakeCompleteTelegramLogin{err: application.ErrIdentityVerificationFailed}, nil)

	response := performRequest(handler.CompleteTelegram, http.MethodGet, "/api/v1/auth/telegram/callback?code=code&state=state", nil)

	require.Equal(t, http.StatusBadRequest, response.Code)
	require.Contains(t, response.Body.String(), `"message":"Authentication failed"`)
}

func TestHandlerCurrentSessionUsesAuthenticatedMiddlewareUser(t *testing.T) {
	userID := uuid.New()
	handler := newTestHandler(nil, nil, nil)

	response := performRequest(func(c *gin.Context) {
		c.Set("user_id", userID)
		handler.CurrentSession(c)
	}, http.MethodGet, "/api/v1/auth/me", nil)

	require.Equal(t, http.StatusOK, response.Code)
	require.JSONEq(t, `{"userId":"`+userID.String()+`","authenticated":true}`, response.Body.String())
}

func TestHandlerCurrentSessionRejectsMissingAuthenticationContext(t *testing.T) {
	handler := newTestHandler(nil, nil, nil)

	response := performRequest(handler.CurrentSession, http.MethodGet, "/api/v1/auth/me", nil)

	require.Equal(t, http.StatusUnauthorized, response.Code)
}

func TestHandlerLogoutRevokesSessionAndClearsCookie(t *testing.T) {
	logout := &fakeLogoutHandler{}
	handler := newTestHandler(nil, nil, logout)
	rawToken := []byte("01234567890123456789012345678901")
	encoded := encodeTestToken(rawToken)

	response := performRequest(handler.Logout, http.MethodPost, "/api/v1/auth/logout", func(request *http.Request) {
		request.AddCookie(&http.Cookie{
			Name:     DefaultCookieName,
			Value:    encoded,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})
	})

	require.Equal(t, http.StatusNoContent, response.Code)
	require.Equal(t, 1, logout.calls)
	require.Equal(t, token.NewGenerator().Hash(rawToken), logout.hash.Bytes())
	setCookie := response.Header().Get("Set-Cookie")
	require.Contains(t, setCookie, "__Host-session=")
	require.Contains(t, setCookie, "Max-Age=0")
	require.Contains(t, setCookie, "HttpOnly")
	require.Contains(t, setCookie, "Secure")
	require.Contains(t, setCookie, "SameSite=Lax")
}

func TestHandlerLogoutIsIdempotentWithoutCookie(t *testing.T) {
	logout := &fakeLogoutHandler{}
	handler := newTestHandler(nil, nil, logout)

	response := performRequest(handler.Logout, http.MethodPost, "/api/v1/auth/logout", nil)

	require.Equal(t, http.StatusNoContent, response.Code)
	require.Equal(t, 0, logout.calls)
	require.Contains(t, response.Header().Get("Set-Cookie"), "Max-Age=0")
}

func TestCookieConfigRejectsInsecureHostCookie(t *testing.T) {
	config := DefaultCookieConfig(60)
	config.Secure = false

	require.Error(t, config.Validate(false))
}

func TestCookieConfigRejectsProductionNonHostCookie(t *testing.T) {
	config := DefaultCookieConfig(60)
	config.Name = "session"

	require.Error(t, config.Validate(true))
}

func TestCookieConfigAllowsExplicitLocalNonHostCookie(t *testing.T) {
	config := DefaultCookieConfig(60)
	config.Name = "session"
	config.Secure = false

	require.NoError(t, config.Validate(false))
}

func TestSessionTokenHashRejectsMalformedCookie(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/", nil)
	request.AddCookie(&http.Cookie{
		Name:     DefaultCookieName,
		Value:    "not-base64***",
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	_, err := sessionTokenHash(request, DefaultCookieName, token.NewGenerator())

	require.ErrorIs(t, err, errInvalidSessionCookie)
}

func encodeTestToken(value []byte) string {
	return base64.RawURLEncoding.EncodeToString(value)
}

var _ ports.TokenGenerator = (*token.Generator)(nil)

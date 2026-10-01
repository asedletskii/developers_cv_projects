package http

import (
	"context"
	"errors"
	"testing"
	"time"

	"encoding/base64"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/auth/application"
	"developers-cv/backend/internal/modules/auth/application/commands"
	"developers-cv/backend/internal/modules/auth/application/ports"
	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	"developers-cv/backend/internal/modules/auth/infrastructure/token"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type fakeStartTelegramLogin struct {
	start commands.TelegramLoginStart
	err   error
}

func (f *fakeStartTelegramLogin) Handle(context.Context) (commands.TelegramLoginStart, error) {
	return f.start, f.err
}

type fakeCompleteTelegramLogin struct {
	session commands.AuthenticatedSession
	err     error
}

func (f *fakeCompleteTelegramLogin) Handle(_ context.Context, _ string, _ string, browserBinding string) (commands.AuthenticatedSession, error) {
	if f.err == nil && browserBinding == "" {
		return commands.AuthenticatedSession{}, application.ErrBrowserBindingInvalid
	}
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
	handler := newTestHandler(&fakeStartTelegramLogin{start: commands.TelegramLoginStart{AuthorizationURL: "https://oauth.example/auth?state=opaque", BrowserBindingToken: "binding", BrowserBindingMaxAge: 300}}, nil, nil)

	response := performRequest(handler.StartTelegram, http.MethodGet, "/api/v1/auth/telegram", nil)

	require.Equal(t, http.StatusFound, response.Code)
	require.Equal(t, "https://oauth.example/auth?state=opaque", response.Header().Get("Location"))
	require.Contains(t, response.Header().Get("Set-Cookie"), "__Host-auth-flow=binding")
	require.Contains(t, response.Header().Get("Set-Cookie"), "Max-Age=300")
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

	response := performRequest(handler.CompleteTelegram, http.MethodGet, "/api/v1/auth/telegram/callback?code=authorization-code&state=state", func(request *http.Request) {
		request.AddCookie(&http.Cookie{Name: DefaultBrowserBindingCookieName, Value: "binding", Secure: true, HttpOnly: true, SameSite: http.SameSiteLaxMode})
	})

	require.Equal(t, http.StatusFound, response.Code)
	require.Equal(t, "/cv/lk", response.Header().Get("Location"))
	cookies := response.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == DefaultCookieName {
			sessionCookie = cookie
			break
		}
	}
	require.NotNil(t, sessionCookie)
	require.Equal(t, "test-session-value", sessionCookie.Value)
	require.Equal(t, "/", sessionCookie.Path)
	require.Equal(t, 2592000, sessionCookie.MaxAge)
	require.True(t, sessionCookie.HttpOnly)
	require.True(t, sessionCookie.Secure)
	require.Equal(t, http.SameSiteLaxMode, sessionCookie.SameSite)
	require.NotContains(t, response.Body.String(), "authorization-code")
	require.NotContains(t, response.Body.String(), userID.String())
}

func TestHandlerCompleteTelegramRejectsMissingCallbackValues(t *testing.T) {
	handler := newTestHandler(nil, &fakeCompleteTelegramLogin{}, nil)

	response := performRequest(handler.CompleteTelegram, http.MethodGet, "/api/v1/auth/telegram/callback?code=authorization-code", nil)

	require.Equal(t, http.StatusFound, response.Code)
	require.Equal(t, "/cv/auth/error?reason=invalid_request", response.Header().Get("Location"))
}

func TestHandlerCompleteTelegramRejectsInvalidState(t *testing.T) {
	handler := newTestHandler(nil, &fakeCompleteTelegramLogin{err: application.ErrStateInvalid}, nil)

	response := performRequest(handler.CompleteTelegram, http.MethodGet, "/api/v1/auth/telegram/callback?code=code&state=state", func(request *http.Request) {
		request.AddCookie(&http.Cookie{Name: DefaultBrowserBindingCookieName, Value: "binding", Secure: true, HttpOnly: true, SameSite: http.SameSiteLaxMode})
	})

	require.Equal(t, http.StatusFound, response.Code)
	require.Equal(t, "/cv/auth/error?reason=state_invalid", response.Header().Get("Location"))
	require.NotContains(t, response.Body.String(), application.ErrStateInvalid.Error())
}

func TestHandlerCompleteTelegramRedirectsExpiredTransaction(t *testing.T) {
	handler := newTestHandler(nil, &fakeCompleteTelegramLogin{err: domain.ErrAuthTransactionExpired}, nil)

	response := performRequest(handler.CompleteTelegram, http.MethodGet, "/api/v1/auth/telegram/callback?code=code&state=state", func(request *http.Request) {
		request.AddCookie(&http.Cookie{Name: DefaultBrowserBindingCookieName, Value: "binding", Secure: true, HttpOnly: true, SameSite: http.SameSiteLaxMode})
	})

	require.Equal(t, http.StatusFound, response.Code)
	require.Equal(t, "/cv/auth/error?reason=transaction_expired", response.Header().Get("Location"))
}

func TestHandlerCompleteTelegramRedirectsConsumedTransaction(t *testing.T) {
	handler := newTestHandler(nil, &fakeCompleteTelegramLogin{err: domain.ErrAuthTransactionConsumed}, nil)

	response := performRequest(handler.CompleteTelegram, http.MethodGet, "/api/v1/auth/telegram/callback?code=code&state=state", func(request *http.Request) {
		request.AddCookie(&http.Cookie{Name: DefaultBrowserBindingCookieName, Value: "binding", Secure: true, HttpOnly: true, SameSite: http.SameSiteLaxMode})
	})

	require.Equal(t, http.StatusFound, response.Code)
	require.Equal(t, "/cv/auth/error?reason=transaction_consumed", response.Header().Get("Location"))
}

func TestHandlerCompleteTelegramProviderDenialRedirects(t *testing.T) {
	handler := newTestHandler(nil, &fakeCompleteTelegramLogin{}, nil)

	response := performRequest(handler.CompleteTelegram, http.MethodGet, "/api/v1/auth/telegram/callback?error=access_denied&state=state", func(request *http.Request) {
		request.AddCookie(&http.Cookie{Name: DefaultBrowserBindingCookieName, Value: "binding", Secure: true, HttpOnly: true, SameSite: http.SameSiteLaxMode})
	})

	require.Equal(t, http.StatusFound, response.Code)
	require.Equal(t, "/cv/auth/error?reason=provider_denied", response.Header().Get("Location"))
}

func TestHandlerCompleteTelegramRejectsUnboundProviderDenial(t *testing.T) {
	handler := newTestHandler(nil, &fakeCompleteTelegramLogin{}, nil)

	response := performRequest(handler.CompleteTelegram, http.MethodGet, "/api/v1/auth/telegram/callback?error=access_denied&state=state", nil)

	require.Equal(t, http.StatusFound, response.Code)
	require.Equal(t, "/cv/auth/error?reason=invalid_request", response.Header().Get("Location"))
}

func TestHandlerCompleteTelegramMapsProviderFailureWithoutDetails(t *testing.T) {
	handler := newTestHandler(nil, &fakeCompleteTelegramLogin{err: errors.New("telegram secret token=should-not-leak")}, nil)

	response := performRequest(handler.CompleteTelegram, http.MethodGet, "/api/v1/auth/telegram/callback?code=code&state=state", func(request *http.Request) {
		request.AddCookie(&http.Cookie{Name: DefaultBrowserBindingCookieName, Value: "binding", Secure: true, HttpOnly: true, SameSite: http.SameSiteLaxMode})
	})

	require.Equal(t, http.StatusFound, response.Code)
	require.Equal(t, "/cv/auth/error?reason=database_error", response.Header().Get("Location"))
	require.NotContains(t, response.Body.String(), "telegram secret")
	require.NotContains(t, response.Body.String(), "should-not-leak")
}

func TestHandlerCompleteTelegramMapsIdentityVerificationFailure(t *testing.T) {
	handler := newTestHandler(nil, &fakeCompleteTelegramLogin{err: application.ErrIdentityVerificationFailed}, nil)

	response := performRequest(handler.CompleteTelegram, http.MethodGet, "/api/v1/auth/telegram/callback?code=code&state=state", func(request *http.Request) {
		request.AddCookie(&http.Cookie{Name: DefaultBrowserBindingCookieName, Value: "binding", Secure: true, HttpOnly: true, SameSite: http.SameSiteLaxMode})
	})

	require.Equal(t, http.StatusFound, response.Code)
	require.Equal(t, "/cv/auth/error?reason=id_token_invalid", response.Header().Get("Location"))
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

func TestBrowserBindingCookieUsesHostPrefixAndSecurityAttributes(t *testing.T) {
	cookie := DefaultCookieConfig(300).BrowserBindingCookie("binding", 300)
	require.Equal(t, DefaultBrowserBindingCookieName, cookie.Name)
	require.Equal(t, "/", cookie.Path)
	require.True(t, cookie.Secure)
	require.True(t, cookie.HttpOnly)
	require.Equal(t, http.SameSiteLaxMode, cookie.SameSite)
	require.Equal(t, 300, cookie.MaxAge)
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

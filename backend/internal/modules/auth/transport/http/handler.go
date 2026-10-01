package http

import (
	"context"
	"errors"
	"strings"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"developers-cv/backend/internal/modules/auth/application"
	"developers-cv/backend/internal/modules/auth/application/commands"
	"developers-cv/backend/internal/modules/auth/application/ports"
	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	apperrors "developers-cv/backend/internal/platform/http/errors"
	"developers-cv/backend/internal/platform/http/middleware"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type startTelegramLogin interface {
	Handle(context.Context) (commands.TelegramLoginStart, error)
}

type completeTelegramLogin interface {
	Handle(context.Context, string, string, string) (commands.AuthenticatedSession, error)
}

type logoutHandler interface {
	Handle(context.Context, valueobjects.SessionTokenHash, time.Time) error
}

type Handler struct {
	startTelegramLogin   startTelegramLogin
	completeTelegram     completeTelegramLogin
	logout               logoutHandler
	observer             ports.AuthObserver
	tokenGenerator       ports.TokenGenerator
	cookie               CookieConfig
	successLocation      string
	failureLocation      string
	browserBindingMaxAge int
	now                  func() time.Time
}

type HandlerDependencies struct {
	StartTelegramLogin   startTelegramLogin
	CompleteTelegram     completeTelegramLogin
	Logout               logoutHandler
	Observer             ports.AuthObserver
	TokenGenerator       ports.TokenGenerator
	Cookie               CookieConfig
	SuccessLocation      string
	FailureLocation      string
	BrowserBindingMaxAge int
}

func NewHandler(deps HandlerDependencies) *Handler {
	successLocation := deps.SuccessLocation
	if successLocation == "" {
		successLocation = "/cv/lk"
	}
	return &Handler{
		startTelegramLogin:   deps.StartTelegramLogin,
		completeTelegram:     deps.CompleteTelegram,
		logout:               deps.Logout,
		observer:             deps.Observer,
		tokenGenerator:       deps.TokenGenerator,
		cookie:               deps.Cookie,
		successLocation:      successLocation,
		failureLocation:      deps.FailureLocation,
		browserBindingMaxAge: deps.BrowserBindingMaxAge,
		now:                  func() time.Time { return time.Now().UTC() },
	}
}

func (h *Handler) StartTelegram(c *gin.Context) {
	if h == nil || h.startTelegramLogin == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("start telegram login handler is not configured")))
		return
	}
	provider, providerErr := domain.NewProvider(string(domain.ProviderTelegram))
	if providerErr != nil {
		apperrors.WriteError(c, apperrors.Internal(providerErr))
		return
	}
	if h.observer != nil {
		h.observer.LoginStarted(c.Request.Context(), provider)
	}
	finishSpan := func(error) {}
	if h.observer != nil {
		spanCtx, finish := h.observer.StartSpan(c.Request.Context(), "auth.start", provider)
		finishSpan = finish
		c.Request = c.Request.WithContext(spanCtx)
	}
	startedAt := h.now()
	start, err := h.startTelegramLogin.Handle(c.Request.Context())
	if err != nil {
		finishSpan(err)
		if h.observer != nil {
			h.observer.LoginFailed(c.Request.Context(), provider, application.FailureReason(err), time.Since(startedAt))
		}
		writeAuthError(c, err)
		return
	}
	finishSpan(nil)
	maxAge := start.BrowserBindingMaxAge
	if maxAge <= 0 {
		maxAge = h.browserBindingMaxAge
	}
	if strings.TrimSpace(start.BrowserBindingToken) == "" || maxAge <= 0 {
		err := application.ErrInvalidAuthFlow
		if h.observer != nil {
			h.observer.LoginFailed(c.Request.Context(), provider, application.FailureReason(err), time.Since(startedAt))
		}
		writeAuthError(c, err)
		return
	}
	http.SetCookie(c.Writer, h.cookie.BrowserBindingCookie(start.BrowserBindingToken, maxAge))
	c.Redirect(http.StatusFound, start.AuthorizationURL)
}

func (h *Handler) CompleteTelegram(c *gin.Context) {
	if h == nil || h.completeTelegram == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("complete telegram login handler is not configured")))
		return
	}
	if err := h.cookie.Validate(false); err != nil {
		apperrors.WriteError(c, apperrors.Internal(err))
		return
	}

	provider, providerErr := domain.NewProvider(string(domain.ProviderTelegram))
	if providerErr != nil {
		apperrors.WriteError(c, apperrors.Internal(providerErr))
		return
	}
	if h.observer != nil {
		h.observer.LoginCallback(c.Request.Context(), provider)
	}
	finishSpan := func(error) {}
	if h.observer != nil {
		spanCtx, finish := h.observer.StartSpan(c.Request.Context(), "auth.callback", provider)
		finishSpan = finish
		c.Request = c.Request.WithContext(spanCtx)
	}
	startedAt := h.now()
	providerError := strings.TrimSpace(c.Query("error"))
	code := strings.TrimSpace(c.Query("code"))
	state := strings.TrimSpace(c.Query("state"))
	browserBinding, bindingErr := browserBindingToken(c.Request, h.cookie.Secure)
	if state == "" || bindingErr != nil {
		err := application.ErrInvalidAuthFlow
		finishSpan(err)
		if h.observer != nil {
			h.observer.LoginFailed(c.Request.Context(), provider, application.FailureReason(err), time.Since(startedAt))
		}
		h.redirectAuthError(c, err)
		return
	}
	if providerError == "access_denied" {
		finishSpan(nil)
		if h.observer != nil {
			h.observer.LoginFailed(c.Request.Context(), provider, application.FailureReasonProviderDenied, time.Since(startedAt))
		}
		http.SetCookie(c.Writer, h.cookie.ClearBrowserBindingCookie())
		c.Redirect(http.StatusFound, h.failureRedirect("provider_denied"))
		return
	}
	if providerError != "" {
		finishSpan(nil)
		if h.observer != nil {
			h.observer.LoginFailed(c.Request.Context(), provider, application.FailureReasonIDTokenInvalid, time.Since(startedAt))
		}
		http.SetCookie(c.Writer, h.cookie.ClearBrowserBindingCookie())
		c.Redirect(http.StatusFound, h.failureRedirect("authentication_failed"))
		return
	}
	if code == "" {
		err := application.ErrInvalidAuthFlow
		finishSpan(err)
		if h.observer != nil {
			h.observer.LoginFailed(c.Request.Context(), provider, application.FailureReason(err), time.Since(startedAt))
		}
		h.redirectAuthError(c, err)
		return
	}

	session, err := h.completeTelegram.Handle(c.Request.Context(), code, state, browserBinding)
	if err != nil {
		finishSpan(err)
		if h.observer != nil {
			h.observer.LoginFailed(c.Request.Context(), provider, application.FailureReason(err), time.Since(startedAt))
		}
		h.redirectAuthError(c, err)
		return
	}
	finishSpan(nil)
	if h.observer != nil {
		h.observer.LoginSucceeded(c.Request.Context(), provider, session.UserID, time.Since(startedAt))
	}
	c.SetSameSite(h.cookie.SameSite)
	http.SetCookie(c.Writer, h.cookie.ClearBrowserBindingCookie())
	http.SetCookie(c.Writer, h.cookie.NewCookie(session.RawSessionToken))
	c.Redirect(http.StatusFound, h.successLocation)
}

func (h *Handler) CurrentSession(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		apperrors.Unauthorized(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"userId":        userID.String(),
		"authenticated": true,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	if h == nil || h.logout == nil || h.tokenGenerator == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("logout handler is not configured")))
		return
	}
	if err := h.cookie.Validate(false); err != nil {
		apperrors.WriteError(c, apperrors.Internal(err))
		return
	}

	hash, err := sessionTokenHash(c.Request, h.cookie.Name, h.tokenGenerator)
	if err != nil && !errors.Is(err, errInvalidSessionCookie) {
		apperrors.WriteError(c, apperrors.Internal(err))
		return
	}
	userID, _ := middleware.GetUserID(c)
	provider, providerErr := domain.NewProvider(string(domain.ProviderTelegram))
	if providerErr != nil {
		apperrors.WriteError(c, apperrors.Internal(providerErr))
		return
	}
	finishSpan := func(error) {}
	if h.observer != nil {
		spanCtx, finish := h.observer.StartSpan(c.Request.Context(), "auth.logout", provider)
		finishSpan = finish
		c.Request = c.Request.WithContext(spanCtx)
	}
	if err == nil {
		if err := h.logout.Handle(c.Request.Context(), hash, h.now().UTC()); err != nil {
			finishSpan(err)
			if h.observer != nil {
				h.observer.Logout(c.Request.Context(), userIDFromUUID(userID), false)
			}
			writeAuthError(c, err)
			return
		}
	}
	finishSpan(nil)
	if h.observer != nil {
		h.observer.Logout(c.Request.Context(), userIDFromUUID(userID), true)
	}

	http.SetCookie(c.Writer, h.cookie.ClearCookie())
	c.Status(http.StatusNoContent)
}

func (h *Handler) failureRedirect(reason string) string {
	base := h.failureLocation
	if base == "" {
		base = "/cv/auth/error"
	}
	return base + "?reason=" + reason
}

func userIDFromUUID(value uuid.UUID) identity.ID {
	id, err := identity.IDFromUUID(value)
	if err != nil {
		return identity.ID{}
	}
	return id
}

func (h *Handler) redirectAuthError(c *gin.Context, err error) {
	reason := application.FailureReason(err)
	if reason == "" {
		reason = application.FailureReasonInvalidRequest
	}
	if !errors.Is(err, application.ErrBrowserBindingInvalid) && !errors.Is(err, application.ErrStateInvalid) {
		http.SetCookie(c.Writer, h.cookie.ClearBrowserBindingCookie())
	}
	c.Redirect(http.StatusFound, h.failureRedirect(reason))
}

func writeAuthError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, middleware.ErrUnauthenticated),
		errors.Is(err, application.ErrSessionInvalid),
		errors.Is(err, application.ErrSessionExpired),
		errors.Is(err, application.ErrSessionIdleExpired),
		errors.Is(err, application.ErrSessionRevoked):
		apperrors.WriteError(c, apperrors.New(apperrors.CodeUnauthorized, http.StatusUnauthorized, "Authentication required"))
	case errors.Is(err, application.ErrProviderUnavailable):
		apperrors.WriteError(c, apperrors.New(apperrors.CodeUnavailable, http.StatusServiceUnavailable, "Authentication provider unavailable"))
	case errors.Is(err, application.ErrInvalidAuthFlow),
		errors.Is(err, application.ErrStateInvalid),
		errors.Is(err, application.ErrBrowserBindingInvalid),
		errors.Is(err, domain.ErrAuthTransactionExpired),
		errors.Is(err, domain.ErrAuthTransactionConsumed),
		errors.Is(err, application.ErrIdentityVerificationFailed),
		errors.Is(err, application.ErrConcurrentIdentityCreation):
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Authentication failed"))
	default:
		apperrors.WriteError(c, apperrors.Internal(err))
	}
}

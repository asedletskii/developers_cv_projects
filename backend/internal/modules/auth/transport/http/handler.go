package http

import (
	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/auth/application/ports"
	"developers-cv/backend/internal/modules/auth/transport/http/handlers"
)

type CookieConfig = handlers.CookieConfig

const (
	DefaultCookieName               = handlers.DefaultCookieName
	DefaultBrowserBindingCookieName = handlers.DefaultBrowserBindingCookieName
)

func DefaultCookieConfig(maxAge int) CookieConfig {
	return handlers.DefaultCookieConfig(maxAge)
}

type Handler struct {
	inner *handlers.AuthHandler
}

type HandlerDependencies struct {
	StartTelegramLogin   handlers.StartTelegramLogin
	CompleteTelegram     handlers.CompleteTelegramLogin
	Logout               handlers.Logout
	Observer             ports.AuthObserver
	TokenGenerator       ports.TokenGenerator
	Cookie               CookieConfig
	SuccessLocation      string
	FailureLocation      string
	BrowserBindingMaxAge int
}

func NewHandler(deps HandlerDependencies) *Handler {
	return &Handler{inner: handlers.NewAuthHandler(handlers.AuthHandlerDependencies{
		StartTelegramLogin:   deps.StartTelegramLogin,
		CompleteTelegram:     deps.CompleteTelegram,
		Logout:               deps.Logout,
		Observer:             deps.Observer,
		TokenGenerator:       deps.TokenGenerator,
		Cookie:               deps.Cookie,
		SuccessLocation:      deps.SuccessLocation,
		FailureLocation:      deps.FailureLocation,
		BrowserBindingMaxAge: deps.BrowserBindingMaxAge,
	})}
}

func (h *Handler) StartTelegram(c *gin.Context)    { h.inner.StartTelegram(c) }
func (h *Handler) CompleteTelegram(c *gin.Context) { h.inner.CompleteTelegram(c) }
func (h *Handler) CurrentSession(c *gin.Context)   { h.inner.CurrentSession(c) }
func (h *Handler) Logout(c *gin.Context)           { h.inner.Logout(c) }

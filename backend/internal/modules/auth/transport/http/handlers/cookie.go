package handlers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"developers-cv/backend/internal/modules/auth/application/ports"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
)

const DefaultCookieName = "__Host-session"
const DefaultBrowserBindingCookieName = "__Host-auth-flow"

var errInvalidSessionCookie = errors.New("invalid auth session cookie")

type CookieConfig struct {
	Name     string
	Path     string
	Domain   string
	Secure   bool
	HTTPOnly bool
	SameSite http.SameSite
	MaxAge   int
}

func DefaultCookieConfig(maxAge int) CookieConfig {
	return CookieConfig{
		Name:     DefaultCookieName,
		Path:     "/",
		Secure:   true,
		HTTPOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
	}
}

func (c CookieConfig) Validate(production bool) error {
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("auth cookie name is empty")
	}
	if c.Path != "/" {
		return errors.New("auth cookie path must be /")
	}
	if c.MaxAge <= 0 {
		return errors.New("auth cookie max age must be positive")
	}
	if !c.HTTPOnly {
		return errors.New("auth cookie must be httpOnly")
	}
	if c.SameSite != http.SameSiteLaxMode {
		return errors.New("auth cookie sameSite must be lax")
	}
	if production && !c.Secure {
		return errors.New("auth cookie must be secure in production")
	}
	if production && c.Name != DefaultCookieName {
		return errors.New("production auth cookie name must be __Host-session")
	}
	if strings.HasPrefix(c.Name, "__Host-") && (c.Domain != "" || c.Path != "/" || !c.Secure) {
		return errors.New("__Host- auth cookie requires Secure, Path=/ and no Domain")
	}
	return nil
}

func (c CookieConfig) NewCookie(value string) *http.Cookie {
	// #nosec G124
	return &http.Cookie{
		Name:     c.Name,
		Value:    value,
		Path:     c.Path,
		Domain:   c.Domain,
		MaxAge:   c.MaxAge,
		Secure:   c.Secure,
		HttpOnly: c.HTTPOnly,
		SameSite: c.SameSite,
	}
}

func (c CookieConfig) ClearCookie() *http.Cookie {
	// #nosec G124
	cookie := c.NewCookie("")
	cookie.MaxAge = -1
	return cookie
}

func (c CookieConfig) BrowserBindingCookie(value string, maxAge int) *http.Cookie {
	name := DefaultBrowserBindingCookieName
	if !c.Secure {
		name = "auth-flow"
	}
	// #nosec G124
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		Secure:   c.Secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func (c CookieConfig) ClearBrowserBindingCookie() *http.Cookie {
	// #nosec G124
	cookie := c.BrowserBindingCookie("", -1)
	cookie.MaxAge = -1
	return cookie
}

func browserBindingCookieName(secure bool) string {
	if secure {
		return DefaultBrowserBindingCookieName
	}
	return "auth-flow"
}

func browserBindingToken(request *http.Request, secure bool) (string, error) {
	if request == nil {
		return "", errInvalidSessionCookie
	}
	cookie, err := request.Cookie(browserBindingCookieName(secure))
	if err != nil {
		return "", errInvalidSessionCookie
	}
	raw := strings.TrimSpace(cookie.Value)
	if raw == "" {
		return "", errInvalidSessionCookie
	}
	return raw, nil
}

func sessionTokenHash(request *http.Request, cookieName string, tokenGenerator ports.TokenGenerator) (valueobjects.SessionTokenHash, error) {
	if request == nil || tokenGenerator == nil {
		return valueobjects.SessionTokenHash{}, errInvalidSessionCookie
	}
	cookie, err := request.Cookie(cookieName)
	if err != nil {
		return valueobjects.SessionTokenHash{}, errInvalidSessionCookie
	}
	raw := strings.TrimSpace(cookie.Value)
	if raw == "" {
		return valueobjects.SessionTokenHash{}, errInvalidSessionCookie
	}
	token, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil || len(token) == 0 {
		return valueobjects.SessionTokenHash{}, errInvalidSessionCookie
	}
	hash, err := valueobjects.NewSessionTokenHash(tokenGenerator.Hash(token))
	if err != nil {
		return valueobjects.SessionTokenHash{}, fmt.Errorf("create session token hash: %w", err)
	}
	return hash, nil
}

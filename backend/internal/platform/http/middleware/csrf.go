package middleware

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	apperrors "developers-cv/backend/internal/platform/http/errors"
)

func CSRFProtection(allowedOrigins []string) gin.HandlerFunc {
	allowed := normalizeAllowedOrigins(allowedOrigins)
	return func(c *gin.Context) {
		if isSafeCSRFMethod(c.Request.Method) {
			c.Next()
			return
		}
		if isValidCSRFRequest(c, allowed) {
			c.Next()
			return
		}
		apperrors.WriteError(c, apperrors.New(
			apperrors.CodeForbidden,
			http.StatusForbidden,
			"CSRF protection: invalid request",
		))
	}
}

func isSafeCSRFMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	default:
		return false
	}
}

func isValidCSRFRequest(c *gin.Context, allowedOrigins map[string]struct{}) bool {
	if !strings.EqualFold(strings.TrimSpace(c.GetHeader("X-Requested-With")), "XMLHttpRequest") {
		return false
	}

	switch strings.ToLower(strings.TrimSpace(c.GetHeader("Sec-Fetch-Site"))) {
	case "cross-site", "none":
		return false
	}

	if origin := strings.TrimSpace(c.GetHeader("Origin")); origin != "" {
		return originIsAllowed(origin, allowedOrigins)
	}

	if referer := strings.TrimSpace(c.GetHeader("Referer")); referer != "" {
		return originIsAllowed(referer, allowedOrigins)
	}

	return true
}

func originIsAllowed(rawOrigin string, allowedOrigins map[string]struct{}) bool {
	parsed, ok := parseOrigin(rawOrigin)
	if !ok {
		return false
	}
	_, ok = allowedOrigins[originKey(parsed)]
	return ok
}

func parseOrigin(raw string) (*url.URL, bool) {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || parsed.User != nil || parsed.Host == "" {
		return nil, false
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, false
	}
	return parsed, true
}

func originKey(parsed *url.URL) string {
	return strings.ToLower(parsed.Scheme + "://" + parsed.Host)
}

func normalizeAllowedOrigins(origins []string) map[string]struct{} {
	result := make(map[string]struct{}, len(origins))
	for _, raw := range origins {
		parsed, ok := parseOrigin(raw)
		if !ok {
			continue
		}
		result[originKey(parsed)] = struct{}{}
	}
	return result
}

package analytics

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"net/netip"
	"strings"
	"time"

	"github.com/google/uuid"

	analyticscontracts "developers-cv/backend/internal/modules/analytics/contracts"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

const (
	CookiePath   = "/api/v1/cv"
	CookieMaxAge = 365 * 24 * 60 * 60
)

type VisitorIdentifier struct {
	secret     []byte
	production bool
}

func NewVisitorIdentifier(secret string, production bool) (*VisitorIdentifier, error) {
	secret = strings.TrimSpace(secret)
	if len(secret) < 32 {
		return nil, errors.New("analytics visitor secret must be at least 32 bytes")
	}
	return &VisitorIdentifier{secret: []byte(secret), production: production}, nil
}

func (i *VisitorIdentifier) Identify(now time.Time, input analyticscontracts.PublicVisitInput) (analyticscontracts.VisitorIdentityResult, error) {
	if i == nil || len(i.secret) < 32 {
		return analyticscontracts.VisitorIdentityResult{}, errors.New("visitor identifier is not configured")
	}
	ipHash := i.dailyHMAC(now, "ip", normalizeIPForAnalytics(input.IP))
	secure := i.production
	if input.HasCookie && input.VisitorID != (sharedvalueobjects.ID{}) {
		return analyticscontracts.VisitorIdentityResult{VisitorID: input.VisitorID, IPHash: ipHash}, nil
	}
	visitorUUID := uuid.New()
	visitorID, err := sharedvalueobjects.IDFromUUID(visitorUUID)
	if err != nil {
		return analyticscontracts.VisitorIdentityResult{}, fmt.Errorf("create fallback visitor id: %w", err)
	}

	result := analyticscontracts.VisitorIdentityResult{VisitorID: visitorID, IPHash: ipHash}
	result.SetCookie = true
	if secure {
		result.Cookie = &http.Cookie{
			Name:     analyticscontracts.ProductionVisitorCookieName,
			Value:    visitorUUID.String(),
			Path:     CookiePath,
			MaxAge:   CookieMaxAge,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}
	} else {
		result.Cookie = &http.Cookie{ // #nosec G124
			Name:     analyticscontracts.LocalVisitorCookieName,
			Value:    visitorUUID.String(),
			Path:     CookiePath,
			MaxAge:   CookieMaxAge,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		}
	}
	return result, nil
}

func normalizeIPForAnalytics(value string) string {
	addr, err := netip.ParseAddr(strings.TrimSpace(value))
	if err != nil {
		return value
	}
	addr = addr.Unmap()
	if addr.Is6() {
		addr = addr.WithZone("")
		prefix, err := addr.Prefix(64)
		if err == nil {
			return prefix.Masked().String()
		}
	}
	return addr.String()
}

func (i *VisitorIdentifier) dailyHMAC(now time.Time, kind, value string) []byte {
	date := now.UTC().Format("2006-01-02")
	master := hmac.New(sha256.New, i.secret)
	master.Write([]byte(date))
	dailySalt := master.Sum(nil)
	mac := hmac.New(sha256.New, dailySalt)
	mac.Write([]byte(kind))
	mac.Write([]byte{0})
	mac.Write([]byte(value))
	return mac.Sum(nil)
}

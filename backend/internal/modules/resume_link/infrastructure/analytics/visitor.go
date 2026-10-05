package analytics

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"developers-cv/backend/internal/modules/resume_link/application/ports"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

const (
	ProductionCookieName = "__Secure-cv_vid"
	LocalCookieName      = "cv_vid"
	CookiePath           = "/api/v1/cv"
	CookieMaxAge         = 365 * 24 * 60 * 60
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

func (i *VisitorIdentifier) Identify(now time.Time, input ports.PublicVisitInput) (ports.VisitorIdentityResult, error) {
	if i == nil || len(i.secret) < 32 {
		return ports.VisitorIdentityResult{}, errors.New("visitor identifier is not configured")
	}
	ipHash := i.dailyHMAC(now, "ip", input.IP)
	var secure bool
	if i.production {
		secure = true
	} else {
		secure = false
	}
	if !input.DNT && !input.GPC && input.HasCookie && input.VisitorID != (sharedvalueobjects.ID{}) {
		return ports.VisitorIdentityResult{VisitorID: input.VisitorID, IPHash: ipHash}, nil
	}
	uaHash := i.dailyHMAC(now, "ua", input.UserAgent)
	name := append(append([]byte{}, ipHash...), uaHash...)
	name = append(name, []byte(now.UTC().Format("2006-01-02"))...)
	visitorUUID := uuid.NewSHA1(uuid.NameSpaceOID, name)
	visitorID, err := sharedvalueobjects.IDFromUUID(visitorUUID)
	if err != nil {
		return ports.VisitorIdentityResult{}, fmt.Errorf("create fallback visitor id: %w", err)
	}

	result := ports.VisitorIdentityResult{VisitorID: visitorID, IPHash: ipHash}
	if input.DNT || input.GPC {
		return result, nil
	}

	result.SetCookie = true
	if secure {
		result.Cookie = &http.Cookie{
			Name:     ProductionCookieName,
			Value:    visitorUUID.String(),
			Path:     CookiePath,
			MaxAge:   CookieMaxAge,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}
	} else {
		result.Cookie = &http.Cookie{ // #nosec G124
			Name:     LocalCookieName,
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

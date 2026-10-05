package analytics

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume_link/application/ports"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestVisitorIdentifierUsesCookieAndSetsNewCookie(t *testing.T) {
	identifier, err := NewVisitorIdentifier("abcdefghijklmnopqrstuvwxyz0123456789", false)
	require.NoError(t, err)
	now := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
	first, err := identifier.Identify(now, ports.PublicVisitInput{IP: "127.0.0.1", UserAgent: "Mozilla/5.0"})
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, first.VisitorID.UUID())
	require.True(t, first.SetCookie)
	require.Equal(t, ports.LocalVisitorCookieName, first.Cookie.Name)
	cookieUUID := uuid.MustParse(first.Cookie.Value)
	require.Equal(t, byte(4), byte(cookieUUID.Version()))
	cookieID, err := sharedvalueobjects.IDFromUUID(cookieUUID)
	require.NoError(t, err)
	require.Equal(t, cookieID, first.VisitorID)
	second, err := identifier.Identify(now, ports.PublicVisitInput{IP: "127.0.0.1", UserAgent: "Mozilla/5.0", VisitorID: cookieID, HasCookie: true})
	require.NoError(t, err)
	require.Equal(t, cookieID, second.VisitorID)
	require.False(t, second.SetCookie)
}

func TestVisitorIdentifierProductionCookieIsSecure(t *testing.T) {
	identifier, err := NewVisitorIdentifier("abcdefghijklmnopqrstuvwxyz0123456789", true)
	require.NoError(t, err)

	result, err := identifier.Identify(
		time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC),
		ports.PublicVisitInput{IP: "127.0.0.1", UserAgent: "Mozilla/5.0"},
	)
	require.NoError(t, err)
	require.True(t, result.SetCookie)
	require.NotNil(t, result.Cookie)
	require.Equal(t, ports.ProductionVisitorCookieName, result.Cookie.Name)
	require.True(t, result.Cookie.Secure)
	require.True(t, result.Cookie.HttpOnly)
	require.Equal(t, http.SameSiteLaxMode, result.Cookie.SameSite)
}

func TestVisitorIdentifierAlwaysSetsCookieRegardlessOfPrivacyHeaders(t *testing.T) {
	identifier, err := NewVisitorIdentifier("abcdefghijklmnopqrstuvwxyz0123456789", false)
	require.NoError(t, err)
	result, err := identifier.Identify(time.Now().UTC(), ports.PublicVisitInput{IP: "127.0.0.1", UserAgent: "Mozilla/5.0"})
	require.NoError(t, err)
	require.True(t, result.SetCookie)
	require.NotNil(t, result.Cookie)
}

func TestVisitorIdentifierFallbackIsRandomWithoutCookie(t *testing.T) {
	identifier, err := NewVisitorIdentifier("abcdefghijklmnopqrstuvwxyz0123456789", false)
	require.NoError(t, err)
	now := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
	input := ports.PublicVisitInput{IP: "127.0.0.1", UserAgent: "Mozilla/5.0"}
	first, err := identifier.Identify(now, input)
	require.NoError(t, err)
	second, err := identifier.Identify(now.Add(5*time.Hour), input)
	require.NoError(t, err)
	require.NotEqual(t, first.VisitorID, second.VisitorID)
	require.NotEqual(t, uuid.Nil, first.VisitorID.UUID())
	require.NotEqual(t, uuid.Nil, second.VisitorID.UUID())

	otherVisitor, err := identifier.Identify(now, ports.PublicVisitInput{IP: "127.0.0.1", UserAgent: "Mozilla/5.0"})
	require.NoError(t, err)
	require.NotEqual(t, first.VisitorID, otherVisitor.VisitorID)
}

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
	require.Equal(t, LocalCookieName, first.Cookie.Name)
	cookieUUID := uuid.MustParse(first.Cookie.Value)
	require.Equal(t, byte(5), byte(cookieUUID.Version()))
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
	require.Equal(t, ProductionCookieName, result.Cookie.Name)
	require.True(t, result.Cookie.Secure)
	require.True(t, result.Cookie.HttpOnly)
	require.Equal(t, http.SameSiteLaxMode, result.Cookie.SameSite)
}

func TestVisitorIdentifierDoesNotSetCookieForDNT(t *testing.T) {
	identifier, err := NewVisitorIdentifier("abcdefghijklmnopqrstuvwxyz0123456789", false)
	require.NoError(t, err)
	result, err := identifier.Identify(time.Now().UTC(), ports.PublicVisitInput{IP: "127.0.0.1", DNT: true})
	require.NoError(t, err)
	require.False(t, result.SetCookie)
	require.NotEqual(t, uuid.Nil, result.VisitorID.UUID())
}

func TestVisitorIdentifierDNTFallbackIsStableWithinDay(t *testing.T) {
	identifier, err := NewVisitorIdentifier("abcdefghijklmnopqrstuvwxyz0123456789", false)
	require.NoError(t, err)
	now := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
	input := ports.PublicVisitInput{IP: "127.0.0.1", UserAgent: "Mozilla/5.0", DNT: true}
	first, err := identifier.Identify(now, input)
	require.NoError(t, err)
	second, err := identifier.Identify(now.Add(5*time.Hour), input)
	require.NoError(t, err)
	require.Equal(t, first.VisitorID, second.VisitorID)

	nextDay, err := identifier.Identify(now.Add(24*time.Hour), input)
	require.NoError(t, err)
	require.NotEqual(t, first.VisitorID, nextDay.VisitorID)
}

func TestVisitorIdentifierFallbackIsStableWithoutCookie(t *testing.T) {
	identifier, err := NewVisitorIdentifier("abcdefghijklmnopqrstuvwxyz0123456789", false)
	require.NoError(t, err)
	now := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
	input := ports.PublicVisitInput{IP: "127.0.0.1", UserAgent: "Mozilla/5.0"}
	first, err := identifier.Identify(now, input)
	require.NoError(t, err)
	second, err := identifier.Identify(now.Add(5*time.Hour), input)
	require.NoError(t, err)
	require.Equal(t, first.VisitorID, second.VisitorID)
	require.NotEqual(t, uuid.Nil, first.VisitorID.UUID())
	require.NotEqual(t, uuid.Nil, second.VisitorID.UUID())
	nextDay, err := identifier.Identify(now.Add(24*time.Hour), input)
	require.NoError(t, err)
	require.NotEqual(t, first.VisitorID, nextDay.VisitorID)
}

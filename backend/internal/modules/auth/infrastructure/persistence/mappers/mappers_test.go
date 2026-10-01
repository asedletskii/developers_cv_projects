package mappers

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestIdentityRoundTrip(t *testing.T) {
	at := time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC)
	subject, err := valueobjects.NewProviderSubject("telegram-user")
	require.NoError(t, err)
	identity, err := domain.NewIdentity(domain.IdentityParams{ID: sharedvalueobjects.NewID(), UserID: sharedvalueobjects.NewID(), Provider: domain.ProviderTelegram, ProviderSubject: subject, CreatedAt: at, UpdatedAt: at, LastLoginAt: at})
	require.NoError(t, err)
	model, err := IdentityToDB(identity)
	require.NoError(t, err)
	mapped, err := IdentityToDomain(model)
	require.NoError(t, err)
	require.Equal(t, identity.ID(), mapped.ID())
	require.Equal(t, identity.UserID(), mapped.UserID())
	require.Equal(t, identity.ProviderSubject(), mapped.ProviderSubject())
}

func TestSessionRoundTrip(t *testing.T) {
	at := time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC)
	hash, err := valueobjects.NewSessionTokenHash(bytes.Repeat([]byte{1}, valueobjects.HashSize))
	require.NoError(t, err)
	session, err := domain.NewSession(domain.SessionParams{ID: sharedvalueobjects.NewID(), UserID: sharedvalueobjects.NewID(), TokenHash: hash, CreatedAt: at, LastSeenAt: at, ExpiresAt: at.Add(time.Hour)})
	require.NoError(t, err)
	model, err := SessionToDB(session)
	require.NoError(t, err)
	mapped, err := SessionToDomain(model)
	require.NoError(t, err)
	require.Equal(t, session.ID(), mapped.ID())
	require.Equal(t, session.UserID(), mapped.UserID())
	require.Equal(t, session.TokenHash(), mapped.TokenHash())
}

func TestAuthTransactionRoundTrip(t *testing.T) {
	at := time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC)
	state, err := valueobjects.NewStateHash(bytes.Repeat([]byte{1}, valueobjects.HashSize))
	require.NoError(t, err)
	nonce, err := valueobjects.NewNonceHash(bytes.Repeat([]byte{2}, valueobjects.HashSize))
	require.NoError(t, err)
	binding, err := valueobjects.NewBrowserBindingHash(bytes.Repeat([]byte{3}, valueobjects.HashSize))
	require.NoError(t, err)
	transaction, err := domain.NewAuthTransaction(domain.AuthTransactionParams{ID: sharedvalueobjects.NewID(), Provider: domain.ProviderTelegram, StateHash: state, NonceHash: nonce, BrowserBindingHash: binding, CodeVerifier: "verifier", RedirectURI: "https://example.test/callback", CreatedAt: at, ExpiresAt: at.Add(5 * time.Minute)})
	require.NoError(t, err)
	model, err := AuthTransactionToDB(transaction)
	require.NoError(t, err)
	mapped, err := AuthTransactionToDomain(model)
	require.NoError(t, err)
	require.Equal(t, transaction.ID(), mapped.ID())
	require.Equal(t, transaction.StateHash(), mapped.StateHash())
	require.Equal(t, transaction.NonceHash(), mapped.NonceHash())
	require.Equal(t, transaction.BrowserBindingHash(), mapped.BrowserBindingHash())
}

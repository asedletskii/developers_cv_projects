package domain

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
)

func newHashes(t *testing.T) (valueobjects.StateHash, valueobjects.NonceHash) {
	t.Helper()
	state, err := valueobjects.NewStateHash(bytes.Repeat([]byte{1}, valueobjects.HashSize))
	require.NoError(t, err)
	nonce, err := valueobjects.NewNonceHash(bytes.Repeat([]byte{2}, valueobjects.HashSize))
	require.NoError(t, err)
	return state, nonce
}

func newTestTransaction(t *testing.T) (*AuthTransaction, time.Time) {
	t.Helper()
	created := time.Date(2026, 9, 30, 10, 0, 0, 0, time.UTC)
	state, nonce := newHashes(t)
	binding, err := valueobjects.NewBrowserBindingHash(bytes.Repeat([]byte{3}, valueobjects.HashSize))
	require.NoError(t, err)
	transaction, err := NewAuthTransaction(AuthTransactionParams{
		Provider: ProviderTelegram, StateHash: state, NonceHash: nonce, BrowserBindingHash: binding,
		CodeVerifier: "verifier", RedirectURI: "https://example.com/callback", CreatedAt: created,
		ExpiresAt: created.Add(5 * time.Minute),
	})
	require.NoError(t, err)
	return transaction, created
}

func TestNewAuthTransaction(t *testing.T) {
	transaction, created := newTestTransaction(t)
	require.False(t, transaction.IsConsumed())
	require.False(t, transaction.IsExpired(created.Add(time.Minute)))
	require.Equal(t, "verifier", transaction.CodeVerifier())
}

func TestNewAuthTransactionGeneratesIDAndRestorePreservesValues(t *testing.T) {
	created := time.Date(2026, 9, 30, 10, 0, 0, 0, time.UTC)
	state, nonce := newHashes(t)
	binding, err := valueobjects.NewBrowserBindingHash(bytes.Repeat([]byte{3}, valueobjects.HashSize))
	require.NoError(t, err)
	transaction, err := NewAuthTransaction(AuthTransactionParams{
		Provider: ProviderTelegram, StateHash: state, NonceHash: nonce, BrowserBindingHash: binding,
		CodeVerifier: " verifier ", RedirectURI: " https://example.com/callback ",
		CreatedAt: created, ExpiresAt: created.Add(5 * time.Minute),
	})
	require.NoError(t, err)
	require.False(t, transaction.ID().IsZero())
	require.Equal(t, "verifier", transaction.CodeVerifier())

	restored, err := RestoreAuthTransaction(AuthTransactionSnapshot{
		ID: transaction.ID(), Provider: transaction.Provider(), StateHash: state, NonceHash: nonce,
		BrowserBindingHash: binding, CodeVerifier: " verifier ", RedirectURI: " https://example.com/callback ",
		CreatedAt: created, ExpiresAt: created.Add(5 * time.Minute),
	})
	require.NoError(t, err)
	require.Equal(t, " verifier ", restored.CodeVerifier())
	require.Equal(t, " https://example.com/callback ", restored.RedirectURI())
}

func TestAuthTransactionConsume(t *testing.T) {
	transaction, created := newTestTransaction(t)
	consumeAt := created.Add(time.Minute)
	require.NoError(t, transaction.Consume(consumeAt))
	require.True(t, transaction.IsConsumed())
	require.Equal(t, consumeAt, *transaction.ConsumedAt())
	require.ErrorIs(t, transaction.Consume(consumeAt.Add(time.Second)), ErrAuthTransactionConsumed)
}

func TestAuthTransactionConsumeRejectsExpired(t *testing.T) {
	transaction, created := newTestTransaction(t)
	require.ErrorIs(t, transaction.Consume(created.Add(5*time.Minute)), ErrAuthTransactionExpired)
}

func TestNewAuthTransactionRejectsInvalidTimestamps(t *testing.T) {
	created := time.Date(2026, 9, 30, 10, 0, 0, 0, time.UTC)
	state, nonce := newHashes(t)
	binding, err := valueobjects.NewBrowserBindingHash(bytes.Repeat([]byte{3}, valueobjects.HashSize))
	require.NoError(t, err)
	_, err = NewAuthTransaction(AuthTransactionParams{
		Provider: ProviderTelegram, StateHash: state, NonceHash: nonce, BrowserBindingHash: binding,
		CodeVerifier: "verifier", RedirectURI: "https://example.com/callback", CreatedAt: created,
		ExpiresAt: created,
	})
	require.ErrorIs(t, err, ErrInvalidAuthTransaction)
}

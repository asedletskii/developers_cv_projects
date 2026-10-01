package domain

import (
	"fmt"
	"strings"
	"time"

	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type AuthTransaction struct {
	id                 identity.ID
	provider           Provider
	stateHash          valueobjects.StateHash
	nonceHash          valueobjects.NonceHash
	browserBindingHash valueobjects.BrowserBindingHash
	codeVerifier       string
	redirectURI        string
	createdAt          time.Time
	expiresAt          time.Time
	consumedAt         *time.Time
}

type AuthTransactionParams struct {
	ID                 identity.ID
	Provider           Provider
	StateHash          valueobjects.StateHash
	NonceHash          valueobjects.NonceHash
	BrowserBindingHash valueobjects.BrowserBindingHash
	CodeVerifier       string
	RedirectURI        string
	CreatedAt          time.Time
	ExpiresAt          time.Time
	ConsumedAt         *time.Time
}

func NewAuthTransaction(params AuthTransactionParams) (*AuthTransaction, error) {
	if params.ID.IsZero() || !params.Provider.IsValid() {
		return nil, fmt.Errorf("%w: invalid transaction identity", ErrInvalidAuthTransaction)
	}
	codeVerifier := strings.TrimSpace(params.CodeVerifier)
	redirectURI := strings.TrimSpace(params.RedirectURI)
	if codeVerifier == "" || redirectURI == "" {
		return nil, fmt.Errorf("%w: code verifier and redirect uri are required", ErrInvalidAuthTransaction)
	}
	createdAt := params.CreatedAt.UTC()
	expiresAt := params.ExpiresAt.UTC()
	if createdAt.IsZero() || expiresAt.IsZero() || !expiresAt.After(createdAt) {
		return nil, fmt.Errorf("%w: invalid transaction timestamps", ErrInvalidAuthTransaction)
	}
	var consumedAt *time.Time
	if params.ConsumedAt != nil {
		value := params.ConsumedAt.UTC()
		if value.IsZero() || value.Before(createdAt) {
			return nil, fmt.Errorf("%w: invalid consumed timestamp", ErrInvalidAuthTransaction)
		}
		consumedAt = &value
	}
	if params.BrowserBindingHash.IsZero() {
		return nil, fmt.Errorf("%w: browser binding hash is required", ErrInvalidAuthTransaction)
	}
	return &AuthTransaction{id: params.ID, provider: params.Provider, stateHash: params.StateHash, nonceHash: params.NonceHash, browserBindingHash: params.BrowserBindingHash, codeVerifier: codeVerifier, redirectURI: redirectURI, createdAt: createdAt, expiresAt: expiresAt, consumedAt: consumedAt}, nil
}

func (a *AuthTransaction) ID() identity.ID                   { return a.id }
func (a *AuthTransaction) Provider() Provider                { return a.provider }
func (a *AuthTransaction) StateHash() valueobjects.StateHash { return a.stateHash }
func (a *AuthTransaction) NonceHash() valueobjects.NonceHash { return a.nonceHash }
func (a *AuthTransaction) BrowserBindingHash() valueobjects.BrowserBindingHash {
	return a.browserBindingHash
}
func (a *AuthTransaction) CodeVerifier() string { return a.codeVerifier }
func (a *AuthTransaction) RedirectURI() string  { return a.redirectURI }
func (a *AuthTransaction) CreatedAt() time.Time { return a.createdAt }
func (a *AuthTransaction) ExpiresAt() time.Time { return a.expiresAt }
func (a *AuthTransaction) ConsumedAt() *time.Time {
	if a == nil || a.consumedAt == nil {
		return nil
	}
	value := *a.consumedAt
	return &value
}
func (a *AuthTransaction) IsExpired(now time.Time) bool {
	return a == nil || !now.UTC().Before(a.expiresAt)
}
func (a *AuthTransaction) IsConsumed() bool { return a != nil && a.consumedAt != nil }

func (a *AuthTransaction) Consume(now time.Time) error {
	if a == nil {
		return fmt.Errorf("%w: transaction is nil", ErrInvalidAuthTransaction)
	}
	if a.IsConsumed() {
		return ErrAuthTransactionConsumed
	}
	now = now.UTC()
	if now.IsZero() || now.Before(a.createdAt) {
		return fmt.Errorf("%w: invalid consume timestamp", ErrInvalidAuthTransaction)
	}
	if a.IsExpired(now) {
		return ErrAuthTransactionExpired
	}
	a.consumedAt = &now
	return nil
}

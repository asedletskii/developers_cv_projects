package domain

import (
	"fmt"
	"strings"
	"time"

	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type AuthTransaction struct {
	id                 sharedvalueobjects.ID
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

type AuthTransactionSnapshot struct {
	ID                 sharedvalueobjects.ID
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
	id := sharedvalueobjects.NewID()
	createdAt := params.CreatedAt.UTC()
	expiresAt := params.ExpiresAt.UTC()
	var consumedAt *time.Time
	if params.ConsumedAt != nil {
		value := params.ConsumedAt.UTC()
		consumedAt = &value
	}
	if err := validateAuthTransaction(
		id, params.Provider, params.StateHash, params.NonceHash, params.BrowserBindingHash,
		strings.TrimSpace(params.CodeVerifier), strings.TrimSpace(params.RedirectURI),
		createdAt, expiresAt, consumedAt,
	); err != nil {
		return nil, err
	}
	return &AuthTransaction{
		id: id, provider: params.Provider, stateHash: params.StateHash, nonceHash: params.NonceHash,
		browserBindingHash: params.BrowserBindingHash, codeVerifier: strings.TrimSpace(params.CodeVerifier),
		redirectURI: strings.TrimSpace(params.RedirectURI), createdAt: createdAt, expiresAt: expiresAt,
		consumedAt: consumedAt,
	}, nil
}

func RestoreAuthTransaction(snapshot AuthTransactionSnapshot) (*AuthTransaction, error) {
	if err := validateAuthTransaction(
		snapshot.ID, snapshot.Provider, snapshot.StateHash, snapshot.NonceHash, snapshot.BrowserBindingHash,
		snapshot.CodeVerifier, snapshot.RedirectURI, snapshot.CreatedAt, snapshot.ExpiresAt, snapshot.ConsumedAt,
	); err != nil {
		return nil, err
	}
	return &AuthTransaction{
		id: snapshot.ID, provider: snapshot.Provider, stateHash: snapshot.StateHash, nonceHash: snapshot.NonceHash,
		browserBindingHash: snapshot.BrowserBindingHash, codeVerifier: snapshot.CodeVerifier,
		redirectURI: snapshot.RedirectURI, createdAt: snapshot.CreatedAt, expiresAt: snapshot.ExpiresAt,
		consumedAt: cloneTime(snapshot.ConsumedAt),
	}, nil
}

func validateAuthTransaction(
	id sharedvalueobjects.ID,
	provider Provider,
	stateHash valueobjects.StateHash,
	nonceHash valueobjects.NonceHash,
	browserBindingHash valueobjects.BrowserBindingHash,
	codeVerifier, redirectURI string,
	createdAt, expiresAt time.Time,
	consumedAt *time.Time,
) error {
	if id.IsZero() || !provider.IsValid() {
		return fmt.Errorf("%w: invalid transaction identity", ErrInvalidAuthTransaction)
	}
	if codeVerifier == "" || redirectURI == "" {
		return fmt.Errorf("%w: code verifier and redirect uri are required", ErrInvalidAuthTransaction)
	}
	if createdAt.IsZero() || expiresAt.IsZero() || !expiresAt.After(createdAt) {
		return fmt.Errorf("%w: invalid transaction timestamps", ErrInvalidAuthTransaction)
	}
	if consumedAt != nil {
		value := *consumedAt
		if value.IsZero() || value.Before(createdAt) {
			return fmt.Errorf("%w: invalid consumed timestamp", ErrInvalidAuthTransaction)
		}
	}
	if browserBindingHash.IsZero() {
		return fmt.Errorf("%w: browser binding hash is required", ErrInvalidAuthTransaction)
	}
	return nil
}

func (a *AuthTransaction) ID() sharedvalueobjects.ID         { return a.id }
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

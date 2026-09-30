package ports

import (
	"context"
	"errors"

	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
)

var (
	ErrTokenExchangeFailed = errors.New("telegram token exchange failed")
	ErrIDTokenInvalid      = errors.New("telegram id token invalid")
	ErrNonceMismatch       = errors.New("telegram nonce mismatch")
)

type TelegramAuthorizationParams struct {
	State         string
	Nonce         string
	CodeChallenge string
	RedirectURI   string
}

type VerifiedIdentity struct {
	Provider        domain.Provider
	ProviderSubject valueobjects.ProviderSubject
}

type TelegramOIDC interface {
	AuthorizationURL(params TelegramAuthorizationParams) string
	Exchange(ctx context.Context, code string, redirectURI string, verifier string, expectedNonceHash valueobjects.NonceHash) (VerifiedIdentity, error)
}

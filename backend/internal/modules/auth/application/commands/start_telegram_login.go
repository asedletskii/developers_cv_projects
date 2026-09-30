package commands

import (
	"context"
	"fmt"
	"time"

	"developers-cv/backend/internal/modules/auth/application"
	"developers-cv/backend/internal/modules/auth/application/ports"
	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type StartTelegramLoginHandler struct {
	repository     ports.AuthTransactionRepository
	oidc           ports.TelegramOIDC
	tokenGenerator ports.TokenGenerator
	uow            ports.UnitOfWork
	redirectURI    string
	transactionTTL time.Duration
	now            func() time.Time
}

func NewStartTelegramLoginHandler(
	repository ports.AuthTransactionRepository,
	oidc ports.TelegramOIDC,
	tokenGenerator ports.TokenGenerator,
	uow ports.UnitOfWork,
	redirectURI string,
	transactionTTL time.Duration,
) *StartTelegramLoginHandler {
	return &StartTelegramLoginHandler{
		repository:     repository,
		oidc:           oidc,
		tokenGenerator: tokenGenerator,
		uow:            uow,
		redirectURI:    redirectURI,
		transactionTTL: transactionTTL,
		now:            func() time.Time { return time.Now().UTC() },
	}
}

func (h *StartTelegramLoginHandler) Handle(ctx context.Context) (string, error) {
	if h == nil || h.repository == nil || h.oidc == nil || h.tokenGenerator == nil || h.uow == nil {
		return "", fmt.Errorf("start telegram login handler dependencies are nil")
	}
	if h.redirectURI == "" || h.transactionTTL <= 0 {
		return "", fmt.Errorf("%w: invalid start telegram login configuration", application.ErrInvalidAuthFlow)
	}

	state, err := h.tokenGenerator.GenerateState()
	if err != nil {
		return "", fmt.Errorf("generate auth state: %w", err)
	}
	nonce, err := h.tokenGenerator.GenerateNonce()
	if err != nil {
		return "", fmt.Errorf("generate auth nonce: %w", err)
	}
	verifier, err := h.tokenGenerator.GeneratePKCEVerifier()
	if err != nil {
		return "", fmt.Errorf("generate pkce verifier: %w", err)
	}

	stateValue := encodeOpaque(state)
	nonceValue := encodeOpaque(nonce)
	verifierValue := encodeOpaque(verifier)
	createdAt := h.now().UTC()

	provider, err := domain.NewProvider(string(domain.ProviderTelegram))
	if err != nil {
		return "", err
	}
	stateHash, err := valueobjects.NewStateHash(h.tokenGenerator.Hash([]byte(stateValue)))
	if err != nil {
		return "", fmt.Errorf("create state hash: %w", err)
	}
	nonceHash, err := valueobjects.NewNonceHash(h.tokenGenerator.Hash([]byte(nonceValue)))
	if err != nil {
		return "", fmt.Errorf("create nonce hash: %w", err)
	}
	transaction, err := domain.NewAuthTransaction(domain.AuthTransactionParams{
		ID: sharedvalueobjects.NewID(), Provider: provider, StateHash: stateHash, NonceHash: nonceHash,
		CodeVerifier: verifierValue, RedirectURI: h.redirectURI, CreatedAt: createdAt, ExpiresAt: createdAt.Add(h.transactionTTL),
	})
	if err != nil {
		return "", err
	}
	if err := h.uow.Do(ctx, func(tx context.Context) error {
		return h.repository.Create(tx, transaction)
	}); err != nil {
		return "", fmt.Errorf("create auth transaction: %w", err)
	}

	return h.oidc.AuthorizationURL(ports.TelegramAuthorizationParams{
		State: stateValue, Nonce: nonceValue, CodeChallenge: encodeOpaque(h.tokenGenerator.Hash([]byte(verifierValue))), RedirectURI: h.redirectURI,
	}), nil
}

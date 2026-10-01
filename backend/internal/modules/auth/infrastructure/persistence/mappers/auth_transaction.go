package mappers

import (
	"fmt"

	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	"developers-cv/backend/internal/modules/auth/infrastructure/persistence/gormmodel"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func AuthTransactionToDB(transaction *domain.AuthTransaction) (gormmodel.AuthTransaction, error) {
	if transaction == nil {
		return gormmodel.AuthTransaction{}, fmt.Errorf("auth transaction is nil")
	}
	return gormmodel.AuthTransaction{
		ID:                 transaction.ID().UUID(),
		Provider:           transaction.Provider().String(),
		StateHash:          transaction.StateHash().Bytes(),
		NonceHash:          transaction.NonceHash().Bytes(),
		BrowserBindingHash: transaction.BrowserBindingHash().Bytes(),
		CodeVerifier:       transaction.CodeVerifier(),
		RedirectURI:        transaction.RedirectURI(),
		CreatedAt:          transaction.CreatedAt(),
		ExpiresAt:          transaction.ExpiresAt(),
		ConsumedAt:         transaction.ConsumedAt(),
	}, nil
}

func AuthTransactionToDomain(model gormmodel.AuthTransaction) (*domain.AuthTransaction, error) {
	id, err := sharedvalueobjects.IDFromUUID(model.ID)
	if err != nil {
		return nil, fmt.Errorf("map auth transaction id: %w", err)
	}
	provider, err := domain.NewProvider(model.Provider)
	if err != nil {
		return nil, err
	}
	state, err := valueobjects.NewStateHash(model.StateHash)
	if err != nil {
		return nil, err
	}
	nonce, err := valueobjects.NewNonceHash(model.NonceHash)
	if err != nil {
		return nil, err
	}
	browserBinding, err := valueobjects.NewBrowserBindingHash(model.BrowserBindingHash)
	if err != nil {
		return nil, err
	}
	return domain.NewAuthTransaction(domain.AuthTransactionParams{
		ID:                 id,
		Provider:           provider,
		StateHash:          state,
		NonceHash:          nonce,
		BrowserBindingHash: browserBinding,
		CodeVerifier:       model.CodeVerifier,
		RedirectURI:        model.RedirectURI,
		CreatedAt:          model.CreatedAt,
		ExpiresAt:          model.ExpiresAt,
		ConsumedAt:         model.ConsumedAt,
	})
}

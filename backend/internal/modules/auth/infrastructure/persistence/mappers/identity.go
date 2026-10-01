package mappers

import (
	"fmt"

	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	"developers-cv/backend/internal/modules/auth/infrastructure/persistence/gormmodel"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func IdentityToDB(identity *domain.Identity) (gormmodel.UserIdentity, error) {
	if identity == nil {
		return gormmodel.UserIdentity{}, fmt.Errorf("identity is nil")
	}
	return gormmodel.UserIdentity{
		ID:              identity.ID().UUID(),
		UserID:          identity.UserID().UUID(),
		Provider:        identity.Provider().String(),
		ProviderSubject: identity.ProviderSubject().String(),
		CreatedAt:       identity.CreatedAt(),
		UpdatedAt:       identity.UpdatedAt(),
		LastLoginAt:     identity.LastLoginAt(),
	}, nil
}

func IdentityToDomain(model gormmodel.UserIdentity) (*domain.Identity, error) {
	id, err := sharedvalueobjects.IDFromUUID(model.ID)
	if err != nil {
		return nil, fmt.Errorf("map identity id: %w", err)
	}
	userID, err := sharedvalueobjects.IDFromUUID(model.UserID)
	if err != nil {
		return nil, fmt.Errorf("map identity user id: %w", err)
	}
	provider, err := domain.NewProvider(model.Provider)
	if err != nil {
		return nil, err
	}
	subject, err := valueobjects.NewProviderSubject(model.ProviderSubject)
	if err != nil {
		return nil, err
	}
	return domain.NewIdentity(domain.IdentityParams{
		ID:              id,
		UserID:          userID,
		Provider:        provider,
		ProviderSubject: subject,
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
		LastLoginAt:     model.LastLoginAt,
	})
}

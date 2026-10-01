package ports

import (
	"context"
	"time"

	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type IdentityRepository interface {
	FindByProviderSubject(ctx context.Context, provider domain.Provider, subject valueobjects.ProviderSubject) (*domain.Identity, error)
	Create(ctx context.Context, identity *domain.Identity) error
	UpdateLastLogin(ctx context.Context, identityID sharedvalueobjects.ID, at time.Time) error
}

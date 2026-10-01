package domain

import (
	"fmt"
	"time"

	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type Identity struct {
	id              identity.ID
	userID          identity.ID
	provider        Provider
	providerSubject valueobjects.ProviderSubject
	createdAt       time.Time
	updatedAt       time.Time
	lastLoginAt     time.Time
}

type IdentityParams struct {
	ID              identity.ID
	UserID          identity.ID
	Provider        Provider
	ProviderSubject valueobjects.ProviderSubject
	CreatedAt       time.Time
	UpdatedAt       time.Time
	LastLoginAt     time.Time
}

func NewIdentity(params IdentityParams) (*Identity, error) {
	if params.ID.IsZero() || params.UserID.IsZero() || !params.Provider.IsValid() || params.ProviderSubject.IsZero() {
		return nil, fmt.Errorf("%w: invalid identity fields", ErrInvalidIdentity)
	}
	createdAt := params.CreatedAt.UTC()
	updatedAt := params.UpdatedAt.UTC()
	lastLoginAt := params.LastLoginAt.UTC()
	if createdAt.IsZero() || updatedAt.IsZero() || lastLoginAt.IsZero() {
		return nil, fmt.Errorf("%w: identity timestamps must be set", ErrInvalidIdentity)
	}
	if updatedAt.Before(createdAt) || lastLoginAt.Before(createdAt) {
		return nil, fmt.Errorf("%w: invalid identity timestamp order", ErrInvalidIdentity)
	}
	return &Identity{id: params.ID, userID: params.UserID, provider: params.Provider, providerSubject: params.ProviderSubject, createdAt: createdAt, updatedAt: updatedAt, lastLoginAt: lastLoginAt}, nil
}

func (i *Identity) ID() identity.ID                               { return i.id }
func (i *Identity) UserID() identity.ID                           { return i.userID }
func (i *Identity) Provider() Provider                            { return i.provider }
func (i *Identity) ProviderSubject() valueobjects.ProviderSubject { return i.providerSubject }
func (i *Identity) CreatedAt() time.Time                          { return i.createdAt }
func (i *Identity) UpdatedAt() time.Time                          { return i.updatedAt }
func (i *Identity) LastLoginAt() time.Time                        { return i.lastLoginAt }

func (i *Identity) TouchLogin(now time.Time) error {
	if i == nil {
		return fmt.Errorf("%w: identity is nil", ErrInvalidIdentity)
	}
	now = now.UTC()
	if now.IsZero() || now.Before(i.createdAt) {
		return fmt.Errorf("%w: invalid login timestamp", ErrInvalidIdentity)
	}
	i.lastLoginAt = now
	if now.After(i.updatedAt) {
		i.updatedAt = now
	}
	return nil
}

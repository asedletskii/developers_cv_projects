package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	"developers-cv/backend/internal/modules/auth/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/platform/database"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type IdentityRepository struct{ db *gorm.DB }

func NewIdentityRepository(db *gorm.DB) *IdentityRepository { return &IdentityRepository{db: db} }

func (r *IdentityRepository) FindByProviderSubject(ctx context.Context, provider domain.Provider, subject valueobjects.ProviderSubject) (*domain.Identity, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("find identity: repository dependencies are nil")
	}
	db := database.DBFromContext(ctx, r.db)
	var model gormmodel.UserIdentity
	err := db.WithContext(ctx).Where("provider = ? AND provider_subject = ?", provider.String(), subject.String()).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrEntityNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find identity: %w", mapDBError(err))
	}
	return identityFromModel(model)
}

func (r *IdentityRepository) Create(ctx context.Context, identity *domain.Identity) error {
	if identity == nil {
		return fmt.Errorf("create identity: identity is nil")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("create identity: %w", err)
	}
	model := gormmodel.UserIdentity{ID: identity.ID().UUID(), UserID: identity.UserID().UUID(), Provider: identity.Provider().String(), ProviderSubject: identity.ProviderSubject().String(), CreatedAt: identity.CreatedAt(), UpdatedAt: identity.UpdatedAt(), LastLoginAt: identity.LastLoginAt()}
	if err := tx.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("create identity: %w", mapDBError(err))
	}
	return nil
}

func (r *IdentityRepository) UpdateLastLogin(ctx context.Context, identityID sharedvalueobjects.ID, at time.Time) error {
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("update identity last login: %w", err)
	}
	at = at.UTC()
	result := tx.WithContext(ctx).Model(&gormmodel.UserIdentity{}).Where("id = ?", identityID.UUID()).Updates(map[string]any{"last_login_at": at, "updated_at": at})
	if result.Error != nil {
		return fmt.Errorf("update identity last login: %w", mapDBError(result.Error))
	}
	if result.RowsAffected != 1 {
		return domain.ErrEntityNotFound
	}
	return nil
}

func identityFromModel(model gormmodel.UserIdentity) (*domain.Identity, error) {
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
	identity, err := domain.NewIdentity(domain.IdentityParams{ID: id, UserID: userID, Provider: provider, ProviderSubject: subject, CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt, LastLoginAt: model.LastLoginAt})
	if err != nil {
		return nil, err
	}
	return identity, nil
}

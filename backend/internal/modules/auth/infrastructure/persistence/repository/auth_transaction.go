package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	"developers-cv/backend/internal/modules/auth/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/platform/database"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type AuthTransactionRepository struct{ db *gorm.DB }

func NewAuthTransactionRepository(db *gorm.DB) *AuthTransactionRepository {
	return &AuthTransactionRepository{db: db}
}

func (r *AuthTransactionRepository) Create(ctx context.Context, transaction *domain.AuthTransaction) error {
	if transaction == nil {
		return fmt.Errorf("create auth transaction: transaction is nil")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("create auth transaction: %w", err)
	}
	model := transactionModel(transaction)
	if err := tx.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("create auth transaction: %w", mapDBError(err))
	}
	return nil
}

func (r *AuthTransactionRepository) FindByStateHash(ctx context.Context, stateHash valueobjects.StateHash) (*domain.AuthTransaction, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("find auth transaction: repository dependencies are nil")
	}
	return r.find(ctx, stateHash, false)
}
func (r *AuthTransactionRepository) FindByStateHashForUpdate(ctx context.Context, stateHash valueobjects.StateHash) (*domain.AuthTransaction, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("find auth transaction: repository dependencies are nil")
	}
	return r.find(ctx, stateHash, true)
}
func (r *AuthTransactionRepository) find(ctx context.Context, stateHash valueobjects.StateHash, lock bool) (*domain.AuthTransaction, error) {
	var db *gorm.DB
	if lock {
		tx, err := database.DBFromTransactionContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("find auth transaction for update: %w", err)
		}
		db = tx
	} else {
		db = database.DBFromContext(ctx, r.db)
	}
	var model gormmodel.AuthTransaction
	query := db.WithContext(ctx)
	if lock {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	err := query.Where("state_hash = ?", stateHash.Bytes()).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrEntityNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find auth transaction: %w", mapDBError(err))
	}
	return transactionFromModel(model)
}

func (r *AuthTransactionRepository) Consume(ctx context.Context, transactionID sharedvalueobjects.ID, at time.Time) error {
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("consume auth transaction: %w", err)
	}
	result := tx.WithContext(ctx).Model(&gormmodel.AuthTransaction{}).Where("id = ? AND consumed_at IS NULL", transactionID.UUID()).Update("consumed_at", at.UTC())
	if result.Error != nil {
		return fmt.Errorf("consume auth transaction: %w", mapDBError(result.Error))
	}
	if result.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
}

func (r *AuthTransactionRepository) DeleteExpiredBatch(ctx context.Context, now time.Time, limit int) (int, error) {
	if limit <= 0 {
		return 0, fmt.Errorf("delete expired auth transactions: limit must be positive")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return 0, fmt.Errorf("delete expired auth transactions: %w", err)
	}
	var ids []uuid.UUID
	if err := tx.WithContext(ctx).Model(&gormmodel.AuthTransaction{}).Where("expires_at <= ?", now.UTC()).Order("expires_at, id").Limit(limit).Pluck("id", &ids).Error; err != nil {
		return 0, fmt.Errorf("find expired auth transactions: %w", mapDBError(err))
	}
	if len(ids) == 0 {
		return 0, nil
	}
	result := tx.WithContext(ctx).Where("id IN ?", ids).Delete(&gormmodel.AuthTransaction{})
	if result.Error != nil {
		return 0, fmt.Errorf("delete expired auth transactions: %w", mapDBError(result.Error))
	}
	return int(result.RowsAffected), nil
}

func transactionModel(t *domain.AuthTransaction) gormmodel.AuthTransaction {
	return gormmodel.AuthTransaction{ID: t.ID().UUID(), Provider: t.Provider().String(), StateHash: t.StateHash().Bytes(), NonceHash: t.NonceHash().Bytes(), CodeVerifier: t.CodeVerifier(), RedirectURI: t.RedirectURI(), CreatedAt: t.CreatedAt(), ExpiresAt: t.ExpiresAt(), ConsumedAt: t.ConsumedAt()}
}
func transactionFromModel(m gormmodel.AuthTransaction) (*domain.AuthTransaction, error) {
	id, err := sharedvalueobjects.IDFromUUID(m.ID)
	if err != nil {
		return nil, err
	}
	provider, err := domain.NewProvider(m.Provider)
	if err != nil {
		return nil, err
	}
	state, err := valueobjects.NewStateHash(m.StateHash)
	if err != nil {
		return nil, err
	}
	nonce, err := valueobjects.NewNonceHash(m.NonceHash)
	if err != nil {
		return nil, err
	}
	return domain.NewAuthTransaction(domain.AuthTransactionParams{ID: id, Provider: provider, StateHash: state, NonceHash: nonce, CodeVerifier: m.CodeVerifier, RedirectURI: m.RedirectURI, CreatedAt: m.CreatedAt, ExpiresAt: m.ExpiresAt, ConsumedAt: m.ConsumedAt})
}

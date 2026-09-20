package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
	"developers-cv/backend/internal/modules/auth/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/modules/auth/infrastructure/persistence/mappers"
	"developers-cv/backend/internal/platform/database"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type SessionRepository struct{ db *gorm.DB }

func NewSessionRepository(db *gorm.DB) *SessionRepository { return &SessionRepository{db: db} }

func (r *SessionRepository) Create(ctx context.Context, session *domain.Session) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("create session: repository dependencies are nil")
	}
	if session == nil {
		return fmt.Errorf("create session: session is nil")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	model, err := mappers.SessionToDB(session)
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	if err := tx.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("create session: %w", mapDBError(err))
	}
	return nil
}

func (r *SessionRepository) FindActiveByTokenHash(ctx context.Context, tokenHash valueobjects.SessionTokenHash) (*domain.Session, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("find active session: repository dependencies are nil")
	}
	db := database.DBFromContext(ctx, r.db)
	var model gormmodel.AuthSession
	err := db.WithContext(ctx).Where("token_hash = ? AND revoked_at IS NULL AND expires_at > ?", tokenHash.Bytes(), time.Now().UTC()).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrEntityNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find active session: %w", mapDBError(err))
	}
	session, err := mappers.SessionToDomain(model)
	if err != nil {
		return nil, fmt.Errorf("find active session: map session: %w", err)
	}
	return session, nil
}

func (r *SessionRepository) Revoke(ctx context.Context, sessionID sharedvalueobjects.ID, at time.Time) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("revoke session: repository dependencies are nil")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}
	result := tx.WithContext(ctx).Model(&gormmodel.AuthSession{}).Where("id = ? AND revoked_at IS NULL", sessionID.UUID()).Updates(map[string]any{"revoked_at": at.UTC()})
	if result.Error != nil {
		return fmt.Errorf("revoke session: %w", mapDBError(result.Error))
	}
	if result.RowsAffected == 0 {
		var exists int64
		if err := tx.WithContext(ctx).Model(&gormmodel.AuthSession{}).Where("id = ?", sessionID.UUID()).Count(&exists).Error; err != nil {
			return fmt.Errorf("revoke session: check session: %w", mapDBError(err))
		}
		if exists == 0 {
			return domain.ErrEntityNotFound
		}
	}
	return nil
}

func (r *SessionRepository) RevokeAllByUserID(ctx context.Context, userID sharedvalueobjects.ID, at time.Time) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("revoke all sessions: repository dependencies are nil")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("revoke all sessions: %w", err)
	}
	if err := tx.WithContext(ctx).Model(&gormmodel.AuthSession{}).Where("user_id = ? AND revoked_at IS NULL", userID.UUID()).Update("revoked_at", at.UTC()).Error; err != nil {
		return fmt.Errorf("revoke all sessions: %w", mapDBError(err))
	}
	return nil
}

func (r *SessionRepository) Touch(ctx context.Context, sessionID sharedvalueobjects.ID, at time.Time) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("touch session: repository dependencies are nil")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("touch session: %w", err)
	}
	result := tx.WithContext(ctx).Model(&gormmodel.AuthSession{}).Where("id = ? AND revoked_at IS NULL", sessionID.UUID()).Update("last_seen_at", at.UTC())
	if result.Error != nil {
		return fmt.Errorf("touch session: %w", mapDBError(result.Error))
	}
	if result.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
}

func (r *SessionRepository) DeleteExpiredBatch(ctx context.Context, now time.Time, limit int) (int, error) {
	if r == nil || r.db == nil {
		return 0, fmt.Errorf("delete expired sessions: repository dependencies are nil")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return 0, fmt.Errorf("delete expired sessions: %w", err)
	}
	var ids []uuid.UUID
	if err := tx.WithContext(ctx).Model(&gormmodel.AuthSession{}).Where("expires_at <= ?", now.UTC()).Order("expires_at, id").Limit(limit).Pluck("id", &ids).Error; err != nil {
		return 0, fmt.Errorf("find expired sessions: %w", mapDBError(err))
	}
	if len(ids) == 0 {
		return 0, nil
	}
	result := tx.WithContext(ctx).Where("id IN ?", ids).Delete(&gormmodel.AuthSession{})
	if result.Error != nil {
		return 0, fmt.Errorf("delete expired sessions: %w", mapDBError(result.Error))
	}
	return int(result.RowsAffected), nil
}

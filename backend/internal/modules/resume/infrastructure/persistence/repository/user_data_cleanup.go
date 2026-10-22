package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"

	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/platform/database"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func (r *ResumeRepository) ListIDsForUpdateByUserID(
	ctx context.Context,
	userID sharedvalueobjects.ID,
) ([]sharedvalueobjects.ID, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("list resume ids for update by user: repository dependencies are nil")
	}
	if userID.IsZero() {
		return nil, fmt.Errorf("list resume ids for update by user: user id is empty")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("list resume ids for update by user: %w", err)
	}
	var resumeIDs []uuid.UUID
	if err := tx.WithContext(ctx).
		Model(&gormmodel.Resume{}).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ?", userID.UUID()).
		Order("id").
		Pluck("id", &resumeIDs).Error; err != nil {
		return nil, fmt.Errorf("list resume ids for update by user: %w", mapDBError(err))
	}
	result := make([]sharedvalueobjects.ID, 0, len(resumeIDs))
	for _, resumeID := range resumeIDs {
		id, err := sharedvalueobjects.IDFromUUID(resumeID)
		if err != nil {
			return nil, fmt.Errorf("list resume ids for update by user: %w", err)
		}
		result = append(result, id)
	}
	return result, nil
}

func (r *ResumeRepository) DeleteByUserID(ctx context.Context, userID sharedvalueobjects.ID) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("delete resumes by user: repository dependencies are nil")
	}
	if userID.IsZero() {
		return fmt.Errorf("delete resumes by user: user id is empty")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("delete resumes by user: %w", err)
	}
	if err := tx.WithContext(ctx).
		Where("user_id = ?", userID.UUID()).
		Delete(&gormmodel.Resume{}).Error; err != nil {
		return fmt.Errorf("delete resumes by user: %w", mapDBError(err))
	}
	return nil
}

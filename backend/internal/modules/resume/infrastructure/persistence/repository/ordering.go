package repository

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/domain"
)

func reorderPositions(tx *gorm.DB, table string, resumeID uuid.UUID, orderedIDs []uuid.UUID) error {
	var existingIDs []uuid.UUID
	if err := tx.Table(table).
		Where("resume_id = ?", resumeID).
		Order("position").
		Pluck("id", &existingIDs).Error; err != nil {
		return err
	}

	if len(existingIDs) != len(orderedIDs) || len(orderedIDs) == 0 {
		return fmt.Errorf("invalid %s order: %w", table, domain.ErrInvalidEntity)
	}

	existing := make(map[uuid.UUID]struct{}, len(existingIDs))
	for _, id := range existingIDs {
		existing[id] = struct{}{}
	}
	requested := make(map[uuid.UUID]struct{}, len(orderedIDs))
	for _, id := range orderedIDs {
		if _, ok := existing[id]; !ok {
			return fmt.Errorf("entity %s does not belong to collection: %w", id, domain.ErrInvalidEntity)
		}
		if _, ok := requested[id]; ok {
			return fmt.Errorf("duplicate entity %s in order: %w", id, domain.ErrInvalidEntity)
		}
		requested[id] = struct{}{}
	}

	offset := len(existingIDs) + 1
	if err := tx.Table(table).
		Where("resume_id = ?", resumeID).
		Update("position", gorm.Expr("position + ?", offset)).Error; err != nil {
		return err
	}

	for position, id := range orderedIDs {
		result := tx.Table(table).
			Where("id = ? AND resume_id = ?", id, resumeID).
			Update("position", position)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return domain.ErrEntityNotFound
		}
	}
	return nil
}

func nextPosition(tx *gorm.DB, table string, resumeID uuid.UUID) (int, error) {
	var maxPosition *int
	if err := tx.Table(table).Where("resume_id = ?", resumeID).Select("MAX(position)").Scan(&maxPosition).Error; err != nil {
		return 0, err
	}
	if maxPosition == nil {
		return 0, nil
	}
	return *maxPosition + 1, nil
}

func nextSortOrder(tx *gorm.DB, table string, resumeID uuid.UUID) (int, error) {
	var maxOrder *int
	if err := tx.Table(table).Where("resume_id = ?", resumeID).Select("MAX(sort_order)").Scan(&maxOrder).Error; err != nil {
		return 0, err
	}
	if maxOrder == nil {
		return 0, nil
	}
	return *maxOrder + 1, nil
}

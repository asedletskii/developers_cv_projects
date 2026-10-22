package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"developers-cv/backend/internal/modules/analytics/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/platform/database"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func (r *AnalyticsRepository) DeleteByLinkIDs(ctx context.Context, linkIDs []sharedvalueobjects.ID) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("delete analytics by link ids: repository dependencies are nil")
	}
	if len(linkIDs) == 0 {
		return nil
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("delete analytics by link ids: %w", err)
	}
	ids := make([]uuid.UUID, 0, len(linkIDs))
	for _, linkID := range linkIDs {
		if linkID.IsZero() {
			return fmt.Errorf("delete analytics by link ids: link id is empty")
		}
		ids = append(ids, linkID.UUID())
	}
	tx = tx.WithContext(ctx)
	if err := tx.Where("link_id IN ?", ids).Delete(&gormmodel.AnalyticsClick{}).Error; err != nil {
		return fmt.Errorf("delete analytics clicks: %w", mapDBError(err))
	}
	if err := tx.Where("link_id IN ?", ids).Delete(&gormmodel.AnalyticsEvent{}).Error; err != nil {
		return fmt.Errorf("delete analytics events: %w", mapDBError(err))
	}
	if err := tx.Where("link_id IN ?", ids).Delete(&gormmodel.AnalyticsView{}).Error; err != nil {
		return fmt.Errorf("delete analytics views: %w", mapDBError(err))
	}
	if err := tx.Where("link_id IN ?", ids).Delete(&gormmodel.AnalyticsVisitor{}).Error; err != nil {
		return fmt.Errorf("delete analytics visitors: %w", mapDBError(err))
	}
	return nil
}

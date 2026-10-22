package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/analytics/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/platform/database"
)

func (r *AnalyticsRepository) DeleteExpiredClicks(
	ctx context.Context,
	cutoff time.Time,
	limit int,
) (int, error) {
	if r == nil || r.db == nil {
		return 0, fmt.Errorf("delete expired analytics clicks: repository dependencies are nil")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return 0, fmt.Errorf("delete expired analytics clicks: %w", err)
	}
	if acquired, err := tryAcquireAnalyticsRetentionLock(tx.WithContext(ctx)); err != nil {
		return 0, fmt.Errorf("acquire analytics retention cleanup lock: %w", mapDBError(err))
	} else if !acquired {
		return 0, nil
	}
	var ids []int64
	if err := tx.WithContext(ctx).Model(&gormmodel.AnalyticsClick{}).
		Where("occurred_at < ?", cutoff.UTC()).
		Order("occurred_at, id").Limit(limit).Pluck("id", &ids).Error; err != nil {
		return 0, fmt.Errorf("find expired analytics clicks: %w", mapDBError(err))
	}
	if len(ids) == 0 {
		return 0, nil
	}
	result := tx.WithContext(ctx).Where("id IN ?", ids).Delete(&gormmodel.AnalyticsClick{})
	if result.Error != nil {
		return 0, fmt.Errorf("delete expired analytics clicks: %w", mapDBError(result.Error))
	}
	return int(result.RowsAffected), nil
}

func (r *AnalyticsRepository) DeleteExpiredViews(
	ctx context.Context,
	cutoff time.Time,
	limit int,
) (int, error) {
	if r == nil || r.db == nil {
		return 0, fmt.Errorf("delete expired analytics views: repository dependencies are nil")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return 0, fmt.Errorf("delete expired analytics views: %w", err)
	}
	if acquired, err := tryAcquireAnalyticsRetentionLock(tx.WithContext(ctx)); err != nil {
		return 0, fmt.Errorf("acquire analytics retention cleanup lock: %w", mapDBError(err))
	} else if !acquired {
		return 0, nil
	}
	var ids []uuid.UUID
	if err := tx.WithContext(ctx).Model(&gormmodel.AnalyticsView{}).
		Where("started_at < ?", cutoff.UTC()).
		Order("started_at, id").Limit(limit).Pluck("id", &ids).Error; err != nil {
		return 0, fmt.Errorf("find expired analytics views: %w", mapDBError(err))
	}
	if len(ids) == 0 {
		return 0, nil
	}
	result := tx.WithContext(ctx).Where("id IN ?", ids).Delete(&gormmodel.AnalyticsView{})
	if result.Error != nil {
		return 0, fmt.Errorf("delete expired analytics views: %w", mapDBError(result.Error))
	}
	return int(result.RowsAffected), nil
}

func tryAcquireAnalyticsRetentionLock(tx *gorm.DB) (bool, error) {
	if tx == nil {
		return false, fmt.Errorf("analytics retention cleanup transaction is nil")
	}
	if tx.Name() != "postgres" {
		return false, fmt.Errorf("analytics retention cleanup lock requires PostgreSQL, got %q", tx.Name())
	}
	const lockKey int64 = 741209882133
	var acquired bool
	if err := tx.Raw("SELECT pg_try_advisory_xact_lock(?)", lockKey).Scan(&acquired).Error; err != nil {
		return false, err
	}
	return acquired, nil
}

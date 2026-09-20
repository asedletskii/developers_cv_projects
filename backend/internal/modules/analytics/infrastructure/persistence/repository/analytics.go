package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"developers-cv/backend/internal/modules/analytics/application/ports"
	analyticscontracts "developers-cv/backend/internal/modules/analytics/contracts"
	"developers-cv/backend/internal/modules/analytics/domain"
	"developers-cv/backend/internal/modules/analytics/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/modules/analytics/infrastructure/persistence/mappers"
	"developers-cv/backend/internal/platform/database"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func (r *AnalyticsRepository) LoadAnalyticsBatchForUpdate(
	ctx context.Context,
	viewID, linkID sharedvalueobjects.ID,
	includeEventCount, includeClickCount, includePrintHistory bool,
) (ports.AnalyticsBatchSnapshot, error) {
	if r == nil || r.db == nil {
		return ports.AnalyticsBatchSnapshot{}, fmt.Errorf("load analytics batch: repository dependencies are nil")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return ports.AnalyticsBatchSnapshot{}, fmt.Errorf("load analytics batch: %w", err)
	}

	var model gormmodel.AnalyticsView
	if err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND link_id = ?", viewID.UUID(), linkID.UUID()).
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ports.AnalyticsBatchSnapshot{}, analyticscontracts.ErrAnalyticsSubjectNotFound
		}
		return ports.AnalyticsBatchSnapshot{}, fmt.Errorf("load analytics batch: lock view: %w", mapDBError(err))
	}

	view, err := mappers.ViewToDomain(model)
	if err != nil {
		return ports.AnalyticsBatchSnapshot{}, fmt.Errorf("load analytics batch: map view: %w", err)
	}

	snapshot := ports.AnalyticsBatchSnapshot{View: view}
	if includeEventCount {
		if err := tx.WithContext(ctx).Model(&gormmodel.AnalyticsEvent{}).
			Where("view_id = ? AND client_event_id IS NOT NULL", viewID.UUID()).
			Count(&snapshot.StoredEventCount).Error; err != nil {
			return ports.AnalyticsBatchSnapshot{}, fmt.Errorf("load analytics batch: count events: %w", mapDBError(err))
		}
	}
	if includeClickCount {
		if err := tx.WithContext(ctx).Model(&gormmodel.AnalyticsClick{}).
			Where("view_id = ?", viewID.UUID()).Count(&snapshot.StoredClickCount).Error; err != nil {
			return ports.AnalyticsBatchSnapshot{}, fmt.Errorf("load analytics batch: count clicks: %w", mapDBError(err))
		}
	}
	if includePrintHistory && snapshot.StoredEventCount < domain.MaxAnalyticsEventsPerView {
		if err := tx.WithContext(ctx).Model(&gormmodel.AnalyticsEvent{}).
			Select("occurred_at").
			Where("view_id = ? AND event_type = ?", viewID.UUID(), string(domain.EventTypePrint)).
			Order("occurred_at ASC").Pluck("occurred_at", &snapshot.ExistingPrints).Error; err != nil {
			return ports.AnalyticsBatchSnapshot{}, fmt.Errorf("load analytics batch: load print history: %w", mapDBError(err))
		}
	}
	return snapshot, nil
}

func (r *AnalyticsRepository) SaveAnalyticsBatch(
	ctx context.Context,
	view *domain.View,
	events []*domain.AnalyticsEvent,
	clicks []*domain.AnalyticsClick,
) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("save analytics batch: repository dependencies are nil")
	}
	if view == nil {
		return fmt.Errorf("save analytics batch: view is nil")
	}
	for _, event := range events {
		if event == nil || event.ViewID() != view.ID() || event.LinkID() != view.LinkID() {
			return fmt.Errorf("%w: analytics event does not belong to batch view", domain.ErrInvalidAnalyticsEvent)
		}
	}
	for _, click := range clicks {
		if click == nil || click.ViewID() != view.ID() || click.LinkID() != view.LinkID() {
			return fmt.Errorf("%w: analytics click does not belong to batch view", domain.ErrInvalidAnalyticsClick)
		}
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("save analytics batch: %w", err)
	}

	if len(events) > 0 {
		models := make([]gormmodel.AnalyticsEvent, 0, len(events))
		for _, event := range events {
			model, mapErr := mappers.EventToDB(event)
			if mapErr != nil {
				return fmt.Errorf("save analytics batch: map event: %w", mapErr)
			}
			models = append(models, model)
		}
		if err := tx.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&models).Error; err != nil {
			return fmt.Errorf("save analytics batch: insert events: %w", mapDBError(err))
		}
	}
	if len(clicks) > 0 {
		models := make([]gormmodel.AnalyticsClick, 0, len(clicks))
		for _, click := range clicks {
			model, mapErr := mappers.ClickToDB(click)
			if mapErr != nil {
				return fmt.Errorf("save analytics batch: map click: %w", mapErr)
			}
			models = append(models, model)
		}
		if err := tx.WithContext(ctx).Create(&models).Error; err != nil {
			return fmt.Errorf("save analytics batch: insert clicks: %w", mapDBError(err))
		}
	}

	state := view.State()
	if err := tx.WithContext(ctx).Model(&gormmodel.AnalyticsView{}).
		Where("id = ? AND link_id = ?", view.ID().UUID(), view.LinkID().UUID()).
		Updates(map[string]any{
			"confirmed_at":   nullableTime(state.ConfirmedAt()),
			"last_seen_at":   nullableTime(state.LastSeenAt()),
			"active_ms":      state.ActiveMS(),
			"last_state_seq": state.LastSeq(),
			"max_section":    string(state.MaxSection()),
			"exit_section":   string(state.ExitSection()),
			"ended_at":       nullableTime(state.EndedAt()),
		}).Error; err != nil {
		return fmt.Errorf("save analytics batch: update state: %w", mapDBError(err))
	}
	return nil
}

func (r *AnalyticsRepository) RecordClicks(ctx context.Context, clicks []*domain.AnalyticsClick) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("record analytics clicks: repository dependencies are nil")
	}
	if len(clicks) == 0 {
		return nil
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("record analytics clicks: %w", err)
	}
	viewID := clicks[0].ViewID().UUID()
	linkID := clicks[0].LinkID().UUID()
	var matchingView gormmodel.AnalyticsView
	if err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Select("id").
		Where("id = ? AND link_id = ?", viewID, linkID).
		First(&matchingView).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return analyticscontracts.ErrAnalyticsSubjectNotFound
		}
		return fmt.Errorf("record analytics clicks: verify and lock view: %w", mapDBError(err))
	}
	return r.recordClicksForLockedView(ctx, tx, viewID, clicks)
}

func (r *AnalyticsRepository) recordClicksForLockedView(
	ctx context.Context,
	tx *gorm.DB,
	viewID uuid.UUID,
	clicks []*domain.AnalyticsClick,
) error {
	var stored int64
	if err := tx.WithContext(ctx).Model(&gormmodel.AnalyticsClick{}).Where("view_id = ?", viewID).Count(&stored).Error; err != nil {
		return fmt.Errorf("record analytics clicks: count: %w", mapDBError(err))
	}
	clicks = domain.SelectAnalyticsClicksForStorage(clicks, stored)
	if len(clicks) == 0 {
		return nil
	}
	models := make([]gormmodel.AnalyticsClick, 0, len(clicks))
	for _, click := range clicks {
		dbModel, err := mappers.ClickToDB(click)
		if err != nil {
			return fmt.Errorf("record analytics clicks: map click: %w", err)
		}
		models = append(models, dbModel)
	}
	if len(models) == 0 {
		return nil
	}
	if err := tx.WithContext(ctx).Create(&models).Error; err != nil {
		return fmt.Errorf("record analytics clicks: insert: %w", mapDBError(err))
	}
	return nil
}

func (r *AnalyticsRepository) RecordEvents(ctx context.Context, events []*domain.AnalyticsEvent) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("record analytics events: repository dependencies are nil")
	}
	if len(events) == 0 {
		return nil
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("record analytics events: %w", err)
	}
	tx = tx.WithContext(ctx)
	first := events[0]
	viewID := first.ViewID().UUID()
	linkID := first.LinkID().UUID()
	models := make([]gormmodel.AnalyticsEvent, 0, len(events))
	for _, event := range events {
		model, err := mappers.EventToDB(event)
		if err != nil {
			return fmt.Errorf("record analytics events: %w", err)
		}
		models = append(models, model)
	}
	var matchingView gormmodel.AnalyticsView
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Select("id").
		Where("id = ? AND link_id = ?", viewID, linkID).
		First(&matchingView).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return analyticscontracts.ErrAnalyticsSubjectNotFound
		}
		return fmt.Errorf("record analytics events: verify and lock view: %w", mapDBError(err))
	}
	if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&models).Error; err != nil {
		return fmt.Errorf("record analytics events: %w", mapDBError(err))
	}
	return nil
}

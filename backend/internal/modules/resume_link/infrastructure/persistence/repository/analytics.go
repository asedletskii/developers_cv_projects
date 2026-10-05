package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/modules/resume_link/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/modules/resume_link/infrastructure/persistence/mappers"
	"developers-cv/backend/internal/platform/database"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func (r *ResumeLinkRepository) CreateView(ctx context.Context, view *domain.View) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("create resume link view: repository dependencies are nil")
	}
	if view == nil || view.ID().IsZero() || view.LinkID().IsZero() || view.VisitorID().IsZero() {
		return fmt.Errorf("create resume link view: invalid view")
	}
	tx := database.DBFromContext(ctx, r.db).WithContext(ctx)
	model, err := mappers.ViewToDB(view)
	if err != nil {
		return fmt.Errorf("create resume link view: %w", err)
	}
	if err := tx.Create(&model).Error; err != nil {
		return fmt.Errorf("create resume link view: %w", mapDBError(err))
	}
	return nil
}

func (r *ResumeLinkRepository) GetView(ctx context.Context, id sharedvalueobjects.ID) (*domain.View, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("get resume link view: repository dependencies are nil")
	}
	if id.IsZero() {
		return nil, fmt.Errorf("get resume link view: id is empty")
	}
	db := database.DBFromContext(ctx, r.db).WithContext(ctx)
	var model gormmodel.ResumeLinkView
	if err := db.Where("id = ?", id.UUID()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, fmt.Errorf("get resume link view: %w", err)
	}
	return mappers.ViewToDomain(model)
}

func (r *ResumeLinkRepository) UpdateViewState(ctx context.Context, id sharedvalueobjects.ID, update domain.ViewStateUpdate) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("update resume link view state: repository dependencies are nil")
	}
	if id.IsZero() {
		return fmt.Errorf("update resume link view state: id is empty")
	}
	tx := database.DBFromContext(ctx, r.db).WithContext(ctx)
	var model gormmodel.ResumeLinkView
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id.UUID()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrEntityNotFound
		}
		return fmt.Errorf("update resume link view state: load view: %w", err)
	}
	view, err := mappers.ViewToDomain(model)
	if err != nil {
		return fmt.Errorf("update resume link view state: map view: %w", err)
	}
	if err := view.ApplyState(update); err != nil {
		return fmt.Errorf("update resume link view state: %w", err)
	}
	state := view.State()
	updates := map[string]any{
		"confirmed_at": nullableTime(state.ConfirmedAt()), "last_seen_at": nullableTime(state.LastSeenAt()),
		"active_ms": state.ActiveMS(), "max_section": string(state.MaxSection()), "exit_section": string(state.ExitSection()),
		"ended_at": nullableTime(state.EndedAt()),
	}
	if err := tx.Model(&gormmodel.ResumeLinkView{}).Where("id = ?", id.UUID()).Updates(updates).Error; err != nil {
		return fmt.Errorf("update resume link view state: persist state: %w", err)
	}
	return nil
}

func (r *ResumeLinkRepository) ApplyBatch(ctx context.Context, viewID, linkID sharedvalueobjects.ID, update domain.ViewStateUpdate, events []*domain.AnalyticsEvent, clicks []*domain.AnalyticsClick) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("apply analytics batch: repository dependencies are nil")
	}
	if viewID.IsZero() || linkID.IsZero() {
		return fmt.Errorf("apply analytics batch: invalid ids")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("apply analytics batch: %w", err)
	}

	var model gormmodel.ResumeLinkView
	if err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", viewID.UUID()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrEntityNotFound
		}
		return fmt.Errorf("apply analytics batch: load view: %w", err)
	}
	if model.LinkID != linkID.UUID() {
		return domain.ErrEntityNotFound
	}

	var link gormmodel.ResumeLink
	if err := tx.WithContext(ctx).Where("id = ?", linkID.UUID()).First(&link).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrEntityNotFound
		}
		return fmt.Errorf("apply analytics batch: load link: %w", err)
	}
	if link.Status != string(domain.ResumeLinkStatusActive) {
		return domain.ErrResumeLinkInactive
	}

	var resume struct {
		ID     uuid.UUID `gorm:"column:id"`
		Status string    `gorm:"column:status"`
	}
	if err := tx.WithContext(ctx).Table("resumes").Select("id, status").Where("id = ?", link.ResumeID).First(&resume).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrEntityNotFound
		}
		return fmt.Errorf("apply analytics batch: load resume: %w", err)
	}
	if resume.Status != "opened" {
		return domain.ErrEntityNotFound
	}
	if model.StartedAt.Before(time.Now().UTC().Add(-12 * time.Hour)) {
		return domain.ErrAnalyticsViewExpired
	}
	view, err := mappers.ViewToDomain(model)
	if err != nil {
		return fmt.Errorf("apply analytics batch: map view: %w", err)
	}
	if view.Traffic() == domain.TrafficBot {
		return nil
	}
	if update.Seq <= view.State().LastSeq() {
		return nil
	}
	if err := view.ApplyState(update); err != nil {
		return fmt.Errorf("apply analytics batch: apply state: %w", err)
	}
	if len(events) > 0 {
		var stored int64
		if err := tx.WithContext(ctx).Model(&gormmodel.ResumeLinkAnalyticsEvent{}).
			Where("view_id = ? AND client_event_id IS NOT NULL", viewID.UUID()).Count(&stored).Error; err != nil {
			return fmt.Errorf("apply analytics batch: count events: %w", err)
		}
		remaining := int64(commandsMaxAnalyticsEventsPerView) - stored
		if remaining > 0 {
			if int64(len(events)) > remaining {
				events = events[:remaining]
			}
			models := make([]gormmodel.ResumeLinkAnalyticsEvent, 0, len(events))
			printAccepted := false
			for _, event := range events {
				if event.Type() == domain.EventTypePrint {
					if printAccepted {
						continue
					}
					var recentPrint int64
					if err := tx.WithContext(ctx).Model(&gormmodel.ResumeLinkAnalyticsEvent{}).
						Where("view_id = ? AND event_type = ? AND occurred_at >= ?", viewID.UUID(), string(domain.EventTypePrint), event.OccurredAt().Add(-2*time.Second)).
						Count(&recentPrint).Error; err != nil {
						return fmt.Errorf("apply analytics batch: check print deduplication: %w", err)
					}
					if recentPrint > 0 {
						continue
					}
					printAccepted = true
				}
				dbModel, mapErr := mappers.AnalyticsEventToDB(event)
				if mapErr != nil {
					return fmt.Errorf("apply analytics batch: map event: %w", mapErr)
				}
				models = append(models, dbModel)
			}
			if len(models) > 0 {
				if err := tx.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&models).Error; err != nil {
					return fmt.Errorf("apply analytics batch: insert events: %w", mapDBError(err))
				}
			}
		}
	}

	if len(clicks) > 0 {
		if err := r.RecordClicks(ctx, clicks); err != nil {
			return fmt.Errorf("apply analytics batch: record clicks: %w", err)
		}
	}

	state := view.State()
	if err := tx.WithContext(ctx).Model(&gormmodel.ResumeLinkView{}).Where("id = ?", viewID.UUID()).Updates(map[string]any{
		"confirmed_at":   nullableTime(state.ConfirmedAt()),
		"last_seen_at":   nullableTime(state.LastSeenAt()),
		"active_ms":      state.ActiveMS(),
		"last_state_seq": state.LastSeq(),
		"max_section":    string(state.MaxSection()),
		"exit_section":   string(state.ExitSection()),
		"ended_at":       nullableTime(state.EndedAt()),
	}).Error; err != nil {
		return fmt.Errorf("apply analytics batch: update state: %w", err)
	}
	return nil
}

const (
	commandsMaxAnalyticsEventsPerView = 1000
	maxAnalyticsClicksPerView         = 1000
)

func (r *ResumeLinkRepository) RecordClicks(ctx context.Context, clicks []*domain.AnalyticsClick) error {
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
	if clicks[0] == nil {
		return fmt.Errorf("record analytics clicks: nil click")
	}
	viewID := clicks[0].ViewID().UUID()
	for _, click := range clicks[1:] {
		if click == nil || click.ViewID().UUID() != viewID {
			return fmt.Errorf("record analytics clicks: mixed view ids")
		}
	}
	var stored int64
	if err := tx.WithContext(ctx).Model(&gormmodel.ResumeLinkClick{}).Where("view_id = ?", viewID).Count(&stored).Error; err != nil {
		return fmt.Errorf("record analytics clicks: count: %w", err)
	}
	remaining := int64(maxAnalyticsClicksPerView) - stored
	if remaining <= 0 {
		return nil
	}
	if int64(len(clicks)) > remaining {
		clicks = clicks[:remaining]
	}
	models := make([]gormmodel.ResumeLinkClick, 0, len(clicks))
	for _, click := range clicks {
		dbModel, err := mappers.AnalyticsClickToDB(click)
		if err != nil {
			continue
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

func (r *ResumeLinkRepository) RecordEvents(ctx context.Context, events []*domain.AnalyticsEvent) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("record analytics events: repository dependencies are nil")
	}
	if len(events) == 0 {
		return nil
	}
	db := database.DBFromContext(ctx, r.db).WithContext(ctx)
	models := make([]gormmodel.ResumeLinkAnalyticsEvent, 0, len(events))
	for _, event := range events {
		model, err := mappers.AnalyticsEventToDB(event)
		if err != nil {
			return fmt.Errorf("record analytics events: %w", err)
		}
		models = append(models, model)
	}
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&models).Error; err != nil {
		return fmt.Errorf("record analytics events: %w", mapDBError(err))
	}
	return nil
}

func (r *ResumeLinkRepository) DeleteExpiredClicks(ctx context.Context, cutoff time.Time, limit int) (int, error) {
	if r == nil || r.db == nil {
		return 0, fmt.Errorf("delete expired analytics clicks: repository dependencies are nil")
	}
	if cutoff.IsZero() || limit <= 0 {
		return 0, fmt.Errorf("delete expired analytics clicks: invalid cutoff or limit")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return 0, fmt.Errorf("delete expired analytics clicks: %w", err)
	}
	var ids []int64
	if err := tx.WithContext(ctx).Model(&gormmodel.ResumeLinkClick{}).
		Where("occurred_at < ?", cutoff.UTC()).
		Order("occurred_at, id").Limit(limit).Pluck("id", &ids).Error; err != nil {
		return 0, fmt.Errorf("find expired analytics clicks: %w", mapDBError(err))
	}
	if len(ids) == 0 {
		return 0, nil
	}
	result := tx.WithContext(ctx).Where("id IN ?", ids).Delete(&gormmodel.ResumeLinkClick{})
	if result.Error != nil {
		return 0, fmt.Errorf("delete expired analytics clicks: %w", mapDBError(result.Error))
	}
	return int(result.RowsAffected), nil
}

func (r *ResumeLinkRepository) DeleteExpiredViews(ctx context.Context, cutoff time.Time, limit int) (int, error) {
	if r == nil || r.db == nil {
		return 0, fmt.Errorf("delete expired analytics views: repository dependencies are nil")
	}
	if cutoff.IsZero() || limit <= 0 {
		return 0, fmt.Errorf("delete expired analytics views: invalid cutoff or limit")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return 0, fmt.Errorf("delete expired analytics views: %w", err)
	}
	var ids []uuid.UUID
	if err := tx.WithContext(ctx).Model(&gormmodel.ResumeLinkView{}).
		Where("started_at < ?", cutoff.UTC()).
		Order("started_at, id").Limit(limit).Pluck("id", &ids).Error; err != nil {
		return 0, fmt.Errorf("find expired analytics views: %w", mapDBError(err))
	}
	if len(ids) == 0 {
		return 0, nil
	}
	result := tx.WithContext(ctx).Where("id IN ?", ids).Delete(&gormmodel.ResumeLinkView{})
	if result.Error != nil {
		return 0, fmt.Errorf("delete expired analytics views: %w", mapDBError(result.Error))
	}
	return int(result.RowsAffected), nil
}

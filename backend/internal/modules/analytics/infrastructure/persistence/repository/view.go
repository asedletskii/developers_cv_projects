package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	analyticscontracts "developers-cv/backend/internal/modules/analytics/contracts"
	"developers-cv/backend/internal/modules/analytics/domain"
	"developers-cv/backend/internal/modules/analytics/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/modules/analytics/infrastructure/persistence/mappers"
	"developers-cv/backend/internal/platform/database"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func (r *AnalyticsRepository) CreateView(ctx context.Context, view *domain.View) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("create analytics view: repository dependencies are nil")
	}
	if view == nil {
		return fmt.Errorf("create analytics view: view is nil")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("create analytics view: %w", err)
	}
	tx = tx.WithContext(ctx)
	model, err := mappers.ViewToDB(view)
	if err != nil {
		return fmt.Errorf("create analytics view: %w", err)
	}
	if err := tx.Create(&model).Error; err != nil {
		return fmt.Errorf("create analytics view: persist view: %w", mapDBError(err))
	}
	if err := registerVisitorFirstView(tx, view); err != nil {
		return fmt.Errorf("create analytics view: register visitor first view: %w", mapDBError(err))
	}
	return nil
}

func (r *AnalyticsRepository) GetView(
	ctx context.Context,
	id sharedvalueobjects.ID,
) (*domain.View, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("get analytics view: repository dependencies are nil")
	}
	db := database.DBFromContext(ctx, r.db).WithContext(ctx)
	var model gormmodel.AnalyticsView
	if err := db.Where("id = ?", id.UUID()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, analyticscontracts.ErrAnalyticsSubjectNotFound
		}
		return nil, fmt.Errorf("get analytics view: %w", mapDBError(err))
	}
	view, err := mappers.ViewToDomain(model)
	if err != nil {
		return nil, fmt.Errorf("get analytics view: map view: %w", err)
	}
	return view, nil
}

func (r *AnalyticsRepository) UpdateViewState(
	ctx context.Context,
	id sharedvalueobjects.ID,
	update domain.ViewStateUpdate,
) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("update analytics view state: repository dependencies are nil")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("update analytics view state: %w", err)
	}
	tx = tx.WithContext(ctx)
	var model gormmodel.AnalyticsView
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id.UUID()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return analyticscontracts.ErrAnalyticsSubjectNotFound
		}
		return fmt.Errorf("update analytics view state: load view: %w", mapDBError(err))
	}
	view, err := mappers.ViewToDomain(model)
	if err != nil {
		return fmt.Errorf("update analytics view state: map view: %w", err)
	}
	if err := view.ApplyState(update); err != nil {
		return fmt.Errorf("update analytics view state: %w", err)
	}
	state := view.State()
	updates := map[string]any{
		"confirmed_at": nullableTime(state.ConfirmedAt()), "last_seen_at": nullableTime(state.LastSeenAt()),
		"active_ms": state.ActiveMS(), "max_section": string(state.MaxSection()), "exit_section": string(state.ExitSection()),
		"ended_at": nullableTime(state.EndedAt()),
	}
	if err := tx.Model(&gormmodel.AnalyticsView{}).Where("id = ?", id.UUID()).Updates(updates).Error; err != nil {
		return fmt.Errorf("update analytics view state: persist state: %w", mapDBError(err))
	}
	return nil
}

func registerVisitorFirstView(tx *gorm.DB, view *domain.View) error {
	linkID := view.LinkID().UUID()
	visitorID := view.VisitorID().UUID()
	viewID := view.ID().UUID()

	result := tx.Exec(`
		INSERT INTO analytics.analytics_visitors (
			link_id, visitor_id, first_seen_at, first_view_id, ordinal
		)
		SELECT link_id, visitor_id, started_at, id, NULL
		FROM analytics.analytics_views
		WHERE id = ?
		ON CONFLICT (link_id, visitor_id) DO NOTHING
	`, viewID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 1 {
		if err := lockVisitorOrdinalScope(tx, linkID); err != nil {
			return err
		}
		var preceding int64
		if err := tx.Raw(`
			SELECT COUNT(*)
			FROM analytics.analytics_visitors AS visitor
			JOIN analytics.analytics_views AS candidate ON candidate.id = ?
			WHERE visitor.link_id = ?
			  AND visitor.ordinal IS NOT NULL
			  AND (visitor.first_seen_at, visitor.first_view_id)
			      < (candidate.started_at, candidate.id)
		`, viewID, linkID).Scan(&preceding).Error; err != nil {
			return err
		}
		newOrdinal := preceding + 1
		if err := tx.Exec(`
			UPDATE analytics.analytics_visitors
			SET ordinal = ordinal + 1
			WHERE link_id = ? AND ordinal >= ?
		`, linkID, newOrdinal).Error; err != nil {
			return err
		}
		if err := tx.Exec(`
			UPDATE analytics.analytics_visitors
			SET ordinal = ?
			WHERE link_id = ? AND visitor_id = ? AND ordinal IS NULL
		`, newOrdinal, linkID, visitorID).Error; err != nil {
			return err
		}
		return nil
	}

	isEarlier := func() (bool, error) {
		var earlier bool
		err := tx.Raw(`
			SELECT (candidate.started_at, candidate.id)
			       < (visitor.first_seen_at, visitor.first_view_id)
			FROM analytics.analytics_views AS candidate
			JOIN analytics.analytics_visitors AS visitor
			  ON visitor.link_id = candidate.link_id
			 AND visitor.visitor_id = candidate.visitor_id
			WHERE candidate.id = ?
		`, viewID).Scan(&earlier).Error
		return earlier, err
	}
	earlier, err := isEarlier()
	if err != nil || !earlier {
		return err
	}

	if err := lockVisitorOrdinalScope(tx, linkID); err != nil {
		return err
	}
	earlier, err = isEarlier()
	if err != nil || !earlier {
		return err
	}
	if err := tx.Exec(`
		UPDATE analytics.analytics_visitors AS visitor
		SET first_seen_at = candidate.started_at,
		    first_view_id = candidate.id
		FROM analytics.analytics_views AS candidate
		WHERE candidate.id = ?
		  AND visitor.link_id = candidate.link_id
		  AND visitor.visitor_id = candidate.visitor_id
	`, viewID).Error; err != nil {
		return err
	}
	return tx.Exec(`
		WITH ranked AS (
			SELECT visitor_id,
			       ROW_NUMBER() OVER (ORDER BY first_seen_at, first_view_id) AS new_ordinal
			FROM analytics.analytics_visitors
			WHERE link_id = ?
		)
		UPDATE analytics.analytics_visitors AS visitor
		SET ordinal = ranked.new_ordinal
		FROM ranked
		WHERE visitor.link_id = ? AND visitor.visitor_id = ranked.visitor_id
	`, linkID, linkID).Error
}

func lockVisitorOrdinalScope(tx *gorm.DB, linkID uuid.UUID) error {
	return tx.Exec(
		`SELECT pg_advisory_xact_lock(hashtextextended('analytics.visitor_ordinal:' || ?::text, 0))`,
		linkID,
	).Error
}

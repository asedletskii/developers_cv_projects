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

type ResumeLinkRepository struct {
	db *gorm.DB
}

func NewResumeLinkRepository(db *gorm.DB) *ResumeLinkRepository {
	return &ResumeLinkRepository{db: db}
}

func (r *ResumeLinkRepository) CountByResumeID(ctx context.Context, resumeID sharedvalueobjects.ID) (int64, error) {
	if r == nil || r.db == nil {
		return 0, fmt.Errorf("count resume links: repository dependencies are nil")
	}
	db, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return 0, fmt.Errorf("count resume links: %w", err)
	}
	var count int64
	if err := db.WithContext(ctx).Model(&gormmodel.ResumeLink{}).
		Where("resume_id = ?", resumeID.UUID()).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count resume links: %w", mapDBError(err))
	}
	return count, nil
}

func (r *ResumeLinkRepository) Create(ctx context.Context, link *domain.ResumeLink) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("create resume link: repository dependencies are nil")
	}
	if link == nil {
		return fmt.Errorf("create resume link: link is nil")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("create resume link: %w", err)
	}

	now := time.Now().UTC()
	model, err := mappers.ResumeLinkToDB(link, now, now)
	if err != nil {
		return fmt.Errorf("create resume link: %w", err)
	}
	if err := tx.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("create resume link: %w", mapDBError(err))
	}
	return nil
}

func (r *ResumeLinkRepository) GetByID(ctx context.Context, id sharedvalueobjects.ID) (*domain.ResumeLink, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("get resume link: repository dependencies are nil")
	}
	db := database.DBFromContext(ctx, r.db).WithContext(ctx)
	var model gormmodel.ResumeLink
	if err := db.Where("id = ?", id.UUID()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, fmt.Errorf("get resume link: %w", mapDBError(err))
	}
	link, err := mappers.ResumeLinkToDomain(model)
	if err != nil {
		return nil, fmt.Errorf("get resume link: map resume link: %w", err)
	}
	return link, nil
}

func (r *ResumeLinkRepository) GetByIDForUpdate(ctx context.Context, id sharedvalueobjects.ID) (*domain.ResumeLink, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("get resume link for update: repository dependencies are nil")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get resume link for update: %w", err)
	}
	var model gormmodel.ResumeLink
	if err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id.UUID()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, fmt.Errorf("get resume link for update: %w", mapDBError(err))
	}
	link, err := mappers.ResumeLinkToDomain(model)
	if err != nil {
		return nil, fmt.Errorf("get resume link for update: map resume link: %w", err)
	}
	return link, nil
}

func (r *ResumeLinkRepository) Rename(ctx context.Context, id sharedvalueobjects.ID, name string, version int64) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("rename resume link: repository dependencies are nil")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("rename resume link: %w", err)
	}
	result := tx.WithContext(ctx).Model(&gormmodel.ResumeLink{}).
		Where("id = ? AND version = ?", id.UUID(), version).
		Updates(map[string]any{
			"name":       name,
			"version":    gorm.Expr("version + 1"),
			"updated_at": time.Now().UTC(),
		})
	if result.Error != nil {
		return fmt.Errorf("rename resume link: %w", mapDBError(result.Error))
	}
	if result.RowsAffected != 1 {
		var exists int64
		if err := tx.WithContext(ctx).Model(&gormmodel.ResumeLink{}).Where("id = ?", id.UUID()).Count(&exists).Error; err != nil {
			return fmt.Errorf("rename resume link: check link: %w", mapDBError(err))
		}
		if exists == 0 {
			return domain.ErrEntityNotFound
		}
		return domain.ErrConcurrencyConflict
	}
	return nil
}

func (r *ResumeLinkRepository) ReplaceStatus(ctx context.Context, id sharedvalueobjects.ID, active bool, version int64) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("replace resume link status: repository dependencies are nil")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("replace resume link status: %w", err)
	}
	status := domain.ResumeLinkStatusInactive
	if active {
		status = domain.ResumeLinkStatusActive
	}
	result := tx.WithContext(ctx).Model(&gormmodel.ResumeLink{}).
		Where("id = ? AND version = ?", id.UUID(), version).
		Updates(map[string]any{
			"status":     status,
			"version":    gorm.Expr("version + 1"),
			"updated_at": time.Now().UTC(),
		})
	if result.Error != nil {
		return fmt.Errorf("replace resume link status: %w", mapDBError(result.Error))
	}
	if result.RowsAffected != 1 {
		var exists int64
		if err := tx.WithContext(ctx).Model(&gormmodel.ResumeLink{}).Where("id = ?", id.UUID()).Count(&exists).Error; err != nil {
			return fmt.Errorf("replace resume link status: check link: %w", mapDBError(err))
		}
		if exists == 0 {
			return domain.ErrEntityNotFound
		}
		return domain.ErrConcurrencyConflict
	}
	return nil
}

func (r *ResumeLinkRepository) Delete(ctx context.Context, id sharedvalueobjects.ID) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("delete resume link: repository dependencies are nil")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("delete resume link: %w", err)
	}
	result := tx.WithContext(ctx).Where("id = ?", id.UUID()).Delete(&gormmodel.ResumeLink{})
	if result.Error != nil {
		return fmt.Errorf("delete resume link: %w", mapDBError(result.Error))
	}
	if result.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
}

func (r *ResumeLinkRepository) ListIDsForUpdateByResumeIDs(
	ctx context.Context,
	resumeIDs []sharedvalueobjects.ID,
) ([]sharedvalueobjects.ID, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("list resume link ids for update: repository dependencies are nil")
	}
	if len(resumeIDs) == 0 {
		return nil, nil
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("list resume link ids for update: %w", err)
	}
	ids := make([]uuid.UUID, 0, len(resumeIDs))
	for _, resumeID := range resumeIDs {
		if resumeID.IsZero() {
			return nil, fmt.Errorf("list resume link ids for update: resume id is empty")
		}
		ids = append(ids, resumeID.UUID())
	}
	var linkIDs []uuid.UUID
	if err := tx.WithContext(ctx).
		Model(&gormmodel.ResumeLink{}).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("resume_id IN ?", ids).
		Order("id").
		Pluck("id", &linkIDs).Error; err != nil {
		return nil, fmt.Errorf("list resume link ids for update: %w", mapDBError(err))
	}
	result := make([]sharedvalueobjects.ID, 0, len(linkIDs))
	for _, linkID := range linkIDs {
		id, err := sharedvalueobjects.IDFromUUID(linkID)
		if err != nil {
			return nil, fmt.Errorf("list resume link ids for update: %w", err)
		}
		result = append(result, id)
	}
	return result, nil
}

func (r *ResumeLinkRepository) DeleteByResumeIDs(ctx context.Context, resumeIDs []sharedvalueobjects.ID) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("delete resume links by resume ids: repository dependencies are nil")
	}
	if len(resumeIDs) == 0 {
		return nil
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("delete resume links by resume ids: %w", err)
	}
	ids := make([]uuid.UUID, 0, len(resumeIDs))
	for _, resumeID := range resumeIDs {
		if resumeID.IsZero() {
			return fmt.Errorf("delete resume links by resume ids: resume id is empty")
		}
		ids = append(ids, resumeID.UUID())
	}
	if err := tx.WithContext(ctx).Where("resume_id IN ?", ids).Delete(&gormmodel.ResumeLink{}).Error; err != nil {
		return fmt.Errorf("delete resume links by resume ids: %w", mapDBError(err))
	}
	return nil
}

package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	resumevalueobjects "developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/modules/resume_link/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume_link/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/platform/database"
)

type ResumeLinkRepository struct {
	db *gorm.DB
}

func NewResumeLinkRepository(db *gorm.DB) *ResumeLinkRepository {
	return &ResumeLinkRepository{db: db}
}

func (r *ResumeLinkRepository) Create(ctx context.Context, link *domain.ResumeLink) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("create resume link: repository dependencies are nil")
	}
	if link == nil || link.ID().IsZero() || link.ResumeID().IsZero() || link.Name() == "" || link.Version() < 1 {
		return fmt.Errorf("create resume link: invalid link")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("create resume link: %w", err)
	}
	var lockedResume struct {
		ID uuid.UUID `gorm:"column:id"`
	}
	if err := tx.WithContext(ctx).Table("resumes").Clauses(clause.Locking{Strength: "UPDATE"}).Select("id").Where("id = ?", link.ResumeID().UUID()).First(&lockedResume).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resumedomain.ErrEntityNotFound
		}
		return fmt.Errorf("create resume link: lock resume: %w", err)
	}

	var count int64
	if err := tx.WithContext(ctx).Model(&gormmodel.ResumeLink{}).Where("resume_id = ?", link.ResumeID().UUID()).Count(&count).Error; err != nil {
		return fmt.Errorf("create resume link: count links: %w", err)
	}
	if count >= domain.MaxResumeLinksPerResume {
		return domain.ErrResumeLinkLimitExceeded
	}

	now := time.Now().UTC()
	model := gormmodel.ResumeLink{
		ID: link.ID().UUID(), ResumeID: link.ResumeID().UUID(), Name: link.Name(),
		Status: statusOf(link), Version: link.Version(), CreatedAt: now, UpdatedAt: now,
	}
	if err := tx.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("create resume link: %w", mapDBError(err))
	}
	return nil
}

func (r *ResumeLinkRepository) GetByID(ctx context.Context, id valueobjects.ID) (*domain.ResumeLink, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("get resume link: repository dependencies are nil")
	}
	if id.IsZero() {
		return nil, fmt.Errorf("get resume link: id is empty")
	}
	db := database.DBFromContext(ctx, r.db).WithContext(ctx)
	var model gormmodel.ResumeLink
	if err := db.Where("id = ?", id.UUID()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, fmt.Errorf("get resume link: %w", err)
	}
	return mapResumeLink(model)
}

func (r *ResumeLinkRepository) ReplaceStatus(ctx context.Context, id valueobjects.ID, active bool, version int64) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("replace resume link status: repository dependencies are nil")
	}
	if id.IsZero() || version < 1 {
		return fmt.Errorf("replace resume link status: invalid link or expected version")
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
			return fmt.Errorf("replace resume link status: check link: %w", err)
		}
		if exists == 0 {
			return domain.ErrEntityNotFound
		}
		return domain.ErrConcurrencyConflict
	}
	return nil
}

func (r *ResumeLinkRepository) Delete(ctx context.Context, id valueobjects.ID) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("delete resume link: repository dependencies are nil")
	}
	if id.IsZero() {
		return fmt.Errorf("delete resume link: id is empty")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("delete resume link: %w", err)
	}
	result := tx.WithContext(ctx).Where("id = ?", id.UUID()).Delete(&gormmodel.ResumeLink{})
	if result.Error != nil {
		return fmt.Errorf("delete resume link: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
}

func (r *ResumeLinkRepository) RecordEvent(ctx context.Context, event *domain.ResumeLinkEvent) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("record resume link event: repository dependencies are nil")
	}
	if event == nil || event.ID().IsZero() || event.LinkID().IsZero() || !event.Type().IsValid() || event.OccurredAt().IsZero() {
		return fmt.Errorf("record resume link event: invalid event")
	}
	db := database.DBFromContext(ctx, r.db)
	model := gormmodel.ResumeLinkEvent{ID: event.ID().UUID(), LinkID: event.LinkID().UUID(), EventType: string(event.Type()), OccurredAt: event.OccurredAt().UTC()}
	if err := db.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("record resume link event: %w", mapDBError(err))
	}
	return nil
}

func statusOf(link *domain.ResumeLink) string {
	if link.IsActive() {
		return domain.ResumeLinkStatusActive
	}
	return domain.ResumeLinkStatusInactive
}

func mapResumeLink(model gormmodel.ResumeLink) (*domain.ResumeLink, error) {
	id, err := valueobjects.IDFromUUID(model.ID)
	if err != nil {
		return nil, fmt.Errorf("map resume link id: %w", err)
	}
	resumeID, err := resumevalueobjects.IDFromUUID(model.ResumeID)
	if err != nil {
		return nil, fmt.Errorf("map resume link resume id: %w", err)
	}
	link, err := domain.NewResumeLink(domain.ResumeLinkParams{
		ID: id, ResumeID: resumeID, Name: model.Name, Version: model.Version,
	})
	if err != nil {
		return nil, fmt.Errorf("map resume link: %w", err)
	}
	if model.Status == domain.ResumeLinkStatusInactive {
		if err := link.Deactivate(); err != nil {
			return nil, fmt.Errorf("map resume link status: %w", err)
		}
	} else if model.Status != domain.ResumeLinkStatusActive {
		return nil, fmt.Errorf("map resume link: invalid status %q", model.Status)
	}
	return link, nil
}

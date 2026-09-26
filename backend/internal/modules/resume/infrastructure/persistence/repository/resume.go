package repository

import (
	"context"
	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/platform/database"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func (r *ResumeRepository) CreateResume(ctx context.Context, c ports.CreateResumeRepositoryCommand) error {
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("create resume: %w", err)
	}
	if c.Resume == nil {
		return fmt.Errorf("resume is nil")
	}
	m := gormmodel.Resume{ID: c.Resume.ID().UUID(), UserID: c.Resume.UserID().UUID(), Title: c.Resume.Title(), Language: string(c.Resume.Language()), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}
	if err := tx.WithContext(ctx).Create(&m).Error; err != nil {
		return fmt.Errorf("create resume: %w", mapDBError(err))
	}
	return nil
}
func (r *ResumeRepository) GetResume(ctx context.Context, q ports.GetResumeRepositoryCommand) (*domain.Resume, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("get resume: repository dependencies are nil")
	}
	tx := database.DBFromContext(ctx, r.db)
	var m gormmodel.Resume
	if err := tx.WithContext(ctx).Where("id = ? AND user_id = ?", q.ResumeID.UUID(), q.UserID.UUID()).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrEntityNotFound
		}
		return nil, fmt.Errorf("get resume: %w", err)
	}
	id, e := valueobjects.IDFromUUID(m.ID)
	if e != nil {
		return nil, e
	}
	uid, e := useridentity.IDFromUUID(m.UserID)
	if e != nil {
		return nil, e
	}
	return domain.NewResume(domain.ResumeParams{ID: id, UserID: uid, Title: m.Title, Language: enums.ResumeLanguage(m.Language)})
}
func (r *ResumeRepository) ReplaceResumeTitle(ctx context.Context, c ports.ReplaceResumeTitleRepositoryCommand) error {
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("replace resume title: %w", err)
	}
	if c.ResumeID.IsZero() || c.UserID.IsZero() {
		return fmt.Errorf("invalid resume identity")
	}
	res := tx.WithContext(ctx).Model(&gormmodel.Resume{}).Where("id = ? AND user_id = ?", c.ResumeID.UUID(), c.UserID.UUID()).Updates(map[string]any{"title": c.Title, "updated_at": time.Now().UTC()})
	if res.Error != nil {
		return fmt.Errorf("replace resume title: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
}
func (r *ResumeRepository) DeleteResume(ctx context.Context, q ports.DeleteResumeRepositoryCommand) error {
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("delete resume: %w", err)
	}
	res := tx.WithContext(ctx).Where("id = ? AND user_id = ?", q.ResumeID.UUID(), q.UserID.UUID()).Delete(&gormmodel.Resume{})
	if res.Error != nil {
		return fmt.Errorf("delete resume: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
}

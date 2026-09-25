package repository

import (
	"context"
	"fmt"
	"time"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/platform/database"

	"gorm.io/gorm"
)

var _ ports.CreateResumeRootRepository = (*ResumeRootRepository)(nil)
var _ ports.GetResumeRootRepository = (*ResumeRootRepository)(nil)
var _ ports.UpdateResumeRootRepository = (*ResumeRootRepository)(nil)
var _ ports.DeleteResumeRootRepository = (*ResumeRootRepository)(nil)

func (r *ResumeRootRepository) CreateRoot(ctx context.Context, c ports.CreateResumeRootRepositoryCommand) error {
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("create resume root: %w", err)
	}
	m, p, sections, err := rootModels(c.Resume)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	m.CreatedAt = now
	m.UpdatedAt = now
	if err := tx.WithContext(ctx).Create(&m).Error; err != nil {
		return fmt.Errorf("create resume root: %w", err)
	}
	if err := tx.WithContext(ctx).Create(&p).Error; err != nil {
		return fmt.Errorf("create resume presentation: %w", err)
	}
	if len(sections) > 0 {
		if err := tx.WithContext(ctx).Create(&sections).Error; err != nil {
			return fmt.Errorf("create resume presentation sections: %w", err)
		}
	}
	return nil
}
func (r *ResumeRootRepository) GetRoot(ctx context.Context, q ports.GetResumeRootRepositoryCommand) (*domain.ResumeRoot, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("get resume root: repository dependencies are nil")
	}
	tx := database.DBFromContext(ctx, r.db)
	var m gormmodel.ResumeRoot
	if err := tx.WithContext(ctx).Where("id = ? AND user_id = ?", q.ResumeID.UUID(), q.UserID.UUID()).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrEntityNotFound
		}
		return nil, fmt.Errorf("get resume root: %w", err)
	}
	var p gormmodel.ResumePresentation
	if err := tx.WithContext(ctx).Where("resume_id = ?", m.ID).First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("get resume presentation: %w", err)
		}
		return nil, fmt.Errorf("get resume presentation: %w", err)
	}
	var ss []gormmodel.ResumePresentationSection
	if err := tx.WithContext(ctx).Where("resume_id = ?", m.ID).Order("position").Find(&ss).Error; err != nil {
		return nil, fmt.Errorf("get resume presentation sections: %w", err)
	}
	return mapRoot(m, p, ss)
}
func (r *ResumeRootRepository) UpdateRoot(ctx context.Context, c ports.UpdateResumeRootRepositoryCommand) error {
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("update resume root: %w", err)
	}
	m, p, sections, err := rootModels(c.Resume)
	if err != nil {
		return err
	}
	result := tx.WithContext(ctx).Model(&gormmodel.ResumeRoot{}).Where("id = ? AND user_id = ?", m.ID, m.UserID).Updates(map[string]any{"title": m.Title, "updated_at": time.Now().UTC()})
	if result.Error != nil {
		return fmt.Errorf("update resume root: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	presentationResult := tx.WithContext(ctx).Model(&gormmodel.ResumePresentation{}).Where("resume_id = ?", m.ID).Updates(map[string]any{
		"greeting_first_ru": p.GreetingFirstRU, "greeting_first_en": p.GreetingFirstEN, "greeting_second_ru": p.GreetingSecondRU, "greeting_second_en": p.GreetingSecondEN, "greeting_highlight_ru": p.GreetingHighlightRU, "greeting_highlight_en": p.GreetingHighlightEN, "greeting_third_ru": p.GreetingThirdRU, "greeting_third_en": p.GreetingThirdEN, "availability_status": p.AvailabilityStatus, "availability_from": p.AvailabilityFrom, "availability_note_ru": p.AvailabilityNoteRU, "availability_note_en": p.AvailabilityNoteEN, "resume_link": p.ResumeLink, "theme_primary": p.ThemePrimary, "theme_secondary": p.ThemeSecondary,
	})
	if presentationResult.Error != nil {
		return fmt.Errorf("update resume presentation: %w", presentationResult.Error)
	}
	if presentationResult.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	if err := tx.WithContext(ctx).Where("resume_id = ?", m.ID).Delete(&gormmodel.ResumePresentationSection{}).Error; err != nil {
		return fmt.Errorf("replace resume presentation sections: %w", err)
	}
	if len(sections) > 0 {
		if err := tx.WithContext(ctx).Create(&sections).Error; err != nil {
			return fmt.Errorf("create resume presentation sections: %w", err)
		}
	}
	return nil
}
func (r *ResumeRootRepository) DeleteRoot(ctx context.Context, q ports.DeleteResumeRootRepositoryCommand) error {
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("delete resume root: %w", err)
	}
	result := tx.WithContext(ctx).Where("id = ? AND user_id = ?", q.ResumeID.UUID(), q.UserID.UUID()).Delete(&gormmodel.ResumeRoot{})
	if result.Error != nil {
		return fmt.Errorf("delete resume root: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
}

package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/mappers"
)

func persistOverview(tx *gorm.DB, resume *domain.Resume) error {
	o := resume.Overview()
	if o == nil {
		if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeOverviewSection{}).Error; err != nil {
			return fmt.Errorf("delete resume overview sections: %w", mapDBError(err))
		}
		if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeOverview{}).Error; err != nil {
			return fmt.Errorf("delete resume overview: %w", mapDBError(err))
		}
		return nil
	}
	p, sections := mappers.OverviewToDBModels(resume.ID(), o)
	result := tx.Model(&gormmodel.ResumeOverview{}).
		Where("resume_id = ?", resume.ID().UUID()).
		Updates(map[string]any{
			"greeting_first":      p.GreetingFirst,
			"greeting_second":     p.GreetingSecond,
			"greeting_highlight":  p.GreetingHighlight,
			"greeting_third":      p.GreetingThird,
			"availability_status": p.AvailabilityStatus,
			"availability_from":   p.AvailabilityFrom,
			"availability_note":   p.AvailabilityNote,
			"resume_link":         p.ResumeLink,
			"theme_primary":       p.ThemePrimary,
			"theme_secondary":     p.ThemeSecondary,
		})
	if result.Error != nil {
		return fmt.Errorf("update resume overview: %w", mapDBError(result.Error))
	}
	if result.RowsAffected == 0 {
		if err := tx.Create(&p).Error; err != nil {
			return fmt.Errorf("create resume overview: %w", mapDBError(err))
		}
	}
	if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeOverviewSection{}).Error; err != nil {
		return fmt.Errorf("replace resume overview sections: %w", mapDBError(err))
	}
	if len(sections) > 0 {
		if err := tx.Create(&sections).Error; err != nil {
			return fmt.Errorf("create resume overview sections: %w", mapDBError(err))
		}
	}
	return nil
}

func loadResumeOverview(ctx context.Context, tx *gorm.DB, resumeID uuid.UUID) (*domain.ResumeOverview, error) {
	var overviewModel gormmodel.ResumeOverview
	err := tx.WithContext(ctx).Where("resume_id = ?", resumeID).First(&overviewModel).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get resume overview: %w", mapDBError(err))
	}

	var sections []gormmodel.ResumeOverviewSection
	if err := tx.WithContext(ctx).Where("resume_id = ?", resumeID).Order("position").Find(&sections).Error; err != nil {
		return nil, fmt.Errorf("get resume overview sections: %w", mapDBError(err))
	}
	return mappers.OverviewToDomain(overviewModel, sections), nil
}

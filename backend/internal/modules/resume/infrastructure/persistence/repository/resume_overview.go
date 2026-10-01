package repository

import (
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/mappers"
)

func persistOverview(tx *gorm.DB, resume *domain.Resume) error {
	o := resume.Overview()
	if o == nil {
		return fmt.Errorf("save resume overview: overview is required")
	}
	p, sections, err := mappers.OverviewToDBModels(resume.ID(), o)
	if err != nil {
		return fmt.Errorf("save resume overview: %w", err)
	}
	result := tx.Model(&gormmodel.ResumeOverview{}).Where("resume_id = ?", resume.ID().UUID()).Updates(map[string]any{
		"greeting_first": p.GreetingFirst, "greeting_second": p.GreetingSecond, "greeting_highlight": p.GreetingHighlight, "greeting_third": p.GreetingThird,
		"availability_status": p.AvailabilityStatus, "availability_from": p.AvailabilityFrom, "availability_note": p.AvailabilityNote, "resume_link": p.ResumeLink, "theme_primary": p.ThemePrimary, "theme_secondary": p.ThemeSecondary,
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

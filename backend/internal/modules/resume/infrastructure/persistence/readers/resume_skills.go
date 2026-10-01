package readers

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

func (r *ResumeReader) loadSkills(ctx context.Context, db *gorm.DB, resumeID uuid.UUID, view *queries.ResumeView) error {
	var model gormmodel.ResumeSkills
	if err := db.WithContext(ctx).Where("resume_id = ?", resumeID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return fmt.Errorf("get skills: %w", err)
	}
	var hard []gormmodel.ResumeSkillHard
	if err := db.WithContext(ctx).Where("skills_id = ?", model.ID).Order("position").Find(&hard).Error; err != nil {
		return fmt.Errorf("get hard skills: %w", err)
	}
	var soft []gormmodel.ResumeSkillSoft
	if err := db.WithContext(ctx).Where("skills_id = ?", model.ID).Order("position").Find(&soft).Error; err != nil {
		return fmt.Errorf("get soft skills: %w", err)
	}
	result := queries.SkillsView{ID: model.ID.String(), HardSkills: make([]string, 0, len(hard)), SoftSkills: make([]string, 0, len(soft))}
	for _, item := range hard {
		result.HardSkills = append(result.HardSkills, item.Value)
	}
	for _, item := range soft {
		result.SoftSkills = append(result.SoftSkills, item.Value)
	}
	view.Skills = &result
	return nil
}

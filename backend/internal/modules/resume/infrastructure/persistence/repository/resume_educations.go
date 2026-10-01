package repository

import (
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/mappers"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func persistEducations(tx *gorm.DB, resume *domain.Resume) error {
	models := make([]gormmodel.ResumeEducation, 0, len(resume.Educations()))
	for position, item := range resume.Educations() {
		if item == nil {
			return fmt.Errorf("%w: nil education", domain.ErrResumeInvariant)
		}
		if item.ResumeID() != resume.ID() {
			return fmt.Errorf("%w: education does not belong to resume", domain.ErrResumeInvariant)
		}
		model, err := mappers.EducationToDB(item)
		if err != nil {
			return err
		}
		model.SortOrder = position
		models = append(models, model)
	}

	if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeEducation{}).Error; err != nil {
		return fmt.Errorf("replace resume educations: %w", mapDBError(err))
	}
	if len(models) > 0 {
		if err := tx.Create(&models).Error; err != nil {
			return fmt.Errorf("create resume educations: %w", mapDBError(err))
		}
	}
	return nil
}
func loadEducations(tx *gorm.DB, resumeID identity.ID) ([]*domain.Education, error) {
	var rows []gormmodel.ResumeEducation
	if err := tx.Where("resume_id = ?", resumeID.UUID()).Order("sort_order").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("get resume educations: %w", mapDBError(err))
	}
	out := make([]*domain.Education, 0, len(rows))
	for _, row := range rows {
		x, err := mappers.EducationToDomain(row)
		if err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	if len(out) == 0 {
		return nil, nil
	}
	return out, nil
}

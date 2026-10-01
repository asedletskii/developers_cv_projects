package repository

import (
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/mappers"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func persistCourses(tx *gorm.DB, resume *domain.Resume) error {
	models := make([]gormmodel.ResumeCourse, 0, len(resume.Courses()))
	for position, item := range resume.Courses() {
		if item == nil {
			return fmt.Errorf("%w: nil course", domain.ErrResumeInvariant)
		}
		if item.ResumeID() != resume.ID() {
			return fmt.Errorf("%w: course does not belong to resume", domain.ErrResumeInvariant)
		}
		model, err := mappers.CourseToDB(item)
		if err != nil {
			return err
		}
		model.SortOrder = position
		models = append(models, model)
	}

	if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeCourse{}).Error; err != nil {
		return fmt.Errorf("replace resume courses: %w", mapDBError(err))
	}
	if len(models) > 0 {
		if err := tx.Create(&models).Error; err != nil {
			return fmt.Errorf("create resume courses: %w", mapDBError(err))
		}
	}
	return nil
}

func loadCourses(tx *gorm.DB, resumeID identity.ID) ([]*domain.Course, error) {
	var rows []gormmodel.ResumeCourse
	if err := tx.Where("resume_id = ?", resumeID.UUID()).Order("sort_order").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("get resume courses: %w", mapDBError(err))
	}
	out := make([]*domain.Course, 0, len(rows))
	for _, row := range rows {
		x, err := mappers.CourseToDomain(row)
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

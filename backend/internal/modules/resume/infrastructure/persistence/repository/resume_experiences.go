package repository

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/mappers"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func loadExperiences(tx *gorm.DB, resumeID sharedvalueobjects.ID) ([]*domain.Experience, error) {
	var models []gormmodel.ResumeExperience
	if err := tx.Where("resume_id = ?", resumeID.UUID()).Order("sort_order").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("load resume experiences: %w", mapDBError(err))
	}
	if len(models) == 0 {
		return make([]*domain.Experience, 0), nil
	}

	ids := make([]uuid.UUID, 0, len(models))
	for _, model := range models {
		ids = append(ids, model.ID)
	}
	achievements, err := loadExperienceAchievements(tx, ids)
	if err != nil {
		return nil, err
	}
	tasks, err := loadExperienceTasks(tx, ids)
	if err != nil {
		return nil, err
	}
	stack, err := loadExperienceStack(tx, ids)
	if err != nil {
		return nil, err
	}

	return mappers.ExperiencesToDomain(models, achievements, tasks, stack)
}

func loadExperienceAchievements(tx *gorm.DB, experienceIDs []uuid.UUID) (map[uuid.UUID][]gormmodel.ResumeExperienceAchievement, error) {
	var rows []gormmodel.ResumeExperienceAchievement
	if err := tx.Where("experience_id IN ?", experienceIDs).Order("experience_id, position").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("load resume experience achievements: %w", mapDBError(err))
	}
	grouped := make(map[uuid.UUID][]gormmodel.ResumeExperienceAchievement, len(experienceIDs))
	for _, row := range rows {
		grouped[row.ExperienceID] = append(grouped[row.ExperienceID], row)
	}
	return grouped, nil
}

func loadExperienceTasks(tx *gorm.DB, experienceIDs []uuid.UUID) (map[uuid.UUID][]gormmodel.ResumeExperienceTask, error) {
	var rows []gormmodel.ResumeExperienceTask
	if err := tx.Where("experience_id IN ?", experienceIDs).Order("experience_id, position").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("load resume experience tasks: %w", mapDBError(err))
	}
	grouped := make(map[uuid.UUID][]gormmodel.ResumeExperienceTask, len(experienceIDs))
	for _, row := range rows {
		grouped[row.ExperienceID] = append(grouped[row.ExperienceID], row)
	}
	return grouped, nil
}

func loadExperienceStack(tx *gorm.DB, experienceIDs []uuid.UUID) (map[uuid.UUID][]gormmodel.ResumeExperienceStack, error) {
	var rows []gormmodel.ResumeExperienceStack
	if err := tx.Where("experience_id IN ?", experienceIDs).Order("experience_id, position").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("load resume experience stack: %w", mapDBError(err))
	}
	grouped := make(map[uuid.UUID][]gormmodel.ResumeExperienceStack, len(experienceIDs))
	for _, row := range rows {
		grouped[row.ExperienceID] = append(grouped[row.ExperienceID], row)
	}
	return grouped, nil
}

func persistExperiences(tx *gorm.DB, resume *domain.Resume) error {
	models, achievements, tasks, stack := mappers.ExperiencesToDBModels(resume.Experiences())
	if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeExperience{}).Error; err != nil {
		return fmt.Errorf("replace resume experiences: %w", mapDBError(err))
	}
	if len(models) > 0 {
		if err := tx.Create(&models).Error; err != nil {
			return fmt.Errorf("create resume experiences: %w", mapDBError(err))
		}
	}
	if len(achievements) > 0 {
		if err := tx.Create(&achievements).Error; err != nil {
			return fmt.Errorf("create resume experience achievements: %w", mapDBError(err))
		}
	}
	if len(tasks) > 0 {
		if err := tx.Create(&tasks).Error; err != nil {
			return fmt.Errorf("create resume experience tasks: %w", mapDBError(err))
		}
	}
	if len(stack) > 0 {
		if err := tx.Create(&stack).Error; err != nil {
			return fmt.Errorf("create resume experience stack: %w", mapDBError(err))
		}
	}
	return nil
}

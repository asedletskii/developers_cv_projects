package repository

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/mappers"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func loadExperiences(tx *gorm.DB, resumeID identity.ID) ([]*domain.Experience, error) {
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

	items := make([]*domain.Experience, 0, len(models))
	for _, model := range models {
		item, err := mappers.ExperienceToDomain(model, achievements[model.ID], tasks[model.ID], stack[model.ID])
		if err != nil {
			return nil, fmt.Errorf("load resume experience: %w", err)
		}
		items = append(items, item)
	}
	return items, nil
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
	models := make([]gormmodel.ResumeExperience, 0, len(resume.Experiences()))
	items := resume.Experiences()
	for position, item := range items {
		if item == nil {
			return fmt.Errorf("%w: nil experience", domain.ErrResumeInvariant)
		}
		if item.ResumeID() != resume.ID() {
			return fmt.Errorf("%w: experience does not belong to resume", domain.ErrResumeInvariant)
		}
		model, err := mappers.ExperienceToDB(item)
		if err != nil {
			return fmt.Errorf("map resume experience: %w", err)
		}
		model.SortOrder = position
		models = append(models, model)
	}

	if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeExperience{}).Error; err != nil {
		return fmt.Errorf("replace resume experiences: %w", mapDBError(err))
	}
	for i, model := range models {
		if err := tx.Create(&model).Error; err != nil {
			return fmt.Errorf("create resume experience: %w", mapDBError(err))
		}
		if err := replaceExperienceChildren(tx, model.ID, items[i]); err != nil {
			return fmt.Errorf("replace resume experience children: %w", err)
		}
	}
	return nil
}

func replaceExperienceChildren(tx *gorm.DB, experienceID uuid.UUID, item *domain.Experience) error {
	if item == nil {
		return fmt.Errorf("experience is required")
	}
	experienceIdentity, err := identity.IDFromUUID(experienceID)
	if err != nil {
		return fmt.Errorf("parse experience id: %w", err)
	}
	if err := tx.Where("experience_id = ?", experienceID).Delete(&gormmodel.ResumeExperienceAchievement{}).Error; err != nil {
		return err
	}
	if err := tx.Where("experience_id = ?", experienceID).Delete(&gormmodel.ResumeExperienceTask{}).Error; err != nil {
		return err
	}
	if err := tx.Where("experience_id = ?", experienceID).Delete(&gormmodel.ResumeExperienceStack{}).Error; err != nil {
		return err
	}
	achievements := mappers.ExperienceAchievementsToDB(experienceIdentity, item.Achievements())
	if len(achievements) > 0 {
		if err := tx.Create(&achievements).Error; err != nil {
			return err
		}
	}
	tasks := mappers.ExperienceTasksToDB(experienceIdentity, item.Tasks())
	if len(tasks) > 0 {
		if err := tx.Create(&tasks).Error; err != nil {
			return err
		}
	}
	stack := mappers.ExperienceStackToDB(experienceIdentity, item.Stack())
	if len(stack) > 0 {
		if err := tx.Create(&stack).Error; err != nil {
			return err
		}
	}
	return nil
}

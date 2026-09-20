package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/mappers"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func persistSkills(tx *gorm.DB, resume *domain.Resume) error {
	skills := resume.Skills()
	if skills == nil {
		if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeSkills{}).Error; err != nil {
			return fmt.Errorf("delete resume skills: %w", mapDBError(err))
		}
		return nil
	}

	m, hardModels, softModels := mappers.SkillsToDBModels(skills)
	result := tx.Model(&gormmodel.ResumeSkills{}).
		Where("id = ? AND resume_id = ?", m.ID, resume.ID().UUID()).UpdateColumn("resume_id", resume.ID().UUID())
	if result.Error != nil {
		return fmt.Errorf("update resume skills: %w", mapDBError(result.Error))
	}
	if result.RowsAffected == 0 {
		if err := tx.Create(&m).Error; err != nil {
			return fmt.Errorf("create resume skills: %w", mapDBError(err))
		}
	}
	if err := tx.Where("skills_id = ?", m.ID).Delete(&gormmodel.ResumeSkillHard{}).Error; err != nil {
		return fmt.Errorf("replace resume hard skills: %w", mapDBError(err))
	}
	if err := tx.Where("skills_id = ?", m.ID).Delete(&gormmodel.ResumeSkillSoft{}).Error; err != nil {
		return fmt.Errorf("replace resume soft skills: %w", mapDBError(err))
	}
	if len(hardModels) > 0 {
		if err := tx.Create(&hardModels).Error; err != nil {
			return fmt.Errorf("create resume hard skills: %w", mapDBError(err))
		}
	}
	if len(softModels) > 0 {
		if err := tx.Create(&softModels).Error; err != nil {
			return fmt.Errorf("create resume soft skills: %w", mapDBError(err))
		}
	}
	return nil
}

func loadSkills(tx *gorm.DB, resumeID sharedvalueobjects.ID) (*domain.Skills, error) {
	var m gormmodel.ResumeSkills
	err := tx.Where("resume_id = ?", resumeID.UUID()).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get resume skills: %w", mapDBError(err))
	}
	var hardRows []gormmodel.ResumeSkillHard
	if err := tx.Where("skills_id = ?", m.ID).Order("position").Find(&hardRows).Error; err != nil {
		return nil, fmt.Errorf("get resume hard skills: %w", mapDBError(err))
	}
	var softRows []gormmodel.ResumeSkillSoft
	if err := tx.Where("skills_id = ?", m.ID).Order("position").Find(&softRows).Error; err != nil {
		return nil, fmt.Errorf("get resume soft skills: %w", mapDBError(err))
	}
	skills, err := mappers.SkillsToDomain(m, hardRows, softRows)
	if err != nil {
		return nil, fmt.Errorf("map resume skills: %w", err)
	}
	return skills, nil
}

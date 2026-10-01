package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/mappers"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func persistSkills(tx *gorm.DB, resume *domain.Resume) error {
	skills := resume.Skills()
	if skills == nil {
		var current gormmodel.ResumeSkills
		err := tx.Where("resume_id = ?", resume.ID().UUID()).First(&current).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("load resume skills: %w", mapDBError(err))
		}
		if err := tx.Where("skills_id = ?", current.ID).Delete(&gormmodel.ResumeSkillHard{}).Error; err != nil {
			return fmt.Errorf("delete resume hard skills: %w", mapDBError(err))
		}
		if err := tx.Where("skills_id = ?", current.ID).Delete(&gormmodel.ResumeSkillSoft{}).Error; err != nil {
			return fmt.Errorf("delete resume soft skills: %w", mapDBError(err))
		}
		if err := tx.Delete(&current).Error; err != nil {
			return fmt.Errorf("delete resume skills: %w", mapDBError(err))
		}
		return nil
	}

	m, hardModels, softModels, err := mappers.SkillsToDBModels(skills)
	if err != nil {
		return fmt.Errorf("map resume skills: %w", err)
	}
	var current gormmodel.ResumeSkills
	err = tx.Where("resume_id = ?", resume.ID().UUID()).First(&current).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := tx.Create(&m).Error; err != nil {
			return fmt.Errorf("create resume skills: %w", mapDBError(err))
		}
	} else if err != nil {
		return fmt.Errorf("load resume skills: %w", mapDBError(err))
	} else if current.ID != m.ID {
		return fmt.Errorf("%w: skills identity mismatch", domain.ErrResumeInvariant)
	}

	if err := tx.Where("skills_id = ?", m.ID).Delete(&gormmodel.ResumeSkillHard{}).Error; err != nil {
		return fmt.Errorf("replace resume hard skills: %w", mapDBError(err))
	}
	if err := tx.Where("skills_id = ?", m.ID).Delete(&gormmodel.ResumeSkillSoft{}).Error; err != nil {
		return fmt.Errorf("replace resume soft skills: %w", mapDBError(err))
	}
	for _, model := range hardModels {
		if err := tx.Create(&model).Error; err != nil {
			return fmt.Errorf("create resume hard skill: %w", mapDBError(err))
		}
	}
	for _, model := range softModels {
		if err := tx.Create(&model).Error; err != nil {
			return fmt.Errorf("create resume soft skill: %w", mapDBError(err))
		}
	}
	return nil
}

func loadSkills(tx *gorm.DB, resumeID identity.ID) (*domain.Skills, error) {
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
	return mappers.SkillsToDomain(m, hardRows, softRows, resumeID)
}

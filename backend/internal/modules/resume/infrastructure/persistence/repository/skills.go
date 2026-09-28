package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
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

	m := gormmodel.ResumeSkills{ID: skills.ID().UUID(), ResumeID: skills.ResumeID().UUID()}
	var current gormmodel.ResumeSkills
	err := tx.Where("resume_id = ?", resume.ID().UUID()).First(&current).Error
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
	for i, value := range skills.HardSkills() {
		if err := tx.Create(&gormmodel.ResumeSkillHard{SkillsID: m.ID, Value: string(value), Position: i}).Error; err != nil {
			return fmt.Errorf("create resume hard skill: %w", mapDBError(err))
		}
	}
	for i, value := range skills.SoftSkills() {
		if err := tx.Create(&gormmodel.ResumeSkillSoft{SkillsID: m.ID, Value: string(value), Position: i}).Error; err != nil {
			return fmt.Errorf("create resume soft skill: %w", mapDBError(err))
		}
	}
	return nil
}

func loadSkills(tx *gorm.DB, resumeID valueobjects.ID) (*domain.Skills, error) {
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
	id, err := valueobjects.IDFromUUID(m.ID)
	if err != nil {
		return nil, err
	}
	hard := make([]enums.HardSkill, len(hardRows))
	for i, row := range hardRows {
		hard[i] = enums.HardSkill(row.Value)
	}
	soft := make([]enums.SoftSkill, len(softRows))
	for i, row := range softRows {
		soft[i] = enums.SoftSkill(row.Value)
	}
	return domain.NewSkills(domain.SkillsParams{ID: id, ResumeID: resumeID, Hard: hard, Soft: soft})
}

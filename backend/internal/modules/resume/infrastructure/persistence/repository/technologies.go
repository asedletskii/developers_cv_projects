package repository

import (
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

func persistTechnologies(tx *gorm.DB, resume *domain.Resume) error {
	technologies := resume.Technologies()
	if technologies == nil {
		if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeTechnology{}).Error; err != nil {
			return fmt.Errorf("delete resume technologies: %w", mapDBError(err))
		}
		return nil
	}

	models := make([]gormmodel.ResumeTechnology, 0, len(technologies.Entries()))
	seenIDs := make(map[valueobjects.ID]struct{}, len(technologies.Entries()))
	for _, entry := range technologies.Entries() {
		id := entry.ID()
		if id.IsZero() {
			return fmt.Errorf("%w: technology id is required", domain.ErrResumeInvariant)
		}
		if _, ok := seenIDs[id]; ok {
			return fmt.Errorf("%w: duplicate technology id", domain.ErrInvalidEntity)
		}
		seenIDs[id] = struct{}{}
		model, err := technologyModel(entry, resume.ID())
		if err != nil {
			return fmt.Errorf("map resume technology: %w", err)
		}
		models = append(models, model)
	}

	if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeTechnology{}).Error; err != nil {
		return fmt.Errorf("replace resume technologies: %w", mapDBError(err))
	}
	if len(models) > 0 {
		if err := tx.Create(&models).Error; err != nil {
			return fmt.Errorf("create resume technologies: %w", mapDBError(err))
		}
	}
	return nil
}

func loadTechnologies(tx *gorm.DB, resumeID valueobjects.ID) (*domain.ResumeTechnologies, error) {
	var rows []gormmodel.ResumeTechnology
	if err := tx.Where("resume_id = ?", resumeID.UUID()).Order("position").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("get resume technologies: %w", mapDBError(err))
	}
	if len(rows) == 0 {
		return nil, nil
	}

	collectionID, err := valueobjects.IDFromUUID(rows[0].ID)
	if err != nil {
		return nil, err
	}
	entries := make([]domain.ResumeTechnologyEntryParams, 0, len(rows))
	for _, row := range rows {
		id, err := valueobjects.IDFromUUID(row.ID)
		if err != nil {
			return nil, err
		}
		var grade *enums.TechnologyGrade
		if row.Grade != nil {
			value := enums.TechnologyGrade(*row.Grade)
			if !value.IsValid() {
				return nil, domain.ErrInvalidEntity
			}
			grade = &value
		}
		entries = append(entries, domain.ResumeTechnologyEntryParams{
			ID: id, Name: row.Name, Grade: grade, Position: row.Position,
		})
	}
	return domain.NewResumeTechnologies(domain.ResumeTechnologiesParams{
		ID: collectionID, ResumeID: resumeID, Entries: entries,
	})
}

func technologyModel(x domain.ResumeTechnologyEntry, resumeID valueobjects.ID) (gormmodel.ResumeTechnology, error) {
	if x.ID().IsZero() {
		return gormmodel.ResumeTechnology{}, fmt.Errorf("technology id is required")
	}
	m := gormmodel.ResumeTechnology{
		ID: x.ID().UUID(), ResumeID: resumeID.UUID(), Name: x.Name(), Position: x.Position(),
	}
	if g := x.Grade(); g != nil {
		if !g.IsValid() {
			return gormmodel.ResumeTechnology{}, fmt.Errorf("invalid technology grade")
		}
		value := string(*g)
		m.Grade = &value
	}
	return m, nil
}

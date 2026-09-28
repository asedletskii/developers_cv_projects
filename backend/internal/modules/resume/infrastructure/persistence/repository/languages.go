package repository

import (
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

func persistLanguages(tx *gorm.DB, resume *domain.Resume) error {
	languages := resume.Languages()
	if languages == nil {
		if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeLanguage{}).Error; err != nil {
			return fmt.Errorf("delete resume languages: %w", mapDBError(err))
		}
		return nil
	}

	entries := languages.Entries()
	models := make([]gormmodel.ResumeLanguage, 0, len(entries))
	seenIDs := make(map[valueobjects.ID]struct{}, len(entries))
	seenLanguages := make(map[enums.SpokenLanguage]struct{}, len(entries))
	for _, entry := range entries {
		id := entry.ID()
		if id.IsZero() {
			return fmt.Errorf("%w: language id is required", domain.ErrResumeInvariant)
		}
		if _, ok := seenIDs[id]; ok {
			return fmt.Errorf("%w: duplicate language id", domain.ErrInvalidEntity)
		}
		seenIDs[id] = struct{}{}
		language := entry.Language()
		if _, ok := seenLanguages[language]; ok {
			return fmt.Errorf("%w: duplicate language", domain.ErrInvalidEntity)
		}
		seenLanguages[language] = struct{}{}
		model, err := languageModel(entry, resume.ID())
		if err != nil {
			return fmt.Errorf("map resume language: %w", err)
		}
		models = append(models, model)
	}

	if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeLanguage{}).Error; err != nil {
		return fmt.Errorf("replace resume languages: %w", mapDBError(err))
	}
	if len(models) > 0 {
		if err := tx.Create(&models).Error; err != nil {
			return fmt.Errorf("create resume languages: %w", mapDBError(err))
		}
	}
	return nil
}

func loadLanguages(tx *gorm.DB, resumeID valueobjects.ID) (*domain.ResumeLanguages, error) {
	var rows []gormmodel.ResumeLanguage
	if err := tx.Where("resume_id = ?", resumeID.UUID()).Order("position").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("get resume languages: %w", mapDBError(err))
	}
	if len(rows) == 0 {
		return nil, nil
	}
	collectionID, err := valueobjects.IDFromUUID(rows[0].ID)
	if err != nil {
		return nil, err
	}
	entries := make([]domain.ResumeLanguageEntryParams, 0, len(rows))
	for _, row := range rows {
		id, err := valueobjects.IDFromUUID(row.ID)
		if err != nil {
			return nil, err
		}
		language := enums.SpokenLanguage(row.Language)
		if !language.IsValid() {
			return nil, domain.ErrInvalidEntity
		}
		var grade *enums.LanguageGrade
		if row.Grade != nil {
			value := enums.LanguageGrade(*row.Grade)
			if !value.IsValid() {
				return nil, domain.ErrInvalidEntity
			}
			grade = &value
		}
		entries = append(entries, domain.ResumeLanguageEntryParams{ID: id, Language: language, Grade: grade, Position: row.Position})
	}
	return domain.NewResumeLanguages(domain.ResumeLanguagesParams{ID: collectionID, ResumeID: resumeID, Entries: entries})
}

func languageModel(x domain.ResumeLanguageEntry, resumeID valueobjects.ID) (gormmodel.ResumeLanguage, error) {
	if x.ID().IsZero() {
		return gormmodel.ResumeLanguage{}, fmt.Errorf("language id is required")
	}
	m := gormmodel.ResumeLanguage{ID: x.ID().UUID(), ResumeID: resumeID.UUID(), Language: string(x.Language()), Position: x.Position()}
	if g := x.Grade(); g != nil {
		value := string(*g)
		m.Grade = &value
	}
	if !enums.SpokenLanguage(m.Language).IsValid() {
		return gormmodel.ResumeLanguage{}, fmt.Errorf("invalid language")
	}
	return m, nil
}

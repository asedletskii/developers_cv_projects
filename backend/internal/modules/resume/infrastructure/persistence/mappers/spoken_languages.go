package mappers

import (
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func LanguagesToDomain(rows []gormmodel.ResumeSpokenLanguage) (*domain.ResumeSpokenLanguages, error) {
	if len(rows) == 0 {
		return nil, nil
	}
	entries := make([]domain.SpokenLanguageEntryParams, 0, len(rows))
	for i, row := range rows {
		id, err := sharedvalueobjects.IDFromUUID(row.ID)
		if err != nil {
			return nil, fmt.Errorf("map spoken language at index %d: %w", i, err)
		}
		language := domain.SpokenLanguage(row.Language)
		var grade *domain.LanguageGrade
		if row.Grade != nil {
			value := domain.LanguageGrade(*row.Grade)
			grade = &value
		}
		entries = append(entries, domain.SpokenLanguageEntryParams{
			ID:       id,
			Language: language,
			Grade:    grade,
			Position: row.Position,
		})
	}
	return domain.RestoreResumeSpokenLanguages(domain.ResumeSpokenLanguagesParams{
		Entries: entries,
	}), nil
}

func LanguageToDB(
	x domain.SpokenLanguageEntry,
	resumeID sharedvalueobjects.ID,
) gormmodel.ResumeSpokenLanguage {
	m := gormmodel.ResumeSpokenLanguage{
		ID:       x.ID().UUID(),
		ResumeID: resumeID.UUID(),
		Language: string(x.Language()),
		Position: x.Position(),
	}
	if g := x.Grade(); g != nil {
		m.Grade = stringPtr(string(*g))
	}
	return m
}

func stringPtr(value string) *string { return &value }

func ResumeSpokenLanguagesToDBModels(
	collection *domain.ResumeSpokenLanguages,
	resumeID sharedvalueobjects.ID,
) []gormmodel.ResumeSpokenLanguage {
	if collection == nil {
		return nil
	}
	entries := collection.Entries()
	models := make([]gormmodel.ResumeSpokenLanguage, 0, len(entries))
	for _, entry := range entries {
		models = append(models, LanguageToDB(entry, resumeID))
	}
	return models
}

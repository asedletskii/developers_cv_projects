package mappers

import (
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func LanguagesToDomain(
	rows []gormmodel.ResumeLanguage,
	resumeID identity.ID,
) (*domain.ResumeLanguages, error) {
	if len(rows) == 0 {
		return nil, nil
	}
	collectionID, err := identity.IDFromUUID(rows[0].ID)
	if err != nil {
		return nil, err
	}
	entries := make([]domain.ResumeLanguageEntryParams, 0, len(rows))
	for _, row := range rows {
		id, err := identity.IDFromUUID(row.ID)
		if err != nil {
			return nil, err
		}
		language := domain.SpokenLanguage(row.Language)
		if !language.IsValid() {
			return nil, domain.ErrInvalidEntity
		}
		var grade *domain.LanguageGrade
		if row.Grade != nil {
			value := domain.LanguageGrade(*row.Grade)
			if !value.IsValid() {
				return nil, domain.ErrInvalidEntity
			}
			grade = &value
		}
		entries = append(entries, domain.ResumeLanguageEntryParams{
			ID:       id,
			Language: language,
			Grade:    grade,
			Position: row.Position,
		})
	}
	return domain.NewResumeLanguages(domain.ResumeLanguagesParams{
		ID:       collectionID,
		ResumeID: resumeID,
		Entries:  entries,
	})
}

func LanguageToDB(
	x domain.ResumeLanguageEntry,
	resumeID identity.ID,
) (gormmodel.ResumeLanguage, error) {
	if x.ID().IsZero() {
		return gormmodel.ResumeLanguage{}, fmt.Errorf("language id is required")
	}
	m := gormmodel.ResumeLanguage{
		ID:       x.ID().UUID(),
		ResumeID: resumeID.UUID(),
		Language: string(x.Language()),
		Position: x.Position(),
	}
	if g := x.Grade(); g != nil {
		m.Grade = stringPtr(string(*g))
	}
	if !domain.SpokenLanguage(m.Language).IsValid() {
		return gormmodel.ResumeLanguage{}, fmt.Errorf("invalid language")
	}
	return m, nil
}

func stringPtr(value string) *string { return &value }

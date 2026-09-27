package domain_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

func TestResumeLanguagesRejectsDuplicateLanguages(t *testing.T) {
	language := enums.SpokenLanguage("en")
	_, err := domain.NewResumeLanguages(domain.ResumeLanguagesParams{
		ID:       valueobjects.NewID(),
		ResumeID: valueobjects.NewID(),
		Entries: []domain.ResumeLanguageEntryParams{
			{ID: valueobjects.NewID(), Language: language, Position: 0},
			{ID: valueobjects.NewID(), Language: language, Position: 1},
		},
	})
	require.ErrorIs(t, err, domain.ErrInvalidEntity)
}

func TestResumeLanguageEntryRejectsInvalidPosition(t *testing.T) {
	_, err := domain.NewResumeLanguageEntry(domain.ResumeLanguageEntryParams{
		ID:       valueobjects.NewID(),
		Language: enums.SpokenLanguage("en"),
		Position: -1,
	})
	require.ErrorIs(t, err, domain.ErrInvalidEntity)
}

func TestResumeLanguagesRejectsDuplicateIDs(t *testing.T) {
	id := valueobjects.NewID()
	_, err := domain.NewResumeLanguages(domain.ResumeLanguagesParams{
		ID:       valueobjects.NewID(),
		ResumeID: valueobjects.NewID(),
		Entries: []domain.ResumeLanguageEntryParams{
			{ID: id, Language: enums.SpokenLanguageEN, Position: 0},
			{ID: id, Language: enums.SpokenLanguageFR, Position: 1},
		},
	})
	require.ErrorIs(t, err, domain.ErrInvalidEntity)
}

func TestResumeLanguagesRejectsDuplicatePositions(t *testing.T) {
	_, err := domain.NewResumeLanguages(domain.ResumeLanguagesParams{
		ID:       valueobjects.NewID(),
		ResumeID: valueobjects.NewID(),
		Entries: []domain.ResumeLanguageEntryParams{
			{ID: valueobjects.NewID(), Language: enums.SpokenLanguageEN, Position: 0},
			{ID: valueobjects.NewID(), Language: enums.SpokenLanguageFR, Position: 0},
		},
	})
	require.ErrorIs(t, err, domain.ErrInvalidEntity)
}

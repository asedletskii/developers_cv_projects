package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type LanguageInput struct {
	ID       sharedvalueobjects.ID
	Language domain.SpokenLanguage
	Grade    *domain.LanguageGrade
	Position int
}

type LanguagesInput struct {
	Entries []LanguageInput
}

func (i LanguagesInput) toDomain() domain.ResumeSpokenLanguagesParams {
	entries := make([]domain.SpokenLanguageEntryParams, 0, len(i.Entries))
	for _, entry := range i.Entries {
		entries = append(entries, domain.SpokenLanguageEntryParams{
			ID: entry.ID, Language: entry.Language, Grade: entry.Grade, Position: entry.Position,
		})
	}
	return domain.ResumeSpokenLanguagesParams{Entries: entries}
}

type ReplaceLanguagesCommand struct {
	ResumeID sharedvalueobjects.ID
	UserID   sharedvalueobjects.ID
	Version  int64
	Input    LanguagesInput
}

type ReplaceLanguagesHandler struct {
	repository ports.ResumeRepository
	uow        ports.UnitOfWork
}

func NewReplaceLanguagesHandler(repository ports.ResumeRepository, uow ports.UnitOfWork) *ReplaceLanguagesHandler {
	return &ReplaceLanguagesHandler{repository: repository, uow: uow}
}

func (h *ReplaceLanguagesHandler) Handle(ctx context.Context, c ReplaceLanguagesCommand) error {
	if h == nil {
		return fmt.Errorf("replace languages handler dependencies are nil")
	}

	languages := c.Input.toDomain()
	return executeResumeMutation(ctx, h.repository, h.uow, c.ResumeID, c.UserID, c.Version, "replace languages", ports.ResumeMutationLanguages, func(resume *domain.Resume) error {
		return resume.ReplaceSpokenLanguages(languages)
	})
}

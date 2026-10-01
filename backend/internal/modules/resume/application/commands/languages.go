package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type LanguageInput struct {
	ID       identity.ID
	Language domain.SpokenLanguage
	Grade    *domain.LanguageGrade
	Position int
}

type LanguagesInput struct {
	Entries []LanguageInput
}

func (i LanguagesInput) toDomain() domain.ResumeLanguagesParams {
	entries := make([]domain.ResumeLanguageEntryParams, 0, len(i.Entries))
	for _, entry := range i.Entries {
		entries = append(entries, domain.ResumeLanguageEntryParams{
			ID: entry.ID, Language: entry.Language, Grade: entry.Grade, Position: entry.Position,
		})
	}
	return domain.ResumeLanguagesParams{Entries: entries}
}

type ReplaceLanguagesCommand struct {
	ResumeID identity.ID
	UserID   identity.ID
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

	return executeResumeMutation(ctx, h.repository, h.uow, c.ResumeID, c.UserID, c.Version, "replace languages", func(resume *domain.Resume) error {
		return resume.ReplaceLanguages(c.Input.toDomain())
	})
}

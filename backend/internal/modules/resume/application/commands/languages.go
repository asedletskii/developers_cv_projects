package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type LanguageInput struct {
	ID       valueobjects.ID
	Language enums.SpokenLanguage
	Grade    *enums.LanguageGrade
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
	ResumeID valueobjects.ID
	UserID   useridentity.ID
	Version  int64
	Input    LanguagesInput
}

type ReplaceLanguagesHandler struct {
	repository ports.ResumeRepository
	uow        UnitOfWork
}

func NewReplaceLanguagesHandler(repository ports.ResumeRepository, uow UnitOfWork) *ReplaceLanguagesHandler {
	return &ReplaceLanguagesHandler{repository: repository, uow: uow}
}

func (h *ReplaceLanguagesHandler) Handle(ctx context.Context, c ReplaceLanguagesCommand) error {
	if h == nil || h.repository == nil || h.uow == nil {
		return fmt.Errorf("replace languages handler dependencies are nil")
	}
	if c.ResumeID.IsZero() || c.UserID.IsZero() || c.Version < 1 {
		return fmt.Errorf("invalid replace languages command")
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		resume, err := h.repository.Get(tx, ports.GetResumeParams{ResumeID: c.ResumeID, UserID: c.UserID})
		if err != nil {
			return err
		}
		if resume.Version() != c.Version {
			return domain.ErrConcurrencyConflict
		}
		if err := resume.ReplaceLanguages(c.Input.toDomain()); err != nil {
			return err
		}
		return h.repository.Save(tx, resume)
	})
}

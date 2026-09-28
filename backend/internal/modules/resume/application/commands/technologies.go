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

type TechnologyInput struct {
	ID       valueobjects.ID
	Name     string
	Grade    *enums.TechnologyGrade
	Position int
}

type TechnologiesInput struct {
	Entries []TechnologyInput
}

func (i TechnologiesInput) toDomain() domain.ResumeTechnologiesParams {
	entries := make([]domain.ResumeTechnologyEntryParams, 0, len(i.Entries))
	for _, entry := range i.Entries {
		entries = append(entries, domain.ResumeTechnologyEntryParams{
			ID: entry.ID, Name: entry.Name, Grade: entry.Grade, Position: entry.Position,
		})
	}
	return domain.ResumeTechnologiesParams{Entries: entries}
}

type ReplaceTechnologiesCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
	Version  int64
	Input    TechnologiesInput
}

type ReplaceTechnologiesHandler struct {
	repository ports.ResumeRepository
	uow        UnitOfWork
}

func NewReplaceTechnologiesHandler(repository ports.ResumeRepository, uow UnitOfWork) *ReplaceTechnologiesHandler {
	return &ReplaceTechnologiesHandler{repository: repository, uow: uow}
}

func (h *ReplaceTechnologiesHandler) Handle(ctx context.Context, c ReplaceTechnologiesCommand) error {
	if h == nil || h.repository == nil || h.uow == nil {
		return fmt.Errorf("replace technologies handler dependencies are nil")
	}
	if c.ResumeID.IsZero() || c.UserID.IsZero() || c.Version < 1 {
		return fmt.Errorf("invalid replace technologies command")
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		resume, err := h.repository.Get(tx, ports.GetResumeParams{ResumeID: c.ResumeID, UserID: c.UserID})
		if err != nil {
			return err
		}
		if resume.Version() != c.Version {
			return domain.ErrConcurrencyConflict
		}
		if err := resume.ReplaceTechnologies(c.Input.toDomain()); err != nil {
			return err
		}
		return h.repository.Save(tx, resume)
	})
}

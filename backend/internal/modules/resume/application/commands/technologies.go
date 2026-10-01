package commands

import (
	"context"
	"developers-cv/backend/internal/modules/resume/domain"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/ports"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type TechnologyInput struct {
	ID       identity.ID
	Name     string
	Grade    *domain.TechnologyGrade
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
	ResumeID identity.ID
	UserID   identity.ID
	Version  int64
	Input    TechnologiesInput
}

type ReplaceTechnologiesHandler struct {
	repository ports.ResumeRepository
	uow        ports.UnitOfWork
}

func NewReplaceTechnologiesHandler(repository ports.ResumeRepository, uow ports.UnitOfWork) *ReplaceTechnologiesHandler {
	return &ReplaceTechnologiesHandler{repository: repository, uow: uow}
}

func (h *ReplaceTechnologiesHandler) Handle(ctx context.Context, c ReplaceTechnologiesCommand) error {
	if h == nil {
		return fmt.Errorf("replace technologies handler dependencies are nil")
	}

	return executeResumeMutation(ctx, h.repository, h.uow, c.ResumeID, c.UserID, c.Version, "replace technologies", func(resume *domain.Resume) error {
		return resume.ReplaceTechnologies(c.Input.toDomain())
	})
}

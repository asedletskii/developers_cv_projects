package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ExperienceStackItemInput struct {
	Name string
}

type ExperienceInput struct {
	ID                sharedvalueobjects.ID
	Company, Position string
	Location          *string
	StartDate         domain.YearMonth
	EndDate           *domain.YearMonth
	Description       *string
	Achievements      []string
	Tasks             []string
	Stack             []ExperienceStackItemInput
}

type ExperiencesInput struct{ Items []ExperienceInput }

func (i ExperienceInput) toDomain() domain.ExperienceParams {
	stack := make([]domain.ExperienceStackParams, 0, len(i.Stack))
	for _, item := range i.Stack {
		stack = append(stack, domain.ExperienceStackParams{Name: item.Name})
	}
	return domain.ExperienceParams{
		ID:           i.ID,
		Company:      i.Company,
		Position:     i.Position,
		Location:     i.Location,
		StartDate:    i.StartDate,
		EndDate:      i.EndDate,
		Description:  i.Description,
		Achievements: append([]string(nil), i.Achievements...),
		Tasks:        append([]string(nil), i.Tasks...),
		Stack:        stack,
	}
}

func (i ExperiencesInput) toDomain() []domain.ExperienceParams {
	experiences := make([]domain.ExperienceParams, 0, len(i.Items))
	for _, input := range i.Items {
		experiences = append(experiences, input.toDomain())
	}
	return experiences
}

type ReplaceExperiencesCommand struct {
	ResumeID sharedvalueobjects.ID
	UserID   sharedvalueobjects.ID
	Version  int64
	Input    ExperiencesInput
}

type ReplaceExperiencesHandler struct {
	repository ports.ResumeRepository
	uow        ports.UnitOfWork
}

func NewReplaceExperiencesHandler(repository ports.ResumeRepository, uow ports.UnitOfWork) *ReplaceExperiencesHandler {
	return &ReplaceExperiencesHandler{repository: repository, uow: uow}
}

func (h *ReplaceExperiencesHandler) Handle(ctx context.Context, c ReplaceExperiencesCommand) error {
	if h == nil {
		return fmt.Errorf("replace experiences handler dependencies are nil")
	}

	experiences := c.Input.toDomain()
	return executeResumeMutation(ctx, h.repository, h.uow, c.ResumeID, c.UserID, c.Version, "replace experiences", ports.ResumeMutationExperiences, func(resume *domain.Resume) error {
		return resume.ReplaceExperiences(experiences)
	})
}

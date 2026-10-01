package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ExperienceStackItemInput struct {
	Name  string
	Grade *domain.TechnologyGrade
}

type ExperienceInput struct {
	ID                identity.ID
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

func (i ExperienceInput) toDomain(resumeID identity.ID) domain.ExperienceParams {
	stack := make([]domain.ExperienceStackParams, 0, len(i.Stack))
	for _, item := range i.Stack {
		stack = append(stack, domain.ExperienceStackParams{Name: item.Name, Grade: item.Grade})
	}
	return domain.ExperienceParams{ID: i.ID, ResumeID: resumeID, Company: i.Company, Position: i.Position, Location: i.Location, StartDate: i.StartDate, EndDate: i.EndDate, Description: i.Description, Achievements: append([]string(nil), i.Achievements...), Tasks: append([]string(nil), i.Tasks...), Stack: stack}
}

type ReplaceExperiencesCommand struct {
	ResumeID identity.ID
	UserID   identity.ID
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

	return executeResumeMutation(ctx, h.repository, h.uow, c.ResumeID, c.UserID, c.Version, "replace experiences", func(resume *domain.Resume) error {
		params := make([]domain.ExperienceParams, 0, len(c.Input.Items))
		for _, input := range c.Input.Items {
			params = append(params, input.toDomain(c.ResumeID))
		}
		return resume.ReplaceExperiences(params)
	})
}

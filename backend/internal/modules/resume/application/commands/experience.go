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

type ExperienceStackItemInput struct {
	Name  string
	Grade *enums.TechnologyGrade
}

type ExperienceInput struct {
	ID                valueobjects.ID
	Company, Position string
	Location          *string
	StartDate         valueobjects.YearMonth
	EndDate           *valueobjects.YearMonth
	Description       *string
	Achievements      []string
	Tasks             []string
	Stack             []ExperienceStackItemInput
}

type ExperiencesInput struct{ Items []ExperienceInput }

func (i ExperienceInput) toDomain(resumeID valueobjects.ID) domain.ExperienceParams {
	stack := make([]domain.ExperienceStackParams, 0, len(i.Stack))
	for _, item := range i.Stack {
		stack = append(stack, domain.ExperienceStackParams{Name: item.Name, Grade: item.Grade})
	}
	return domain.ExperienceParams{ID: i.ID, ResumeID: resumeID, Company: i.Company, Position: i.Position, Location: i.Location, StartDate: i.StartDate, EndDate: i.EndDate, Description: i.Description, Achievements: append([]string(nil), i.Achievements...), Tasks: append([]string(nil), i.Tasks...), Stack: stack}
}

type ReplaceExperiencesCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
	Version  int64
	Input    ExperiencesInput
}

type ReplaceExperiencesHandler struct {
	repository ports.ResumeRepository
	uow        UnitOfWork
}

func NewReplaceExperiencesHandler(repository ports.ResumeRepository, uow UnitOfWork) *ReplaceExperiencesHandler {
	return &ReplaceExperiencesHandler{repository: repository, uow: uow}
}

func (h *ReplaceExperiencesHandler) Handle(ctx context.Context, c ReplaceExperiencesCommand) error {
	if h == nil || h.repository == nil || h.uow == nil {
		return fmt.Errorf("replace experiences handler dependencies are nil")
	}
	if c.ResumeID.IsZero() || c.UserID.IsZero() || c.Version < 1 {
		return fmt.Errorf("invalid replace experiences command")
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		resume, err := h.repository.Get(tx, ports.GetResumeParams{ResumeID: c.ResumeID, UserID: c.UserID})
		if err != nil {
			return err
		}
		if resume.Version() != c.Version {
			return domain.ErrConcurrencyConflict
		}
		params := make([]domain.ExperienceParams, 0, len(c.Input.Items))
		for _, input := range c.Input.Items {
			params = append(params, input.toDomain(c.ResumeID))
		}
		if err := resume.ReplaceExperiences(params); err != nil {
			return err
		}
		return h.repository.Save(tx, resume)
	})
}

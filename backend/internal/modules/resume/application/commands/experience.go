package commands

import (
	"context"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
	"fmt"
)

type ExperienceTechnologyInput struct {
	Name  string
	Grade *enums.TechnologyGrade
}
type ExperienceInput struct {
	ID, ResumeID      valueobjects.ID
	Company, Position string
	Location          *string
	StartDate         valueobjects.YearMonth
	EndDate           *valueobjects.YearMonth
	Description       *string
	Achievements      []string
	Technologies      []ExperienceTechnologyInput
}
type CreateExperienceCommand struct {
	ID, ResumeID valueobjects.ID
	UserID       useridentity.ID
	Input        ExperienceInput
}
type GetExperienceCommand struct {
	ResumeID, ExperienceID valueobjects.ID
	UserID                 useridentity.ID
}
type UpdateExperienceCommand struct {
	ResumeID, ExperienceID valueobjects.ID
	UserID                 useridentity.ID
	Input                  ExperienceInput
}
type DeleteExperienceCommand struct {
	ResumeID, ExperienceID valueobjects.ID
	UserID                 useridentity.ID
}
type CreateExperienceRepositoryCommand struct {
	Experience *domain.Experience
	ResumeID   valueobjects.ID
	UserID     useridentity.ID
}
type UpdateExperienceRepositoryCommand struct {
	Experience *domain.Experience
	ResumeID   valueobjects.ID
	UserID     useridentity.ID
}
type CreateExperienceRepository interface {
	CreateExperience(context.Context, CreateExperienceRepositoryCommand) error
}
type GetExperienceRepository interface {
	GetExperience(context.Context, GetExperienceCommand) (*domain.Experience, error)
}
type UpdateExperienceRepository interface {
	UpdateExperience(context.Context, UpdateExperienceRepositoryCommand) error
}
type DeleteExperienceRepository interface {
	DeleteExperience(context.Context, DeleteExperienceCommand) error
}
type CreateExperienceHandler struct {
	r CreateExperienceRepository
	u UnitOfWork
}

func NewCreateExperienceHandler(r CreateExperienceRepository, u UnitOfWork) *CreateExperienceHandler {
	return &CreateExperienceHandler{r: r, u: u}
}
func (h *CreateExperienceHandler) Handle(ctx context.Context, c CreateExperienceCommand) (valueobjects.ID, error) {
	if h == nil || h.r == nil || h.u == nil {
		return valueobjects.ID{}, fmt.Errorf("create experience dependencies are nil")
	}
	if c.ID.IsZero() {
		c.ID = valueobjects.NewID()
	}
	c.Input.ID, c.Input.ResumeID = c.ID, c.ResumeID
	a, err := experienceFromInput(c.Input)
	if err != nil {
		return valueobjects.ID{}, err
	}
	err = h.u.Do(ctx, func(tx context.Context) error {
		return h.r.CreateExperience(tx, CreateExperienceRepositoryCommand{Experience: a, ResumeID: c.ResumeID, UserID: c.UserID})
	})
	if err != nil {
		return valueobjects.ID{}, err
	}
	return a.ID(), nil
}

type UpdateExperienceHandler struct {
	r   GetExperienceRepository
	u   UpdateExperienceRepository
	uow UnitOfWork
}

func NewUpdateExperienceHandler(r GetExperienceRepository, u UpdateExperienceRepository, uow UnitOfWork) *UpdateExperienceHandler {
	return &UpdateExperienceHandler{r: r, u: u, uow: uow}
}
func (h *UpdateExperienceHandler) Handle(ctx context.Context, c UpdateExperienceCommand) error {
	if h == nil || h.r == nil || h.u == nil || h.uow == nil {
		return fmt.Errorf("update experience dependencies are nil")
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		current, err := h.r.GetExperience(tx, GetExperienceCommand{ResumeID: c.ResumeID, ExperienceID: c.ExperienceID, UserID: c.UserID})
		if err != nil {
			return err
		}
		c.Input.ID, c.Input.ResumeID = current.ID(), c.ResumeID
		updated, err := experienceFromInput(c.Input)
		if err != nil {
			return err
		}
		return h.u.UpdateExperience(tx, UpdateExperienceRepositoryCommand{Experience: updated, ResumeID: c.ResumeID, UserID: c.UserID})
	})
}

type DeleteExperienceHandler struct {
	r DeleteExperienceRepository
	u UnitOfWork
}

func NewDeleteExperienceHandler(r DeleteExperienceRepository, u UnitOfWork) *DeleteExperienceHandler {
	return &DeleteExperienceHandler{r: r, u: u}
}
func (h *DeleteExperienceHandler) Handle(ctx context.Context, c DeleteExperienceCommand) error {
	if h == nil || h.r == nil || h.u == nil {
		return fmt.Errorf("delete experience dependencies are nil")
	}
	return h.u.Do(ctx, func(tx context.Context) error { return h.r.DeleteExperience(tx, c) })
}
func experienceFromInput(input ExperienceInput) (*domain.Experience, error) {
	ts := make([]domain.ExperienceTechnologyParams, len(input.Technologies))
	for i, t := range input.Technologies {
		ts[i] = domain.ExperienceTechnologyParams{Name: t.Name, Grade: t.Grade}
	}
	return domain.NewExperience(domain.ExperienceParams{ID: input.ID, ResumeID: input.ResumeID, Company: input.Company, Position: input.Position, Location: input.Location, StartDate: input.StartDate, EndDate: input.EndDate, Description: input.Description, Achievements: input.Achievements, Technologies: ts})
}

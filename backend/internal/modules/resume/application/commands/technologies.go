package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type TechnologyInput struct {
	ID       valueobjects.ID
	Name     string
	Grade    *enums.TechnologyGrade
	Position int
}

type CreateTechnologyCommand struct {
	ID, ResumeID valueobjects.ID
	UserID       useridentity.ID
	Input        TechnologyInput
}

type GetTechnologyCommand struct {
	ResumeID, TechnologyID valueobjects.ID
	UserID                 useridentity.ID
}

type UpdateTechnologyCommand struct {
	ResumeID, TechnologyID valueobjects.ID
	UserID                 useridentity.ID
	Input                  TechnologyInput
}

type DeleteTechnologyCommand struct {
	ResumeID, TechnologyID valueobjects.ID
	UserID                 useridentity.ID
}

type CreateTechnologyRepositoryCommand struct {
	Technology domain.ResumeTechnologyEntry
	ResumeID   valueobjects.ID
	UserID     useridentity.ID
}

type UpdateTechnologyRepositoryCommand struct {
	Technology domain.ResumeTechnologyEntry
	ResumeID   valueobjects.ID
	UserID     useridentity.ID
}

type CreateTechnologyRepository interface {
	CreateTechnology(context.Context, CreateTechnologyRepositoryCommand) error
}

type GetTechnologyRepository interface {
	GetTechnology(context.Context, GetTechnologyCommand) (domain.ResumeTechnologyEntry, error)
}

type UpdateTechnologyRepository interface {
	UpdateTechnology(context.Context, UpdateTechnologyRepositoryCommand) error
}

type DeleteTechnologyRepository interface {
	DeleteTechnology(context.Context, DeleteTechnologyCommand) error
}

type CreateTechnologyHandler struct {
	r CreateTechnologyRepository
	u UnitOfWork
}

func NewCreateTechnologyHandler(r CreateTechnologyRepository, u UnitOfWork) *CreateTechnologyHandler {
	return &CreateTechnologyHandler{r: r, u: u}
}

func (h *CreateTechnologyHandler) Handle(ctx context.Context, c CreateTechnologyCommand) (valueobjects.ID, error) {
	if h == nil || h.r == nil || h.u == nil {
		return valueobjects.ID{}, fmt.Errorf("create technology dependencies are nil")
	}
	if c.ID.IsZero() {
		c.ID = valueobjects.NewID()
	}
	c.Input.ID = c.ID
	a, err := domain.NewResumeTechnologyEntry(domain.ResumeTechnologyEntryParams{ID: c.Input.ID, Name: c.Input.Name, Grade: c.Input.Grade, Position: c.Input.Position})
	if err != nil {
		return valueobjects.ID{}, err
	}
	err = h.u.Do(ctx, func(tx context.Context) error {
		return h.r.CreateTechnology(tx, CreateTechnologyRepositoryCommand{Technology: a, ResumeID: c.ResumeID, UserID: c.UserID})
	})
	if err != nil {
		return valueobjects.ID{}, err
	}
	return a.ID(), nil
}

type UpdateTechnologyHandler struct {
	r   GetTechnologyRepository
	u   UpdateTechnologyRepository
	uow UnitOfWork
}

func NewUpdateTechnologyHandler(r GetTechnologyRepository, u UpdateTechnologyRepository, uow UnitOfWork) *UpdateTechnologyHandler {
	return &UpdateTechnologyHandler{r: r, u: u, uow: uow}
}

func (h *UpdateTechnologyHandler) Handle(ctx context.Context, c UpdateTechnologyCommand) error {
	if h == nil || h.r == nil || h.u == nil || h.uow == nil {
		return fmt.Errorf("update technology dependencies are nil")
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		current, err := h.r.GetTechnology(tx, GetTechnologyCommand{ResumeID: c.ResumeID, TechnologyID: c.TechnologyID, UserID: c.UserID})
		if err != nil {
			return err
		}
		if err := current.Replace(domain.ResumeTechnologyEntryParams{ID: current.ID(), Name: c.Input.Name, Grade: c.Input.Grade, Position: c.Input.Position}); err != nil {
			return err
		}
		return h.u.UpdateTechnology(tx, UpdateTechnologyRepositoryCommand{Technology: current, ResumeID: c.ResumeID, UserID: c.UserID})
	})
}

type DeleteTechnologyHandler struct {
	r DeleteTechnologyRepository
	u UnitOfWork
}

func NewDeleteTechnologyHandler(r DeleteTechnologyRepository, u UnitOfWork) *DeleteTechnologyHandler {
	return &DeleteTechnologyHandler{r: r, u: u}
}

func (h *DeleteTechnologyHandler) Handle(ctx context.Context, c DeleteTechnologyCommand) error {
	if h == nil || h.r == nil || h.u == nil {
		return fmt.Errorf("delete technology dependencies are nil")
	}
	return h.u.Do(ctx, func(tx context.Context) error { return h.r.DeleteTechnology(tx, c) })
}

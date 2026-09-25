package commands

import (
	"context"
	"developers-cv/backend/internal/modules/resume/application/ports"
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

type UpdateTechnologyCommand struct {
	ResumeID, TechnologyID valueobjects.ID
	UserID                 useridentity.ID
	Input                  TechnologyInput
}

type DeleteTechnologyCommand struct {
	ResumeID, TechnologyID valueobjects.ID
	UserID                 useridentity.ID
}

type CreateTechnologyHandler struct {
	r ports.CreateTechnologyRepository
	u UnitOfWork
}

func NewCreateTechnologyHandler(r ports.CreateTechnologyRepository, u UnitOfWork) *CreateTechnologyHandler {
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
		return h.r.CreateTechnology(tx, ports.CreateTechnologyRepositoryCommand{Technology: a, ResumeID: c.ResumeID, UserID: c.UserID})
	})
	if err != nil {
		return valueobjects.ID{}, err
	}
	return a.ID(), nil
}

type UpdateTechnologyHandler struct {
	u   ports.UpdateTechnologyRepository
	uow UnitOfWork
}

func NewUpdateTechnologyHandler(u ports.UpdateTechnologyRepository, uow UnitOfWork) *UpdateTechnologyHandler {
	return &UpdateTechnologyHandler{u: u, uow: uow}
}

func (h *UpdateTechnologyHandler) Handle(ctx context.Context, c UpdateTechnologyCommand) error {
	if h == nil || h.u == nil || h.uow == nil {
		return fmt.Errorf("update technology dependencies are nil")
	}
	entity, err := domain.NewResumeTechnologyEntry(domain.ResumeTechnologyEntryParams{ID: c.TechnologyID, Name: c.Input.Name, Grade: c.Input.Grade, Position: c.Input.Position})
	if err != nil {
		return err
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		return h.u.UpdateTechnology(tx, ports.UpdateTechnologyRepositoryCommand{Technology: entity, ResumeID: c.ResumeID, UserID: c.UserID})
	})
}

type DeleteTechnologyHandler struct {
	r ports.DeleteTechnologyRepository
	u UnitOfWork
}

func NewDeleteTechnologyHandler(r ports.DeleteTechnologyRepository, u UnitOfWork) *DeleteTechnologyHandler {
	return &DeleteTechnologyHandler{r: r, u: u}
}

func (h *DeleteTechnologyHandler) Handle(ctx context.Context, c DeleteTechnologyCommand) error {
	if h == nil || h.r == nil || h.u == nil {
		return fmt.Errorf("delete technology dependencies are nil")
	}
	return h.u.Do(ctx, func(tx context.Context) error {
		return h.r.DeleteTechnology(tx, ports.DeleteTechnologyRepositoryCommand{ResumeID: c.ResumeID, TechnologyID: c.TechnologyID, UserID: c.UserID})
	})
}

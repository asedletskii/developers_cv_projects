package commands

import (
	"context"

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
type ListTechnologiesQuery struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type GetTechnologyQuery struct {
	ResumeID, TechnologyID valueobjects.ID
	UserID                 useridentity.ID
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
type TechnologiesRepository interface {
	CreateTechnology(context.Context, TechnologyInput, valueobjects.ID, useridentity.ID) error
	ListTechnologies(context.Context, ListTechnologiesQuery) ([]domain.ResumeTechnologyEntry, error)
	GetTechnology(context.Context, GetTechnologyQuery) (domain.ResumeTechnologyEntry, error)
	UpdateTechnology(context.Context, UpdateTechnologyCommand) error
	DeleteTechnology(context.Context, DeleteTechnologyCommand) error
}
type CreateTechnologyHandler struct {
	r TechnologiesRepository
	u UnitOfWork
}

func NewCreateTechnologyHandler(r TechnologiesRepository, u UnitOfWork) *CreateTechnologyHandler {
	return &CreateTechnologyHandler{r, u}
}
func (h *CreateTechnologyHandler) Handle(ctx context.Context, c CreateTechnologyCommand) (valueobjects.ID, error) {
	if c.ID.IsZero() {
		c.ID = valueobjects.NewID()
	}
	c.Input.ID = c.ID
	if _, err := domain.NewResumeTechnologyEntry(domain.ResumeTechnologyEntryParams{
		ID: c.Input.ID, Name: c.Input.Name, Grade: c.Input.Grade, Position: c.Input.Position,
	}); err != nil {
		return valueobjects.ID{}, err
	}
	return c.ID, h.u.Do(ctx, func(tx context.Context) error { return h.r.CreateTechnology(tx, c.Input, c.ResumeID, c.UserID) })
}

type ListTechnologiesHandler struct{ r TechnologiesRepository }

func NewListTechnologiesHandler(r TechnologiesRepository) *ListTechnologiesHandler {
	return &ListTechnologiesHandler{r}
}
func (h *ListTechnologiesHandler) Handle(ctx context.Context, q ListTechnologiesQuery) ([]domain.ResumeTechnologyEntry, error) {
	return h.r.ListTechnologies(ctx, q)
}

type GetTechnologyHandler struct{ r TechnologiesRepository }

func NewGetTechnologyHandler(r TechnologiesRepository) *GetTechnologyHandler {
	return &GetTechnologyHandler{r}
}
func (h *GetTechnologyHandler) Handle(ctx context.Context, q GetTechnologyQuery) (domain.ResumeTechnologyEntry, error) {
	return h.r.GetTechnology(ctx, q)
}

type UpdateTechnologyHandler struct {
	r TechnologiesRepository
	u UnitOfWork
}

func NewUpdateTechnologyHandler(r TechnologiesRepository, u UnitOfWork) *UpdateTechnologyHandler {
	return &UpdateTechnologyHandler{r, u}
}
func (h *UpdateTechnologyHandler) Handle(ctx context.Context, c UpdateTechnologyCommand) error {
	c.Input.ID = c.TechnologyID
	if _, err := domain.NewResumeTechnologyEntry(domain.ResumeTechnologyEntryParams{
		ID: c.Input.ID, Name: c.Input.Name, Grade: c.Input.Grade, Position: c.Input.Position,
	}); err != nil {
		return err
	}
	return h.u.Do(ctx, func(tx context.Context) error { return h.r.UpdateTechnology(tx, c) })
}

type DeleteTechnologyHandler struct {
	r TechnologiesRepository
	u UnitOfWork
}

func NewDeleteTechnologyHandler(r TechnologiesRepository, u UnitOfWork) *DeleteTechnologyHandler {
	return &DeleteTechnologyHandler{r, u}
}
func (h *DeleteTechnologyHandler) Handle(ctx context.Context, c DeleteTechnologyCommand) error {
	return h.u.Do(ctx, func(tx context.Context) error { return h.r.DeleteTechnology(tx, c) })
}

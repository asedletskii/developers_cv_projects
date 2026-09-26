package commands

import (
	"context"
	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
	"fmt"
)

type ActivityInput struct {
	ID, ResumeID valueobjects.ID
	Type         enums.ActivityType
	Title        string
	Description  *string
	Date         *valueobjects.YearMonth
	Link         *valueobjects.ExternalURL
}
type CreateActivityCommand struct {
	ID, ResumeID valueobjects.ID
	UserID       useridentity.ID
	Input        ActivityInput
}
type UpdateActivityCommand struct {
	ResumeID, ActivityID valueobjects.ID
	UserID               useridentity.ID
	Input                ActivityInput
}
type DeleteActivityCommand struct {
	ResumeID, ActivityID valueobjects.ID
	UserID               useridentity.ID
}

type CreateActivityHandler struct {
	r ports.CreateActivityRepository
	u UnitOfWork
}

func NewCreateActivityHandler(r ports.CreateActivityRepository, u UnitOfWork) *CreateActivityHandler {
	return &CreateActivityHandler{r: r, u: u}
}
func (h *CreateActivityHandler) Handle(ctx context.Context, c CreateActivityCommand) (valueobjects.ID, error) {
	if h == nil || h.r == nil || h.u == nil {
		return valueobjects.ID{}, fmt.Errorf("create activity dependencies are nil")
	}
	if c.ID.IsZero() {
		c.ID = valueobjects.NewID()
	}
	c.Input.ID, c.Input.ResumeID = c.ID, c.ResumeID
	a, err := activityFromInput(c.Input)
	if err != nil {
		return valueobjects.ID{}, err
	}
	err = h.u.Do(ctx, func(tx context.Context) error {
		return h.r.CreateActivity(tx, ports.CreateActivityRepositoryCommand{Activity: a, ResumeID: c.ResumeID, UserID: c.UserID})
	})
	if err != nil {
		return valueobjects.ID{}, err
	}
	return a.ID(), nil
}

type UpdateActivityHandler struct {
	u   ports.UpdateActivityRepository
	uow UnitOfWork
}

func NewUpdateActivityHandler(u ports.UpdateActivityRepository, uow UnitOfWork) *UpdateActivityHandler {
	return &UpdateActivityHandler{u: u, uow: uow}
}

func (h *UpdateActivityHandler) Handle(ctx context.Context, c UpdateActivityCommand) error {
	if h == nil || h.u == nil || h.uow == nil {
		return fmt.Errorf("update activity dependencies are nil")
	}
	c.Input.ID, c.Input.ResumeID = c.ActivityID, c.ResumeID
	entity, err := activityFromInput(c.Input)
	if err != nil {
		return err
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		return h.u.UpdateActivity(tx, ports.UpdateActivityRepositoryCommand{Activity: entity, ResumeID: c.ResumeID, UserID: c.UserID})
	})
}

type DeleteActivityHandler struct {
	r ports.DeleteActivityRepository
	u UnitOfWork
}

func NewDeleteActivityHandler(r ports.DeleteActivityRepository, u UnitOfWork) *DeleteActivityHandler {
	return &DeleteActivityHandler{r: r, u: u}
}
func (h *DeleteActivityHandler) Handle(ctx context.Context, c DeleteActivityCommand) error {
	if h == nil || h.r == nil || h.u == nil {
		return fmt.Errorf("delete activity dependencies are nil")
	}
	return h.u.Do(ctx, func(tx context.Context) error {
		return h.r.DeleteActivity(tx, ports.DeleteActivityRepositoryCommand{ResumeID: c.ResumeID, ActivityID: c.ActivityID, UserID: c.UserID})
	})
}
func activityFromInput(input ActivityInput) (*domain.Activity, error) {
	return domain.NewActivity(domain.ActivityParams{ID: input.ID, ResumeID: input.ResumeID, Type: input.Type, Title: input.Title, Description: input.Description, Date: input.Date, Link: input.Link})
}

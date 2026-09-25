package commands

import (
	"context"
	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
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
type GetActivityCommand struct {
	ResumeID, ActivityID valueobjects.ID
	UserID               useridentity.ID
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
type CreateActivityRepository = ports.CreateActivityRepository
type CreateActivityRepositoryCommand = ports.CreateActivityRepositoryCommand
type GetActivityRepository = ports.GetActivityRepository
type UpdateActivityRepository = ports.UpdateActivityRepository
type UpdateActivityRepositoryCommand = ports.UpdateActivityRepositoryCommand
type DeleteActivityRepository = ports.DeleteActivityRepository

type CreateActivityHandler struct {
	r CreateActivityRepository
	u UnitOfWork
}

func NewCreateActivityHandler(r CreateActivityRepository, u UnitOfWork) *CreateActivityHandler {
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
	r   GetActivityRepository
	u   UpdateActivityRepository
	uow UnitOfWork
}

func NewUpdateActivityHandler(r GetActivityRepository, u UpdateActivityRepository, uow UnitOfWork) *UpdateActivityHandler {
	return &UpdateActivityHandler{r: r, u: u, uow: uow}
}
func (h *UpdateActivityHandler) Handle(ctx context.Context, c UpdateActivityCommand) error {
	if h == nil || h.r == nil || h.u == nil || h.uow == nil {
		return fmt.Errorf("update activity dependencies are nil")
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		current, err := h.r.GetActivity(tx, ports.GetActivityRepositoryCommand{ResumeID: c.ResumeID, ActivityID: c.ActivityID, UserID: c.UserID})
		if err != nil {
			return err
		}
		c.Input.ID, c.Input.ResumeID = current.ID(), c.ResumeID
		if err := current.Replace(domain.ActivityParams{ID: c.Input.ID, ResumeID: c.Input.ResumeID, Type: c.Input.Type, Title: c.Input.Title, Description: c.Input.Description, Date: c.Input.Date, Link: c.Input.Link}); err != nil {
			return err
		}
		return h.u.UpdateActivity(tx, ports.UpdateActivityRepositoryCommand{Activity: current, ResumeID: c.ResumeID, UserID: c.UserID})
	})
}

type DeleteActivityHandler struct {
	r DeleteActivityRepository
	u UnitOfWork
}

func NewDeleteActivityHandler(r DeleteActivityRepository, u UnitOfWork) *DeleteActivityHandler {
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

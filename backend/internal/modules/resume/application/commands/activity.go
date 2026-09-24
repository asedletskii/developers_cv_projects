package commands

import (
	"context"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
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
type ListActivitiesQuery struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type GetActivityQuery struct {
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

type ActivitiesRepository interface {
	CreateActivity(context.Context, CreateActivityCommand) error
	ListActivities(context.Context, ListActivitiesQuery) ([]*domain.Activity, error)
	GetActivity(context.Context, GetActivityQuery) (*domain.Activity, error)
	UpdateActivity(context.Context, UpdateActivityCommand) error
	DeleteActivity(context.Context, DeleteActivityCommand) error
}
type CreateActivityHandler struct {
	repository ActivitiesRepository
	uow        UnitOfWork
}

func NewCreateActivityHandler(repository ActivitiesRepository, uow UnitOfWork) *CreateActivityHandler {
	return &CreateActivityHandler{repository: repository, uow: uow}
}
func (h *CreateActivityHandler) Handle(ctx context.Context, command CreateActivityCommand) (valueobjects.ID, error) {
	if command.ID.IsZero() {
		command.ID = valueobjects.NewID()
	}
	command.Input.ID, command.Input.ResumeID = command.ID, command.ResumeID
	if _, err := activityFromInput(command.Input); err != nil {
		return valueobjects.ID{}, err
	}
	return command.ID, h.uow.Do(ctx, func(tx context.Context) error { return h.repository.CreateActivity(tx, command) })
}

type ListActivitiesHandler struct{ repository ActivitiesRepository }

func NewListActivitiesHandler(repository ActivitiesRepository) *ListActivitiesHandler {
	return &ListActivitiesHandler{repository: repository}
}
func (h *ListActivitiesHandler) Handle(ctx context.Context, query ListActivitiesQuery) ([]*domain.Activity, error) {
	return h.repository.ListActivities(ctx, query)
}

type GetActivityHandler struct{ repository ActivitiesRepository }

func NewGetActivityHandler(repository ActivitiesRepository) *GetActivityHandler {
	return &GetActivityHandler{repository: repository}
}
func (h *GetActivityHandler) Handle(ctx context.Context, query GetActivityQuery) (*domain.Activity, error) {
	return h.repository.GetActivity(ctx, query)
}

type UpdateActivityHandler struct {
	repository ActivitiesRepository
	uow        UnitOfWork
}

func NewUpdateActivityHandler(repository ActivitiesRepository, uow UnitOfWork) *UpdateActivityHandler {
	return &UpdateActivityHandler{repository: repository, uow: uow}
}
func (h *UpdateActivityHandler) Handle(ctx context.Context, command UpdateActivityCommand) error {
	command.Input.ID, command.Input.ResumeID = command.ActivityID, command.ResumeID
	if _, err := activityFromInput(command.Input); err != nil {
		return err
	}
	return h.uow.Do(ctx, func(tx context.Context) error { return h.repository.UpdateActivity(tx, command) })
}

type DeleteActivityHandler struct {
	repository ActivitiesRepository
	uow        UnitOfWork
}

func NewDeleteActivityHandler(repository ActivitiesRepository, uow UnitOfWork) *DeleteActivityHandler {
	return &DeleteActivityHandler{repository: repository, uow: uow}
}
func (h *DeleteActivityHandler) Handle(ctx context.Context, command DeleteActivityCommand) error {
	return h.uow.Do(ctx, func(tx context.Context) error { return h.repository.DeleteActivity(tx, command) })
}
func activityFromInput(input ActivityInput) (*domain.Activity, error) {
	return domain.NewActivity(domain.ActivityParams{ID: input.ID, ResumeID: input.ResumeID, Type: input.Type, Title: input.Title, Description: input.Description, Date: input.Date, Link: input.Link})
}

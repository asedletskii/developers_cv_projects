package commands

import (
	"context"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type EducationInput struct {
	ID          valueobjects.ID
	ResumeID    valueobjects.ID
	Location    string
	University  string
	Grade       enums.EducationLevel
	ProgramName string
	ProfileName *string
	Faculty     string
	Course      *int
	EndDate     *valueobjects.Year
}

type CreateEducationCommand struct {
	ID, ResumeID valueobjects.ID
	UserID       useridentity.ID
	Input        EducationInput
}
type ListEducationsQuery struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type GetEducationQuery struct {
	ResumeID, EducationID valueobjects.ID
	UserID                useridentity.ID
}
type UpdateEducationCommand struct {
	ResumeID, EducationID valueobjects.ID
	UserID                useridentity.ID
	Input                 EducationInput
}
type DeleteEducationCommand struct {
	ResumeID, EducationID valueobjects.ID
	UserID                useridentity.ID
}

type EducationsRepository interface {
	CreateEducation(context.Context, CreateEducationCommand) error
	ListEducations(context.Context, ListEducationsQuery) ([]*domain.Education, error)
	GetEducation(context.Context, GetEducationQuery) (*domain.Education, error)
	UpdateEducation(context.Context, UpdateEducationCommand) error
	DeleteEducation(context.Context, DeleteEducationCommand) error
}

type CreateEducationHandler struct {
	repository EducationsRepository
	uow        UnitOfWork
}

func NewCreateEducationHandler(repository EducationsRepository, uow UnitOfWork) *CreateEducationHandler {
	return &CreateEducationHandler{repository: repository, uow: uow}
}
func (h *CreateEducationHandler) Handle(ctx context.Context, command CreateEducationCommand) (valueobjects.ID, error) {
	if command.ID.IsZero() {
		command.ID = valueobjects.NewID()
	}
	command.Input.ID, command.Input.ResumeID = command.ID, command.ResumeID
	if _, err := educationFromInput(command.Input); err != nil {
		return valueobjects.ID{}, err
	}
	return command.ID, h.uow.Do(ctx, func(tx context.Context) error { return h.repository.CreateEducation(tx, command) })
}

type ListEducationsHandler struct{ repository EducationsRepository }

func NewListEducationsHandler(repository EducationsRepository) *ListEducationsHandler {
	return &ListEducationsHandler{repository: repository}
}
func (h *ListEducationsHandler) Handle(ctx context.Context, query ListEducationsQuery) ([]*domain.Education, error) {
	return h.repository.ListEducations(ctx, query)
}

type GetEducationHandler struct{ repository EducationsRepository }

func NewGetEducationHandler(repository EducationsRepository) *GetEducationHandler {
	return &GetEducationHandler{repository: repository}
}
func (h *GetEducationHandler) Handle(ctx context.Context, query GetEducationQuery) (*domain.Education, error) {
	return h.repository.GetEducation(ctx, query)
}

type UpdateEducationHandler struct {
	repository EducationsRepository
	uow        UnitOfWork
}

func NewUpdateEducationHandler(repository EducationsRepository, uow UnitOfWork) *UpdateEducationHandler {
	return &UpdateEducationHandler{repository: repository, uow: uow}
}
func (h *UpdateEducationHandler) Handle(ctx context.Context, command UpdateEducationCommand) error {
	command.Input.ID, command.Input.ResumeID = command.EducationID, command.ResumeID
	if _, err := educationFromInput(command.Input); err != nil {
		return err
	}
	return h.uow.Do(ctx, func(tx context.Context) error { return h.repository.UpdateEducation(tx, command) })
}

type DeleteEducationHandler struct {
	repository EducationsRepository
	uow        UnitOfWork
}

func NewDeleteEducationHandler(repository EducationsRepository, uow UnitOfWork) *DeleteEducationHandler {
	return &DeleteEducationHandler{repository: repository, uow: uow}
}
func (h *DeleteEducationHandler) Handle(ctx context.Context, command DeleteEducationCommand) error {
	return h.uow.Do(ctx, func(tx context.Context) error { return h.repository.DeleteEducation(tx, command) })
}

func educationFromInput(input EducationInput) (*domain.Education, error) {
	return domain.NewEducation(domain.EducationParams{ID: input.ID, ResumeID: input.ResumeID, Location: input.Location, University: input.University, Grade: input.Grade, ProgramName: input.ProgramName, ProfileName: input.ProfileName, Faculty: input.Faculty, Course: input.Course, EndDate: input.EndDate})
}

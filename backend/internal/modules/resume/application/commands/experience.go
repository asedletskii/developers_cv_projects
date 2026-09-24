package commands

import (
	"context"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type ExperienceTechnologyInput struct {
	Name  string
	Grade *enums.TechnologyGrade
}

type ExperienceInput struct {
	ID           valueobjects.ID
	ResumeID     valueobjects.ID
	Company      string
	Position     string
	Location     *string
	StartDate    valueobjects.YearMonth
	EndDate      *valueobjects.YearMonth
	Description  *string
	Achievements []string
	Technologies []ExperienceTechnologyInput
}

type CreateExperienceCommand struct {
	ID, ResumeID valueobjects.ID
	UserID       useridentity.ID
	Input        ExperienceInput
}

type ListExperiencesQuery struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}

type GetExperienceQuery struct {
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

type ExperiencesRepository interface {
	CreateExperience(context.Context, CreateExperienceCommand) error
	ListExperiences(context.Context, ListExperiencesQuery) ([]*domain.Experience, error)
	GetExperience(context.Context, GetExperienceQuery) (*domain.Experience, error)
	UpdateExperience(context.Context, UpdateExperienceCommand) error
	DeleteExperience(context.Context, DeleteExperienceCommand) error
}

type CreateExperienceHandler struct {
	repository ExperiencesRepository
	uow        UnitOfWork
}

func NewCreateExperienceHandler(repository ExperiencesRepository, uow UnitOfWork) *CreateExperienceHandler {
	return &CreateExperienceHandler{repository: repository, uow: uow}
}

func (h *CreateExperienceHandler) Handle(ctx context.Context, command CreateExperienceCommand) (valueobjects.ID, error) {
	if command.ID.IsZero() {
		command.ID = valueobjects.NewID()
	}
	command.Input.ID = command.ID
	command.Input.ResumeID = command.ResumeID
	if _, err := experienceFromInput(command.Input); err != nil {
		return valueobjects.ID{}, err
	}
	return command.ID, h.uow.Do(ctx, func(tx context.Context) error {
		return h.repository.CreateExperience(tx, command)
	})
}

type ListExperiencesHandler struct{ repository ExperiencesRepository }

func NewListExperiencesHandler(repository ExperiencesRepository) *ListExperiencesHandler {
	return &ListExperiencesHandler{repository: repository}
}

func (h *ListExperiencesHandler) Handle(ctx context.Context, query ListExperiencesQuery) ([]*domain.Experience, error) {
	return h.repository.ListExperiences(ctx, query)
}

type GetExperienceHandler struct{ repository ExperiencesRepository }

func NewGetExperienceHandler(repository ExperiencesRepository) *GetExperienceHandler {
	return &GetExperienceHandler{repository: repository}
}

func (h *GetExperienceHandler) Handle(ctx context.Context, query GetExperienceQuery) (*domain.Experience, error) {
	return h.repository.GetExperience(ctx, query)
}

type UpdateExperienceHandler struct {
	repository ExperiencesRepository
	uow        UnitOfWork
}

func NewUpdateExperienceHandler(repository ExperiencesRepository, uow UnitOfWork) *UpdateExperienceHandler {
	return &UpdateExperienceHandler{repository: repository, uow: uow}
}

func (h *UpdateExperienceHandler) Handle(ctx context.Context, command UpdateExperienceCommand) error {
	command.Input.ID = command.ExperienceID
	command.Input.ResumeID = command.ResumeID
	if _, err := experienceFromInput(command.Input); err != nil {
		return err
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		return h.repository.UpdateExperience(tx, command)
	})
}

type DeleteExperienceHandler struct {
	repository ExperiencesRepository
	uow        UnitOfWork
}

func NewDeleteExperienceHandler(repository ExperiencesRepository, uow UnitOfWork) *DeleteExperienceHandler {
	return &DeleteExperienceHandler{repository: repository, uow: uow}
}

func (h *DeleteExperienceHandler) Handle(ctx context.Context, command DeleteExperienceCommand) error {
	return h.uow.Do(ctx, func(tx context.Context) error {
		return h.repository.DeleteExperience(tx, command)
	})
}

func experienceFromInput(input ExperienceInput) (*domain.Experience, error) {
	technologies := make([]domain.ExperienceTechnologyParams, len(input.Technologies))
	for i, technology := range input.Technologies {
		technologies[i] = domain.ExperienceTechnologyParams{Name: technology.Name, Grade: technology.Grade}
	}
	return domain.NewExperience(domain.ExperienceParams{
		ID:           input.ID,
		ResumeID:     input.ResumeID,
		Company:      input.Company,
		Position:     input.Position,
		Location:     input.Location,
		StartDate:    input.StartDate,
		EndDate:      input.EndDate,
		Description:  input.Description,
		Achievements: input.Achievements,
		Technologies: technologies,
	})
}

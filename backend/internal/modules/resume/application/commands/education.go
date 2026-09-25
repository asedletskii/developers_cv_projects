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

type EducationInput struct {
	ID, ResumeID         valueobjects.ID
	Location, University string
	Grade                enums.EducationLevel
	ProgramName, Faculty string
	ProfileName          *string
	Course               *int
	EndDate              *valueobjects.Year
}
type CreateEducationCommand struct {
	ID, ResumeID valueobjects.ID
	UserID       useridentity.ID
	Input        EducationInput
}
type GetEducationCommand struct {
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
type CreateEducationRepository = ports.CreateEducationRepository
type CreateEducationRepositoryCommand = ports.CreateEducationRepositoryCommand
type GetEducationRepository = ports.GetEducationRepository
type UpdateEducationRepository = ports.UpdateEducationRepository
type UpdateEducationRepositoryCommand = ports.UpdateEducationRepositoryCommand
type DeleteEducationRepository = ports.DeleteEducationRepository

type CreateEducationHandler struct {
	r CreateEducationRepository
	u UnitOfWork
}

func NewCreateEducationHandler(r CreateEducationRepository, u UnitOfWork) *CreateEducationHandler {
	return &CreateEducationHandler{r: r, u: u}
}
func (h *CreateEducationHandler) Handle(ctx context.Context, c CreateEducationCommand) (valueobjects.ID, error) {
	if h == nil || h.r == nil || h.u == nil {
		return valueobjects.ID{}, fmt.Errorf("create education dependencies are nil")
	}
	if c.ID.IsZero() {
		c.ID = valueobjects.NewID()
	}
	c.Input.ID, c.Input.ResumeID = c.ID, c.ResumeID
	a, err := educationFromInput(c.Input)
	if err != nil {
		return valueobjects.ID{}, err
	}
	err = h.u.Do(ctx, func(tx context.Context) error {
		return h.r.CreateEducation(tx, ports.CreateEducationRepositoryCommand{Education: a, ResumeID: c.ResumeID, UserID: c.UserID})
	})
	if err != nil {
		return valueobjects.ID{}, err
	}
	return a.ID(), nil
}

type UpdateEducationHandler struct {
	r   GetEducationRepository
	u   UpdateEducationRepository
	uow UnitOfWork
}

func NewUpdateEducationHandler(r GetEducationRepository, u UpdateEducationRepository, uow UnitOfWork) *UpdateEducationHandler {
	return &UpdateEducationHandler{r: r, u: u, uow: uow}
}
func (h *UpdateEducationHandler) Handle(ctx context.Context, c UpdateEducationCommand) error {
	if h == nil || h.r == nil || h.u == nil || h.uow == nil {
		return fmt.Errorf("update education dependencies are nil")
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		current, err := h.r.GetEducation(tx, ports.GetEducationRepositoryCommand{ResumeID: c.ResumeID, EducationID: c.EducationID, UserID: c.UserID})
		if err != nil {
			return err
		}
		c.Input.ID, c.Input.ResumeID = current.ID(), c.ResumeID
		if err := current.Replace(domain.EducationParams{ID: c.Input.ID, ResumeID: c.Input.ResumeID, Location: c.Input.Location, University: c.Input.University, Grade: c.Input.Grade, ProgramName: c.Input.ProgramName, ProfileName: c.Input.ProfileName, Faculty: c.Input.Faculty, Course: c.Input.Course, EndDate: c.Input.EndDate}); err != nil {
			return err
		}
		return h.u.UpdateEducation(tx, ports.UpdateEducationRepositoryCommand{Education: current, ResumeID: c.ResumeID, UserID: c.UserID})
	})
}

type DeleteEducationHandler struct {
	r DeleteEducationRepository
	u UnitOfWork
}

func NewDeleteEducationHandler(r DeleteEducationRepository, u UnitOfWork) *DeleteEducationHandler {
	return &DeleteEducationHandler{r: r, u: u}
}
func (h *DeleteEducationHandler) Handle(ctx context.Context, c DeleteEducationCommand) error {
	if h == nil || h.r == nil || h.u == nil {
		return fmt.Errorf("delete education dependencies are nil")
	}
	return h.u.Do(ctx, func(tx context.Context) error {
		return h.r.DeleteEducation(tx, ports.DeleteEducationRepositoryCommand{ResumeID: c.ResumeID, EducationID: c.EducationID, UserID: c.UserID})
	})
}
func educationFromInput(input EducationInput) (*domain.Education, error) {
	return domain.NewEducation(domain.EducationParams{ID: input.ID, ResumeID: input.ResumeID, Location: input.Location, University: input.University, Grade: input.Grade, ProgramName: input.ProgramName, ProfileName: input.ProfileName, Faculty: input.Faculty, Course: input.Course, EndDate: input.EndDate})
}

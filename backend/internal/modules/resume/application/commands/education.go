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

type EducationInput struct {
	ID                   valueobjects.ID
	Location, University string
	Grade                enums.EducationLevel
	ProgramName, Faculty string
	ProfileName          *string
	Course               *int
	EndDate              *valueobjects.Year
}
type EducationsInput struct{ Items []EducationInput }

func (i EducationsInput) toDomain() []domain.EducationParams {
	out := make([]domain.EducationParams, 0, len(i.Items))
	for _, x := range i.Items {
		out = append(out, domain.EducationParams{ID: x.ID, Location: x.Location, University: x.University, Grade: x.Grade, ProgramName: x.ProgramName, ProfileName: x.ProfileName, Faculty: x.Faculty, Course: x.Course, EndDate: x.EndDate})
	}
	return out
}

type ReplaceEducationsCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
	Version  int64
	Input    EducationsInput
}
type ReplaceEducationsHandler struct {
	repository ports.ResumeRepository
	uow        UnitOfWork
}

func NewReplaceEducationsHandler(repository ports.ResumeRepository, uow UnitOfWork) *ReplaceEducationsHandler {
	return &ReplaceEducationsHandler{repository: repository, uow: uow}
}
func (h *ReplaceEducationsHandler) Handle(ctx context.Context, c ReplaceEducationsCommand) error {
	if h == nil || h.repository == nil || h.uow == nil {
		return fmt.Errorf("replace educations handler dependencies are nil")
	}
	if c.ResumeID.IsZero() || c.UserID.IsZero() || c.Version < 1 {
		return fmt.Errorf("invalid replace educations command")
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		r, err := h.repository.Get(tx, ports.GetResumeParams{ResumeID: c.ResumeID, UserID: c.UserID})
		if err != nil {
			return err
		}
		if r.Version() != c.Version {
			return domain.ErrConcurrencyConflict
		}
		if err := r.ReplaceEducations(c.Input.toDomain()); err != nil {
			return err
		}
		return h.repository.Save(tx, r)
	})
}

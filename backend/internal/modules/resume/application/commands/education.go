package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type EducationInput struct {
	ID                   sharedvalueobjects.ID
	Location, University string
	Grade                domain.EducationLevel
	ProgramName, Faculty string
	ProfileName          *string
	Course               *int
	EndDate              *domain.Year
}
type EducationsInput struct{ Items []EducationInput }

func (i EducationsInput) toDomain() []domain.EducationParams {
	out := make([]domain.EducationParams, 0, len(i.Items))
	for _, x := range i.Items {
		out = append(out, domain.EducationParams{
			ID:          x.ID,
			Location:    x.Location,
			University:  x.University,
			Grade:       x.Grade,
			ProgramName: x.ProgramName,
			ProfileName: x.ProfileName,
			Faculty:     x.Faculty,
			Course:      x.Course,
			EndDate:     x.EndDate,
		})
	}
	return out
}

type ReplaceEducationsCommand struct {
	ResumeID sharedvalueobjects.ID
	UserID   sharedvalueobjects.ID
	Version  int64
	Input    EducationsInput
}
type ReplaceEducationsHandler struct {
	repository ports.ResumeRepository
	uow        ports.UnitOfWork
}

func NewReplaceEducationsHandler(repository ports.ResumeRepository, uow ports.UnitOfWork) *ReplaceEducationsHandler {
	return &ReplaceEducationsHandler{repository: repository, uow: uow}
}
func (h *ReplaceEducationsHandler) Handle(ctx context.Context, c ReplaceEducationsCommand) error {
	if h == nil {
		return fmt.Errorf("replace educations handler dependencies are nil")
	}

	educations := c.Input.toDomain()
	return executeResumeMutation(ctx, h.repository, h.uow, c.ResumeID, c.UserID, c.Version, "replace educations", ports.ResumeMutationEducations, func(resume *domain.Resume) error {
		return resume.ReplaceEducations(educations)
	})
}

package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type CourseInput struct {
	ID              sharedvalueobjects.ID
	Name            string
	Organization    string
	EndDate         domain.YearMonth
	Hours           *int
	CertificateLink *domain.ExternalURL
}
type CoursesInput struct{ Items []CourseInput }

func (i CoursesInput) toDomain() []domain.CourseParams {
	out := make([]domain.CourseParams, 0, len(i.Items))
	for _, x := range i.Items {
		out = append(out, domain.CourseParams{
			ID:              x.ID,
			Name:            x.Name,
			Organization:    x.Organization,
			EndDate:         x.EndDate,
			Hours:           x.Hours,
			CertificateLink: x.CertificateLink,
		})
	}
	return out
}

type ReplaceCoursesCommand struct {
	ResumeID sharedvalueobjects.ID
	UserID   sharedvalueobjects.ID
	Version  int64
	Input    CoursesInput
}
type ReplaceCoursesHandler struct {
	repository ports.ResumeRepository
	uow        ports.UnitOfWork
}

func NewReplaceCoursesHandler(repository ports.ResumeRepository, uow ports.UnitOfWork) *ReplaceCoursesHandler {
	return &ReplaceCoursesHandler{repository: repository, uow: uow}
}
func (h *ReplaceCoursesHandler) Handle(ctx context.Context, c ReplaceCoursesCommand) error {
	if h == nil {
		return fmt.Errorf("replace courses handler dependencies are nil")
	}

	courses := c.Input.toDomain()
	return executeResumeMutation(ctx, h.repository, h.uow, c.ResumeID, c.UserID, c.Version, "replace courses", ports.ResumeMutationCourses, func(resume *domain.Resume) error {
		return resume.ReplaceCourses(courses)
	})
}

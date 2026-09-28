package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type CourseInput struct {
	ID              valueobjects.ID
	Name            string
	Organization    string
	EndDate         valueobjects.YearMonth
	Hours           *int
	CertificateLink *valueobjects.ExternalURL
}
type CoursesInput struct{ Items []CourseInput }

func (i CoursesInput) toDomain() []domain.CourseParams {
	out := make([]domain.CourseParams, 0, len(i.Items))
	for _, x := range i.Items {
		out = append(out, domain.CourseParams{ID: x.ID, Name: x.Name, Organization: x.Organization, EndDate: x.EndDate, Hours: x.Hours, CertificateLink: x.CertificateLink})
	}
	return out
}

type ReplaceCoursesCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
	Version  int64
	Input    CoursesInput
}
type ReplaceCoursesHandler struct {
	repository ports.ResumeRepository
	uow        UnitOfWork
}

func NewReplaceCoursesHandler(repository ports.ResumeRepository, uow UnitOfWork) *ReplaceCoursesHandler {
	return &ReplaceCoursesHandler{repository: repository, uow: uow}
}
func (h *ReplaceCoursesHandler) Handle(ctx context.Context, c ReplaceCoursesCommand) error {
	if h == nil || h.repository == nil || h.uow == nil {
		return fmt.Errorf("replace courses handler dependencies are nil")
	}
	if c.ResumeID.IsZero() || c.UserID.IsZero() || c.Version < 1 {
		return fmt.Errorf("invalid replace courses command")
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		r, err := h.repository.Get(tx, ports.GetResumeParams{ResumeID: c.ResumeID, UserID: c.UserID})
		if err != nil {
			return err
		}
		if r.Version() != c.Version {
			return domain.ErrConcurrencyConflict
		}
		if err := r.ReplaceCourses(c.Input.toDomain()); err != nil {
			return err
		}
		return h.repository.Save(tx, r)
	})
}

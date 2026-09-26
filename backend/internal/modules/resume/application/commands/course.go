package commands

import (
	"context"
	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
	"fmt"
)

type CourseInput struct {
	ID, ResumeID       valueobjects.ID
	Name, Organization string
	EndDate            valueobjects.YearMonth
	Hours              *int
	CertificateLink    *valueobjects.ExternalURL
}
type CreateCourseCommand struct {
	ID, ResumeID valueobjects.ID
	UserID       useridentity.ID
	Input        CourseInput
}
type UpdateCourseCommand struct {
	ResumeID, CourseID valueobjects.ID
	UserID             useridentity.ID
	Input              CourseInput
}
type DeleteCourseCommand struct {
	ResumeID, CourseID valueobjects.ID
	UserID             useridentity.ID
}

type CreateCourseHandler struct {
	r ports.CreateCourseRepository
	u UnitOfWork
}

func NewCreateCourseHandler(r ports.CreateCourseRepository, u UnitOfWork) *CreateCourseHandler {
	return &CreateCourseHandler{r: r, u: u}
}
func (h *CreateCourseHandler) Handle(ctx context.Context, c CreateCourseCommand) (valueobjects.ID, error) {
	if h == nil || h.r == nil || h.u == nil {
		return valueobjects.ID{}, fmt.Errorf("create course dependencies are nil")
	}
	if c.ID.IsZero() {
		c.ID = valueobjects.NewID()
	}
	c.Input.ID, c.Input.ResumeID = c.ID, c.ResumeID
	a, err := courseFromInput(c.Input)
	if err != nil {
		return valueobjects.ID{}, err
	}
	err = h.u.Do(ctx, func(tx context.Context) error {
		return h.r.CreateCourse(tx, ports.CreateCourseRepositoryCommand{Course: a, ResumeID: c.ResumeID, UserID: c.UserID})
	})
	if err != nil {
		return valueobjects.ID{}, err
	}
	return a.ID(), nil
}

type UpdateCourseHandler struct {
	u   ports.UpdateCourseRepository
	uow UnitOfWork
}

func NewUpdateCourseHandler(u ports.UpdateCourseRepository, uow UnitOfWork) *UpdateCourseHandler {
	return &UpdateCourseHandler{u: u, uow: uow}
}
func (h *UpdateCourseHandler) Handle(ctx context.Context, c UpdateCourseCommand) error {
	if h == nil || h.u == nil || h.uow == nil {
		return fmt.Errorf("update course dependencies are nil")
	}
	c.Input.ID, c.Input.ResumeID = c.CourseID, c.ResumeID
	course, err := courseFromInput(c.Input)
	if err != nil {
		return err
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		return h.u.UpdateCourse(tx, ports.UpdateCourseRepositoryCommand{Course: course, ResumeID: c.ResumeID, UserID: c.UserID})
	})
}

type DeleteCourseHandler struct {
	r ports.DeleteCourseRepository
	u UnitOfWork
}

func NewDeleteCourseHandler(r ports.DeleteCourseRepository, u UnitOfWork) *DeleteCourseHandler {
	return &DeleteCourseHandler{r: r, u: u}
}
func (h *DeleteCourseHandler) Handle(ctx context.Context, c DeleteCourseCommand) error {
	if h == nil || h.r == nil || h.u == nil {
		return fmt.Errorf("delete course dependencies are nil")
	}
	return h.u.Do(ctx, func(tx context.Context) error {
		return h.r.DeleteCourse(tx, ports.DeleteCourseRepositoryCommand{ResumeID: c.ResumeID, CourseID: c.CourseID, UserID: c.UserID})
	})
}
func courseFromInput(input CourseInput) (*domain.Course, error) {
	return domain.NewCourse(domain.CourseParams{ID: input.ID, ResumeID: input.ResumeID, Name: input.Name, Organization: input.Organization, EndDate: input.EndDate, Hours: input.Hours, CertificateLink: input.CertificateLink})
}

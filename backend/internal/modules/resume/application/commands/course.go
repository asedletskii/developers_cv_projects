package commands

import (
	"context"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
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
type ListCoursesQuery struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type GetCourseQuery struct {
	ResumeID, CourseID valueobjects.ID
	UserID             useridentity.ID
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

type CoursesRepository interface {
	CreateCourse(context.Context, CreateCourseCommand) error
	ListCourses(context.Context, ListCoursesQuery) ([]*domain.Course, error)
	GetCourse(context.Context, GetCourseQuery) (*domain.Course, error)
	UpdateCourse(context.Context, UpdateCourseCommand) error
	DeleteCourse(context.Context, DeleteCourseCommand) error
}
type CreateCourseHandler struct {
	repository CoursesRepository
	uow        UnitOfWork
}

func NewCreateCourseHandler(repository CoursesRepository, uow UnitOfWork) *CreateCourseHandler {
	return &CreateCourseHandler{repository: repository, uow: uow}
}
func (h *CreateCourseHandler) Handle(ctx context.Context, command CreateCourseCommand) (valueobjects.ID, error) {
	if command.ID.IsZero() {
		command.ID = valueobjects.NewID()
	}
	command.Input.ID, command.Input.ResumeID = command.ID, command.ResumeID
	if _, err := courseFromInput(command.Input); err != nil {
		return valueobjects.ID{}, err
	}
	return command.ID, h.uow.Do(ctx, func(tx context.Context) error { return h.repository.CreateCourse(tx, command) })
}

type ListCoursesHandler struct{ repository CoursesRepository }

func NewListCoursesHandler(repository CoursesRepository) *ListCoursesHandler {
	return &ListCoursesHandler{repository: repository}
}
func (h *ListCoursesHandler) Handle(ctx context.Context, query ListCoursesQuery) ([]*domain.Course, error) {
	return h.repository.ListCourses(ctx, query)
}

type GetCourseHandler struct{ repository CoursesRepository }

func NewGetCourseHandler(repository CoursesRepository) *GetCourseHandler {
	return &GetCourseHandler{repository: repository}
}
func (h *GetCourseHandler) Handle(ctx context.Context, query GetCourseQuery) (*domain.Course, error) {
	return h.repository.GetCourse(ctx, query)
}

type UpdateCourseHandler struct {
	repository CoursesRepository
	uow        UnitOfWork
}

func NewUpdateCourseHandler(repository CoursesRepository, uow UnitOfWork) *UpdateCourseHandler {
	return &UpdateCourseHandler{repository: repository, uow: uow}
}
func (h *UpdateCourseHandler) Handle(ctx context.Context, command UpdateCourseCommand) error {
	command.Input.ID, command.Input.ResumeID = command.CourseID, command.ResumeID
	if _, err := courseFromInput(command.Input); err != nil {
		return err
	}
	return h.uow.Do(ctx, func(tx context.Context) error { return h.repository.UpdateCourse(tx, command) })
}

type DeleteCourseHandler struct {
	repository CoursesRepository
	uow        UnitOfWork
}

func NewDeleteCourseHandler(repository CoursesRepository, uow UnitOfWork) *DeleteCourseHandler {
	return &DeleteCourseHandler{repository: repository, uow: uow}
}
func (h *DeleteCourseHandler) Handle(ctx context.Context, command DeleteCourseCommand) error {
	return h.uow.Do(ctx, func(tx context.Context) error { return h.repository.DeleteCourse(tx, command) })
}
func courseFromInput(input CourseInput) (*domain.Course, error) {
	return domain.NewCourse(domain.CourseParams{ID: input.ID, ResumeID: input.ResumeID, Name: input.Name, Organization: input.Organization, EndDate: input.EndDate, Hours: input.Hours, CertificateLink: input.CertificateLink})
}

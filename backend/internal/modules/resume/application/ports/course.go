package ports

import (
	"context"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type CreateCourseRepositoryCommand struct {
	Course   *domain.Course
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}

type GetCourseRepositoryCommand struct {
	ResumeID, CourseID valueobjects.ID
	UserID             useridentity.ID
}

type UpdateCourseRepositoryCommand struct {
	Course   *domain.Course
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}

type DeleteCourseRepositoryCommand struct {
	ResumeID, CourseID valueobjects.ID
	UserID             useridentity.ID
}

type CreateCourseRepository interface {
	CreateCourse(context.Context, CreateCourseRepositoryCommand) error
}

type GetCourseRepository interface {
	GetCourse(context.Context, GetCourseRepositoryCommand) (*domain.Course, error)
}

type UpdateCourseRepository interface {
	UpdateCourse(context.Context, UpdateCourseRepositoryCommand) error
}

type DeleteCourseRepository interface {
	DeleteCourse(context.Context, DeleteCourseRepositoryCommand) error
}

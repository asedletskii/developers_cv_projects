package ports

import (
	"context"

	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type GetResumeParams struct {
	ResumeID sharedvalueobjects.ID
	UserID   sharedvalueobjects.ID
}

type DeleteResumeParams struct {
	ResumeID sharedvalueobjects.ID
	UserID   sharedvalueobjects.ID
	Version  int64
}

type ResumeMutationSection string

const (
	ResumeMutationTitle        ResumeMutationSection = "title"
	ResumeMutationStatus       ResumeMutationSection = "status"
	ResumeMutationOverview     ResumeMutationSection = "overview"
	ResumeMutationProfile      ResumeMutationSection = "profile"
	ResumeMutationContacts     ResumeMutationSection = "contacts"
	ResumeMutationSkills       ResumeMutationSection = "skills"
	ResumeMutationLanguages    ResumeMutationSection = "languages"
	ResumeMutationTechnologies ResumeMutationSection = "technologies"
	ResumeMutationExperiences  ResumeMutationSection = "experiences"
	ResumeMutationEducations   ResumeMutationSection = "educations"
	ResumeMutationCourses      ResumeMutationSection = "courses"
	ResumeMutationActivities   ResumeMutationSection = "activities"
	ResumeMutationProjects     ResumeMutationSection = "projects"
)

type ResumeRepository interface {
	Get(context.Context, GetResumeParams) (*domain.Resume, error)
	SaveSection(context.Context, *domain.Resume, ResumeMutationSection) error
	Create(context.Context, *domain.Resume) error
	LockForDelete(context.Context, DeleteResumeParams) error
	Delete(context.Context, DeleteResumeParams) error
}

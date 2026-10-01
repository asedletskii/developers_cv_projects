package mappers

import (
	"fmt"
	"time"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func ResumeToDB(resume *domain.Resume, createdAt, updatedAt time.Time) (gormmodel.Resume, error) {
	if resume == nil {
		return gormmodel.Resume{}, fmt.Errorf("resume is nil")
	}
	return gormmodel.Resume{
		ID:        resume.ID().UUID(),
		UserID:    resume.UserID().UUID(),
		Title:     resume.Title(),
		Language:  string(resume.Language()),
		Status:    string(resume.Status()),
		Version:   resume.Version(),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func ResumeToDomainSnapshot(
	model gormmodel.Resume,
	overview *domain.ResumeOverview,
	profile *domain.Profile,
	contacts *domain.Contacts,
	skills *domain.Skills,
	languages *domain.ResumeLanguages,
	technologies *domain.ResumeTechnologies,
	experiences []*domain.Experience,
	educations []*domain.Education,
	courses []*domain.Course,
	activities []*domain.Activity,
	projects []*domain.Project,
) (domain.ResumeSnapshot, error) {
	id, err := sharedvalueobjects.IDFromUUID(model.ID)
	if err != nil {
		return domain.ResumeSnapshot{}, err
	}
	userID, err := sharedvalueobjects.IDFromUUID(model.UserID)
	if err != nil {
		return domain.ResumeSnapshot{}, err
	}
	return domain.ResumeSnapshot{
		ID:           id,
		UserID:       userID,
		Title:        model.Title,
		Language:     domain.ResumeLanguage(model.Language),
		Status:       domain.ResumeStatus(model.Status),
		Version:      model.Version,
		Overview:     overview,
		Profile:      profile,
		Contacts:     contacts,
		Skills:       skills,
		Languages:    languages,
		Technologies: technologies,
		Experiences:  experiences,
		Educations:   educations,
		Courses:      courses,
		Activities:   activities,
		Projects:     projects,
	}, nil
}

package ports

import (
	"context"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type CreateResumeRootRepositoryCommand struct{ Resume *domain.ResumeRoot }
type GetResumeRootRepositoryCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type UpdateResumeRootRepositoryCommand struct{ Resume *domain.ResumeRoot }
type DeleteResumeRootRepositoryCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type CreateResumeRootRepository interface {
	CreateRoot(context.Context, CreateResumeRootRepositoryCommand) error
}
type GetResumeRootRepository interface {
	GetRoot(context.Context, GetResumeRootRepositoryCommand) (*domain.ResumeRoot, error)
}
type UpdateResumeRootRepository interface {
	UpdateRoot(context.Context, UpdateResumeRootRepositoryCommand) error
}
type DeleteResumeRootRepository interface {
	DeleteRoot(context.Context, DeleteResumeRootRepositoryCommand) error
}

type CreateResumeProfileRepositoryCommand struct {
	Profile *domain.Profile
	UserID  useridentity.ID
}
type GetResumeProfileRepositoryCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type UpdateResumeProfileRepositoryCommand struct {
	Profile *domain.Profile
	UserID  useridentity.ID
}
type CreateResumeProfileRepository interface {
	CreateProfile(context.Context, CreateResumeProfileRepositoryCommand) error
}
type GetResumeProfileRepository interface {
	GetProfile(context.Context, GetResumeProfileRepositoryCommand) (*domain.Profile, error)
}
type UpdateResumeProfileRepository interface {
	UpdateProfile(context.Context, UpdateResumeProfileRepositoryCommand) error
}

type CreateResumeContactsRepositoryCommand struct {
	Contacts *domain.Contacts
	UserID   useridentity.ID
}
type GetResumeContactsRepositoryCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type UpdateResumeContactsRepositoryCommand struct {
	Contacts *domain.Contacts
	UserID   useridentity.ID
}
type CreateResumeContactsRepository interface {
	CreateContacts(context.Context, CreateResumeContactsRepositoryCommand) error
}
type GetResumeContactsRepository interface {
	GetContacts(context.Context, GetResumeContactsRepositoryCommand) (*domain.Contacts, error)
}
type UpdateResumeContactsRepository interface {
	UpdateContacts(context.Context, UpdateResumeContactsRepositoryCommand) error
}

type CreateSkillsRepositoryCommand struct {
	Skills *domain.Skills
	UserID useridentity.ID
}
type GetSkillsRepositoryCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type UpdateSkillsRepositoryCommand struct {
	Skills *domain.Skills
	UserID useridentity.ID
}
type CreateSkillsRepository interface {
	CreateSkills(context.Context, CreateSkillsRepositoryCommand) error
}
type GetSkillsRepository interface {
	GetSkills(context.Context, GetSkillsRepositoryCommand) (*domain.Skills, error)
}
type UpdateSkillsRepository interface {
	UpdateSkills(context.Context, UpdateSkillsRepositoryCommand) error
}

type CreateLanguageRepositoryCommand struct {
	Language domain.ResumeLanguageEntry
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type GetLanguageRepositoryCommand struct {
	ResumeID, LanguageID valueobjects.ID
	UserID               useridentity.ID
}
type UpdateLanguageRepositoryCommand struct {
	Language domain.ResumeLanguageEntry
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type DeleteLanguageRepositoryCommand struct {
	ResumeID, LanguageID valueobjects.ID
	UserID               useridentity.ID
}
type CreateLanguageRepository interface {
	CreateLanguage(context.Context, CreateLanguageRepositoryCommand) error
}
type GetLanguageRepository interface {
	GetLanguage(context.Context, GetLanguageRepositoryCommand) (domain.ResumeLanguageEntry, error)
}
type UpdateLanguageRepository interface {
	UpdateLanguage(context.Context, UpdateLanguageRepositoryCommand) error
}
type DeleteLanguageRepository interface {
	DeleteLanguage(context.Context, DeleteLanguageRepositoryCommand) error
}

type CreateTechnologyRepositoryCommand struct {
	Technology domain.ResumeTechnologyEntry
	ResumeID   valueobjects.ID
	UserID     useridentity.ID
}
type GetTechnologyRepositoryCommand struct {
	ResumeID, TechnologyID valueobjects.ID
	UserID                 useridentity.ID
}
type UpdateTechnologyRepositoryCommand struct {
	Technology domain.ResumeTechnologyEntry
	ResumeID   valueobjects.ID
	UserID     useridentity.ID
}
type DeleteTechnologyRepositoryCommand struct {
	ResumeID, TechnologyID valueobjects.ID
	UserID                 useridentity.ID
}
type CreateTechnologyRepository interface {
	CreateTechnology(context.Context, CreateTechnologyRepositoryCommand) error
}
type GetTechnologyRepository interface {
	GetTechnology(context.Context, GetTechnologyRepositoryCommand) (domain.ResumeTechnologyEntry, error)
}
type UpdateTechnologyRepository interface {
	UpdateTechnology(context.Context, UpdateTechnologyRepositoryCommand) error
}
type DeleteTechnologyRepository interface {
	DeleteTechnology(context.Context, DeleteTechnologyRepositoryCommand) error
}

type CreateExperienceRepositoryCommand struct {
	Experience *domain.Experience
	ResumeID   valueobjects.ID
	UserID     useridentity.ID
}
type GetExperienceRepositoryCommand struct {
	ResumeID, ExperienceID valueobjects.ID
	UserID                 useridentity.ID
}
type UpdateExperienceRepositoryCommand struct {
	Experience *domain.Experience
	ResumeID   valueobjects.ID
	UserID     useridentity.ID
}
type DeleteExperienceRepositoryCommand struct {
	ResumeID, ExperienceID valueobjects.ID
	UserID                 useridentity.ID
}
type CreateExperienceRepository interface {
	CreateExperience(context.Context, CreateExperienceRepositoryCommand) error
}
type GetExperienceRepository interface {
	GetExperience(context.Context, GetExperienceRepositoryCommand) (*domain.Experience, error)
}
type UpdateExperienceRepository interface {
	UpdateExperience(context.Context, UpdateExperienceRepositoryCommand) error
}
type DeleteExperienceRepository interface {
	DeleteExperience(context.Context, DeleteExperienceRepositoryCommand) error
}

type CreateEducationRepositoryCommand struct {
	Education *domain.Education
	ResumeID  valueobjects.ID
	UserID    useridentity.ID
}
type GetEducationRepositoryCommand struct {
	ResumeID, EducationID valueobjects.ID
	UserID                useridentity.ID
}
type UpdateEducationRepositoryCommand struct {
	Education *domain.Education
	ResumeID  valueobjects.ID
	UserID    useridentity.ID
}
type DeleteEducationRepositoryCommand struct {
	ResumeID, EducationID valueobjects.ID
	UserID                useridentity.ID
}
type CreateEducationRepository interface {
	CreateEducation(context.Context, CreateEducationRepositoryCommand) error
}
type GetEducationRepository interface {
	GetEducation(context.Context, GetEducationRepositoryCommand) (*domain.Education, error)
}
type UpdateEducationRepository interface {
	UpdateEducation(context.Context, UpdateEducationRepositoryCommand) error
}
type DeleteEducationRepository interface {
	DeleteEducation(context.Context, DeleteEducationRepositoryCommand) error
}

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

type CreateActivityRepositoryCommand struct {
	Activity *domain.Activity
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type GetActivityRepositoryCommand struct {
	ResumeID, ActivityID valueobjects.ID
	UserID               useridentity.ID
}
type UpdateActivityRepositoryCommand struct {
	Activity *domain.Activity
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type DeleteActivityRepositoryCommand struct {
	ResumeID, ActivityID valueobjects.ID
	UserID               useridentity.ID
}
type CreateActivityRepository interface {
	CreateActivity(context.Context, CreateActivityRepositoryCommand) error
}
type GetActivityRepository interface {
	GetActivity(context.Context, GetActivityRepositoryCommand) (*domain.Activity, error)
}
type UpdateActivityRepository interface {
	UpdateActivity(context.Context, UpdateActivityRepositoryCommand) error
}
type DeleteActivityRepository interface {
	DeleteActivity(context.Context, DeleteActivityRepositoryCommand) error
}

type CreateProjectRepositoryCommand struct {
	Project  *domain.Project
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type GetProjectRepositoryCommand struct {
	ResumeID, ProjectID valueobjects.ID
	UserID              useridentity.ID
}
type UpdateProjectRepositoryCommand struct {
	Project   *domain.Project
	ResumeID  valueobjects.ID
	ProjectID valueobjects.ID
	UserID    useridentity.ID
}
type DeleteProjectRepositoryCommand struct {
	ResumeID, ProjectID valueobjects.ID
	UserID              useridentity.ID
}
type CreateProjectRepository interface {
	CreateProject(context.Context, CreateProjectRepositoryCommand) error
}
type GetProjectRepository interface {
	GetProject(context.Context, GetProjectRepositoryCommand) (*domain.Project, error)
}
type UpdateProjectRepository interface {
	UpdateProject(context.Context, UpdateProjectRepositoryCommand) error
}
type DeleteProjectRepository interface {
	DeleteProject(context.Context, DeleteProjectRepositoryCommand) error
}

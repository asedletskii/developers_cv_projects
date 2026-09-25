//go:build wireinject

package app

import (
	"context"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"

	resumecommands "developers-cv/backend/internal/modules/resume/application/commands"
	resumeports "developers-cv/backend/internal/modules/resume/application/ports"
	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumereaders "developers-cv/backend/internal/modules/resume/infrastructure/persistence/readers"
	resumerepository "developers-cv/backend/internal/modules/resume/infrastructure/persistence/repository"
	resumehttp "developers-cv/backend/internal/modules/resume/transport/http"
	activitieshttp "developers-cv/backend/internal/modules/resume/transport/http/activities"
	courseshttp "developers-cv/backend/internal/modules/resume/transport/http/courses"
	educationshttp "developers-cv/backend/internal/modules/resume/transport/http/educations"
	experienceshttp "developers-cv/backend/internal/modules/resume/transport/http/experiences"
	projectshttp "developers-cv/backend/internal/modules/resume/transport/http/projects"
	publichandler "developers-cv/backend/internal/modules/resume/transport/http/public"
	userqueries "developers-cv/backend/internal/modules/user/application/queries"
	userauth "developers-cv/backend/internal/modules/user/infrastructure/auth"
	userreaders "developers-cv/backend/internal/modules/user/infrastructure/persistence/readers"
	userhttp "developers-cv/backend/internal/modules/user/transport/http"

	"developers-cv/backend/internal/platform/config"
	"developers-cv/backend/internal/platform/database"
	httpserver "developers-cv/backend/internal/platform/http"
	"developers-cv/backend/internal/platform/logger"
	"developers-cv/backend/internal/platform/telemetry"
)

func Initialize() (*Application, error) {
	panic(wire.Build(
		config.Load,

		func(cfg config.Config) (*pgxpool.Pool, error) {
			return newPostgresPool(context.Background(), cfg)
		},

		logger.New,
		newHealthService,
		telemetry.New,
		telemetry.NewMetrics,
		newPublicResumeReader,
		newListResumesReader,
		resumereaders.NewResumeRootReader,
		resumereaders.NewProfileReader,
		resumereaders.NewContactsReader,
		resumereaders.NewSkillsReader,
		resumereaders.NewLanguagesReader,
		resumereaders.NewTechnologiesReader,
		resumereaders.NewExperiencesReader,
		resumereaders.NewEducationReader,
		resumereaders.NewCourseReader,
		resumereaders.NewActivityReader,
		resumereaders.NewProjectsReader,
		database.NewGORMFromConfig,
		newGetProfileReader,
		database.NewUnitOfWork,
		resumerepository.NewResumeRootRepository,
		resumerepository.NewProfileRepository,
		resumerepository.NewContactsRepository,
		resumerepository.NewSkillsRepository,
		resumerepository.NewLanguagesRepository,
		resumerepository.NewTechnologiesRepository,
		resumerepository.NewExperiencesRepository,
		resumerepository.NewEducationRepository,
		resumerepository.NewCourseRepository,
		resumerepository.NewActivityRepository,
		resumerepository.NewProjectRepository,
		wire.Bind(new(resumeports.CreateResumeRootRepository), new(*resumerepository.ResumeRootRepository)),
		wire.Bind(new(resumeports.GetResumeRootRepository), new(*resumerepository.ResumeRootRepository)),
		wire.Bind(new(resumeports.UpdateResumeRootRepository), new(*resumerepository.ResumeRootRepository)),
		wire.Bind(new(resumeports.DeleteResumeRootRepository), new(*resumerepository.ResumeRootRepository)),
		wire.Bind(new(resumeports.CreateResumeProfileRepository), new(*resumerepository.ProfileRepository)),
		wire.Bind(new(resumeports.GetResumeProfileRepository), new(*resumerepository.ProfileRepository)),
		wire.Bind(new(resumeports.UpdateResumeProfileRepository), new(*resumerepository.ProfileRepository)),
		wire.Bind(new(resumeports.CreateResumeContactsRepository), new(*resumerepository.ContactsRepository)),
		wire.Bind(new(resumeports.GetResumeContactsRepository), new(*resumerepository.ContactsRepository)),
		wire.Bind(new(resumeports.UpdateResumeContactsRepository), new(*resumerepository.ContactsRepository)),
		wire.Bind(new(resumeports.CreateLanguageRepository), new(*resumerepository.LanguagesRepository)),
		wire.Bind(new(resumeports.GetLanguageRepository), new(*resumerepository.LanguagesRepository)),
		wire.Bind(new(resumeports.UpdateLanguageRepository), new(*resumerepository.LanguagesRepository)),
		wire.Bind(new(resumeports.DeleteLanguageRepository), new(*resumerepository.LanguagesRepository)),
		wire.Bind(new(resumeports.CreateTechnologyRepository), new(*resumerepository.TechnologiesRepository)),
		wire.Bind(new(resumeports.GetTechnologyRepository), new(*resumerepository.TechnologiesRepository)),
		wire.Bind(new(resumeports.UpdateTechnologyRepository), new(*resumerepository.TechnologiesRepository)),
		wire.Bind(new(resumeports.DeleteTechnologyRepository), new(*resumerepository.TechnologiesRepository)),
		wire.Bind(new(resumeports.CreateExperienceRepository), new(*resumerepository.ExperiencesRepository)),
		wire.Bind(new(resumeports.GetExperienceRepository), new(*resumerepository.ExperiencesRepository)),
		wire.Bind(new(resumeports.UpdateExperienceRepository), new(*resumerepository.ExperiencesRepository)),
		wire.Bind(new(resumeports.DeleteExperienceRepository), new(*resumerepository.ExperiencesRepository)),
		wire.Bind(new(resumeports.CreateEducationRepository), new(*resumerepository.EducationRepository)),
		wire.Bind(new(resumeports.GetEducationRepository), new(*resumerepository.EducationRepository)),
		wire.Bind(new(resumeports.UpdateEducationRepository), new(*resumerepository.EducationRepository)),
		wire.Bind(new(resumeports.DeleteEducationRepository), new(*resumerepository.EducationRepository)),
		wire.Bind(new(resumeports.CreateCourseRepository), new(*resumerepository.CourseRepository)),
		wire.Bind(new(resumeports.GetCourseRepository), new(*resumerepository.CourseRepository)),
		wire.Bind(new(resumeports.UpdateCourseRepository), new(*resumerepository.CourseRepository)),
		wire.Bind(new(resumeports.DeleteCourseRepository), new(*resumerepository.CourseRepository)),
		wire.Bind(new(resumeports.CreateActivityRepository), new(*resumerepository.ActivityRepository)),
		wire.Bind(new(resumeports.GetActivityRepository), new(*resumerepository.ActivityRepository)),
		wire.Bind(new(resumeports.UpdateActivityRepository), new(*resumerepository.ActivityRepository)),
		wire.Bind(new(resumeports.DeleteActivityRepository), new(*resumerepository.ActivityRepository)),
		wire.Bind(new(resumeports.CreateProjectRepository), new(*resumerepository.ProjectRepository)),
		wire.Bind(new(resumeports.GetProjectRepository), new(*resumerepository.ProjectRepository)),
		wire.Bind(new(resumeports.UpdateProjectRepository), new(*resumerepository.ProjectRepository)),
		wire.Bind(new(resumeports.DeleteProjectRepository), new(*resumerepository.ProjectRepository)),
		wire.Bind(new(resumequeries.ResumeRootReader), new(*resumereaders.ResumeRootReader)),
		wire.Bind(new(resumequeries.ResumeProfileReader), new(*resumereaders.ProfileReader)),
		wire.Bind(new(resumequeries.ResumeContactsReader), new(*resumereaders.ContactsReader)),
		wire.Bind(new(resumequeries.SkillsReader), new(*resumereaders.SkillsReader)),
		wire.Bind(new(resumequeries.LanguagesReader), new(*resumereaders.LanguagesReader)),
		wire.Bind(new(resumequeries.TechnologiesReader), new(*resumereaders.TechnologiesReader)),
		wire.Bind(new(resumequeries.ExperiencesReader), new(*resumereaders.ExperiencesReader)),
		wire.Bind(new(resumequeries.EducationsReader), new(*resumereaders.EducationReader)),
		wire.Bind(new(resumequeries.CoursesReader), new(*resumereaders.CourseReader)),
		wire.Bind(new(resumequeries.ActivitiesReader), new(*resumereaders.ActivityReader)),
		wire.Bind(new(resumequeries.ProjectsReader), new(*resumereaders.ProjectsReader)),
		wire.Bind(new(resumeports.CreateSkillsRepository), new(*resumerepository.SkillsRepository)),
		wire.Bind(new(resumeports.GetSkillsRepository), new(*resumerepository.SkillsRepository)),
		wire.Bind(new(resumeports.UpdateSkillsRepository), new(*resumerepository.SkillsRepository)),
		wire.Bind(new(resumecommands.UnitOfWork), new(*database.UnitOfWork)),
		userreaders.NewUserReader,
		wire.Bind(new(userqueries.UserReader), new(*userreaders.UserReader)),
		userauth.NewMockAuthenticator,
		newAuthenticate,
		resumequeries.NewPublicResumeQueryHandler,
		resumequeries.NewListResumesQueryHandler,
		resumequeries.NewGetResumeRootQueryHandler,
		resumequeries.NewGetResumeProfileQueryHandler,
		resumequeries.NewGetResumeContactsQueryHandler,
		resumequeries.NewGetSkillsQueryHandler,
		userqueries.NewGetProfileHandler,
		userhttp.NewHandler,
		resumecommands.NewCreateResumeRootHandler,
		resumecommands.NewUpdateResumeRootHandler,
		resumecommands.NewDeleteResumeRootHandler,
		resumecommands.NewCreateResumeProfileHandler,
		resumecommands.NewUpdateResumeProfileHandler,
		resumecommands.NewCreateResumeContactsHandler,
		resumecommands.NewUpdateResumeContactsHandler,
		resumecommands.NewCreateSkillsHandler,
		resumecommands.NewUpdateSkillsHandler,
		resumecommands.NewCreateLanguageHandler,
		resumequeries.NewListLanguagesQueryHandler,
		resumequeries.NewGetLanguageQueryHandler,
		resumecommands.NewUpdateLanguageHandler,
		resumecommands.NewDeleteLanguageHandler,
		resumecommands.NewCreateTechnologyHandler,
		resumequeries.NewListTechnologiesQueryHandler,
		resumequeries.NewGetTechnologyQueryHandler,
		resumecommands.NewUpdateTechnologyHandler,
		resumecommands.NewDeleteTechnologyHandler,
		resumecommands.NewCreateExperienceHandler,
		resumequeries.NewListExperiencesQueryHandler,
		resumequeries.NewGetExperienceQueryHandler,
		resumecommands.NewUpdateExperienceHandler,
		resumecommands.NewDeleteExperienceHandler,
		resumecommands.NewCreateEducationHandler,
		resumequeries.NewListEducationsQueryHandler,
		resumequeries.NewGetEducationQueryHandler,
		resumecommands.NewUpdateEducationHandler,
		resumecommands.NewDeleteEducationHandler,
		resumecommands.NewCreateCourseHandler,
		resumequeries.NewListCoursesQueryHandler,
		resumequeries.NewGetCourseQueryHandler,
		resumecommands.NewUpdateCourseHandler,
		resumecommands.NewDeleteCourseHandler,
		resumecommands.NewCreateActivityHandler,
		resumequeries.NewListActivitiesQueryHandler,
		resumequeries.NewGetActivityQueryHandler,
		resumecommands.NewUpdateActivityHandler,
		resumecommands.NewDeleteActivityHandler,
		resumecommands.NewCreateProjectHandler,
		resumequeries.NewListProjectsQueryHandler,
		resumequeries.NewGetProjectQueryHandler,
		resumecommands.NewUpdateProjectHandler,
		resumecommands.NewDeleteProjectHandler,
		experienceshttp.NewHandler,
		educationshttp.NewHandler,
		courseshttp.NewHandler,
		activitieshttp.NewHandler,
		projectshttp.NewHandler,
		publichandler.NewHandler,
		resumehttp.NewHandler,
		httpserver.NewRouter,
		httpserver.NewServer,
		NewApplication,
	))
}

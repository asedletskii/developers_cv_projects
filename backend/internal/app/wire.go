//go:build wireinject

package app

import (
	"context"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"

	resumecommands "developers-cv/backend/internal/modules/resume/application/commands"
	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
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
		wire.Bind(new(resumecommands.CreateResumeRootRepository), new(*resumerepository.ResumeRootRepository)),
		wire.Bind(new(resumecommands.GetResumeRootRepository), new(*resumerepository.ResumeRootRepository)),
		wire.Bind(new(resumecommands.UpdateResumeRootRepository), new(*resumerepository.ResumeRootRepository)),
		wire.Bind(new(resumecommands.DeleteResumeRootRepository), new(*resumerepository.ResumeRootRepository)),
		wire.Bind(new(resumecommands.CreateResumeProfileRepository), new(*resumerepository.ProfileRepository)),
		wire.Bind(new(resumecommands.GetResumeProfileRepository), new(*resumerepository.ProfileRepository)),
		wire.Bind(new(resumecommands.UpdateResumeProfileRepository), new(*resumerepository.ProfileRepository)),
		wire.Bind(new(resumecommands.CreateResumeContactsRepository), new(*resumerepository.ContactsRepository)),
		wire.Bind(new(resumecommands.GetResumeContactsRepository), new(*resumerepository.ContactsRepository)),
		wire.Bind(new(resumecommands.UpdateResumeContactsRepository), new(*resumerepository.ContactsRepository)),
		wire.Bind(new(resumecommands.SkillsRepository), new(*resumerepository.SkillsRepository)),
		wire.Bind(new(resumecommands.LanguagesRepository), new(*resumerepository.LanguagesRepository)),
		wire.Bind(new(resumecommands.TechnologiesRepository), new(*resumerepository.TechnologiesRepository)),
		wire.Bind(new(resumecommands.ExperiencesRepository), new(*resumerepository.ExperiencesRepository)),
		wire.Bind(new(resumecommands.EducationsRepository), new(*resumerepository.EducationRepository)),
		wire.Bind(new(resumecommands.CoursesRepository), new(*resumerepository.CourseRepository)),
		wire.Bind(new(resumecommands.ActivitiesRepository), new(*resumerepository.ActivityRepository)),
		wire.Bind(new(resumecommands.ProjectsRepository), new(*resumerepository.ProjectRepository)),
		wire.Bind(new(resumecommands.UnitOfWork), new(*database.UnitOfWork)),
		userreaders.NewUserReader,
		wire.Bind(new(userqueries.UserReader), new(*userreaders.UserReader)),
		userauth.NewMockAuthenticator,
		newAuthenticate,
		resumequeries.NewPublicResumeQueryHandler,
		resumequeries.NewListResumesHandler,
		userqueries.NewGetProfileHandler,
		userhttp.NewHandler,
		resumecommands.NewCreateResumeRootHandler,
		resumecommands.NewGetResumeRootHandler,
		resumecommands.NewUpdateResumeRootHandler,
		resumecommands.NewDeleteResumeRootHandler,
		resumecommands.NewCreateResumeProfileHandler,
		resumecommands.NewGetResumeProfileHandler,
		resumecommands.NewUpdateResumeProfileHandler,
		resumecommands.NewCreateResumeContactsHandler,
		resumecommands.NewGetResumeContactsHandler,
		resumecommands.NewUpdateResumeContactsHandler,
		resumecommands.NewCreateSkillsHandler,
		resumecommands.NewGetSkillsHandler,
		resumecommands.NewUpdateSkillsHandler,
		resumecommands.NewCreateLanguageHandler,
		resumecommands.NewListLanguagesHandler,
		resumecommands.NewGetLanguageHandler,
		resumecommands.NewUpdateLanguageHandler,
		resumecommands.NewDeleteLanguageHandler,
		resumecommands.NewCreateTechnologyHandler,
		resumecommands.NewListTechnologiesHandler,
		resumecommands.NewGetTechnologyHandler,
		resumecommands.NewUpdateTechnologyHandler,
		resumecommands.NewDeleteTechnologyHandler,
		resumecommands.NewCreateExperienceHandler,
		resumecommands.NewListExperiencesHandler,
		resumecommands.NewGetExperienceHandler,
		resumecommands.NewUpdateExperienceHandler,
		resumecommands.NewDeleteExperienceHandler,
		resumecommands.NewCreateEducationHandler,
		resumecommands.NewListEducationsHandler,
		resumecommands.NewGetEducationHandler,
		resumecommands.NewUpdateEducationHandler,
		resumecommands.NewDeleteEducationHandler,
		resumecommands.NewCreateCourseHandler,
		resumecommands.NewListCoursesHandler,
		resumecommands.NewGetCourseHandler,
		resumecommands.NewUpdateCourseHandler,
		resumecommands.NewDeleteCourseHandler,
		resumecommands.NewCreateActivityHandler,
		resumecommands.NewListActivitiesHandler,
		resumecommands.NewGetActivityHandler,
		resumecommands.NewUpdateActivityHandler,
		resumecommands.NewDeleteActivityHandler,
		resumecommands.NewCreateProjectHandler,
		resumecommands.NewListProjectsHandler,
		resumecommands.NewGetProjectHandler,
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

//go:build wireinject

package app

import (
	"github.com/google/wire"

	resumecommands "developers-cv/backend/internal/modules/resume/application/commands"
	resumeports "developers-cv/backend/internal/modules/resume/application/ports"
	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumereaders "developers-cv/backend/internal/modules/resume/infrastructure/persistence/readers"
	resumerepository "developers-cv/backend/internal/modules/resume/infrastructure/persistence/repository"
	resumehttp "developers-cv/backend/internal/modules/resume/transport/http"
	activitieshttp "developers-cv/backend/internal/modules/resume/transport/http/activities"
	contactshandler "developers-cv/backend/internal/modules/resume/transport/http/contacts"
	courseshttp "developers-cv/backend/internal/modules/resume/transport/http/courses"
	educationshttp "developers-cv/backend/internal/modules/resume/transport/http/educations"
	experienceshttp "developers-cv/backend/internal/modules/resume/transport/http/experiences"
	languageshttp "developers-cv/backend/internal/modules/resume/transport/http/languages"
	overviewhandler "developers-cv/backend/internal/modules/resume/transport/http/overview"
	profilehandler "developers-cv/backend/internal/modules/resume/transport/http/profile"
	projectshttp "developers-cv/backend/internal/modules/resume/transport/http/projects"
	publichandler "developers-cv/backend/internal/modules/resume/transport/http/public"
	skillshttp "developers-cv/backend/internal/modules/resume/transport/http/skills"
	technologieshttp "developers-cv/backend/internal/modules/resume/transport/http/technologies"
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

var resumeReadSet = wire.NewSet(
	newListResumesReader,
	resumereaders.NewResumeReader,
	newGetProfileReader,
	wire.Bind(new(resumequeries.ResumeReader), new(*resumereaders.ResumeReader)),
	wire.Bind(new(resumequeries.PublicResumeReader), new(*resumereaders.ResumeReader)),
	wire.Bind(new(resumequeries.ContactViewReader), new(*resumereaders.ResumeReader)),
)

var resumeCommandSet = wire.NewSet(
	resumecommands.NewCreateResumeHandler,
	resumecommands.NewReplaceResumeTitleHandler,
	resumecommands.NewReplaceOverviewHandler,
	resumecommands.NewReplaceProfileHandler,
	resumecommands.NewDeleteResumeHandler,
	resumecommands.NewReplaceContactsHandler,
	resumecommands.NewReplaceSkillsHandler,
	resumecommands.NewReplaceLanguagesHandler,
	resumecommands.NewReplaceTechnologiesHandler,
	resumecommands.NewReplaceEducationsHandler,
	resumecommands.NewReplaceCoursesHandler,
	resumecommands.NewReplaceActivitiesHandler,
	resumecommands.NewReplaceProjectsHandler,
	wire.Bind(new(resumecommands.UnitOfWork), new(*database.UnitOfWork)),
)

var resumeHTTPSet = wire.NewSet(
	activitieshttp.NewHandler,
	contactshandler.NewHandler,
	educationshttp.NewHandler,
	courseshttp.NewHandler,
	experienceshttp.NewHandler,
	projectshttp.NewHandler,
	skillshttp.NewHandler,
	languageshttp.NewHandler,
	technologieshttp.NewHandler,
	overviewhandler.NewHandler,
	profilehandler.NewHandler,
	resumehttp.NewHandler,
	publichandler.NewHandler,
	publichandler.NewContactViewHandler,
	wire.Struct(new(resumehttp.HandlerDependencies), "*"),
)

func Initialize() (*Application, error) {
	panic(wire.Build(
		config.Load,
		logger.New,
		newHealthService,
		telemetry.New,
		telemetry.NewMetrics,

		database.NewGORMFromConfig,
		database.NewUnitOfWork,
		resumerepository.NewResumeRepository,
		wire.Bind(new(resumeports.ResumeRepository), new(*resumerepository.ResumeRepository)),

		resumeReadSet,
		resumeCommandSet,
		resumeHTTPSet,

		resumequeries.NewPublicResumeQueryHandler,
		resumequeries.NewListResumesQueryHandler,
		resumequeries.NewGetResumeQueryHandler,
		resumequeries.NewContactViewQueryHandler,

		userreaders.NewUserReader,
		wire.Bind(new(userqueries.UserReader), new(*userreaders.UserReader)),
		userauth.NewMockAuthenticator,
		newAuthenticate,
		userqueries.NewGetProfileHandler,
		userhttp.NewHandler,

		httpserver.NewRouter,
		httpserver.NewServer,
		NewApplication,
	))
}

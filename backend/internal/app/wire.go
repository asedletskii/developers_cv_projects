//go:build wireinject

package app

import (
	"github.com/google/wire"

	authcommands "developers-cv/backend/internal/modules/auth/application/commands"
	authports "developers-cv/backend/internal/modules/auth/application/ports"
	authqueries "developers-cv/backend/internal/modules/auth/application/queries"
	cleanup "developers-cv/backend/internal/modules/auth/infrastructure/cleanup"
	authrepository "developers-cv/backend/internal/modules/auth/infrastructure/persistence/repository"
	authtoken "developers-cv/backend/internal/modules/auth/infrastructure/token"
	authhttp "developers-cv/backend/internal/modules/auth/transport/http"
	resumecommands "developers-cv/backend/internal/modules/resume/application/commands"
	resumeports "developers-cv/backend/internal/modules/resume/application/ports"
	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumereaders "developers-cv/backend/internal/modules/resume/infrastructure/persistence/readers"
	resumerepository "developers-cv/backend/internal/modules/resume/infrastructure/persistence/repository"
	resumehttp "developers-cv/backend/internal/modules/resume/transport/http"
	resumelinkcommands "developers-cv/backend/internal/modules/resume_link/application/commands"
	resumelinkports "developers-cv/backend/internal/modules/resume_link/application/ports"
	resumelinkqueries "developers-cv/backend/internal/modules/resume_link/application/queries"
	resumelinkanalytics "developers-cv/backend/internal/modules/resume_link/infrastructure/analytics"
	resumelinkreaders "developers-cv/backend/internal/modules/resume_link/infrastructure/persistence/readers"
	resumelinkrepository "developers-cv/backend/internal/modules/resume_link/infrastructure/persistence/repository"
	resumelinkhttp "developers-cv/backend/internal/modules/resume_link/transport/http"
	resumelinkpublichandler "developers-cv/backend/internal/modules/resume_link/transport/http/public"
	usercommands "developers-cv/backend/internal/modules/user/application/commands"
	userqueries "developers-cv/backend/internal/modules/user/application/queries"
	userreaders "developers-cv/backend/internal/modules/user/infrastructure/persistence/readers"
	userrepository "developers-cv/backend/internal/modules/user/infrastructure/persistence/repository"
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
	resumecommands.NewChangeResumeStatusHandler,
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
	wire.Bind(new(resumeports.UnitOfWork), new(*database.UnitOfWork)),
)

var resumeLinkReadSet = wire.NewSet(
	resumelinkreaders.NewResumeLinkReader,
	resumelinkqueries.NewPublicResumeLinkQueryHandler,
	resumelinkqueries.NewGetResumeLinkStatisticsQueryHandler,
	wire.Bind(new(resumelinkqueries.ResumeLinkReader), new(*resumelinkreaders.ResumeLinkReader)),
	wire.Bind(new(resumelinkqueries.ResumeLinkPublicReader), new(*resumelinkreaders.ResumeLinkReader)),
	wire.Bind(new(resumelinkqueries.ResumeLinkStatisticsReader), new(*resumelinkreaders.ResumeLinkReader)),
	wire.Bind(new(resumelinkqueries.ResumeOwnershipReader), new(*resumereaders.ResumeReader)),
	wire.Bind(new(resumelinkqueries.ResumePublicQuery), new(*resumequeries.PublicResumeQueryHandler)),
)

var resumeLinkAnalyticsSet = wire.NewSet(
	resumelinkanalytics.NewRecorder,
	wire.Bind(new(resumelinkports.LinkEventRecorder), new(*resumelinkanalytics.Recorder)),
)

var resumeLinkCommandSet = wire.NewSet(
	resumelinkcommands.NewCreateResumeLinkHandler,
	resumelinkcommands.NewDeleteResumeLinkHandler,
	resumelinkcommands.NewReplaceStatusResumeLinkHandler,
	wire.Bind(new(resumelinkcommands.UnitOfWork), new(*database.UnitOfWork)),
)

var resumeHTTPSet = wire.NewSet(
	resumehttp.NewHandler,
	wire.Struct(new(resumehttp.HandlerDependencies), "*"),
)

var resumeLinkHTTPSet = wire.NewSet(
	resumelinkpublichandler.NewHandler,
	resumelinkpublichandler.NewContactViewHandler,
	resumelinkhttp.NewHandler,
	wire.Struct(new(resumelinkhttp.HandlerDependencies), "*"),
)

var authRepositorySet = wire.NewSet(
	authrepository.NewAuthTransactionRepository,
	authrepository.NewIdentityRepository,
	authrepository.NewSessionRepository,
	wire.Bind(new(authports.AuthTransactionRepository), new(*authrepository.AuthTransactionRepository)),
	wire.Bind(new(authports.IdentityRepository), new(*authrepository.IdentityRepository)),
	wire.Bind(new(authports.SessionRepository), new(*authrepository.SessionRepository)),
)

var authSet = wire.NewSet(
	authRepositorySet,
	authtoken.NewGenerator,
	wire.Bind(new(authports.TokenGenerator), new(*authtoken.Generator)),
	usercommands.NewCreateUserHandler,
	wire.Bind(new(authports.UserProvisioner), new(*usercommands.CreateUserHandler)),
	authqueries.NewCurrentSessionHandler,
	authcommands.NewStartTelegramLoginHandler,
	newCompleteTelegramLoginHandler,
	authcommands.NewLogoutHandler,
	authcommands.NewCleanupExpiredHandler,
	authhttp.NewHandler,
	wire.Struct(new(authhttp.HandlerDependencies), "*"),
)

func Initialize() (*Application, error) {
	panic(wire.Build(
		config.Load, logger.New, newHealthService, telemetry.New, telemetry.NewMetrics,
		database.NewGORMFromConfig, database.NewUnitOfWork,
		resumerepository.NewResumeRepository,
		wire.Bind(new(resumeports.ResumeRepository), new(*resumerepository.ResumeRepository)),
		resumelinkrepository.NewResumeLinkRepository,
		wire.Bind(new(resumelinkports.ResumeLinkRepository), new(*resumelinkrepository.ResumeLinkRepository)),
		wire.Bind(new(resumelinkports.ResumeOwnershipReader), new(*resumereaders.ResumeReader)),
		resumeReadSet, resumeCommandSet, resumeLinkReadSet, resumeLinkCommandSet, resumeLinkAnalyticsSet,
		resumeHTTPSet, resumeLinkHTTPSet,
		resumequeries.NewPublicResumeQueryHandler, resumequeries.NewListResumesQueryHandler,
		resumequeries.NewGetResumeQueryHandler, resumequeries.NewContactViewQueryHandler,
		userreaders.NewUserReader,
		wire.Bind(new(userqueries.UserReader), new(*userreaders.UserReader)),
		userrepository.NewUserRepository,
		userrepository.NewProfileRepository,
		usercommands.NewReplaceProfileHandler,
		userqueries.NewGetProfileHandler,
		userhttp.NewHandler,
		authSet,
		newTelegramHTTPClient,
		newTelegramOIDC,
		newAuthObserver,
		newRuntimeAuthenticator,
		newAuthenticate,
		cookieConfigFromRuntime,
		httpserver.NewRouter, httpserver.NewServer, cleanup.NewWorker, NewApplication,
	))
}

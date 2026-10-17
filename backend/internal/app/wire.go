//go:build wireinject

package app

import (
	"github.com/google/wire"

	authcommands "developers-cv/backend/internal/modules/auth/application/commands"
	authports "developers-cv/backend/internal/modules/auth/application/ports"
	cleanup "developers-cv/backend/internal/modules/auth/infrastructure/cleanup"
	authrepository "developers-cv/backend/internal/modules/auth/infrastructure/persistence/repository"
	authtoken "developers-cv/backend/internal/modules/auth/infrastructure/token"
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
	resumelinkcleanup "developers-cv/backend/internal/modules/resume_link/infrastructure/cleanup"
	resumelinkreaders "developers-cv/backend/internal/modules/resume_link/infrastructure/persistence/readers"
	resumelinkrepository "developers-cv/backend/internal/modules/resume_link/infrastructure/persistence/repository"
	resumelinkhttp "developers-cv/backend/internal/modules/resume_link/transport/http"
	resumelinkhandlers "developers-cv/backend/internal/modules/resume_link/transport/http/handlers"
	usercommands "developers-cv/backend/internal/modules/user/application/commands"
	userports "developers-cv/backend/internal/modules/user/application/ports"
	userqueries "developers-cv/backend/internal/modules/user/application/queries"
	userreaders "developers-cv/backend/internal/modules/user/infrastructure/persistence/readers"
	userrepository "developers-cv/backend/internal/modules/user/infrastructure/persistence/repository"
	userhttp "developers-cv/backend/internal/modules/user/transport/http"
	"developers-cv/backend/internal/platform/config"
	"developers-cv/backend/internal/platform/database"
)

var resumeReadSet = wire.NewSet(
	newListResumesReader,
	resumereaders.NewResumeReader,
	newGetProfileReader,
	wire.Bind(new(resumequeries.ResumeReader), new(*resumereaders.ResumeReader)),
	wire.Bind(new(resumequeries.PublicResumeReader), new(*resumereaders.ResumeReader)),
	wire.Bind(new(resumequeries.ContactViewReader), new(*resumereaders.ResumeReader)),
	wire.Bind(new(resumequeries.DownloadContactsReader), new(*resumereaders.ResumeReader)),
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
	resumecommands.NewReplaceExperiencesHandler,
)

var resumeLinkReadSet = wire.NewSet(
	resumelinkreaders.NewResumeLinkReader,
	resumelinkqueries.NewListResumeLinksQueryHandler,
	resumelinkqueries.NewPublicResumeLinkQueryHandlerWithVisit,
	resumelinkqueries.NewLinkStatisticsHandler,
	resumelinkqueries.NewDailyStatisticsHandler,
	resumelinkqueries.NewVisitsHandler,
	resumelinkqueries.NewSectionsHandler,
	resumelinkqueries.NewTechnologiesHandler,
	resumelinkqueries.NewResumeStatisticsHandler,
	resumelinkqueries.NewContactViewThroughLinkQueryHandlerWithAnalytics,
	wire.Bind(new(resumelinkqueries.ResumeLinkReader), new(*resumelinkreaders.ResumeLinkReader)),
	wire.Bind(new(resumelinkqueries.ResumeLinkPublicReader), new(*resumelinkreaders.ResumeLinkReader)),
	wire.Bind(new(resumelinkqueries.LinkStatisticsReader), new(*resumelinkreaders.ResumeLinkReader)),
	wire.Bind(new(resumelinkqueries.DailyStatisticsReader), new(*resumelinkreaders.ResumeLinkReader)),
	wire.Bind(new(resumelinkqueries.VisitsReader), new(*resumelinkreaders.ResumeLinkReader)),
	wire.Bind(new(resumelinkqueries.SectionsReader), new(*resumelinkreaders.ResumeLinkReader)),
	wire.Bind(new(resumelinkqueries.TechnologiesReader), new(*resumelinkreaders.ResumeLinkReader)),
	wire.Bind(new(resumelinkqueries.ResumeStatisticsReader), new(*resumelinkreaders.ResumeLinkReader)),
)

var resumeLinkCommandSet = wire.NewSet(
	resumelinkcommands.NewCreatePublicVisitHandler,
	resumelinkcommands.NewRecordPublicInteractionEventHandler,
	newIngestAnalyticsBatchHandler,
	newCleanupAnalyticsHandler,
	resumelinkcleanup.NewWorker,
	resumelinkqueries.NewAnalyticsClickQueryHandler,
	resumelinkhandlers.NewResumeLinkAnalyticsEventsHandler,
	resumelinkanalytics.NewLayoutSnapshotBuilder,
	resumelinkcommands.NewCreateResumeLinkHandler,
	resumelinkcommands.NewDeleteResumeLinkHandler,
	resumelinkcommands.NewReplaceStatusResumeLinkHandler,
	wire.Bind(new(resumelinkports.PublicVisitCreator), new(*resumelinkcommands.CreatePublicVisitHandler)),
	wire.Bind(new(resumelinkports.PublicInteractionEventRecorder), new(*resumelinkcommands.RecordPublicInteractionEventHandler)),
)

var resumeHTTPSet = wire.NewSet(
	resumehttp.NewHandler,
	wire.Struct(new(resumehttp.HandlerDependencies), "*"),
)

var resumeLinkHTTPSet = wire.NewSet(
	resumelinkhandlers.NewPublicResumeLinkHandler,
	resumelinkhandlers.NewPublicContactViewHandler,
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
	newCurrentSessionHandler,
	newStartTelegramLoginHandler,
	newCompleteTelegramLoginHandler,
	authcommands.NewLogoutHandler,
	authcommands.NewCleanupExpiredHandler,
	newAuthHandler,
)

func initialize() (*Application, func(), error) {
	panic(wire.Build(
		config.Load, wire.FieldsOf(new(config.Config), "Auth"), newLogger, newHealthService, newTelemetry, newMetrics,
		newContactViewAnalyticsErrorReporter, newPublicAnalyticsErrorReporter,
		newDatabase, database.NewUnitOfWork,
		wire.Bind(new(resumeports.UnitOfWork), new(*database.UnitOfWork)),
		wire.Bind(new(userports.UnitOfWork), new(*database.UnitOfWork)),
		wire.Bind(new(resumelinkports.UnitOfWork), new(*database.UnitOfWork)),
		wire.Bind(new(authports.UnitOfWork), new(*database.UnitOfWork)),
		wire.Bind(new(resumelinkqueries.ResumePublicQuery), new(*resumequeries.PublicResumeQueryHandler)),
		wire.Bind(new(resumelinkqueries.ResumeContactViewQuery), new(*resumequeries.ContactViewQueryHandler)),
		wire.Bind(new(resumelinkqueries.ResumeContactsDownloadQuery), new(*resumequeries.DownloadContactsQueryHandler)),
		wire.Bind(new(resumelinkports.VisitorIdentity), new(*resumelinkanalytics.VisitorIdentifier)),
		wire.Bind(new(resumelinkports.TrafficClassifier), new(*resumelinkanalytics.TrafficClassifier)),
		wire.Bind(new(resumelinkhandlers.PublicResumeLinkQuery), new(*resumelinkqueries.PublicResumeLinkQueryHandler)),
		wire.Bind(new(resumelinkhandlers.ContactViewQuery), new(*resumelinkqueries.ContactViewThroughLinkQueryHandler)),
		wire.Bind(new(resumelinkqueries.AnalyticsClickReader), new(*resumelinkreaders.ResumeLinkReader)),
		wire.Bind(new(resumelinkports.LayoutSnapshotBuilder), new(*resumelinkanalytics.LayoutSnapshotBuilder)),
		resumerepository.NewResumeRepository,
		wire.Bind(new(resumeports.ResumeRepository), new(*resumerepository.ResumeRepository)),
		resumelinkrepository.NewResumeLinkRepository,
		newAnalyticsTrafficClassifier,
		newAnalyticsVisitorIdentifier,
		newOptionalAuthenticate,
		wire.Bind(new(resumelinkports.ResumeLinkRepository), new(*resumelinkrepository.ResumeLinkRepository)),
		wire.Bind(new(resumelinkports.ViewRepository), new(*resumelinkrepository.ResumeLinkRepository)),
		wire.Bind(new(resumelinkports.AnalyticsEventRepository), new(*resumelinkrepository.ResumeLinkRepository)),
		wire.Bind(new(resumelinkports.AnalyticsBatchRepository), new(*resumelinkrepository.ResumeLinkRepository)),
		wire.Bind(new(resumelinkports.AnalyticsRetentionRepository), new(*resumelinkrepository.ResumeLinkRepository)),
		wire.Bind(new(resumelinkports.ResumeOwnershipReader), new(*resumereaders.ResumeReader)),
		wire.Bind(new(resumelinkqueries.ResumeOverviewReader), new(*resumereaders.ResumeReader)),
		resumeReadSet, resumeCommandSet, resumeLinkReadSet, resumeLinkCommandSet,
		resumeHTTPSet, resumeLinkHTTPSet,
		resumequeries.NewPublicResumeQueryHandler, resumequeries.NewListResumesQueryHandler,
		resumequeries.NewGetResumeQueryHandler, resumequeries.NewContactViewQueryHandler,
		resumequeries.NewDownloadContactsQueryHandler,
		userreaders.NewUserReader,
		wire.Bind(new(userqueries.UserReader), new(*userreaders.UserReader)),
		userrepository.NewUserRepository,
		wire.Bind(new(userports.UserRepository), new(*userrepository.UserRepository)),
		userrepository.NewProfileRepository,
		wire.Bind(new(userports.ProfileRepository), new(*userrepository.ProfileRepository)),
		usercommands.NewReplaceProfileHandler,
		userqueries.NewGetProfileHandler,
		userhttp.NewHandler,
		wire.Struct(new(userhttp.HandlerDependencies), "*"),
		authSet,
		newTurnstileVerifier,
		newTelegramHTTPClient,
		newTelegramOIDC,
		newAuthObserver,
		newRuntimeAuthenticator,
		newAuthenticate,
		newRouter, newHTTPHandler, newHTTPServer, cleanup.NewWorker, newBackgroundWorkers, newApplication,
	))
}

package app

import (
	"gorm.io/gorm"

	analyticscommands "developers-cv/backend/internal/modules/analytics/application/commands"
	analyticsports "developers-cv/backend/internal/modules/analytics/application/ports"
	analyticscontracts "developers-cv/backend/internal/modules/analytics/contracts"
	analyticsreaders "developers-cv/backend/internal/modules/analytics/infrastructure/persistence/readers"
	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
	"developers-cv/backend/internal/modules/resume_link/infrastructure/analyticsaccess"
	"developers-cv/backend/internal/platform/config"
	"developers-cv/backend/internal/platform/telemetry"
)

func newAnalyticsAccessReader(
	links analyticsaccess.LinkLifecycleReader,
	resumes resumecontracts.ResumeOwnershipChecker,
	publication resumecontracts.ResumePublicationLocker,
	overview resumecontracts.OwnedResumeSectionsQueryHandler,
	publicResume resumecontracts.PublicResumeQueryHandler,
) *analyticsaccess.Reader {
	return analyticsaccess.NewReader(links, resumes, publication, overview, publicResume)
}

func newAnalyticsReader(db *gorm.DB, subjects analyticscontracts.AnalyticsSubjectReader) *analyticsreaders.AnalyticsReader {
	return analyticsreaders.NewAnalyticsReader(db, subjects)
}

func newCleanupAnalyticsHandler(
	repository analyticsports.RetentionRepository,
	uow analyticsports.UnitOfWork,
	cfg config.Config,
) *analyticscommands.CleanupAnalyticsHandler {
	return analyticscommands.NewCleanupAnalyticsHandler(repository, uow, analyticscommands.AnalyticsRetentionPolicy{
		ClickRetentionDays: cfg.Analytics.ClickRetentionDays,
		ViewRetentionDays:  cfg.Analytics.ViewRetentionDays,
	})
}

func newIngestAnalyticsBatchHandler(
	repository analyticsports.AnalyticsBatchRepository,
	uow analyticsports.UnitOfWork,
	subjects analyticscontracts.AnalyticsSubjectReader,
	metrics *telemetry.Metrics,
) *analyticscommands.IngestAnalyticsBatchHandler {
	return analyticscommands.NewIngestAnalyticsBatchHandler(repository, uow, subjects, metrics)
}

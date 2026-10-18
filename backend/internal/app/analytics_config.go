package app

import (
	analyticscommands "developers-cv/backend/internal/modules/analytics/application/commands"
	analyticsports "developers-cv/backend/internal/modules/analytics/application/ports"
	resumelinkports "developers-cv/backend/internal/modules/resume_link/application/ports"
	resumelinkqueries "developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/platform/config"
	"developers-cv/backend/internal/platform/telemetry"
)

func newAnalyticsAccessReader(
	links resumelinkqueries.ResumeLinkPublicReader,
	resumes resumelinkports.ResumeOwnershipReader,
	overview resumelinkqueries.AnalyticsResumeOverviewReader,
	publicResume resumelinkqueries.ResumePublicQuery,
) analyticsports.AnalyticsAccessReader {
	return resumelinkqueries.NewAnalyticsAccessReader(links, resumes, overview, publicResume)
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
	metrics *telemetry.Metrics,
) *analyticscommands.IngestAnalyticsBatchHandler {
	return analyticscommands.NewIngestAnalyticsBatchHandler(repository, uow, metrics)
}

package app

import (
	"developers-cv/backend/internal/modules/resume_link/application/commands"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/platform/config"
	"developers-cv/backend/internal/platform/telemetry"
)

func newCleanupAnalyticsHandler(
	repository ports.AnalyticsRetentionRepository,
	uow ports.UnitOfWork,
	cfg config.Config,
) *commands.CleanupAnalyticsHandler {
	return commands.NewCleanupAnalyticsHandler(repository, uow, commands.AnalyticsRetentionPolicy{
		ClickRetentionDays: cfg.Analytics.ClickRetentionDays,
		ViewRetentionDays:  cfg.Analytics.ViewRetentionDays,
	})
}

func newIngestAnalyticsBatchHandler(
	repository ports.AnalyticsBatchRepository,
	uow ports.UnitOfWork,
	metrics *telemetry.Metrics,
) *commands.IngestAnalyticsBatchHandler {
	return commands.NewIngestAnalyticsBatchHandler(repository, uow, metrics)
}

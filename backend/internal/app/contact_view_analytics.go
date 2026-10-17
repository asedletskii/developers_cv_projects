package app

import (
	"go.uber.org/zap"

	resumelinkqueries "developers-cv/backend/internal/modules/resume_link/application/queries"
)

func newContactViewAnalyticsErrorReporter(logger *zap.Logger) resumelinkqueries.ContactViewAnalyticsErrorReporter {
	return func(err error) {
		logger.Error("failed to record contact view analytics", zap.Error(err))
	}
}

func newPublicAnalyticsErrorReporter(logger *zap.Logger) resumelinkqueries.PublicAnalyticsErrorReporter {
	return func(failure resumelinkqueries.PublicAnalyticsFailure) {
		fields := []zap.Field{
			zap.String("analytics_stage", string(failure.Stage)),
			zap.Error(failure.Err),
		}
		if !failure.LinkID.IsZero() {
			fields = append(fields, zap.String("link_id", failure.LinkID.String()))
		}
		if !failure.ResumeID.IsZero() {
			fields = append(fields, zap.String("resume_id", failure.ResumeID.String()))
		}
		logger.Error("public resume analytics failed", fields...)
	}
}

package app

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"developers-cv/backend/internal/platform/telemetry"
)

const (
	schemaOrphanCheckInterval = 5 * time.Minute
	schemaOrphanCheckTimeout  = 30 * time.Second
)

type orphanChecker struct {
	db       *gorm.DB
	logger   *zap.Logger
	metrics  *telemetry.Metrics
	interval time.Duration
}

func newOrphanChecker(db *gorm.DB, logger *zap.Logger, metrics *telemetry.Metrics) *orphanChecker {
	return &orphanChecker{db: db, logger: logger, metrics: metrics, interval: schemaOrphanCheckInterval}
}

func (c *orphanChecker) Run(ctx context.Context) {
	if c == nil || c.db == nil {
		return
	}
	c.check(ctx)
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.check(ctx)
		}
	}
}

type schemaOrphanCounts struct {
	AnalyticsWithoutLink      int64 `gorm:"column:analytics_without_link"`
	LinksWithoutResume        int64 `gorm:"column:links_without_resume"`
	ResumesWithoutUser        int64 `gorm:"column:resumes_without_user"`
	AuthIdentitiesWithoutUser int64 `gorm:"column:auth_identities_without_user"`
	AuthSessionsWithoutUser   int64 `gorm:"column:auth_sessions_without_user"`
}

func (c *schemaOrphanCounts) total() int64 {
	return c.AnalyticsWithoutLink + c.LinksWithoutResume + c.ResumesWithoutUser +
		c.AuthIdentitiesWithoutUser + c.AuthSessionsWithoutUser
}

func (c *orphanChecker) check(parent context.Context) {
	ctx, cancel := context.WithTimeout(parent, schemaOrphanCheckTimeout)
	defer cancel()

	var counts schemaOrphanCounts
	err := c.db.WithContext(ctx).Raw(`
		SELECT
			(SELECT count(*) FROM analytics.analytics_views v
			 LEFT JOIN resume_link.resume_links l ON l.id = v.link_id
			 WHERE l.id IS NULL) AS analytics_without_link,
			(SELECT count(*) FROM resume_link.resume_links l
			 LEFT JOIN resume.resumes r ON r.id = l.resume_id
			 WHERE r.id IS NULL) AS links_without_resume,
			(SELECT count(*) FROM resume.resumes r
			 LEFT JOIN users.users u ON u.id = r.user_id
			 WHERE u.id IS NULL) AS resumes_without_user,
			(SELECT count(*) FROM auth.user_identities i
			 LEFT JOIN users.users u ON u.id = i.user_id
			 WHERE u.id IS NULL) AS auth_identities_without_user,
			(SELECT count(*) FROM auth.auth_sessions s
			 LEFT JOIN users.users u ON u.id = s.user_id
			 WHERE u.id IS NULL) AS auth_sessions_without_user
	`).Scan(&counts).Error
	if err != nil {
		if c.logger != nil {
			c.logger.Warn("schema orphan check failed", zap.Error(err))
		}
		return
	}

	c.metrics.RecordSchemaOrphanRows("analytics_without_link", counts.AnalyticsWithoutLink)
	c.metrics.RecordSchemaOrphanRows("links_without_resume", counts.LinksWithoutResume)
	c.metrics.RecordSchemaOrphanRows("resumes_without_user", counts.ResumesWithoutUser)
	c.metrics.RecordSchemaOrphanRows("auth_identities_without_user", counts.AuthIdentitiesWithoutUser)
	c.metrics.RecordSchemaOrphanRows("auth_sessions_without_user", counts.AuthSessionsWithoutUser)

	if counts.total() > 0 && c.logger != nil {
		c.logger.Warn("schema orphan rows detected",
			zap.Int64("analytics_without_link", counts.AnalyticsWithoutLink),
			zap.Int64("links_without_resume", counts.LinksWithoutResume),
			zap.Int64("resumes_without_user", counts.ResumesWithoutUser),
			zap.Int64("auth_identities_without_user", counts.AuthIdentitiesWithoutUser),
			zap.Int64("auth_sessions_without_user", counts.AuthSessionsWithoutUser),
		)
	}
}

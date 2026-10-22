package readers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"developers-cv/backend/internal/modules/analytics/application/queries"
	"developers-cv/backend/internal/platform/database"
)

const resumeStatisticsQuery = `
SELECT
    COUNT(*) FILTER (WHERE v.traffic = 'human') AS views,
    COUNT(*) FILTER (
        WHERE v.traffic = 'human' AND v.confirmed_at IS NOT NULL
    ) AS reads,
    COUNT(DISTINCT v.visitor_id) FILTER (
        WHERE v.traffic = 'human'
    ) AS unique_visitors,
    COUNT(DISTINCT v.visitor_id) FILTER (
        WHERE v.traffic = 'human'
          AND EXISTS (
              SELECT 1
              FROM analytics.analytics_views AS previous
              WHERE previous.resume_id = @resumeID
                AND previous.visitor_id = v.visitor_id
                AND previous.traffic = 'human'
                AND previous.started_at >= @periodFrom
                AND previous.started_at < @periodTo
                AND (previous.started_at AT TIME ZONE @periodTimezone)::date
                    <> (v.started_at AT TIME ZONE @periodTimezone)::date
          )
    ) AS returning_visitors,
    COUNT(*) FILTER (
        WHERE v.traffic = 'human'
          AND EXISTS (
              SELECT 1
              FROM analytics.analytics_events AS event
              WHERE event.view_id = v.id
                AND event.event_type = 'contact_view'
                AND event.occurred_at >= @periodFrom
                AND event.occurred_at < @periodTo
          )
    ) AS contact_views,
    COUNT(*) FILTER (
        WHERE v.traffic = 'human'
          AND EXISTS (
              SELECT 1
              FROM analytics.analytics_events AS event
              WHERE event.view_id = v.id
                AND event.event_type = 'download'
                AND event.occurred_at >= @periodFrom
                AND event.occurred_at < @periodTo
          )
    ) AS downloads,
    COALESCE((
        SELECT percentile_cont(0.5) WITHIN GROUP (ORDER BY median_view.active_ms)
        FROM analytics.analytics_views AS median_view
        WHERE median_view.resume_id = @resumeID
          AND median_view.traffic = 'human'
          AND median_view.confirmed_at IS NOT NULL
          AND median_view.started_at >= @periodFrom
          AND median_view.started_at < @periodTo
    ), 0) AS median
FROM analytics.analytics_views AS v
WHERE v.resume_id = @resumeID
  AND v.started_at >= @periodFrom
  AND v.started_at < @periodTo`

type resumeLinkViewRow struct {
	LinkID     uuid.UUID `gorm:"column:link_id"`
	Views      int64
	Reads      int64
	LastViewAt *time.Time `gorm:"column:last_view_at"`
}

func (r *AnalyticsReader) GetResumeStatistics(
	ctx context.Context,
	query queries.ResumeAnalyticsQuery,
) (queries.ResumeStatistics, error) {
	if r == nil || r.db == nil {
		return queries.ResumeStatistics{}, fmt.Errorf("get resume analytics statistics: reader dependencies are nil")
	}
	var result queries.ResumeStatistics
	if err := database.WithReadSnapshot(ctx, r.db, func(snapshotCtx context.Context) error {
		var err error
		result, err = r.getResumeStatistics(snapshotCtx, query)
		return err
	}); err != nil {
		return queries.ResumeStatistics{}, fmt.Errorf("read resume analytics statistics snapshot: %w", err)
	}
	return result, nil
}

func (r *AnalyticsReader) getResumeStatistics(
	ctx context.Context,
	query queries.ResumeAnalyticsQuery,
) (queries.ResumeStatistics, error) {
	db := database.DBFromContext(ctx, r.db).WithContext(ctx)
	base := queries.LinkStatistics{}
	var firstLast firstLastViewRow
	if err := db.Table("analytics.analytics_views AS view").
		Select("MIN(view.started_at) AS first_view_at, MAX(view.started_at) AS last_view_at").
		Where("view.resume_id = ?", query.ResumeID.UUID()).
		Where("view.traffic = ?", "human").
		Scan(&firstLast).Error; err != nil {
		return queries.ResumeStatistics{}, fmt.Errorf("get resume analytics first/last view: %w", err)
	}
	base.FirstViewAt, base.LastViewAt = firstLast.First, firstLast.Last

	var row linkStatisticsRow
	if err := db.Raw(
		resumeStatisticsQuery,
		sql.Named("resumeID", query.ResumeID.UUID()),
		sql.Named("periodFrom", query.Period.From),
		sql.Named("periodTo", query.Period.To),
		sql.Named("periodTimezone", query.Period.TZ),
	).Scan(&row).Error; err != nil {
		return queries.ResumeStatistics{}, fmt.Errorf("get resume analytics statistics: %w", err)
	}
	applyStatistics(&base, row)

	var err error
	base.Devices, err = readResumeDeviceStatistics(db, query.ResumeID.UUID(), query.Period.From, query.Period.To)
	if err != nil {
		return queries.ResumeStatistics{}, fmt.Errorf("get resume analytics devices: %w", err)
	}
	base.Referrers, err = readResumeReferrerStatistics(db, query.ResumeID.UUID(), query.Period.From, query.Period.To)
	if err != nil {
		return queries.ResumeStatistics{}, fmt.Errorf("get resume analytics referrers: %w", err)
	}
	base.Filtered, err = readResumeFilteredStatistics(db, query.ResumeID.UUID(), query.Period.From, query.Period.To)
	if err != nil {
		return queries.ResumeStatistics{}, fmt.Errorf("get resume analytics filtered: %w", err)
	}

	var rows []resumeLinkViewRow
	err = db.Table("analytics.analytics_views AS view").
		Select(`
			view.link_id AS link_id,
			COUNT(*) FILTER (
				WHERE view.traffic = 'human'
				  AND view.started_at >= ? AND view.started_at < ?
			) AS views,
			COUNT(*) FILTER (
				WHERE view.traffic = 'human'
				  AND view.confirmed_at IS NOT NULL
				  AND view.started_at >= ? AND view.started_at < ?
			) AS reads,
			MAX(view.started_at) FILTER (
				WHERE view.traffic = 'human'
			) AS last_view_at
		`, query.Period.From, query.Period.To, query.Period.From, query.Period.To).
		Where("view.resume_id = ?", query.ResumeID.UUID()).
		Group("view.link_id").
		Scan(&rows).Error
	if err != nil {
		return queries.ResumeStatistics{}, fmt.Errorf("get resume analytics link views: %w", err)
	}

	if r.subjects == nil {
		return queries.ResumeStatistics{}, fmt.Errorf("get resume analytics links: subject reader is nil")
	}
	refs, err := r.subjects.ListResumeLinks(ctx, query.ResumeID)
	if err != nil {
		return queries.ResumeStatistics{}, fmt.Errorf("get resume analytics links: %w", err)
	}
	byLinkID := make(map[uuid.UUID]resumeLinkViewRow, len(rows))
	for _, row := range rows {
		byLinkID[row.LinkID] = row
	}
	links := make([]queries.ResumeLinkStatistic, 0, len(refs))
	for _, ref := range refs {
		link := queries.ResumeLinkStatistic{ID: ref.ID.String(), Name: ref.Name}
		if row, ok := byLinkID[ref.ID.UUID()]; ok {
			link.Views = row.Views
			link.Reads = row.Reads
			link.LastViewAt = row.LastViewAt
		}
		links = append(links, link)
	}

	return queries.ResumeStatistics{
		LinkStatistics: base,
		Links:          links,
	}, nil
}

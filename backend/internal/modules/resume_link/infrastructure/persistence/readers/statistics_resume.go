package readers

import (
	"context"
	"database/sql"
	"fmt"

	"developers-cv/backend/internal/modules/resume_link/application/queries"
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
              FROM resume_link_views AS previous
              JOIN resume_links AS previous_link
                ON previous_link.id = previous.link_id
              WHERE previous_link.resume_id = @resumeID
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
              FROM resume_link_events AS event
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
              FROM resume_link_events AS event
              WHERE event.view_id = v.id
                AND event.event_type = 'download'
                AND event.occurred_at >= @periodFrom
                AND event.occurred_at < @periodTo
          )
    ) AS downloads,
    COALESCE((
        SELECT percentile_cont(0.5) WITHIN GROUP (ORDER BY median_view.active_ms)
        FROM resume_link_views AS median_view
        JOIN resume_links AS median_link ON median_link.id = median_view.link_id
        WHERE median_link.resume_id = @resumeID
          AND median_view.traffic = 'human'
          AND median_view.confirmed_at IS NOT NULL
          AND median_view.started_at >= @periodFrom
          AND median_view.started_at < @periodTo
    ), 0) AS median
FROM resume_link_views AS v
JOIN resume_links AS link ON link.id = v.link_id
WHERE link.resume_id = @resumeID
  AND v.started_at >= @periodFrom
  AND v.started_at < @periodTo`

func (r *ResumeLinkReader) GetResumeStatistics(
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

func (r *ResumeLinkReader) getResumeStatistics(
	ctx context.Context,
	query queries.ResumeAnalyticsQuery,
) (queries.ResumeStatistics, error) {
	db := database.DBFromContext(ctx, r.db).WithContext(ctx)
	base := queries.LinkStatistics{}
	var firstLast firstLastViewRow
	if err := db.Table("resume_link_views AS view").
		Select("MIN(view.started_at) AS first_view_at, MAX(view.started_at) AS last_view_at").
		Joins("JOIN resume_links AS link ON link.id = view.link_id").
		Where("link.resume_id = ?", query.ResumeID.UUID()).
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

	var links []queries.ResumeLinkStatistic
	err = db.Table("resume_links AS link").
		Select(`
			link.id,
			link.name,
			COUNT(view.id) FILTER (
				WHERE view.traffic = 'human'
				  AND view.started_at >= ? AND view.started_at < ?
			) AS views,
			COUNT(view.id) FILTER (
				WHERE view.traffic = 'human'
				  AND view.confirmed_at IS NOT NULL
				  AND view.started_at >= ? AND view.started_at < ?
			) AS reads,
			MAX(view.started_at) FILTER (
				WHERE view.traffic = 'human'
			) AS last_view_at
		`, query.Period.From, query.Period.To, query.Period.From, query.Period.To).
		Joins("LEFT JOIN resume_link_views AS view ON view.link_id = link.id").
		Where("link.resume_id = ?", query.ResumeID.UUID()).
		Group("link.id, link.name, link.created_at").
		Order("link.created_at ASC, link.id ASC").
		Scan(&links).Error
	if err != nil {
		return queries.ResumeStatistics{}, fmt.Errorf("get resume analytics links: %w", err)
	}

	return queries.ResumeStatistics{
		LinkStatistics: base,
		Links:          append(make([]queries.ResumeLinkStatistic, 0, len(links)), links...),
	}, nil
}

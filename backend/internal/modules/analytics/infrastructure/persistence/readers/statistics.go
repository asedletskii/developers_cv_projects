package readers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"developers-cv/backend/internal/modules/analytics/application/queries"
	"developers-cv/backend/internal/modules/analytics/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/platform/database"
)

const linkStatisticsQuery = `
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
              FROM analytics_views AS previous
              WHERE previous.link_id = v.link_id
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
              FROM analytics_events AS event
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
              FROM analytics_events AS event
              WHERE event.view_id = v.id
                AND event.event_type = 'download'
                AND event.occurred_at >= @periodFrom
                AND event.occurred_at < @periodTo
          )
    ) AS downloads,
    COALESCE((
        SELECT percentile_cont(0.5) WITHIN GROUP (ORDER BY median_view.active_ms)
        FROM analytics_views AS median_view
        WHERE median_view.link_id = @linkID
          AND median_view.traffic = 'human'
          AND median_view.confirmed_at IS NOT NULL
          AND median_view.started_at >= @periodFrom
          AND median_view.started_at < @periodTo
    ), 0) AS median
FROM analytics_views AS v
WHERE v.link_id = @linkID
  AND v.started_at >= @periodFrom
  AND v.started_at < @periodTo`

func (r *AnalyticsReader) GetLinkStatistics(
	ctx context.Context,
	query queries.LinkAnalyticsQuery,
) (queries.LinkStatistics, error) {
	if r == nil || r.db == nil {
		return queries.LinkStatistics{}, fmt.Errorf("get analytics statistics: reader dependencies are nil")
	}
	var result queries.LinkStatistics
	if err := database.WithReadSnapshot(ctx, r.db, func(snapshotCtx context.Context) error {
		var err error
		result, err = r.getLinkStatistics(snapshotCtx, query)
		return err
	}); err != nil {
		return queries.LinkStatistics{}, fmt.Errorf("read analytics statistics snapshot: %w", err)
	}
	return result, nil
}

func (r *AnalyticsReader) getLinkStatistics(
	ctx context.Context,
	query queries.LinkAnalyticsQuery,
) (queries.LinkStatistics, error) {
	db := database.DBFromContext(ctx, r.db).WithContext(ctx)
	out := queries.LinkStatistics{}
	var firstLast firstLastViewRow
	if err := db.Model(&gormmodel.AnalyticsView{}).
		Select("MIN(started_at) AS first_view_at, MAX(started_at) AS last_view_at").
		Where("link_id = ?", query.LinkID.UUID()).
		Where("traffic = ?", "human").
		Scan(&firstLast).Error; err != nil {
		return out, fmt.Errorf("get analytics first/last view: %w", err)
	}
	out.FirstViewAt, out.LastViewAt = firstLast.First, firstLast.Last

	var base linkStatisticsRow
	if err := db.Raw(
		linkStatisticsQuery,
		sql.Named("linkID", query.LinkID.UUID()),
		sql.Named("periodFrom", query.Period.From),
		sql.Named("periodTo", query.Period.To),
		sql.Named("periodTimezone", query.Period.TZ),
	).Scan(&base).Error; err != nil {
		return out, fmt.Errorf("get analytics statistics: %w", err)
	}
	applyStatistics(&out, base)

	var err error
	out.Devices, err = readDeviceStatistics(db, query.LinkID.UUID(), query.Period.From, query.Period.To)
	if err != nil {
		return out, fmt.Errorf("get analytics devices: %w", err)
	}
	out.Referrers, err = readReferrerStatistics(db, query.LinkID.UUID(), query.Period.From, query.Period.To)
	if err != nil {
		return out, fmt.Errorf("get analytics referrers: %w", err)
	}
	out.Filtered, err = readFilteredStatistics(db, query.LinkID.UUID(), query.Period.From, query.Period.To)
	if err != nil {
		return out, fmt.Errorf("get analytics filtered: %w", err)
	}
	return out, nil
}

type firstLastViewRow struct {
	First *time.Time `gorm:"column:first_view_at"`
	Last  *time.Time `gorm:"column:last_view_at"`
}

type linkStatisticsRow struct {
	Views, Reads, UniqueVisitors, ReturningVisitors int64
	ContactViews, Downloads                         int64
	Median                                          float64
}

func applyStatistics(out *queries.LinkStatistics, row linkStatisticsRow) {
	out.Views = row.Views
	out.Reads = row.Reads
	out.UniqueVisitors = row.UniqueVisitors
	out.ReturningVisitors = row.ReturningVisitors
	out.ContactViews = row.ContactViews
	out.Downloads = row.Downloads
	out.MedianActiveSeconds = int64(row.Median / 1000)
}

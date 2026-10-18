package readers

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/analytics/application/queries"
	"developers-cv/backend/internal/modules/analytics/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/platform/database"
)

func (r *AnalyticsReader) GetDaily(
	ctx context.Context,
	query queries.LinkAnalyticsQuery,
) (queries.LinkDaily, error) {
	if r == nil || r.db == nil {
		return queries.LinkDaily{}, fmt.Errorf("get analytics daily: reader dependencies are nil")
	}
	db := database.DBFromContext(ctx, r.db).WithContext(ctx)
	var rows []dailyStatisticsRow
	err := db.Model(&gormmodel.AnalyticsView{}).
		Select(`
			(started_at AT TIME ZONE ?)::date::text AS date,
			COUNT(*) AS views,
			COUNT(*) FILTER (WHERE confirmed_at IS NOT NULL) AS reads,
			COUNT(DISTINCT visitor_id) AS unique_visitors
		`, query.Period.TZ).
		Where("link_id = ?", query.LinkID.UUID()).
		Where("traffic = ?", "human").
		Where("started_at >= ? AND started_at < ?", query.Period.From, query.Period.To).
		Group("date").
		Order("date ASC").
		Scan(&rows).Error
	if err != nil {
		return queries.LinkDaily{}, fmt.Errorf("get analytics daily: %w", err)
	}

	byDate := make(map[string]dailyStatisticsRow, len(rows))
	for _, row := range rows {
		byDate[row.Date] = row
	}
	days := query.Period.Days()
	out := queries.LinkDaily{Days: make([]queries.DailyPoint, 0, len(days))}
	for _, day := range days {
		row := byDate[day]
		out.Days = append(out.Days, queries.DailyPoint{
			Date: day, Views: row.Views, Reads: row.Reads,
			UniqueVisitors: row.UniqueVisitors,
		})
	}
	return out, nil
}

package readers

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/analytics/application/queries"
	"developers-cv/backend/internal/modules/analytics/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/platform/database"
)

func (r *AnalyticsReader) GetSections(
	ctx context.Context,
	query queries.LinkSectionsQuery,
) (queries.SectionsStatistics, error) {
	if r == nil || r.db == nil {
		return queries.SectionsStatistics{}, fmt.Errorf("get analytics sections: reader dependencies are nil")
	}
	var result queries.SectionsStatistics
	if err := database.WithReadSnapshot(ctx, r.db, func(snapshotCtx context.Context) error {
		var err error
		result, err = r.getSections(snapshotCtx, query)
		return err
	}); err != nil {
		return queries.SectionsStatistics{}, fmt.Errorf("read analytics sections snapshot: %w", err)
	}
	return result, nil
}

func (r *AnalyticsReader) getSections(
	ctx context.Context,
	query queries.LinkSectionsQuery,
) (queries.SectionsStatistics, error) {
	db := database.DBFromContext(ctx, r.db).WithContext(ctx)
	var reads int64
	if err := db.Model(&gormmodel.AnalyticsView{}).
		Where("link_id = ?", query.LinkID.UUID()).
		Where("traffic = ?", "human").
		Where("confirmed_at IS NOT NULL").
		Where("started_at >= ? AND started_at < ?", query.Period.From, query.Period.To).
		Count(&reads).Error; err != nil {
		return queries.SectionsStatistics{}, fmt.Errorf("get analytics section reads: %w", err)
	}

	var rows []sectionStatisticsRow
	err := db.Table("analytics_events AS event").
		Select(`
			event.payload->>'section' AS section,
			COUNT(DISTINCT event.view_id)
				FILTER (WHERE view.confirmed_at IS NOT NULL) AS reached,
			COUNT(DISTINCT view.id)
				FILTER (
					WHERE view.confirmed_at IS NOT NULL
					  AND view.exit_section = event.payload->>'section'
				) AS exits
		`).
		Joins("JOIN analytics_views AS view ON view.id = event.view_id").
		Where("event.link_id = ?", query.LinkID.UUID()).
		Where("event.event_type = ?", "section_view").
		Where("view.traffic = ?", "human").
		Where("view.started_at >= ? AND view.started_at < ?", query.Period.From, query.Period.To).
		Group("event.payload->>'section'").
		Scan(&rows).Error
	if err != nil {
		return queries.SectionsStatistics{}, fmt.Errorf("get analytics sections: %w", err)
	}

	counts := make(map[string]queries.SectionStatistics, len(rows))
	for _, row := range rows {
		counts[row.Section] = queries.SectionStatistics{
			Section: row.Section, Reached: row.Reached, Exits: row.Exits,
		}
	}

	out := queries.SectionsStatistics{
		Reads:    reads,
		Sections: make([]queries.SectionStatistics, 0, len(query.SectionOrder)+len(rows)),
	}
	seen := make(map[string]bool, len(query.SectionOrder))
	for _, section := range query.SectionOrder {
		if section == "techStack" || seen[section] {
			continue
		}
		statistic := counts[section]
		statistic.Section = section
		out.Sections = append(out.Sections, statistic)
		seen[section] = true
	}
	for _, row := range rows {
		if seen[row.Section] {
			continue
		}
		out.Sections = append(out.Sections, counts[row.Section])
		seen[row.Section] = true
	}
	return out, nil
}

func (r *AnalyticsReader) GetTechnologies(
	ctx context.Context,
	query queries.LinkAnalyticsQuery,
) (queries.TechnologiesStatistics, error) {
	if r == nil || r.db == nil {
		return queries.TechnologiesStatistics{}, fmt.Errorf("get analytics technologies: reader dependencies are nil")
	}
	db := database.DBFromContext(ctx, r.db).WithContext(ctx)
	var rows []technologyStatisticsRow
	err := db.Table("analytics_events AS event").
		Select(`
			event.payload->>'technology' AS name,
			COUNT(*) AS selections,
			COUNT(DISTINCT view.visitor_id) AS visitors
		`).
		Joins("JOIN analytics_views AS view ON view.id = event.view_id").
		Where("event.link_id = ?", query.LinkID.UUID()).
		Where("event.event_type = ?", "tech_filter").
		Where("event.payload->>'action' = ?", "select").
		Where("view.traffic = ?", "human").
		Where("view.started_at >= ? AND view.started_at < ?", query.Period.From, query.Period.To).
		Group("event.payload->>'technology'").
		Order("selections DESC, name ASC").
		Limit(20).
		Scan(&rows).Error
	if err != nil {
		return queries.TechnologiesStatistics{}, fmt.Errorf("get analytics technologies: %w", err)
	}

	out := queries.TechnologiesStatistics{
		Technologies: make([]queries.TechnologyInterest, 0, len(rows)),
	}
	for _, row := range rows {
		out.Technologies = append(out.Technologies, queries.TechnologyInterest{
			Name: row.Name, Selections: row.Selections, Visitors: row.Visitors,
		})
	}
	return out, nil
}

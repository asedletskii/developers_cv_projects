package readers

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/analytics/application/queries"
	"developers-cv/backend/internal/platform/database"
)

func (r *AnalyticsReader) GetClicks(
	ctx context.Context,
	query queries.LinkClicksQuery,
) (queries.ClickGrid, error) {
	if r == nil || r.db == nil {
		return queries.ClickGrid{}, fmt.Errorf("get analytics clicks: reader dependencies are nil")
	}
	var result queries.ClickGrid
	if err := database.WithReadSnapshot(ctx, r.db, func(snapshotCtx context.Context) error {
		var err error
		result, err = r.getClicks(snapshotCtx, query)
		return err
	}); err != nil {
		return queries.ClickGrid{}, fmt.Errorf("read analytics clicks snapshot: %w", err)
	}
	return result, nil
}

func (r *AnalyticsReader) getClicks(
	ctx context.Context,
	query queries.LinkClicksQuery,
) (queries.ClickGrid, error) {
	breakpoint, err := clickBreakpointCondition(query.Breakpoint)
	if err != nil {
		return queries.ClickGrid{}, err
	}

	db := database.DBFromContext(ctx, r.db).WithContext(ctx)
	cellQuery := clickQuery(db, query, breakpoint)
	if !query.IncludeStale {
		cellQuery = cellQuery.Where("c.layout_hash = ?", query.CurrentHash)
	}

	var rows []clickCellRow
	cellQuery = cellQuery.
		Select(`
			LEAST((c.x_pct * 20 / 10000), 19) AS x,
			LEAST((c.y_pct * 20 / 10000), 19) AS y,
			COUNT(*) AS count
		`).
		Group("LEAST((c.x_pct * 20 / 10000), 19), LEAST((c.y_pct * 20 / 10000), 19)").
		Order("x ASC, y ASC")
	if err := cellQuery.Scan(&rows).Error; err != nil {
		return queries.ClickGrid{}, fmt.Errorf("get analytics click cells: %w", err)
	}

	totalQuery := clickQuery(db, query, breakpoint)
	if !query.IncludeStale {
		totalQuery = totalQuery.Where("c.layout_hash = ?", query.CurrentHash)
	}

	var total int64
	if err := totalQuery.Count(&total).Error; err != nil {
		return queries.ClickGrid{}, fmt.Errorf("get analytics click total: %w", err)
	}

	var staleTotal int64
	if !query.IncludeStale {
		staleQuery := clickQuery(db, query, breakpoint).
			Where("c.layout_hash <> ?", query.CurrentHash)
		if err := staleQuery.Count(&staleTotal).Error; err != nil {
			return queries.ClickGrid{}, fmt.Errorf("get analytics stale click total: %w", err)
		}
	}

	cells := make([]queries.ClickCell, 0, len(rows))
	for _, row := range rows {
		cells = append(cells, queries.ClickCell{
			X: row.X, Y: row.Y, Count: row.Count,
		})
	}

	return queries.ClickGrid{
		Section:    query.Section,
		Grid:       queries.ClickGridSize,
		Cells:      cells,
		Total:      total,
		StaleTotal: staleTotal,
	}, nil
}

func clickQuery(
	db *gorm.DB,
	query queries.LinkClicksQuery,
	breakpoint string,
) *gorm.DB {
	return db.
		Table("analytics.analytics_clicks AS c").
		Joins("JOIN analytics.analytics_views AS v ON v.id = c.view_id").
		Where("c.link_id = ?", query.LinkID.UUID()).
		Where("c.section = ?", query.Section).
		Where("v.traffic = ?", "human").
		Where(breakpoint)
}

func clickBreakpointCondition(breakpoint string) (string, error) {
	switch breakpoint {
	case "xs":
		return "c.viewport_w < 600", nil
	case "md":
		return "c.viewport_w >= 600 AND c.viewport_w < 1200", nil
	case "lg":
		return "c.viewport_w >= 1200", nil
	default:
		return "", fmt.Errorf("get analytics clicks: invalid breakpoint")
	}
}

type clickCellRow struct {
	X     int64 `gorm:"column:x"`
	Y     int64 `gorm:"column:y"`
	Count int64 `gorm:"column:count"`
}

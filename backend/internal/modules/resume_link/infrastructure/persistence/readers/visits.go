package readers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/platform/database"
)

const visitsQuery = `
SELECT
    v.id,
    v.started_at,
    v.visitor_id,
    v.device,
    v.browser,
    v.os,
    v.referrer_host,
    v.active_ms,
    v.confirmed_at,
    v.max_section,
    v.exit_section,
    v.traffic,
    visitor.ordinal,
    v.started_at > visitor.first_seen_at AS is_returning
FROM resume_link_views AS v
JOIN resume_link_visitors AS visitor
  ON visitor.link_id = v.link_id
 AND visitor.visitor_id = v.visitor_id
WHERE v.link_id = ?
  AND (? OR v.traffic = 'human')
  AND (
      NOT ?
      OR v.started_at < ?
      OR (v.started_at = ? AND v.id < ?)
  )
ORDER BY v.started_at DESC, v.id DESC
LIMIT ?`

func (r *ResumeLinkReader) GetVisits(
	ctx context.Context,
	query queries.VisitsReadQuery,
) (queries.LinkVisitsPage, error) {
	if r == nil || r.db == nil {
		return queries.LinkVisitsPage{}, fmt.Errorf("get analytics visits: reader dependencies are nil")
	}
	var result queries.LinkVisitsPage
	if err := database.WithReadSnapshot(ctx, r.db, func(snapshotCtx context.Context) error {
		var err error
		result, err = r.getVisits(snapshotCtx, query)
		return err
	}); err != nil {
		return queries.LinkVisitsPage{}, fmt.Errorf("read analytics visits snapshot: %w", err)
	}
	return result, nil
}

func (r *ResumeLinkReader) getVisits(
	ctx context.Context,
	query queries.VisitsReadQuery,
) (queries.LinkVisitsPage, error) {
	hasCursor := !query.CursorTime.IsZero()
	var cursorTimeArg any
	var cursorIDArg any
	if hasCursor {
		cursorTimeArg = query.CursorTime
		cursorIDArg = query.CursorID
	}

	db := database.DBFromContext(ctx, r.db).WithContext(ctx)
	args := []any{
		query.LinkID.UUID(),
		query.IncludeFiltered,
		hasCursor,
		cursorTimeArg,
		cursorTimeArg,
		cursorIDArg,
		query.Limit + 1,
	}

	var rows []visitRow
	if err := db.Raw(visitsQuery, args...).Scan(&rows).Error; err != nil {
		return queries.LinkVisitsPage{}, fmt.Errorf("get analytics visits: %w", err)
	}

	hasMore := len(rows) > query.Limit
	if hasMore {
		rows = rows[:query.Limit]
	}

	out := queries.LinkVisitsPage{Items: make([]queries.LinkVisit, 0, len(rows))}
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		out.Items = append(out.Items, row.toQuery())
		ids = append(ids, row.ID)
	}

	if err := r.attachVisitEvents(db, ids, out.Items); err != nil {
		return queries.LinkVisitsPage{}, err
	}
	if hasMore {
		last := out.Items[len(out.Items)-1]
		next, err := queries.EncodeVisitCursor(last.StartedAt, last.ID)
		if err != nil {
			return queries.LinkVisitsPage{}, err
		}
		out.NextCursor = &next
	}

	return out, nil
}

type visitRow struct {
	ID, VisitorID, Device, Browser, OS, Traffic string
	StartedAt                                   time.Time
	ReferrerHost                                *string
	ActiveMS                                    int64
	ConfirmedAt                                 *time.Time
	MaxSection, ExitSection                     *string
	Ordinal                                     int64
	IsReturning                                 bool
}

func (row visitRow) toQuery() queries.LinkVisit {
	return queries.LinkVisit{
		ID:        row.ID,
		StartedAt: row.StartedAt,
		Visitor: queries.VisitorInfo{
			Ordinal:     row.Ordinal,
			IsReturning: row.IsReturning,
		},
		Device:        row.Device,
		Browser:       row.Browser,
		OS:            row.OS,
		ReferrerHost:  row.ReferrerHost,
		IsRead:        row.ConfirmedAt != nil,
		ActiveSeconds: row.ActiveMS / 1000,
		MaxSection:    row.MaxSection,
		ExitSection:   row.ExitSection,
		Traffic:       row.Traffic,
		Events:        []queries.VisitEvent{},
	}
}

func (r *ResumeLinkReader) attachVisitEvents(
	db *gorm.DB,
	ids []string,
	visits []queries.LinkVisit,
) error {
	if len(ids) == 0 {
		return nil
	}

	var events []visitEventRow
	if err := db.Table("resume_link_events").
		Select("view_id, event_type, occurred_at, payload").
		Where("view_id IN ?", ids).
		Where("event_type <> ?", "section_view").
		Order("occurred_at ASC, id ASC").
		Scan(&events).Error; err != nil {
		return fmt.Errorf("get analytics visit events: %w", err)
	}

	byID := make(map[string][]queries.VisitEvent, len(ids))
	for _, event := range events {
		var payload map[string]string
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			return fmt.Errorf("get analytics visit events: decode payload for view %q: %w", event.ViewID, err)
		}
		byID[event.ViewID] = append(byID[event.ViewID], queries.VisitEvent{
			Type:       event.EventType,
			At:         event.OccurredAt,
			Contact:    payload["contact"],
			Technology: payload["technology"],
			Kind:       payload["kind"],
			Target:     payload["target"],
		})
	}
	for i := range visits {
		visits[i].Events = append(make([]queries.VisitEvent, 0, len(byID[visits[i].ID])), byID[visits[i].ID]...)
	}
	return nil
}

type visitEventRow struct {
	ViewID, EventType string
	OccurredAt        time.Time
	Payload           []byte
}

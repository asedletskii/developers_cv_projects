package readers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/modules/resume_link/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/platform/database"
)

type ResumeLinkReader struct {
	db *gorm.DB
}

func NewResumeLinkReader(db *gorm.DB) *ResumeLinkReader {
	return &ResumeLinkReader{db: db}
}

func (r *ResumeLinkReader) GetByID(ctx context.Context, query queries.PublicResumeLinkQuery) (queries.ResumeLinkPublicView, error) {
	if r == nil || r.db == nil {
		return queries.ResumeLinkPublicView{}, fmt.Errorf("get resume link: reader dependencies are nil")
	}
	if query.LinkID.IsZero() {
		return queries.ResumeLinkPublicView{}, fmt.Errorf("get resume link: id is empty")
	}

	db := database.DBFromContext(ctx, r.db).WithContext(ctx)
	var model gormmodel.ResumeLink
	if err := db.Where("id = ?", query.LinkID.UUID()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return queries.ResumeLinkPublicView{}, domain.ErrEntityNotFound
		}
		return queries.ResumeLinkPublicView{}, fmt.Errorf("get resume link: %w", err)
	}
	return mapView(model), nil
}

func mapView(model gormmodel.ResumeLink) queries.ResumeLinkPublicView {
	return queries.ResumeLinkPublicView{
		ID: model.ID.String(), ResumeID: model.ResumeID.String(), Status: model.Status,
	}
}

func (r *ResumeLinkReader) ListByResumeID(ctx context.Context, query queries.ListResumeLinksQuery) ([]queries.ResumeLinkListItem, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("list resume links: reader dependencies are nil")
	}
	if query.ResumeID.IsZero() {
		return nil, fmt.Errorf("list resume links: resume id is empty")
	}

	db := database.DBFromContext(ctx, r.db).WithContext(ctx)
	var models []gormmodel.ResumeLink
	if err := db.Where("resume_id = ?", query.ResumeID.UUID()).Order("created_at, id").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("list resume links: %w", err)
	}
	result := make([]queries.ResumeLinkListItem, 0, len(models))
	for _, model := range models {
		result = append(result, queries.ResumeLinkListItem{ID: model.ID.String(), Name: model.Name, Status: model.Status, Version: model.Version, CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt})
	}
	return result, nil
}

func (r *ResumeLinkReader) GetLinkStatistics(ctx context.Context, query queries.LinkAnalyticsQuery) (queries.LinkStatistics, error) {
	if r == nil || r.db == nil {
		return queries.LinkStatistics{}, fmt.Errorf("get analytics statistics: reader dependencies are nil")
	}
	if err := query.Period.Validate(); err != nil {
		return queries.LinkStatistics{}, err
	}
	db := database.DBFromContext(ctx, r.db).WithContext(ctx)
	var out queries.LinkStatistics
	var firstLast struct {
		First *time.Time `gorm:"column:first_view_at"`
		Last  *time.Time `gorm:"column:last_view_at"`
	}
	if err := db.Raw(`SELECT MIN(started_at) first_view_at, MAX(started_at) last_view_at FROM resume_link_views WHERE link_id = ? AND traffic = 'human'`, query.LinkID.UUID()).Scan(&firstLast).Error; err != nil {
		return out, fmt.Errorf("get analytics first/last view: %w", err)
	}
	out.FirstViewAt, out.LastViewAt = firstLast.First, firstLast.Last
	var base struct {
		Views, Reads, UniqueVisitors, ReturningVisitors, ContactViews, Downloads int64
		Median                                                                   float64
	}
	if err := db.Raw(`
SELECT
 COUNT(*) FILTER (WHERE v.traffic='human') views,
 COUNT(*) FILTER (WHERE v.traffic='human' AND v.active_ms >= 5000) reads,
 COUNT(DISTINCT v.visitor_id) FILTER (WHERE v.traffic='human') unique_visitors,
 COUNT(DISTINCT v.visitor_id) FILTER (WHERE v.traffic='human' AND EXISTS (SELECT 1 FROM resume_link_views p WHERE p.link_id=v.link_id AND p.visitor_id=v.visitor_id AND p.traffic='human' AND p.started_at>=? AND p.started_at<? AND (p.started_at AT TIME ZONE ?)::date <> (v.started_at AT TIME ZONE ?)::date)) returning_visitors,
 COUNT(*) FILTER (WHERE v.traffic='human' AND EXISTS (SELECT 1 FROM resume_link_events e WHERE e.view_id=v.id AND e.event_type='contact_view' AND e.occurred_at >= ? AND e.occurred_at < ?)) contact_views,
 COUNT(*) FILTER (WHERE v.traffic='human' AND EXISTS (SELECT 1 FROM resume_link_events e WHERE e.view_id=v.id AND e.event_type='download' AND e.occurred_at >= ? AND e.occurred_at < ?)) downloads,
 COALESCE((SELECT percentile_cont(0.5) WITHIN GROUP (ORDER BY active_ms) FROM resume_link_views m WHERE m.link_id=? AND m.traffic='human' AND m.active_ms >= 5000 AND m.started_at >= ? AND m.started_at < ?),0) median
FROM resume_link_views v WHERE v.link_id=? AND v.started_at >= ? AND v.started_at < ?`, query.Period.From, query.Period.To, query.Period.TZ, query.Period.TZ, query.Period.From, query.Period.To, query.Period.From, query.Period.To, query.LinkID.UUID(), query.Period.From, query.Period.To, query.LinkID.UUID(), query.Period.From, query.Period.To).Scan(&base).Error; err != nil {
		return out, fmt.Errorf("get analytics statistics: %w", err)
	}
	out.Views, out.Reads, out.UniqueVisitors, out.ReturningVisitors, out.ContactViews, out.Downloads = base.Views, base.Reads, base.UniqueVisitors, base.ReturningVisitors, base.ContactViews, base.Downloads
	out.MedianActiveSeconds = int64(base.Median / 1000.0)
	out.Devices = map[string]int64{"desktop": 0, "mobile": 0, "tablet": 0, "unknown": 0}
	var deviceRows []struct {
		Device string
		Count  int64
	}
	if err := db.Raw(`SELECT device, COUNT(*) count FROM resume_link_views WHERE link_id=? AND traffic='human' AND started_at>=? AND started_at<? GROUP BY device`, query.LinkID.UUID(), query.Period.From, query.Period.To).Scan(&deviceRows).Error; err != nil {
		return out, fmt.Errorf("get analytics devices: %w", err)
	}
	for _, row := range deviceRows {
		out.Devices[row.Device] = row.Count
	}
	var refs []struct {
		Host  *string
		Count int64
	}
	if err := db.Raw(`SELECT NULLIF(referrer_host,'') host, COUNT(*) count FROM resume_link_views WHERE link_id=? AND traffic='human' AND started_at>=? AND started_at<? GROUP BY referrer_host ORDER BY count DESC`, query.LinkID.UUID(), query.Period.From, query.Period.To).Scan(&refs).Error; err != nil {
		return out, fmt.Errorf("get analytics referrers: %w", err)
	}
	out.Referrers = make([]queries.ReferrerStat, 0, len(refs))
	for _, x := range refs {
		out.Referrers = append(out.Referrers, queries.ReferrerStat{Host: x.Host, Views: x.Count})
	}
	var filtered struct{ Bot, Owner, Suspect int64 }
	if err := db.Raw(`SELECT COUNT(*) FILTER(WHERE traffic='bot') bot, COUNT(*) FILTER(WHERE traffic='owner') owner, COUNT(*) FILTER(WHERE traffic='suspect') suspect FROM resume_link_views WHERE link_id=? AND started_at>=? AND started_at<?`, query.LinkID.UUID(), query.Period.From, query.Period.To).Scan(&filtered).Error; err != nil {
		return out, fmt.Errorf("get analytics filtered: %w", err)
	}
	out.Filtered = queries.FilteredStatistics{Bot: filtered.Bot, Owner: filtered.Owner, Suspect: filtered.Suspect}
	return out, nil
}

func (r *ResumeLinkReader) GetDaily(ctx context.Context, query queries.LinkAnalyticsQuery) (queries.LinkDaily, error) {
	if r == nil || r.db == nil {
		return queries.LinkDaily{}, fmt.Errorf("get analytics daily: reader dependencies are nil")
	}
	if err := query.Period.Validate(); err != nil {
		return queries.LinkDaily{}, err
	}
	db := database.DBFromContext(ctx, r.db).WithContext(ctx)
	var rows []struct {
		Date                         string
		Views, Reads, UniqueVisitors int64
	}
	if err := db.Raw(`SELECT (started_at AT TIME ZONE ?)::date::text date, COUNT(*) views, COUNT(*) FILTER(WHERE active_ms>=5000) reads, COUNT(DISTINCT visitor_id) unique_visitors FROM resume_link_views WHERE link_id=? AND traffic='human' AND started_at>=? AND started_at<? GROUP BY 1 ORDER BY 1`, query.Period.TZ, query.LinkID.UUID(), query.Period.From, query.Period.To).Scan(&rows).Error; err != nil {
		return queries.LinkDaily{}, fmt.Errorf("get analytics daily: %w", err)
	}
	byDate := make(map[string]struct{ Views, Reads, UniqueVisitors int64 }, len(rows))
	for _, x := range rows {
		byDate[x.Date] = struct{ Views, Reads, UniqueVisitors int64 }{x.Views, x.Reads, x.UniqueVisitors}
	}
	days := query.Period.Days()
	out := queries.LinkDaily{Days: make([]queries.DailyPoint, 0, len(days))}
	for _, d := range days {
		x := byDate[d]
		out.Days = append(out.Days, queries.DailyPoint{Date: d, Views: x.Views, Reads: x.Reads, UniqueVisitors: x.UniqueVisitors})
	}
	return out, nil
}

func (r *ResumeLinkReader) GetVisits(ctx context.Context, query queries.LinkVisitsQuery) (queries.LinkVisitsPage, error) {
	if r == nil || r.db == nil {
		return queries.LinkVisitsPage{}, fmt.Errorf("get analytics visits: reader dependencies are nil")
	}
	limit, err := queries.NormalizeVisitLimit(query.Limit)
	if err != nil {
		return queries.LinkVisitsPage{}, err
	}
	cursorTime, cursorID, err := queries.DecodeVisitCursor(query.Cursor)
	if err != nil {
		return queries.LinkVisitsPage{}, err
	}
	db := database.DBFromContext(ctx, r.db).WithContext(ctx)
	trafficFilter := "traffic='human'"
	if query.IncludeFiltered {
		trafficFilter = "TRUE"
	}
	args := []any{query.LinkID.UUID()}
	cursorSQL := ""
	if !cursorTime.IsZero() {
		cursorSQL = " AND (v.started_at < ? OR (v.started_at = ? AND v.id < ?))"
		args = append(args, cursorTime, cursorTime, cursorID)
	}
	args = append(args, limit+1)
	sql := fmt.Sprintf(`WITH first_visitors AS (SELECT DISTINCT ON (visitor_id) visitor_id, id first_id, started_at first_at FROM resume_link_views WHERE link_id=? ORDER BY visitor_id, started_at, id), ranked AS (SELECT visitor_id, first_id, ROW_NUMBER() OVER (ORDER BY first_at, first_id) ordinal FROM first_visitors) SELECT v.id,v.started_at,v.visitor_id,v.device,v.browser,v.os,v.referrer_host,v.active_ms,v.max_section,v.exit_section,v.traffic,r.ordinal,EXISTS(SELECT 1 FROM resume_link_views p WHERE p.link_id=v.link_id AND p.visitor_id=v.visitor_id AND p.started_at<v.started_at) is_returning FROM resume_link_views v JOIN ranked r ON r.visitor_id=v.visitor_id WHERE v.link_id=? AND %s%s ORDER BY v.started_at DESC,v.id DESC LIMIT ?`, trafficFilter, cursorSQL)
	var rows []struct {
		ID, VisitorID, Device, Browser, OS, Traffic string
		StartedAt                                   time.Time
		ReferrerHost                                *string
		ActiveMS                                    int64
		MaxSection, ExitSection                     *string
		Ordinal                                     int64
		IsReturning                                 bool
	}
	if err := db.Raw(sql, args...).Scan(&rows).Error; err != nil {
		return queries.LinkVisitsPage{}, fmt.Errorf("get analytics visits: %w", err)
	}
	hasMore := len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}
	out := queries.LinkVisitsPage{Items: make([]queries.LinkVisit, 0, len(rows))}
	ids := make([]string, 0, len(rows))
	for _, x := range rows {
		out.Items = append(out.Items, queries.LinkVisit{ID: x.ID, StartedAt: x.StartedAt, Visitor: queries.VisitorInfo{Ordinal: x.Ordinal, IsReturning: x.IsReturning}, Device: x.Device, Browser: x.Browser, OS: x.OS, ReferrerHost: x.ReferrerHost, IsRead: x.ActiveMS >= 5000, ActiveSeconds: x.ActiveMS / 1000, MaxSection: x.MaxSection, ExitSection: x.ExitSection, Traffic: x.Traffic, Events: []queries.VisitEvent{}})
		ids = append(ids, x.ID)
	}
	if len(ids) > 0 {
		var evs []struct {
			ViewID, EventType string
			OccurredAt        time.Time
			Payload           []byte
		}
		if err := db.Raw(`SELECT view_id,event_type,occurred_at,payload FROM resume_link_events WHERE view_id IN ? AND event_type <> 'section_view' ORDER BY occurred_at,id`, ids).Scan(&evs).Error; err != nil {
			return queries.LinkVisitsPage{}, fmt.Errorf("get analytics visit events: %w", err)
		}
		byID := map[string][]queries.VisitEvent{}
		for _, e := range evs {
			var p map[string]string
			if err := json.Unmarshal(e.Payload, &p); err != nil {
				continue
			}
			byID[e.ViewID] = append(byID[e.ViewID], queries.VisitEvent{Type: e.EventType, At: e.OccurredAt, Contact: p["contact"], Technology: p["technology"], Kind: p["kind"], Target: p["target"]})
		}
		for i := range out.Items {
			out.Items[i].Events = byID[out.Items[i].ID]
		}
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

func (r *ResumeLinkReader) GetSections(ctx context.Context, query queries.LinkSectionsQuery) (queries.SectionsStatistics, error) {
	if r == nil || r.db == nil {
		return queries.SectionsStatistics{}, fmt.Errorf("get analytics sections: reader dependencies are nil")
	}
	if err := query.Period.Validate(); err != nil {
		return queries.SectionsStatistics{}, err
	}
	db := database.DBFromContext(ctx, r.db).WithContext(ctx)
	var reads int64
	if err := db.Raw(`SELECT COUNT(*) FROM resume_link_views WHERE link_id=? AND traffic='human' AND active_ms>=5000 AND started_at>=? AND started_at<?`, query.LinkID.UUID(), query.Period.From, query.Period.To).Scan(&reads).Error; err != nil {
		return queries.SectionsStatistics{}, fmt.Errorf("get analytics section reads: %w", err)
	}
	var rows []struct {
		Section        string
		Reached, Exits int64
	}
	if err := db.Raw(`SELECT e.payload->>'section' section, COUNT(DISTINCT e.view_id) FILTER(WHERE v.active_ms>=5000) reached, COUNT(DISTINCT v.id) FILTER(WHERE v.active_ms>=5000 AND v.exit_section=e.payload->>'section') exits FROM resume_link_events e JOIN resume_link_views v ON v.id=e.view_id WHERE e.link_id=? AND e.event_type='section_view' AND v.traffic='human' AND v.started_at>=? AND v.started_at<? GROUP BY 1`, query.LinkID.UUID(), query.Period.From, query.Period.To).Scan(&rows).Error; err != nil {
		return queries.SectionsStatistics{}, fmt.Errorf("get analytics sections: %w", err)
	}
	counts := map[string]queries.SectionStatistics{}
	for _, x := range rows {
		counts[x.Section] = queries.SectionStatistics{Section: x.Section, Reached: x.Reached, Exits: x.Exits}
	}
	out := queries.SectionsStatistics{Reads: reads}
	seen := map[string]bool{}
	for _, s := range query.SectionOrder {
		if s == "techStack" || seen[s] {
			continue
		}
		out.Sections = append(out.Sections, counts[s])
		out.Sections[len(out.Sections)-1].Section = s
		seen[s] = true
	}
	for _, x := range rows {
		if !seen[x.Section] {
			out.Sections = append(out.Sections, counts[x.Section])
		}
	}
	return out, nil
}

func (r *ResumeLinkReader) GetTechnologies(ctx context.Context, query queries.LinkAnalyticsQuery) (queries.TechnologiesStatistics, error) {
	if r == nil || r.db == nil {
		return queries.TechnologiesStatistics{}, fmt.Errorf("get analytics technologies: reader dependencies are nil")
	}
	if err := query.Period.Validate(); err != nil {
		return queries.TechnologiesStatistics{}, err
	}
	db := database.DBFromContext(ctx, r.db).WithContext(ctx)
	var rows []struct {
		Name                 string
		Selections, Visitors int64
	}
	if err := db.Raw(`SELECT e.payload->>'technology' name, COUNT(*) selections, COUNT(DISTINCT v.visitor_id) visitors FROM resume_link_events e JOIN resume_link_views v ON v.id=e.view_id WHERE e.link_id=? AND e.event_type='tech_filter' AND e.payload->>'action'='select' AND v.traffic='human' AND v.started_at>=? AND v.started_at<? GROUP BY 1 ORDER BY selections DESC,name LIMIT 20`, query.LinkID.UUID(), query.Period.From, query.Period.To).Scan(&rows).Error; err != nil {
		return queries.TechnologiesStatistics{}, fmt.Errorf("get analytics technologies: %w", err)
	}
	out := queries.TechnologiesStatistics{Technologies: make([]queries.TechnologyInterest, 0, len(rows))}
	for _, x := range rows {
		out.Technologies = append(out.Technologies, queries.TechnologyInterest{Name: x.Name, Selections: x.Selections, Visitors: x.Visitors})
	}
	return out, nil
}

func (r *ResumeLinkReader) GetResumeStatistics(ctx context.Context, query queries.ResumeAnalyticsQuery) (queries.ResumeStatistics, error) {
	if r == nil || r.db == nil {
		return queries.ResumeStatistics{}, fmt.Errorf("get resume analytics statistics: reader dependencies are nil")
	}
	if err := query.Period.Validate(); err != nil {
		return queries.ResumeStatistics{}, err
	}
	db := database.DBFromContext(ctx, r.db).WithContext(ctx)
	var base queries.LinkStatistics
	var firstLast struct {
		First *time.Time `gorm:"column:first_view_at"`
		Last  *time.Time `gorm:"column:last_view_at"`
	}
	if err := db.Raw(`SELECT MIN(v.started_at) first_view_at, MAX(v.started_at) last_view_at FROM resume_link_views v JOIN resume_links l ON l.id=v.link_id WHERE l.resume_id=? AND v.traffic='human'`, query.ResumeID.UUID()).Scan(&firstLast).Error; err != nil {
		return queries.ResumeStatistics{}, fmt.Errorf("get resume analytics first/last view: %w", err)
	}
	base.FirstViewAt, base.LastViewAt = firstLast.First, firstLast.Last
	var row struct {
		Views, Reads, UniqueVisitors, ReturningVisitors, ContactViews, Downloads int64
		Median                                                                   float64
	}
	if err := db.Raw(`SELECT COUNT(*) FILTER(WHERE v.traffic='human') views, COUNT(*) FILTER(WHERE v.traffic='human' AND v.active_ms>=5000) reads, COUNT(DISTINCT v.visitor_id) FILTER(WHERE v.traffic='human') unique_visitors, COUNT(DISTINCT v.visitor_id) FILTER(WHERE v.traffic='human' AND EXISTS(SELECT 1 FROM resume_link_views p JOIN resume_links pl ON pl.id=p.link_id WHERE pl.resume_id=? AND p.visitor_id=v.visitor_id AND p.traffic='human' AND p.started_at>=? AND p.started_at<? AND (p.started_at AT TIME ZONE ?)::date <> (v.started_at AT TIME ZONE ?)::date)) returning_visitors, COUNT(*) FILTER(WHERE v.traffic='human' AND EXISTS(SELECT 1 FROM resume_link_events e WHERE e.view_id=v.id AND e.event_type='contact_view' AND e.occurred_at>=? AND e.occurred_at<?)) contact_views, COUNT(*) FILTER(WHERE v.traffic='human' AND EXISTS(SELECT 1 FROM resume_link_events e WHERE e.view_id=v.id AND e.event_type='download' AND e.occurred_at>=? AND e.occurred_at<?)) downloads, COALESCE((SELECT percentile_cont(0.5) WITHIN GROUP(ORDER BY active_ms) FROM resume_link_views m JOIN resume_links lm ON lm.id=m.link_id WHERE lm.resume_id=? AND m.traffic='human' AND m.active_ms>=5000 AND m.started_at>=? AND m.started_at<?),0) median FROM resume_link_views v JOIN resume_links l ON l.id=v.link_id WHERE l.resume_id=? AND v.started_at>=? AND v.started_at<?`, query.ResumeID.UUID(), query.Period.From, query.Period.To, query.Period.TZ, query.Period.TZ, query.Period.From, query.Period.To, query.Period.From, query.Period.To, query.ResumeID.UUID(), query.Period.From, query.Period.To, query.ResumeID.UUID(), query.Period.From, query.Period.To).Scan(&row).Error; err != nil {
		return queries.ResumeStatistics{}, fmt.Errorf("get resume analytics statistics: %w", err)
	}
	base.Views, base.Reads, base.UniqueVisitors, base.ReturningVisitors, base.ContactViews, base.Downloads = row.Views, row.Reads, row.UniqueVisitors, row.ReturningVisitors, row.ContactViews, row.Downloads
	base.MedianActiveSeconds = int64(row.Median / 1000)
	base.Devices = map[string]int64{"desktop": 0, "mobile": 0, "tablet": 0, "unknown": 0}
	var dev []struct {
		Device string
		Count  int64
	}
	if err := db.Raw(`SELECT device,COUNT(*) count FROM resume_link_views v JOIN resume_links l ON l.id=v.link_id WHERE l.resume_id=? AND v.traffic='human' AND v.started_at>=? AND v.started_at<? GROUP BY device`, query.ResumeID.UUID(), query.Period.From, query.Period.To).Scan(&dev).Error; err != nil {
		return queries.ResumeStatistics{}, fmt.Errorf("get resume analytics devices: %w", err)
	}
	for _, x := range dev {
		base.Devices[x.Device] = x.Count
	}
	var refs []struct {
		Host  *string
		Count int64
	}
	if err := db.Raw(`SELECT NULLIF(v.referrer_host,'') host,COUNT(*) count FROM resume_link_views v JOIN resume_links l ON l.id=v.link_id WHERE l.resume_id=? AND v.traffic='human' AND v.started_at>=? AND v.started_at<? GROUP BY v.referrer_host ORDER BY count DESC`, query.ResumeID.UUID(), query.Period.From, query.Period.To).Scan(&refs).Error; err != nil {
		return queries.ResumeStatistics{}, fmt.Errorf("get resume analytics referrers: %w", err)
	}
	for _, x := range refs {
		base.Referrers = append(base.Referrers, queries.ReferrerStat{Host: x.Host, Views: x.Count})
	}
	var filtered queries.FilteredStatistics
	if err := db.Raw(`SELECT COUNT(*) FILTER(WHERE traffic='bot') bot,COUNT(*) FILTER(WHERE traffic='owner') owner,COUNT(*) FILTER(WHERE traffic='suspect') suspect FROM resume_link_views v JOIN resume_links l ON l.id=v.link_id WHERE l.resume_id=? AND v.started_at>=? AND v.started_at<?`, query.ResumeID.UUID(), query.Period.From, query.Period.To).Scan(&filtered).Error; err != nil {
		return queries.ResumeStatistics{}, fmt.Errorf("get resume analytics filtered: %w", err)
	}
	base.Filtered = filtered
	var links []queries.ResumeLinkStatistic
	if err := db.Raw(`SELECT l.id,l.name,COUNT(v.id) FILTER(WHERE v.traffic='human' AND v.started_at>=? AND v.started_at<?) views,COUNT(v.id) FILTER(WHERE v.traffic='human' AND v.active_ms>=5000 AND v.started_at>=? AND v.started_at<?) reads,MAX(v.started_at) FILTER(WHERE v.traffic='human') last_view_at FROM resume_links l LEFT JOIN resume_link_views v ON v.link_id=l.id WHERE l.resume_id=? GROUP BY l.id,l.name,l.created_at ORDER BY l.created_at,l.id`, query.Period.From, query.Period.To, query.Period.From, query.Period.To, query.ResumeID.UUID()).Scan(&links).Error; err != nil {
		return queries.ResumeStatistics{}, fmt.Errorf("get resume analytics links: %w", err)
	}
	baseLinks := append(make([]queries.ResumeLinkStatistic, 0, len(links)), links...)
	return queries.ResumeStatistics{LinkStatistics: base, Links: baseLinks}, nil
}

func (r *ResumeLinkReader) GetClicks(ctx context.Context, query queries.LinkClicksQuery) (queries.ClickGrid, error) {
	if r == nil || r.db == nil {
		return queries.ClickGrid{}, fmt.Errorf("get analytics clicks: reader dependencies are nil")
	}
	if query.LinkID.IsZero() || query.Section == "" || query.CurrentHash == "" {
		return queries.ClickGrid{}, fmt.Errorf("get analytics clicks: invalid query")
	}
	var breakpointSQL string
	switch query.Breakpoint {
	case "xs":
		breakpointSQL = "c.viewport_w < 600"
	case "md":
		breakpointSQL = "c.viewport_w >= 600 AND c.viewport_w < 1200"
	case "lg":
		breakpointSQL = "c.viewport_w >= 1200"
	default:
		return queries.ClickGrid{}, fmt.Errorf("get analytics clicks: invalid breakpoint")
	}

	db := database.DBFromContext(ctx, r.db).WithContext(ctx)
	layoutFilter := "c.layout_hash = ?"
	args := []any{query.LinkID.UUID(), query.Section, query.CurrentHash}
	if query.IncludeStale {
		layoutFilter = "TRUE"
		args = []any{query.LinkID.UUID(), query.Section}
	}
	where := fmt.Sprintf(`c.link_id = ? AND c.section = ? AND v.traffic = 'human' AND %s AND %s`, layoutFilter, breakpointSQL)

	var rows []struct {
		X     int64 `gorm:"column:x"`
		Y     int64 `gorm:"column:y"`
		Count int64 `gorm:"column:count"`
	}
	cellSQL := fmt.Sprintf(`
SELECT LEAST((c.x_pct * 20 / 10000), 19) AS x,
       LEAST((c.y_pct * 20 / 10000), 19) AS y,
       COUNT(*) AS count
FROM resume_link_clicks c
JOIN resume_link_views v ON v.id = c.view_id
WHERE %s
GROUP BY 1, 2
ORDER BY 1, 2`, where)
	cellArgs := args
	if err := db.Raw(cellSQL, cellArgs...).Scan(&rows).Error; err != nil {
		return queries.ClickGrid{}, fmt.Errorf("get analytics click cells: %w", err)
	}

	var total int64
	if err := db.Raw(fmt.Sprintf(`SELECT COUNT(*) FROM resume_link_clicks c JOIN resume_link_views v ON v.id=c.view_id WHERE %s`, where), args...).Scan(&total).Error; err != nil {
		return queries.ClickGrid{}, fmt.Errorf("get analytics click total: %w", err)
	}

	var staleTotal int64
	if !query.IncludeStale {
		staleWhere := fmt.Sprintf(`c.link_id = ? AND c.section = ? AND v.traffic = 'human' AND c.layout_hash <> ? AND %s`, breakpointSQL)
		if err := db.Raw(fmt.Sprintf(`SELECT COUNT(*) FROM resume_link_clicks c JOIN resume_link_views v ON v.id=c.view_id WHERE %s`, staleWhere), query.LinkID.UUID(), query.Section, query.CurrentHash).Scan(&staleTotal).Error; err != nil {
			return queries.ClickGrid{}, fmt.Errorf("get analytics stale click total: %w", err)
		}
	}

	cells := make([]queries.ClickCell, 0, len(rows))
	for _, row := range rows {
		cells = append(cells, queries.ClickCell{X: row.X, Y: row.Y, Count: row.Count})
	}
	return queries.ClickGrid{Section: query.Section, Grid: queries.ClickGridSize, Cells: cells, Total: total, StaleTotal: staleTotal}, nil
}

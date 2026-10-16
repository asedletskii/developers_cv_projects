package queries

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

var (
	ErrInvalidAnalyticsPeriod = fmt.Errorf("invalid analytics period")
	ErrInvalidAnalyticsCursor = fmt.Errorf("invalid analytics cursor")
	ErrInvalidAnalyticsLimit  = fmt.Errorf("invalid analytics limit")
)

type AnalyticsPeriod struct {
	From time.Time
	To   time.Time
	TZ   string
}

func NewAnalyticsPeriod(from, to, tz string, now time.Time) (AnalyticsPeriod, error) {
	if tz == "" {
		tz = "UTC"
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return AnalyticsPeriod{}, fmt.Errorf("%w: unknown timezone", ErrInvalidAnalyticsPeriod)
	}
	localNow := now.In(loc)
	today := time.Date(localNow.Year(), localNow.Month(), localNow.Day(), 0, 0, 0, 0, loc)
	if to == "" {
		to = today.Format("2006-01-02")
	}
	if from == "" {
		from = today.AddDate(0, 0, -29).Format("2006-01-02")
	}
	fromDate, err := time.ParseInLocation("2006-01-02", from, loc)
	if err != nil {
		return AnalyticsPeriod{}, fmt.Errorf("%w: invalid from", ErrInvalidAnalyticsPeriod)
	}
	toDate, err := time.ParseInLocation("2006-01-02", to, loc)
	if err != nil {
		return AnalyticsPeriod{}, fmt.Errorf("%w: invalid to", ErrInvalidAnalyticsPeriod)
	}
	if fromDate.After(toDate) {
		return AnalyticsPeriod{}, fmt.Errorf("%w: invalid range", ErrInvalidAnalyticsPeriod)
	}
	dayCount := 0
	for d := fromDate; d.Before(toDate); d = d.AddDate(0, 0, 1) {
		dayCount++
		if dayCount > 365 {
			return AnalyticsPeriod{}, fmt.Errorf("%w: period exceeds 366 calendar days", ErrInvalidAnalyticsPeriod)
		}
	}
	return AnalyticsPeriod{From: fromDate, To: toDate.AddDate(0, 0, 1), TZ: tz}, nil
}

func (p AnalyticsPeriod) Days() []string {
	loc, err := time.LoadLocation(p.TZ)
	if err != nil {
		return nil
	}
	start := p.From.In(loc)
	end := p.To.In(loc)
	result := make([]string, 0, calendarDayCount(start, end))
	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		result = append(result, d.Format("2006-01-02"))
	}
	return result
}

func calendarDayCount(start, end time.Time) int {
	count := 0
	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		count++
	}
	return count
}

type LinkAnalyticsQuery struct {
	UserID sharedvalueobjects.ID
	LinkID sharedvalueobjects.ID
	Period AnalyticsPeriod
}

type ResumeAnalyticsQuery struct {
	UserID   sharedvalueobjects.ID
	ResumeID sharedvalueobjects.ID
	Period   AnalyticsPeriod
}

type LinkVisitsQuery struct {
	UserID          sharedvalueobjects.ID
	LinkID          sharedvalueobjects.ID
	Cursor          string
	Limit           int
	IncludeFiltered bool
}

type LinkSectionsQuery struct {
	UserID       sharedvalueobjects.ID
	LinkID       sharedvalueobjects.ID
	Period       AnalyticsPeriod
	SectionOrder []string
}

type LinkStatistics struct {
	Views               int64              `json:"views"`
	Reads               int64              `json:"reads"`
	UniqueVisitors      int64              `json:"uniqueVisitors"`
	ReturningVisitors   int64              `json:"returningVisitors"`
	ContactViews        int64              `json:"contactViews"`
	Downloads           int64              `json:"downloads"`
	FirstViewAt         *time.Time         `json:"firstViewAt"`
	LastViewAt          *time.Time         `json:"lastViewAt"`
	MedianActiveSeconds int64              `json:"medianActiveSeconds"`
	Devices             map[string]int64   `json:"devices"`
	Referrers           []ReferrerStat     `json:"referrers"`
	Filtered            FilteredStatistics `json:"filtered"`
}

func (s *LinkStatistics) NormalizeResponse() {
	if s == nil {
		return
	}
	if s.Devices == nil {
		s.Devices = make(map[string]int64, 4)
	}
	for _, device := range []string{"desktop", "mobile", "tablet", "unknown"} {
		if _, ok := s.Devices[device]; !ok {
			s.Devices[device] = 0
		}
	}
	if s.Referrers == nil {
		s.Referrers = make([]ReferrerStat, 0)
	}
}

type FilteredStatistics struct {
	Bot     int64 `json:"bot"`
	Owner   int64 `json:"owner"`
	Suspect int64 `json:"suspect"`
}

type ReferrerStat struct {
	Host  *string `json:"host"`
	Views int64   `json:"views"`
}
type DailyPoint struct {
	Date           string `json:"date"`
	Views          int64  `json:"views"`
	Reads          int64  `json:"reads"`
	UniqueVisitors int64  `json:"uniqueVisitors"`
}
type LinkDaily struct {
	Days []DailyPoint `json:"days"`
}
type VisitEvent struct {
	Type       string    `json:"type"`
	At         time.Time `json:"at"`
	Contact    string    `json:"contact,omitempty"`
	Technology string    `json:"technology,omitempty"`
	Kind       string    `json:"kind,omitempty"`
	Target     string    `json:"target,omitempty"`
}
type LinkVisit struct {
	ID            string       `json:"id"`
	StartedAt     time.Time    `json:"startedAt"`
	Visitor       VisitorInfo  `json:"visitor"`
	Device        string       `json:"device"`
	Browser       string       `json:"browser"`
	OS            string       `json:"os"`
	ReferrerHost  *string      `json:"referrerHost"`
	IsRead        bool         `json:"isRead"`
	ActiveSeconds int64        `json:"activeSeconds"`
	MaxSection    *string      `json:"maxSection"`
	ExitSection   *string      `json:"exitSection"`
	Events        []VisitEvent `json:"events"`
	Traffic       string       `json:"traffic"`
}
type VisitorInfo struct {
	Ordinal     int64 `json:"ordinal"`
	IsReturning bool  `json:"isReturning"`
}
type LinkVisitsPage struct {
	Items      []LinkVisit `json:"items"`
	NextCursor *string     `json:"nextCursor"`
}
type SectionsStatistics struct {
	Reads    int64               `json:"reads"`
	Sections []SectionStatistics `json:"sections"`
}
type SectionStatistics struct {
	Section string `json:"section"`
	Reached int64  `json:"reached"`
	Exits   int64  `json:"exits"`
}
type TechnologyInterest struct {
	Name       string `json:"name"`
	Selections int64  `json:"selections"`
	Visitors   int64  `json:"visitors"`
}
type TechnologiesStatistics struct {
	Technologies []TechnologyInterest `json:"technologies"`
}
type ResumeStatistics struct {
	LinkStatistics
	Links []ResumeLinkStatistic `json:"links"`
}

func (s *ResumeStatistics) NormalizeResponse() {
	if s == nil {
		return
	}
	s.LinkStatistics.NormalizeResponse()
	if s.Links == nil {
		s.Links = make([]ResumeLinkStatistic, 0)
	}
}

type ResumeLinkStatistic struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Views      int64      `json:"views"`
	Reads      int64      `json:"reads"`
	LastViewAt *time.Time `json:"lastViewAt"`
}

type AnalyticsReadReader interface {
	GetLinkStatistics(context.Context, LinkAnalyticsQuery) (LinkStatistics, error)
	GetDaily(context.Context, LinkAnalyticsQuery) (LinkDaily, error)
	GetVisits(context.Context, LinkVisitsQuery) (LinkVisitsPage, error)
	GetSections(context.Context, LinkSectionsQuery) (SectionsStatistics, error)
	GetTechnologies(context.Context, LinkAnalyticsQuery) (TechnologiesStatistics, error)
	GetResumeStatistics(context.Context, ResumeAnalyticsQuery) (ResumeStatistics, error)
}

type AnalyticsReadHandler struct {
	links   ResumeLinkPublicReader
	resumes ports.ResumeOwnershipReader
	reader  AnalyticsReadReader
}

func NewAnalyticsReadHandler(links ResumeLinkPublicReader, resumes ports.ResumeOwnershipReader, reader AnalyticsReadReader) *AnalyticsReadHandler {
	return &AnalyticsReadHandler{links: links, resumes: resumes, reader: reader}
}

func (h *AnalyticsReadHandler) authorizeLink(ctx context.Context, userID, linkID sharedvalueobjects.ID) (ResumeLinkPublicView, resumequeries.ResumeView, error) {
	if h == nil || h.links == nil || h.resumes == nil || h.reader == nil {
		return ResumeLinkPublicView{}, resumequeries.ResumeView{}, fmt.Errorf("analytics read handler dependencies are nil")
	}
	if userID.IsZero() || linkID.IsZero() {
		return ResumeLinkPublicView{}, resumequeries.ResumeView{}, fmt.Errorf("invalid analytics query")
	}
	link, err := h.links.GetByID(ctx, PublicResumeLinkQuery{LinkID: linkID})
	if err != nil {
		return ResumeLinkPublicView{}, resumequeries.ResumeView{}, err
	}
	resumeID, err := sharedvalueobjects.ParseID(link.ResumeID)
	if err != nil {
		return ResumeLinkPublicView{}, resumequeries.ResumeView{}, fmt.Errorf("resolve resume link: %w", err)
	}
	resume, err := h.resumes.GetOwned(ctx, resumequeries.GetResumeQuery{ResumeID: resumeID, UserID: userID})
	if err != nil {
		return ResumeLinkPublicView{}, resumequeries.ResumeView{}, err
	}
	return link, resume, nil
}

func (h *AnalyticsReadHandler) Handle(ctx context.Context, q LinkAnalyticsQuery) (LinkStatistics, error) {
	return h.Statistics(ctx, q)
}

func (h *AnalyticsReadHandler) Statistics(ctx context.Context, q LinkAnalyticsQuery) (LinkStatistics, error) {
	_, _, err := h.authorizeLink(ctx, q.UserID, q.LinkID)
	if err != nil {
		return LinkStatistics{}, err
	}
	result, err := h.reader.GetLinkStatistics(ctx, q)
	if err != nil {
		return LinkStatistics{}, err
	}
	result.NormalizeResponse()
	return result, nil
}
func (h *AnalyticsReadHandler) Daily(ctx context.Context, q LinkAnalyticsQuery) (LinkDaily, error) {
	_, _, err := h.authorizeLink(ctx, q.UserID, q.LinkID)
	if err != nil {
		return LinkDaily{}, err
	}
	return h.reader.GetDaily(ctx, q)
}
func (h *AnalyticsReadHandler) Visits(ctx context.Context, q LinkVisitsQuery) (LinkVisitsPage, error) {
	_, _, err := h.authorizeLink(ctx, q.UserID, q.LinkID)
	if err != nil {
		return LinkVisitsPage{}, err
	}
	return h.reader.GetVisits(ctx, q)
}
func (h *AnalyticsReadHandler) Sections(ctx context.Context, q LinkSectionsQuery) (SectionsStatistics, error) {
	_, resume, err := h.authorizeLink(ctx, q.UserID, q.LinkID)
	if err != nil {
		return SectionsStatistics{}, err
	}
	if len(q.SectionOrder) == 0 {
		q.SectionOrder = currentSectionOrder(resume)
	}
	return h.reader.GetSections(ctx, q)
}
func (h *AnalyticsReadHandler) Technologies(ctx context.Context, q LinkAnalyticsQuery) (TechnologiesStatistics, error) {
	_, _, err := h.authorizeLink(ctx, q.UserID, q.LinkID)
	if err != nil {
		return TechnologiesStatistics{}, err
	}
	return h.reader.GetTechnologies(ctx, q)
}
func (h *AnalyticsReadHandler) ResumeStatistics(ctx context.Context, q ResumeAnalyticsQuery) (ResumeStatistics, error) {
	if h == nil || h.resumes == nil || h.reader == nil {
		return ResumeStatistics{}, fmt.Errorf("analytics read handler dependencies are nil")
	}
	if q.UserID.IsZero() || q.ResumeID.IsZero() {
		return ResumeStatistics{}, fmt.Errorf("invalid analytics query")
	}
	if _, err := h.resumes.GetOwned(ctx, resumequeries.GetResumeQuery{ResumeID: q.ResumeID, UserID: q.UserID}); err != nil {
		return ResumeStatistics{}, err
	}
	result, err := h.reader.GetResumeStatistics(ctx, q)
	if err != nil {
		return ResumeStatistics{}, err
	}
	result.NormalizeResponse()
	return result, nil
}

func currentSectionOrder(r resumequeries.ResumeView) []string {
	out := []string{"about"}
	seen := map[string]bool{"about": true}
	for _, s := range r.Overview.Sections {
		if s.Enabled && s.ID != "techStack" && !seen[s.ID] {
			out = append(out, s.ID)
			seen[s.ID] = true
		}
	}
	return out
}

func EncodeVisitCursor(t time.Time, id string) (string, error) {
	b, err := json.Marshal(struct {
		T  time.Time `json:"t"`
		ID string    `json:"id"`
	}{t, id})
	if err != nil {
		return "", fmt.Errorf("encode visit cursor: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
func DecodeVisitCursor(value string) (time.Time, string, error) {
	if value == "" {
		return time.Time{}, "", nil
	}
	b, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return time.Time{}, "", ErrInvalidAnalyticsCursor
	}
	var v struct {
		T  time.Time `json:"t"`
		ID string    `json:"id"`
	}
	if json.Unmarshal(b, &v) != nil || v.T.IsZero() || v.ID == "" {
		return time.Time{}, "", ErrInvalidAnalyticsCursor
	}
	return v.T, v.ID, nil
}
func NormalizeVisitLimit(v int) (int, error) {
	if v == 0 {
		return 20, nil
	}
	if v < 1 || v > 50 {
		return 0, ErrInvalidAnalyticsLimit
	}
	return v, nil
}
func NormalizeSectionOrder(values []string) []string {
	out := append([]string(nil), values...)
	sort.Strings(out)
	return out
}
func (p AnalyticsPeriod) Validate() error {
	if p.From.IsZero() || p.To.IsZero() || p.From.After(p.To) || strings.TrimSpace(p.TZ) == "" {
		return ErrInvalidAnalyticsPeriod
	}
	return nil
}

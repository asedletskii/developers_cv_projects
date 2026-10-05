package domain

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type TrafficType string

const (
	TrafficHuman   TrafficType = "human"
	TrafficBot     TrafficType = "bot"
	TrafficOwner   TrafficType = "owner"
	TrafficSuspect TrafficType = "suspect"
)

func (v TrafficType) IsValid() bool {
	switch v {
	case TrafficHuman, TrafficBot, TrafficOwner, TrafficSuspect:
		return true
	default:
		return false
	}
}

type DeviceType string

const (
	DeviceDesktop DeviceType = "desktop"
	DeviceMobile  DeviceType = "mobile"
	DeviceTablet  DeviceType = "tablet"
	DeviceUnknown DeviceType = "unknown"
)

func (v DeviceType) IsValid() bool {
	switch v {
	case DeviceDesktop, DeviceMobile, DeviceTablet, DeviceUnknown:
		return true
	default:
		return false
	}
}

type BrowserType string

const (
	BrowserChrome  BrowserType = "chrome"
	BrowserSafari  BrowserType = "safari"
	BrowserFirefox BrowserType = "firefox"
	BrowserEdge    BrowserType = "edge"
	BrowserYandex  BrowserType = "yandex"
	BrowserOther   BrowserType = "other"
)

func (v BrowserType) IsValid() bool {
	switch v {
	case BrowserChrome, BrowserSafari, BrowserFirefox, BrowserEdge, BrowserYandex, BrowserOther:
		return true
	default:
		return false
	}
}

type OSType string

const (
	OSWindows OSType = "windows"
	OSMacOS   OSType = "macos"
	OSLinux   OSType = "linux"
	OSIOS     OSType = "ios"
	OSAndroid OSType = "android"
	OSOther   OSType = "other"
)

func (v OSType) IsValid() bool {
	switch v {
	case OSWindows, OSMacOS, OSLinux, OSIOS, OSAndroid, OSOther:
		return true
	default:
		return false
	}
}

type LayoutSnapshot struct {
	order  []resumedomain.SectionID
	hashes map[resumedomain.SectionID]string
}

func NewLayoutSnapshot(order []resumedomain.SectionID, hashes map[resumedomain.SectionID]string) (LayoutSnapshot, error) {
	if len(order) == 0 {
		return LayoutSnapshot{}, fmt.Errorf("%w: layout order is empty", ErrInvalidLayoutSnapshot)
	}
	seen := make(map[resumedomain.SectionID]struct{}, len(order))
	normalizedOrder := make([]resumedomain.SectionID, 0, len(order))
	for _, section := range order {
		if !section.IsValid() || section == resumedomain.SectionTechStack {
			return LayoutSnapshot{}, fmt.Errorf("%w: invalid section %q", ErrInvalidLayoutSnapshot, section)
		}
		if _, ok := seen[section]; ok {
			return LayoutSnapshot{}, fmt.Errorf("%w: duplicate section %q", ErrInvalidLayoutSnapshot, section)
		}
		seen[section] = struct{}{}
		normalizedOrder = append(normalizedOrder, section)
	}

	normalizedHashes := make(map[resumedomain.SectionID]string, len(hashes))
	for _, section := range normalizedOrder {
		if _, ok := hashes[section]; !ok {
			return LayoutSnapshot{}, fmt.Errorf("%w: missing hash for section %q", ErrInvalidLayoutSnapshot, section)
		}
	}
	for section, hash := range hashes {
		if _, ok := seen[section]; !ok {
			return LayoutSnapshot{}, fmt.Errorf("%w: hash for hidden section %q", ErrInvalidLayoutSnapshot, section)
		}
		if len(hash) != 16 || strings.Trim(hash, "0123456789abcdef") != "" {
			return LayoutSnapshot{}, fmt.Errorf("%w: invalid hash for section %q", ErrInvalidLayoutSnapshot, section)
		}
		normalizedHashes[section] = hash
	}

	return LayoutSnapshot{order: normalizedOrder, hashes: normalizedHashes}, nil
}

func (s LayoutSnapshot) Order() []resumedomain.SectionID {
	return append([]resumedomain.SectionID(nil), s.order...)
}

func (s LayoutSnapshot) Hashes() map[resumedomain.SectionID]string {
	result := make(map[resumedomain.SectionID]string, len(s.hashes))
	for key, value := range s.hashes {
		result[key] = value
	}
	return result
}

func (s LayoutSnapshot) IndexOf(section resumedomain.SectionID) int {
	for i, value := range s.order {
		if value == section {
			return i
		}
	}
	return -1
}

func (s LayoutSnapshot) HashFor(section resumedomain.SectionID) (string, bool) {
	value, ok := s.hashes[section]
	return value, ok
}

func (s LayoutSnapshot) MarshalJSON() ([]byte, error) {
	type snapshot struct {
		Order  []resumedomain.SectionID          `json:"order"`
		Hashes map[resumedomain.SectionID]string `json:"hashes"`
	}
	return json.Marshal(snapshot{Order: s.Order(), Hashes: s.Hashes()})
}

type ViewState struct {
	lastSeq     int64
	confirmedAt time.Time
	lastSeenAt  time.Time
	activeMS    int64
	maxSection  resumedomain.SectionID
	exitSection resumedomain.SectionID
	endedAt     time.Time
}

type ViewStateUpdate struct {
	Seq         int64
	ActiveMS    int64
	MaxSection  resumedomain.SectionID
	ExitSection resumedomain.SectionID
	SeenAt      time.Time
	Final       bool
}

func (s *ViewState) Apply(update ViewStateUpdate, layout LayoutSnapshot) error {
	if s == nil {
		return fmt.Errorf("%w: state is nil", ErrInvalidViewState)
	}
	if update.Seq <= 0 {
		return fmt.Errorf("%w: seq must be positive", ErrInvalidViewState)
	}
	if update.ActiveMS < 0 {
		return fmt.Errorf("%w: active_ms is negative", ErrInvalidViewState)
	}
	if update.Seq <= s.lastSeq {
		return nil
	}
	if update.SeenAt.IsZero() {
		return fmt.Errorf("%w: seen_at is empty", ErrInvalidViewState)
	}
	if update.MaxSection != "" && layout.IndexOf(update.MaxSection) < 0 {
		return fmt.Errorf("%w: max_section is not in layout", ErrInvalidViewState)
	}
	if update.ExitSection != "" && layout.IndexOf(update.ExitSection) < 0 {
		return fmt.Errorf("%w: exit_section is not in layout", ErrInvalidViewState)
	}

	s.lastSeq = update.Seq
	if update.ActiveMS > s.activeMS {
		s.activeMS = update.ActiveMS
	}
	if s.confirmedAt.IsZero() && s.activeMS >= 5000 {
		s.confirmedAt = update.SeenAt.UTC()
	}
	if s.lastSeenAt.IsZero() || update.SeenAt.After(s.lastSeenAt) {
		s.lastSeenAt = update.SeenAt.UTC()
	}
	if update.MaxSection != "" && (s.maxSection == "" || layout.IndexOf(update.MaxSection) > layout.IndexOf(s.maxSection)) {
		s.maxSection = update.MaxSection
	}
	if update.ExitSection != "" {
		s.exitSection = update.ExitSection
	}
	if update.Final {
		s.endedAt = update.SeenAt.UTC()
	} else if !s.endedAt.IsZero() {
		s.endedAt = time.Time{}
	}
	return nil
}

func (s ViewState) LastSeq() int64                      { return s.lastSeq }
func (s ViewState) ConfirmedAt() time.Time              { return s.confirmedAt }
func (s ViewState) LastSeenAt() time.Time               { return s.lastSeenAt }
func (s ViewState) ActiveMS() int64                     { return s.activeMS }
func (s ViewState) MaxSection() resumedomain.SectionID  { return s.maxSection }
func (s ViewState) ExitSection() resumedomain.SectionID { return s.exitSection }
func (s ViewState) EndedAt() time.Time                  { return s.endedAt }

type ViewParams struct {
	ID            sharedvalueobjects.ID
	LinkID        sharedvalueobjects.ID
	VisitorID     sharedvalueobjects.ID
	IPHash        []byte
	Device        DeviceType
	Browser       BrowserType
	OS            OSType
	ReferrerHost  string
	Traffic       TrafficType
	TrafficReason string
	Layout        LayoutSnapshot
	StartedAt     time.Time
	CreatedAt     time.Time
}

type View struct {
	id            sharedvalueobjects.ID
	linkID        sharedvalueobjects.ID
	visitorID     sharedvalueobjects.ID
	ipHash        []byte
	device        DeviceType
	browser       BrowserType
	os            OSType
	referrerHost  string
	traffic       TrafficType
	trafficReason string
	layout        LayoutSnapshot
	startedAt     time.Time
	createdAt     time.Time
	state         ViewState
}

func NewView(params ViewParams) (*View, error) {
	switch {
	case params.ID.IsZero(), params.LinkID.IsZero(), params.VisitorID.IsZero():
		return nil, fmt.Errorf("%w: identity is empty", ErrInvalidView)
	case len(params.IPHash) != 32:
		return nil, fmt.Errorf("%w: ip hash must be 32 bytes", ErrInvalidView)
	case !params.Device.IsValid(), !params.Browser.IsValid(), !params.OS.IsValid(), !params.Traffic.IsValid():
		return nil, fmt.Errorf("%w: invalid classification", ErrInvalidView)
	case params.StartedAt.IsZero():
		return nil, fmt.Errorf("%w: started_at is empty", ErrInvalidView)
	}
	if _, err := NewLayoutSnapshot(params.Layout.Order(), params.Layout.Hashes()); err != nil {
		return nil, err
	}
	createdAt := params.CreatedAt.UTC()
	if params.CreatedAt.IsZero() {
		createdAt = time.Now().UTC()
	}
	return &View{
		id: params.ID, linkID: params.LinkID, visitorID: params.VisitorID,
		ipHash: append([]byte(nil), params.IPHash...), device: params.Device,
		browser: params.Browser, os: params.OS, referrerHost: strings.TrimSpace(params.ReferrerHost),
		traffic: params.Traffic, trafficReason: strings.TrimSpace(params.TrafficReason),
		layout: params.Layout, startedAt: params.StartedAt.UTC(), createdAt: createdAt,
	}, nil
}

func (v *View) ID() sharedvalueobjects.ID {
	if v == nil {
		return sharedvalueobjects.ID{}
	}
	return v.id
}
func (v *View) LinkID() sharedvalueobjects.ID {
	if v == nil {
		return sharedvalueobjects.ID{}
	}
	return v.linkID
}
func (v *View) VisitorID() sharedvalueobjects.ID {
	if v == nil {
		return sharedvalueobjects.ID{}
	}
	return v.visitorID
}
func (v *View) IPHash() []byte {
	if v == nil {
		return nil
	}
	return append([]byte(nil), v.ipHash...)
}
func (v *View) Device() DeviceType {
	if v == nil {
		return ""
	}
	return v.device
}
func (v *View) Browser() BrowserType {
	if v == nil {
		return ""
	}
	return v.browser
}
func (v *View) OS() OSType {
	if v == nil {
		return ""
	}
	return v.os
}
func (v *View) ReferrerHost() string {
	if v == nil {
		return ""
	}
	return v.referrerHost
}
func (v *View) Traffic() TrafficType {
	if v == nil {
		return ""
	}
	return v.traffic
}
func (v *View) TrafficReason() string {
	if v == nil {
		return ""
	}
	return v.trafficReason
}
func (v *View) Layout() LayoutSnapshot {
	if v == nil {
		return LayoutSnapshot{}
	}
	return v.layout
}
func (v *View) StartedAt() time.Time {
	if v == nil {
		return time.Time{}
	}
	return v.startedAt
}
func (v *View) CreatedAt() time.Time {
	if v == nil {
		return time.Time{}
	}
	return v.createdAt
}
func (v *View) State() ViewState {
	if v == nil {
		return ViewState{}
	}
	return v.state
}

func (v *View) RestoreState(state ViewState) error {
	if v == nil {
		return fmt.Errorf("%w: view is nil", ErrInvalidView)
	}
	if state.activeMS < 0 {
		return fmt.Errorf("%w: active_ms is negative", ErrInvalidViewState)
	}
	if state.maxSection != "" && v.layout.IndexOf(state.maxSection) < 0 {
		return fmt.Errorf("%w: max_section is not in layout", ErrInvalidViewState)
	}
	if state.exitSection != "" && v.layout.IndexOf(state.exitSection) < 0 {
		return fmt.Errorf("%w: exit_section is not in layout", ErrInvalidViewState)
	}
	v.state = state
	return nil
}

func (v *View) ApplyState(update ViewStateUpdate) error {
	if v == nil {
		return fmt.Errorf("%w: view is nil", ErrInvalidView)
	}
	return v.state.Apply(update, v.layout)
}

type AnalyticsClickParams struct {
	ID         sharedvalueobjects.ID
	ViewID     sharedvalueobjects.ID
	LinkID     sharedvalueobjects.ID
	Section    resumedomain.SectionID
	XPct       int
	YPct       int
	ViewportW  int
	LayoutHash string
	Target     string
	OccurredAt time.Time
}

type AnalyticsClick struct {
	id         sharedvalueobjects.ID
	viewID     sharedvalueobjects.ID
	linkID     sharedvalueobjects.ID
	section    resumedomain.SectionID
	xPct       int
	yPct       int
	viewportW  int
	layoutHash string
	target     string
	occurredAt time.Time
}

func NewAnalyticsClick(params AnalyticsClickParams) (*AnalyticsClick, error) {
	if params.ID.IsZero() || params.ViewID.IsZero() || params.LinkID.IsZero() {
		return nil, fmt.Errorf("%w: invalid ids", ErrInvalidAnalyticsClick)
	}
	if !params.Section.IsValid() || params.Section == resumedomain.SectionTechStack {
		return nil, fmt.Errorf("%w: invalid section", ErrInvalidAnalyticsClick)
	}
	if params.XPct < 0 || params.XPct > 10000 || params.YPct < 0 || params.YPct > 10000 {
		return nil, fmt.Errorf("%w: coordinates out of range", ErrInvalidAnalyticsClick)
	}
	if params.ViewportW <= 0 || params.ViewportW > 32767 {
		return nil, fmt.Errorf("%w: invalid viewport width", ErrInvalidAnalyticsClick)
	}
	if len(params.LayoutHash) != 16 || strings.Trim(params.LayoutHash, "0123456789abcdef") != "" {
		return nil, fmt.Errorf("%w: invalid layout hash", ErrInvalidAnalyticsClick)
	}
	if len(params.Target) > 64 {
		return nil, fmt.Errorf("%w: target is too long", ErrInvalidAnalyticsClick)
	}
	if params.OccurredAt.IsZero() {
		return nil, fmt.Errorf("%w: occurred_at is empty", ErrInvalidAnalyticsClick)
	}
	return &AnalyticsClick{
		id: params.ID, viewID: params.ViewID, linkID: params.LinkID, section: params.Section,
		xPct: params.XPct, yPct: params.YPct, viewportW: params.ViewportW, layoutHash: params.LayoutHash,
		target: params.Target, occurredAt: params.OccurredAt.UTC(),
	}, nil
}

func (c *AnalyticsClick) ID() sharedvalueobjects.ID       { return c.id }
func (c *AnalyticsClick) ViewID() sharedvalueobjects.ID   { return c.viewID }
func (c *AnalyticsClick) LinkID() sharedvalueobjects.ID   { return c.linkID }
func (c *AnalyticsClick) Section() resumedomain.SectionID { return c.section }
func (c *AnalyticsClick) XPct() int                       { return c.xPct }
func (c *AnalyticsClick) YPct() int                       { return c.yPct }
func (c *AnalyticsClick) ViewportW() int                  { return c.viewportW }
func (c *AnalyticsClick) LayoutHash() string              { return c.layoutHash }
func (c *AnalyticsClick) Target() string                  { return c.target }
func (c *AnalyticsClick) OccurredAt() time.Time           { return c.occurredAt }

type AnalyticsEventType string

const (
	EventTypeContactView   AnalyticsEventType = "contact_view"
	EventTypeSectionView   AnalyticsEventType = "section_view"
	EventTypeDownload      AnalyticsEventType = "download"
	EventTypeTechFilter    AnalyticsEventType = "tech_filter"
	EventTypeExternalClick AnalyticsEventType = "external_click"
	EventTypePrint         AnalyticsEventType = "print"
	EventTypeContactCopy   AnalyticsEventType = "contact_copy"
	EventTypeContactClick  AnalyticsEventType = "contact_click"
)

func (v AnalyticsEventType) IsValid() bool {
	switch v {
	case EventTypeContactView, EventTypeSectionView, EventTypeDownload, EventTypeTechFilter, EventTypeExternalClick, EventTypePrint, EventTypeContactCopy, EventTypeContactClick:
		return true
	default:
		return false
	}
}

type AnalyticsEventParams struct {
	ID            sharedvalueobjects.ID
	ViewID        sharedvalueobjects.ID
	LinkID        sharedvalueobjects.ID
	ClientEventID *sharedvalueobjects.ID
	Type          AnalyticsEventType
	Payload       map[string]string
	OccurredAt    time.Time
}

type AnalyticsEvent struct {
	id            sharedvalueobjects.ID
	viewID        sharedvalueobjects.ID
	linkID        sharedvalueobjects.ID
	clientEventID *sharedvalueobjects.ID
	eventType     AnalyticsEventType
	payload       map[string]string
	occurredAt    time.Time
}

func NewAnalyticsEvent(params AnalyticsEventParams) (*AnalyticsEvent, error) {
	if params.ID.IsZero() || params.ViewID.IsZero() || params.LinkID.IsZero() || !params.Type.IsValid() || params.OccurredAt.IsZero() {
		return nil, fmt.Errorf("%w: invalid identity, type or timestamp", ErrInvalidAnalyticsEvent)
	}
	payload, err := normalizeEventPayload(params.Type, params.Payload)
	if err != nil {
		return nil, err
	}
	var clientEventID *sharedvalueobjects.ID
	if params.ClientEventID != nil {
		if params.ClientEventID.IsZero() {
			return nil, fmt.Errorf("%w: client event id is empty", ErrInvalidAnalyticsEvent)
		}
		id := *params.ClientEventID
		clientEventID = &id
	}
	return &AnalyticsEvent{id: params.ID, viewID: params.ViewID, linkID: params.LinkID, clientEventID: clientEventID, eventType: params.Type, payload: payload, occurredAt: params.OccurredAt.UTC()}, nil
}

func normalizeEventPayload(eventType AnalyticsEventType, input map[string]string) (map[string]string, error) {
	allowed := map[AnalyticsEventType]map[string]bool{
		EventTypeContactView: {"contact": true}, EventTypeSectionView: {"section": true},
		EventTypeDownload: {}, EventTypeTechFilter: {"technology": true, "action": true},
		EventTypeExternalClick: {"kind": true, "target": true}, EventTypePrint: {},
		EventTypeContactCopy: {"contact": true}, EventTypeContactClick: {"contact": true},
	}[eventType]
	result := make(map[string]string)
	for key, value := range input {
		if allowed[key] {
			result[key] = truncateEventValue(value)
		}
	}
	switch eventType {
	case EventTypeContactView, EventTypeContactCopy, EventTypeContactClick:
		if !resumedomain.ContactKey(result["contact"]).IsValid() {
			return nil, fmt.Errorf("%w: invalid contact", ErrInvalidAnalyticsEvent)
		}
	case EventTypeSectionView:
		section := resumedomain.SectionID(result["section"])
		if !section.IsValid() || section == resumedomain.SectionTechStack {
			return nil, fmt.Errorf("%w: invalid section", ErrInvalidAnalyticsEvent)
		}
	case EventTypeTechFilter:
		if strings.TrimSpace(result["technology"]) == "" || (result["action"] != "select" && result["action"] != "clear") {
			return nil, fmt.Errorf("%w: invalid technology filter", ErrInvalidAnalyticsEvent)
		}
	case EventTypeExternalClick:
		if strings.TrimSpace(result["kind"]) == "" || strings.TrimSpace(result["target"]) == "" {
			return nil, fmt.Errorf("%w: invalid external click", ErrInvalidAnalyticsEvent)
		}
	case EventTypeDownload, EventTypePrint:
		if len(result) != 0 {
			return nil, fmt.Errorf("%w: payload is not allowed", ErrInvalidAnalyticsEvent)
		}
	}
	return result, nil
}

func truncateEventValue(value string) string {
	value = strings.TrimSpace(value)
	if len(value) <= 64 {
		return value
	}
	return value[:64]
}

func (e *AnalyticsEvent) ID() sharedvalueobjects.ID {
	if e == nil {
		return sharedvalueobjects.ID{}
	}
	return e.id
}
func (e *AnalyticsEvent) ViewID() sharedvalueobjects.ID {
	if e == nil {
		return sharedvalueobjects.ID{}
	}
	return e.viewID
}
func (e *AnalyticsEvent) LinkID() sharedvalueobjects.ID {
	if e == nil {
		return sharedvalueobjects.ID{}
	}
	return e.linkID
}
func (e *AnalyticsEvent) ClientEventID() *sharedvalueobjects.ID {
	if e == nil || e.clientEventID == nil {
		return nil
	}
	id := *e.clientEventID
	return &id
}
func (e *AnalyticsEvent) Type() AnalyticsEventType {
	if e == nil {
		return ""
	}
	return e.eventType
}
func (e *AnalyticsEvent) Payload() map[string]string {
	if e == nil {
		return nil
	}
	result := make(map[string]string, len(e.payload))
	for k, v := range e.payload {
		result[k] = v
	}
	return result
}
func (e *AnalyticsEvent) OccurredAt() time.Time {
	if e == nil {
		return time.Time{}
	}
	return e.occurredAt
}

func ViewStateFromPersistence(lastSeq int64, confirmedAt, lastSeenAt *time.Time, activeMS int64, maxSection, exitSection string, endedAt *time.Time) ViewState {
	state := ViewState{lastSeq: lastSeq, activeMS: activeMS, maxSection: resumedomain.SectionID(maxSection), exitSection: resumedomain.SectionID(exitSection)}
	if confirmedAt != nil {
		state.confirmedAt = confirmedAt.UTC()
	}
	if lastSeenAt != nil {
		state.lastSeenAt = lastSeenAt.UTC()
	}
	if endedAt != nil {
		state.endedAt = endedAt.UTC()
	}
	return state
}

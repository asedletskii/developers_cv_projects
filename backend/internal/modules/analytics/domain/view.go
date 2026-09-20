package domain

import (
	"fmt"
	"strings"
	"time"

	analyticscontracts "developers-cv/backend/internal/modules/analytics/contracts"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ViewParams struct {
	LinkID        sharedvalueobjects.ID
	ResumeID      sharedvalueobjects.ID
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
	resumeID      sharedvalueobjects.ID
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

type ViewSnapshot struct {
	ID            sharedvalueobjects.ID
	LinkID        sharedvalueobjects.ID
	ResumeID      sharedvalueobjects.ID
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
	State         ViewStateSnapshot
}

func NewView(params ViewParams) (*View, error) {
	id := sharedvalueobjects.NewID()
	params.ReferrerHost = strings.TrimSpace(params.ReferrerHost)
	params.TrafficReason = strings.TrimSpace(params.TrafficReason)
	params.StartedAt = params.StartedAt.UTC()
	if params.CreatedAt.IsZero() {
		params.CreatedAt = params.StartedAt
	} else {
		params.CreatedAt = params.CreatedAt.UTC()
	}
	if err := validateView(id, params); err != nil {
		return nil, err
	}
	return &View{
		id:            id,
		linkID:        params.LinkID,
		resumeID:      params.ResumeID,
		visitorID:     params.VisitorID,
		ipHash:        append([]byte(nil), params.IPHash...),
		device:        params.Device,
		browser:       params.Browser,
		os:            params.OS,
		referrerHost:  params.ReferrerHost,
		traffic:       params.Traffic,
		trafficReason: params.TrafficReason,
		layout:        params.Layout,
		startedAt:     params.StartedAt,
		createdAt:     params.CreatedAt,
	}, nil
}

func RestoreView(snapshot ViewSnapshot) (*View, error) {
	params := ViewParams{
		LinkID: snapshot.LinkID, ResumeID: snapshot.ResumeID, VisitorID: snapshot.VisitorID,
		IPHash: snapshot.IPHash, Device: snapshot.Device, Browser: snapshot.Browser, OS: snapshot.OS,
		ReferrerHost: snapshot.ReferrerHost, Traffic: snapshot.Traffic, TrafficReason: snapshot.TrafficReason,
		Layout: snapshot.Layout, StartedAt: snapshot.StartedAt, CreatedAt: snapshot.CreatedAt,
	}
	if err := validateView(snapshot.ID, params); err != nil {
		return nil, err
	}
	view := &View{
		id: snapshot.ID, linkID: snapshot.LinkID, resumeID: snapshot.ResumeID, visitorID: snapshot.VisitorID,
		ipHash: append([]byte(nil), snapshot.IPHash...), device: snapshot.Device,
		browser: snapshot.Browser, os: snapshot.OS, referrerHost: snapshot.ReferrerHost,
		traffic: snapshot.Traffic, trafficReason: snapshot.TrafficReason, layout: snapshot.Layout,
		startedAt: snapshot.StartedAt, createdAt: snapshot.CreatedAt,
		state: ViewState{
			lastSeq: snapshot.State.LastSeq, activeMS: snapshot.State.ActiveMS,
			maxSection: snapshot.State.MaxSection, exitSection: snapshot.State.ExitSection,
			confirmedAt: snapshot.State.ConfirmedAt.UTC(), lastSeenAt: snapshot.State.LastSeenAt.UTC(),
			endedAt: snapshot.State.EndedAt.UTC(),
		},
	}
	if err := view.RestoreState(view.state); err != nil {
		return nil, err
	}
	return view, nil
}

func validateView(id sharedvalueobjects.ID, params ViewParams) error {
	switch {
	case id.IsZero(), params.LinkID.IsZero(), params.ResumeID.IsZero(), params.VisitorID.IsZero():
		return fmt.Errorf("%w: identity is empty", ErrInvalidView)
	case len(params.IPHash) != 32:
		return fmt.Errorf("%w: ip hash must be 32 bytes", ErrInvalidView)
	case !params.Device.IsValid(), !params.Browser.IsValid(), !params.OS.IsValid(), !params.Traffic.IsValid():
		return fmt.Errorf("%w: invalid classification", ErrInvalidView)
	case params.StartedAt.IsZero():
		return fmt.Errorf("%w: started_at is empty", ErrInvalidView)
	case params.CreatedAt.IsZero():
		return fmt.Errorf("%w: created_at is empty", ErrInvalidView)
	case len(params.ReferrerHost) > MaxReferrerHostLength:
		return fmt.Errorf("%w: referrer_host exceeds %d bytes", ErrInvalidView, MaxReferrerHostLength)
	case len(params.TrafficReason) > MaxTrafficReasonLength:
		return fmt.Errorf("%w: traffic_reason exceeds %d bytes", ErrInvalidView, MaxTrafficReasonLength)
	}
	if params.Layout.HasTechnologyCatalog() {
		if _, err := NewLayoutSnapshot(params.Layout.Order(), params.Layout.Hashes(), params.Layout.Technologies()); err != nil {
			return err
		}
	} else if _, err := NewLayoutSnapshot(params.Layout.Order(), params.Layout.Hashes()); err != nil {
		return err
	}
	return nil
}

func (v *View) ID() sharedvalueobjects.ID        { return v.id }
func (v *View) LinkID() sharedvalueobjects.ID    { return v.linkID }
func (v *View) ResumeID() sharedvalueobjects.ID  { return v.resumeID }
func (v *View) VisitorID() sharedvalueobjects.ID { return v.visitorID }
func (v *View) IPHash() []byte                   { return append([]byte(nil), v.ipHash...) }
func (v *View) Device() DeviceType               { return v.device }
func (v *View) Browser() BrowserType             { return v.browser }
func (v *View) OS() OSType                       { return v.os }
func (v *View) ReferrerHost() string             { return v.referrerHost }
func (v *View) Traffic() TrafficType             { return v.traffic }
func (v *View) TrafficReason() string            { return v.trafficReason }
func (v *View) Layout() LayoutSnapshot           { return v.layout }
func (v *View) StartedAt() time.Time             { return v.startedAt }
func (v *View) CreatedAt() time.Time             { return v.createdAt }
func (v *View) State() ViewState                 { return v.state }

func (v *View) RestoreState(state ViewState) error {
	if v == nil {
		return fmt.Errorf("%w: view is nil", ErrInvalidView)
	}
	if err := state.validate(v.layout); err != nil {
		return err
	}
	v.state = state
	return nil
}

const (
	maxActiveMSClockSkew   = 30 * time.Second
	analyticsViewTTL       = 12 * time.Hour
	MaxReferrerHostLength  = 253
	MaxTrafficReasonLength = 100
	MaxActiveMS            = int64((analyticsViewTTL + maxActiveMSClockSkew) / time.Millisecond)
)

type AnalyticsBatchDecision struct {
	Ignore  bool
	SameSeq bool
}

func (v *View) EvaluateAnalyticsBatch(
	linkActive bool,
	resumePublished bool,
	now time.Time,
	update ViewStateUpdate,
) (AnalyticsBatchDecision, error) {
	if v == nil {
		return AnalyticsBatchDecision{}, fmt.Errorf("%w: view is nil", ErrInvalidView)
	}
	if !linkActive {
		return AnalyticsBatchDecision{}, ErrAnalyticsLinkInactive
	}
	if !resumePublished {
		return AnalyticsBatchDecision{}, analyticscontracts.ErrAnalyticsSubjectNotFound
	}
	if now.IsZero() {
		return AnalyticsBatchDecision{}, fmt.Errorf("%w: current time is empty", ErrInvalidViewState)
	}
	if v.startedAt.Before(now.UTC().Add(-analyticsViewTTL)) {
		return AnalyticsBatchDecision{}, ErrAnalyticsViewExpired
	}
	if v.traffic == TrafficBot {
		return AnalyticsBatchDecision{Ignore: true}, nil
	}
	lastSeq := v.state.LastSeq()
	if update.Seq < lastSeq || (update.Seq == lastSeq && !update.Final) {
		return AnalyticsBatchDecision{Ignore: true}, nil
	}
	return AnalyticsBatchDecision{SameSeq: update.Seq == lastSeq}, nil
}

func (v *View) ApplyState(update ViewStateUpdate) error {
	if v == nil {
		return fmt.Errorf("%w: view is nil", ErrInvalidView)
	}
	if update.Seq >= v.state.LastSeq() && !update.SeenAt.IsZero() {
		elapsed := update.SeenAt.Sub(v.startedAt)
		if elapsed < 0 {
			elapsed = 0
		}
		maxActiveMS := elapsed.Milliseconds() + maxActiveMSClockSkew.Milliseconds()
		if update.ActiveMS > maxActiveMS {
			return fmt.Errorf("%w: active_ms exceeds elapsed view time", ErrInvalidViewState)
		}
	}
	return v.state.Apply(update, v.layout)
}

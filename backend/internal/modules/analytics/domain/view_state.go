package domain

import (
	"fmt"
	"time"
)

const ViewConfirmationActiveMS int64 = 5000

type ViewState struct {
	lastSeq     int64
	confirmedAt time.Time
	lastSeenAt  time.Time
	activeMS    int64
	maxSection  AnalyticsSectionID
	exitSection AnalyticsSectionID
	endedAt     time.Time
}

type ViewStateUpdate struct {
	Seq         int64
	ActiveMS    int64
	MaxSection  AnalyticsSectionID
	ExitSection AnalyticsSectionID
	SeenAt      time.Time
	Final       bool
}

type ViewStateSnapshot struct {
	LastSeq     int64
	ConfirmedAt time.Time
	LastSeenAt  time.Time
	ActiveMS    int64
	MaxSection  AnalyticsSectionID
	ExitSection AnalyticsSectionID
	EndedAt     time.Time
}

func (s ViewState) validate(layout LayoutSnapshot) error {
	if s.lastSeq < 0 {
		return fmt.Errorf("%w: last sequence is negative", ErrInvalidViewState)
	}
	if s.activeMS < 0 || s.activeMS > MaxActiveMS {
		return fmt.Errorf("%w: active_ms is outside the allowed range", ErrInvalidViewState)
	}
	if s.maxSection != "" && layout.IndexOf(s.maxSection) < 0 {
		return fmt.Errorf("%w: max_section is not in layout", ErrInvalidViewState)
	}
	if s.exitSection != "" && layout.IndexOf(s.exitSection) < 0 {
		return fmt.Errorf("%w: exit_section is not in layout", ErrInvalidViewState)
	}
	return nil
}

func (s *ViewState) Apply(update ViewStateUpdate, layout LayoutSnapshot) error {
	if s == nil {
		return fmt.Errorf("%w: state is nil", ErrInvalidViewState)
	}
	if update.Seq <= 0 {
		return fmt.Errorf("%w: seq must be positive", ErrInvalidViewState)
	}
	if update.ActiveMS < 0 || update.ActiveMS > MaxActiveMS {
		return fmt.Errorf("%w: active_ms is outside the allowed range", ErrInvalidViewState)
	}
	if update.Seq < s.lastSeq {
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
	if update.Seq == s.lastSeq {
		if update.Final && (s.endedAt.IsZero() || update.SeenAt.After(s.endedAt)) {
			s.endedAt = update.SeenAt.UTC()
		}
		return nil
	}

	s.lastSeq = update.Seq
	if update.ActiveMS > s.activeMS {
		s.activeMS = update.ActiveMS
	}
	if s.confirmedAt.IsZero() && s.activeMS >= ViewConfirmationActiveMS {
		s.confirmedAt = update.SeenAt.UTC()
	}
	if s.lastSeenAt.IsZero() || update.SeenAt.After(s.lastSeenAt) {
		s.lastSeenAt = update.SeenAt.UTC()
	}
	if update.MaxSection != "" && (s.maxSection == "" ||
		layout.IndexOf(update.MaxSection) > layout.IndexOf(s.maxSection)) {
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

func (s ViewState) Snapshot() ViewStateSnapshot {
	return ViewStateSnapshot{
		LastSeq: s.lastSeq, ConfirmedAt: s.confirmedAt, LastSeenAt: s.lastSeenAt,
		ActiveMS: s.activeMS, MaxSection: s.maxSection, ExitSection: s.exitSection,
		EndedAt: s.endedAt,
	}
}

func (s ViewState) LastSeq() int64                  { return s.lastSeq }
func (s ViewState) ConfirmedAt() time.Time          { return s.confirmedAt }
func (s ViewState) LastSeenAt() time.Time           { return s.lastSeenAt }
func (s ViewState) ActiveMS() int64                 { return s.activeMS }
func (s ViewState) MaxSection() AnalyticsSectionID  { return s.maxSection }
func (s ViewState) ExitSection() AnalyticsSectionID { return s.exitSection }
func (s ViewState) EndedAt() time.Time              { return s.endedAt }

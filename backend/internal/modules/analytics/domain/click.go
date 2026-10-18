package domain

import (
	"fmt"
	"strings"
	"time"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type AnalyticsClickParams struct {
	ViewID     sharedvalueobjects.ID
	LinkID     sharedvalueobjects.ID
	Section    AnalyticsSectionID
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
	section    AnalyticsSectionID
	xPct       int
	yPct       int
	viewportW  int
	layoutHash string
	target     string
	occurredAt time.Time
}

func NewAnalyticsClick(params AnalyticsClickParams) (*AnalyticsClick, error) {
	if params.ViewID.IsZero() || params.LinkID.IsZero() {
		return nil, fmt.Errorf("%w: invalid ids", ErrInvalidAnalyticsClick)
	}
	section, err := ParseAnalyticsSectionID(string(params.Section))
	if err != nil || section == AnalyticsSectionTechStack {
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
		id:         sharedvalueobjects.NewID(),
		viewID:     params.ViewID,
		linkID:     params.LinkID,
		section:    section,
		xPct:       params.XPct,
		yPct:       params.YPct,
		viewportW:  params.ViewportW,
		layoutHash: params.LayoutHash,
		target:     params.Target,
		occurredAt: params.OccurredAt.UTC(),
	}, nil
}

func (c *AnalyticsClick) ID() sharedvalueobjects.ID     { return c.id }
func (c *AnalyticsClick) ViewID() sharedvalueobjects.ID { return c.viewID }
func (c *AnalyticsClick) LinkID() sharedvalueobjects.ID { return c.linkID }
func (c *AnalyticsClick) Section() AnalyticsSectionID   { return c.section }
func (c *AnalyticsClick) XPct() int                     { return c.xPct }
func (c *AnalyticsClick) YPct() int                     { return c.yPct }
func (c *AnalyticsClick) ViewportW() int                { return c.viewportW }
func (c *AnalyticsClick) LayoutHash() string            { return c.layoutHash }
func (c *AnalyticsClick) Target() string                { return c.target }
func (c *AnalyticsClick) OccurredAt() time.Time         { return c.occurredAt }

package queries

import (
	"context"
	"fmt"
	"time"

	"developers-cv/backend/internal/modules/resume_link/application/ports"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type LinkVisitsQuery struct {
	UserID          sharedvalueobjects.ID
	LinkID          sharedvalueobjects.ID
	Cursor          string
	Limit           int
	IncludeFiltered bool
}

type VisitsReadQuery struct {
	LinkID          sharedvalueobjects.ID
	CursorTime      time.Time
	CursorID        string
	Limit           int
	IncludeFiltered bool
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

type VisitsReader interface {
	GetVisits(context.Context, VisitsReadQuery) (LinkVisitsPage, error)
}

type VisitsHandler struct {
	links   ResumeLinkPublicReader
	resumes ports.ResumeOwnershipReader
	reader  VisitsReader
}

func NewVisitsHandler(
	links ResumeLinkPublicReader,
	resumes ports.ResumeOwnershipReader,
	reader VisitsReader,
) *VisitsHandler {
	return &VisitsHandler{links: links, resumes: resumes, reader: reader}
}

func (h *VisitsHandler) Handle(ctx context.Context, q LinkVisitsQuery) (LinkVisitsPage, error) {
	if h == nil || h.reader == nil {
		return LinkVisitsPage{}, fmt.Errorf("visits handler dependencies are nil")
	}
	limit, err := NormalizeVisitLimit(q.Limit)
	if err != nil {
		return LinkVisitsPage{}, err
	}
	cursorTime, cursorID, err := DecodeVisitCursor(q.Cursor)
	if err != nil {
		return LinkVisitsPage{}, err
	}
	if _, _, err := authorizeAnalyticsLink(ctx, h.links, h.resumes, q.UserID, q.LinkID); err != nil {
		return LinkVisitsPage{}, err
	}
	return h.reader.GetVisits(ctx, VisitsReadQuery{
		LinkID: q.LinkID, CursorTime: cursorTime, CursorID: cursorID,
		Limit: limit, IncludeFiltered: q.IncludeFiltered,
	})
}

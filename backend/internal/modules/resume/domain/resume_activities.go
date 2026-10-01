package domain

import (
	"fmt"
	"strings"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ActivityType string

const (
	ActivityTypeOpenSource   ActivityType = "openSource"
	ActivityTypeHackathon    ActivityType = "hackathon"
	ActivityTypeSpeaker      ActivityType = "speaker"
	ActivityTypeVolunteering ActivityType = "volunteering"
	ActivityTypeMentoring    ActivityType = "mentoring"
	ActivityTypeCommunity    ActivityType = "community"
)

func (v ActivityType) IsValid() bool {
	switch v {
	case ActivityTypeOpenSource, ActivityTypeHackathon, ActivityTypeSpeaker,
		ActivityTypeVolunteering, ActivityTypeMentoring, ActivityTypeCommunity:
		return true
	default:
		return false
	}
}

type ActivityParams struct {
	ID          sharedvalueobjects.ID
	ResumeID    sharedvalueobjects.ID
	Type        ActivityType
	Title       string
	Description *string
	Date        *YearMonth
	Link        *ExternalURL
}

type Activity struct {
	id          sharedvalueobjects.ID
	resumeID    sharedvalueobjects.ID
	typeValue   ActivityType
	title       string
	description *string
	date        *YearMonth
	link        *ExternalURL
}

func NewActivity(params ActivityParams) (*Activity, error) {
	title := strings.TrimSpace(params.Title)
	if params.ID.IsZero() || params.ResumeID.IsZero() || !params.Type.IsValid() || title == "" {
		return nil, fmt.Errorf("%w: invalid activity fields", ErrInvalidEntity)
	}
	return &Activity{
		id: params.ID, resumeID: params.ResumeID, typeValue: params.Type, title: title,
		description: cloneString(params.Description), date: cloneYearMonth(params.Date), link: cloneExternalURL(params.Link),
	}, nil
}

func (a *Activity) Replace(params ActivityParams) error {
	params.ID, params.ResumeID = a.id, a.resumeID
	updated, err := NewActivity(params)
	if err != nil {
		return err
	}
	*a = *updated
	return nil
}

func (a *Activity) ID() sharedvalueobjects.ID       { return a.id }
func (a *Activity) ResumeID() sharedvalueobjects.ID { return a.resumeID }
func (a *Activity) Type() ActivityType              { return a.typeValue }
func (a *Activity) Title() string                   { return a.title }
func (a *Activity) Description() *string            { return cloneString(a.description) }
func (a *Activity) Date() *YearMonth                { return cloneYearMonth(a.date) }
func (a *Activity) Link() *ExternalURL              { return cloneExternalURL(a.link) }

func (r *Resume) ReplaceActivities(params []ActivityParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	existing := make(map[sharedvalueobjects.ID]*Activity, len(r.activities))
	for _, item := range r.activities {
		if item == nil {
			return fmt.Errorf("%w: nil activity", ErrResumeInvariant)
		}
		existing[item.ID()] = item
	}
	seen := make(map[sharedvalueobjects.ID]struct{}, len(params))
	replaced := make([]*Activity, 0, len(params))
	for _, item := range params {
		if item.ID.IsZero() {
			item.ID = sharedvalueobjects.NewID()
		} else if _, ok := existing[item.ID]; !ok {
			return fmt.Errorf("%w: activity %s does not belong to resume", ErrResumeInvariant, item.ID.String())
		}
		if _, ok := seen[item.ID]; ok {
			return fmt.Errorf("%w: duplicate activity %s", ErrInvalidEntity, item.ID.String())
		}
		seen[item.ID] = struct{}{}
		item.ResumeID = r.id
		activity, err := NewActivity(item)
		if err != nil {
			return err
		}
		replaced = append(replaced, activity)
	}
	r.activities = replaced
	return nil
}

func (r *Resume) Activities() []*Activity {
	if r == nil {
		return nil
	}
	return append([]*Activity(nil), r.activities...)
}

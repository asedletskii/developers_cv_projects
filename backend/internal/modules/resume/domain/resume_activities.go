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

func (v ActivityType) isValid() bool {
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
	Type        ActivityType
	Title       string
	Description *string
	Date        *YearMonth
	Link        *ExternalURL
}

type ActivitySnapshot struct {
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

func newActivity(id, resumeID sharedvalueobjects.ID, params ActivityParams) (*Activity, error) {
	activity := &Activity{
		id: id, resumeID: resumeID, typeValue: params.Type,
		title: strings.TrimSpace(params.Title), description: normalizeOptionalString(params.Description),
		date: cloneYearMonth(params.Date), link: cloneExternalURL(params.Link),
	}
	if err := activity.validate(); err != nil {
		return nil, err
	}
	return activity, nil
}

func RestoreActivity(snapshot ActivitySnapshot) *Activity {
	return &Activity{
		id: snapshot.ID, resumeID: snapshot.ResumeID, typeValue: snapshot.Type, title: snapshot.Title,
		description: copyString(snapshot.Description), date: cloneYearMonth(snapshot.Date), link: cloneExternalURL(snapshot.Link),
	}
}

func (a *Activity) validate() error {
	if a == nil || a.id.IsZero() || a.resumeID.IsZero() || !a.typeValue.isValid() || a.title == "" {
		return fmt.Errorf("%w: invalid activity fields", ErrInvalidEntity)
	}
	if !textLengthWithin(a.title, MaxShortTextLength) || !optionalTextLengthWithin(a.description, MaxDescriptionLength) {
		return fmt.Errorf("%w: activity text exceeds configured limits", ErrInvalidEntity)
	}
	if a.date != nil {
		if err := a.date.validate(); err != nil {
			return fmt.Errorf("%w: invalid activity date", ErrInvalidEntity)
		}
	}
	if a.link != nil {
		if err := a.link.validate(); err != nil {
			return fmt.Errorf("%w: invalid activity URL", ErrInvalidEntity)
		}
	}
	return nil
}

func (a *Activity) ID() sharedvalueobjects.ID       { return a.id }
func (a *Activity) ResumeID() sharedvalueobjects.ID { return a.resumeID }
func (a *Activity) Type() ActivityType              { return a.typeValue }
func (a *Activity) Title() string                   { return a.title }
func (a *Activity) Description() *string            { return copyString(a.description) }
func (a *Activity) Date() *YearMonth                { return cloneYearMonth(a.date) }
func (a *Activity) Link() *ExternalURL              { return cloneExternalURL(a.link) }

func (r *Resume) ReplaceActivities(params []ActivityParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	if len(params) > MaxActivitiesPerResume {
		return fmt.Errorf("%w: resume cannot contain more than %d activities", ErrInvalidEntity, MaxActivitiesPerResume)
	}
	existing := make(map[sharedvalueobjects.ID]struct{}, len(r.activities))
	for _, current := range r.activities {
		if current == nil {
			return fmt.Errorf("%w: nil activity", ErrResumeInvariant)
		}
		existing[current.ID()] = struct{}{}
	}
	seen := make(map[sharedvalueobjects.ID]struct{}, len(params))
	replaced := make([]*Activity, 0, len(params))
	for _, item := range params {
		resolvedID, err := resolveCollectionItemID(item.ID, existing, "activity")
		if err != nil {
			return err
		}
		if _, ok := seen[resolvedID]; ok {
			return fmt.Errorf("%w: duplicate activity %s", ErrInvalidEntity, item.ID.String())
		}
		seen[resolvedID] = struct{}{}
		activity, err := newActivity(resolvedID, r.id, item)
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
	return cloneActivities(r.activities)
}

func cloneActivities(values []*Activity) []*Activity {
	if values == nil {
		return nil
	}
	result := make([]*Activity, len(values))
	for i, value := range values {
		if value == nil {
			continue
		}
		copy := *value
		copy.description = copyString(value.description)
		copy.date = cloneYearMonth(value.date)
		copy.link = cloneExternalURL(value.link)
		result[i] = &copy
	}
	return result
}

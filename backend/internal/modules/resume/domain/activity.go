package domain

import (
	"fmt"
	"strings"

	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

type ActivityParams struct {
	ID          valueobjects.ID
	ResumeID    valueobjects.ID
	Type        enums.ActivityType
	Title       string
	Description *string
	Date        *valueobjects.YearMonth
	Link        *valueobjects.ExternalURL
}

type Activity struct {
	id          valueobjects.ID
	resumeID    valueobjects.ID
	typeValue   enums.ActivityType
	title       string
	description *string
	date        *valueobjects.YearMonth
	link        *valueobjects.ExternalURL
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

func (a *Activity) ID() valueobjects.ID             { return a.id }
func (a *Activity) ResumeID() valueobjects.ID       { return a.resumeID }
func (a *Activity) Type() enums.ActivityType        { return a.typeValue }
func (a *Activity) Title() string                   { return a.title }
func (a *Activity) Description() *string            { return cloneString(a.description) }
func (a *Activity) Date() *valueobjects.YearMonth   { return cloneYearMonth(a.date) }
func (a *Activity) Link() *valueobjects.ExternalURL { return cloneExternalURL(a.link) }

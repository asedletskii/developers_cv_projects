package requests

import (
	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ActivityRequest struct {
	ID          string  `json:"id,omitempty"`
	Type        string  `json:"type"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Date        *string `json:"date,omitempty"`
	Link        *string `json:"link,omitempty"`
}

type ReplaceActivitiesRequest struct {
	Version int64             `json:"version"`
	Items   []ActivityRequest `json:"items"`
}

func (r ReplaceActivitiesRequest) ToInput() (commands.ActivitiesInput, error) {
	out := make([]commands.ActivityInput, 0, len(r.Items))
	for _, item := range r.Items {
		var id sharedvalueobjects.ID
		if item.ID != "" {
			parsed, err := sharedvalueobjects.ParseID(item.ID)
			if err != nil {
				return commands.ActivitiesInput{}, err
			}
			id = parsed
		}
		var date *domain.YearMonth
		if item.Date != nil {
			parsed, err := domain.ParseYearMonth(*item.Date)
			if err != nil {
				return commands.ActivitiesInput{}, err
			}
			date = &parsed
		}
		var link *domain.ExternalURL
		if item.Link != nil {
			parsed, err := domain.NewExternalURL(*item.Link)
			if err != nil {
				return commands.ActivitiesInput{}, err
			}
			link = &parsed
		}
		out = append(out, commands.ActivityInput{
			ID: id, Type: domain.ActivityType(item.Type), Title: item.Title,
			Description: item.Description, Date: date, Link: link,
		})
	}
	return commands.ActivitiesInput{Items: out}, nil
}

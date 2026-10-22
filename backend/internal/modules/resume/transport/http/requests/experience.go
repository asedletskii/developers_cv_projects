package requests

import (
	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ExperienceStackItemRequest struct {
	Name string `json:"name"`
}

type ExperienceRequest struct {
	ID           string                       `json:"id,omitempty"`
	Company      string                       `json:"company"`
	Position     string                       `json:"position"`
	Location     *string                      `json:"location,omitempty"`
	StartDate    string                       `json:"startDate"`
	EndDate      *string                      `json:"endDate,omitempty"`
	Description  *string                      `json:"description,omitempty"`
	Achievements []string                     `json:"achievements"`
	Tasks        []string                     `json:"tasks"`
	Stack        []ExperienceStackItemRequest `json:"stack"`
}

type ReplaceExperiencesRequest struct {
	Version int64               `json:"version"`
	Items   []ExperienceRequest `json:"items"`
}

func (r ReplaceExperiencesRequest) ToInput() (commands.ExperiencesInput, error) {
	items := make([]commands.ExperienceInput, 0, len(r.Items))
	for _, item := range r.Items {
		var id sharedvalueobjects.ID
		if item.ID != "" {
			parsed, err := sharedvalueobjects.ParseID(item.ID)
			if err != nil {
				return commands.ExperiencesInput{}, err
			}
			id = parsed
		}
		startDate, err := domain.ParseYearMonth(item.StartDate)
		if err != nil {
			return commands.ExperiencesInput{}, err
		}
		var endDate *domain.YearMonth
		if item.EndDate != nil {
			parsed, err := domain.ParseYearMonth(*item.EndDate)
			if err != nil {
				return commands.ExperiencesInput{}, err
			}
			endDate = &parsed
		}
		stack := make([]commands.ExperienceStackItemInput, len(item.Stack))
		for i, value := range item.Stack {
			stack[i] = commands.ExperienceStackItemInput{Name: value.Name}
		}
		items = append(items, commands.ExperienceInput{
			ID: id, Company: item.Company, Position: item.Position, Location: item.Location,
			StartDate: startDate, EndDate: endDate, Description: item.Description,
			Achievements: append([]string(nil), item.Achievements...),
			Tasks:        append([]string(nil), item.Tasks...), Stack: stack,
		})
	}
	return commands.ExperiencesInput{Items: items}, nil
}

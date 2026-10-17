package requests

import (
	"developers-cv/backend/internal/modules/resume/application/commands"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type TechnologyRequest struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Position int    `json:"position"`
}

type ResumeTechnologiesRequest struct {
	Version int64               `json:"version"`
	Items   []TechnologyRequest `json:"items"`
}

func (r ResumeTechnologiesRequest) ToInput() (commands.TechnologiesInput, error) {
	items := make([]commands.TechnologyInput, 0, len(r.Items))
	for _, item := range r.Items {
		var id sharedvalueobjects.ID
		if item.ID != "" {
			parsed, err := sharedvalueobjects.ParseID(item.ID)
			if err != nil {
				return commands.TechnologiesInput{}, err
			}
			id = parsed
		}
		items = append(items, commands.TechnologyInput{ID: id, Name: item.Name, Position: item.Position})
	}
	return commands.TechnologiesInput{Entries: items}, nil
}

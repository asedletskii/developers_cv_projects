package requests

import (
	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type SpokenLanguageRequest struct {
	ID       string  `json:"id,omitempty"`
	Language string  `json:"language"`
	Grade    *string `json:"grade,omitempty"`
	Position int     `json:"position"`
}

type ResumeSpokenLanguagesRequest struct {
	Version int64                   `json:"version"`
	Items   []SpokenLanguageRequest `json:"items"`
}

func (r ResumeSpokenLanguagesRequest) ToInput() (commands.LanguagesInput, error) {
	items := make([]commands.LanguageInput, 0, len(r.Items))
	for _, item := range r.Items {
		var id sharedvalueobjects.ID
		if item.ID != "" {
			parsed, err := sharedvalueobjects.ParseID(item.ID)
			if err != nil {
				return commands.LanguagesInput{}, err
			}
			id = parsed
		}
		var grade *domain.LanguageGrade
		if item.Grade != nil {
			value := domain.LanguageGrade(*item.Grade)
			grade = &value
		}
		items = append(items, commands.LanguageInput{
			ID: id, Language: domain.SpokenLanguage(item.Language), Grade: grade, Position: item.Position,
		})
	}
	return commands.LanguagesInput{Entries: items}, nil
}

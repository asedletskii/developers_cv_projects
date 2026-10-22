package requests

import (
	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type EducationRequest struct {
	ID          string  `json:"id,omitempty"`
	Location    string  `json:"location"`
	University  string  `json:"university"`
	Grade       int     `json:"grade"`
	ProgramName string  `json:"programName"`
	ProfileName *string `json:"profileName,omitempty"`
	Faculty     string  `json:"faculty"`
	Course      *int    `json:"course,omitempty"`
	EndDate     *string `json:"endDate,omitempty"`
}

type ReplaceEducationsRequest struct {
	Version int64              `json:"version"`
	Items   []EducationRequest `json:"items"`
}

func (r ReplaceEducationsRequest) ToInput() (commands.EducationsInput, error) {
	out := make([]commands.EducationInput, 0, len(r.Items))
	for _, item := range r.Items {
		var id sharedvalueobjects.ID
		if item.ID != "" {
			parsed, err := sharedvalueobjects.ParseID(item.ID)
			if err != nil {
				return commands.EducationsInput{}, err
			}
			id = parsed
		}
		var endDate *domain.Year
		if item.EndDate != nil {
			parsed, err := domain.ParseYear(*item.EndDate)
			if err != nil {
				return commands.EducationsInput{}, err
			}
			endDate = &parsed
		}
		out = append(out, commands.EducationInput{
			ID: id, Location: item.Location, University: item.University,
			Grade: domain.EducationLevel(item.Grade), ProgramName: item.ProgramName,
			ProfileName: item.ProfileName, Faculty: item.Faculty, Course: item.Course,
			EndDate: endDate,
		})
	}
	return commands.EducationsInput{Items: out}, nil
}

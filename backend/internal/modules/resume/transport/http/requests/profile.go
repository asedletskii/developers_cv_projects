package requests

import (
	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ProfileRequest struct {
	Version     int64    `json:"version"`
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name"`
	Surname     string   `json:"surname"`
	Patronymic  *string  `json:"patronymic,omitempty"`
	Position    string   `json:"position"`
	PhotoURL    string   `json:"photoPath"`
	DateOfBirth string   `json:"dateOfBirth"`
	Location    string   `json:"location"`
	Relocation  bool     `json:"relocation"`
	WorkFormat  []string `json:"workFormat"`
}

func (r ProfileRequest) ToInput() (commands.ProfileInput, error) {
	var id sharedvalueobjects.ID
	if r.ID != "" {
		parsed, err := sharedvalueobjects.ParseID(r.ID)
		if err != nil {
			return commands.ProfileInput{}, err
		}
		id = parsed
	}
	photo, err := domain.NewExternalURL(r.PhotoURL)
	if err != nil {
		return commands.ProfileInput{}, err
	}
	dob, err := domain.ParseDate(r.DateOfBirth)
	if err != nil {
		return commands.ProfileInput{}, err
	}
	formats := make([]domain.WorkFormat, len(r.WorkFormat))
	for i, value := range r.WorkFormat {
		formats[i] = domain.WorkFormat(value)
	}
	return commands.ProfileInput{
		ID: id, Name: r.Name, Surname: r.Surname, Patronymic: r.Patronymic,
		Position: r.Position, PhotoURL: photo, DateOfBirth: dob, Location: r.Location,
		Relocation: r.Relocation, WorkFormats: formats,
	}, nil
}

package request

import (
	"github.com/google/uuid"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

func (r ResumeRequest) ToUpdateCommand(userID uuid.UUID, resumeID valueobjects.ID) (commands.UpdateResumeCommand, error) {
	createCommand, err := r.ToCreateCommand(userID)
	if err != nil {
		return commands.UpdateResumeCommand{}, err
	}
	return commands.UpdateResumeCommand{
		ResumeID:       resumeID,
		UserID:         createCommand.UserID,
		Name:           createCommand.Name,
		Surname:        createCommand.Surname,
		Patronymic:     createCommand.Patronymic,
		Position:       createCommand.Position,
		Availability:   createCommand.Availability,
		ResumeLink:     createCommand.ResumeLink,
		Theme:          createCommand.Theme,
		Sections:       createCommand.Sections,
		HiddenContacts: createCommand.HiddenContacts,
		ContactsText:   createCommand.ContactsText,
		Greeting:       createCommand.Greeting,
		PhotoPath:      createCommand.PhotoPath,
		DateOfBirth:    createCommand.DateOfBirth,
		Location:       createCommand.Location,
		Relocation:     createCommand.Relocation,
		WorkFormats:    createCommand.WorkFormats,
		Contacts:       createCommand.Contacts,
		Languages:      createCommand.Languages,
		Experience:     createCommand.Experience,
		Education:      createCommand.Education,
		Courses:        createCommand.Courses,
		Activities:     createCommand.Activities,
		HardSkills:     createCommand.HardSkills,
		SoftSkills:     createCommand.SoftSkills,
		Technologies:   createCommand.Technologies,
		Projects:       createCommand.Projects,
	}, nil
}

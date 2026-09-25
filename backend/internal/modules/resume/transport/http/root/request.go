package root

import (
	"github.com/google/uuid"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type Request struct {
	Title        string              `json:"title"`
	Presentation PresentationRequest `json:"presentation"`
}
type PresentationRequest struct {
	Theme        *ThemeRequest        `json:"theme,omitempty"`
	Greeting     GreetingRequest      `json:"greeting"`
	Availability *AvailabilityRequest `json:"availability,omitempty"`
	ResumeLink   *string              `json:"resumeLink,omitempty"`
	Sections     []SectionRequest     `json:"sections"`
}
type ThemeRequest struct {
	Primary   *string `json:"primary,omitempty"`
	Secondary *string `json:"secondary,omitempty"`
}
type LocalizedRequest struct {
	RU string `json:"ru"`
	EN string `json:"en"`
}
type GreetingRequest struct {
	First           LocalizedRequest `json:"first"`
	Second          LocalizedRequest `json:"second"`
	SecondHighlight LocalizedRequest `json:"secondHighlight"`
	Third           LocalizedRequest `json:"third"`
}
type AvailabilityRequest struct {
	Status string            `json:"status"`
	From   *string           `json:"from,omitempty"`
	Note   *LocalizedRequest `json:"note,omitempty"`
}
type SectionRequest struct {
	ID        string `json:"id"`
	Enabled   bool   `json:"enabled"`
	Collapsed bool   `json:"collapsed"`
}

func (r Request) presentation() commands.PresentationInput {
	var theme *commands.ThemeInput
	if r.Presentation.Theme != nil {
		theme = &commands.ThemeInput{Primary: r.Presentation.Theme.Primary, Secondary: r.Presentation.Theme.Secondary}
	}
	var av *commands.AvailabilityInput
	if r.Presentation.Availability != nil {
		var note *commands.LocalizedInput
		if r.Presentation.Availability.Note != nil {
			note = &commands.LocalizedInput{RU: r.Presentation.Availability.Note.RU, EN: r.Presentation.Availability.Note.EN}
		}
		av = &commands.AvailabilityInput{Status: r.Presentation.Availability.Status, From: r.Presentation.Availability.From, Note: note}
	}
	sections := make([]commands.SectionInput, len(r.Presentation.Sections))
	for i, v := range r.Presentation.Sections {
		sections[i] = commands.SectionInput{ID: v.ID, Enabled: v.Enabled, Collapsed: v.Collapsed}
	}
	return commands.PresentationInput{Theme: theme, Greeting: commands.GreetingInput{First: commands.LocalizedInput{RU: r.Presentation.Greeting.First.RU, EN: r.Presentation.Greeting.First.EN}, Second: commands.LocalizedInput{RU: r.Presentation.Greeting.Second.RU, EN: r.Presentation.Greeting.Second.EN}, SecondHighlight: commands.LocalizedInput{RU: r.Presentation.Greeting.SecondHighlight.RU, EN: r.Presentation.Greeting.SecondHighlight.EN}, Third: commands.LocalizedInput{RU: r.Presentation.Greeting.Third.RU, EN: r.Presentation.Greeting.Third.EN}}, Availability: av, ResumeLink: r.Presentation.ResumeLink, Sections: sections}
}

func (r Request) ToCreateCommand(userID uuid.UUID) (commands.CreateResumeRootCommand, error) {
	uid, err := useridentity.IDFromUUID(userID)
	if err != nil {
		return commands.CreateResumeRootCommand{}, err
	}
	return commands.CreateResumeRootCommand{UserID: uid, Title: r.Title, Presentation: r.presentation()}, nil
}
func (r Request) ToUpdateCommand(userID uuid.UUID, resumeID valueobjects.ID) (commands.UpdateResumeRootCommand, error) {
	uid, err := useridentity.IDFromUUID(userID)
	if err != nil {
		return commands.UpdateResumeRootCommand{}, err
	}
	return commands.UpdateResumeRootCommand{UserID: uid, ResumeID: resumeID, Title: r.Title, Presentation: r.presentation()}, nil
}

package root

import (
	"github.com/google/uuid"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type ResumeRootRequest struct {
	Title        string              `json:"title"`
	Language     string              `json:"language"`
	Presentation PresentationRequest `json:"presentation"`
}
type ResumeRootUpdateRequest struct {
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
type GreetingRequest struct {
	First           string `json:"first"`
	Second          string `json:"second"`
	SecondHighlight string `json:"secondHighlight"`
	Third           string `json:"third"`
}
type AvailabilityRequest struct {
	Status string  `json:"status"`
	From   *string `json:"from,omitempty"`
	Note   *string `json:"note,omitempty"`
}
type SectionRequest struct {
	ID        string `json:"id"`
	Enabled   bool   `json:"enabled"`
	Collapsed bool   `json:"collapsed"`
}

func presentationInput(p PresentationRequest) commands.PresentationInput {
	var theme *commands.ThemeInput
	if p.Theme != nil {
		theme = &commands.ThemeInput{Primary: p.Theme.Primary, Secondary: p.Theme.Secondary}
	}
	var av *commands.AvailabilityInput
	if p.Availability != nil {
		av = &commands.AvailabilityInput{Status: p.Availability.Status, From: p.Availability.From, Note: p.Availability.Note}
	}
	sections := make([]commands.SectionInput, len(p.Sections))
	for i, v := range p.Sections {
		sections[i] = commands.SectionInput{ID: v.ID, Enabled: v.Enabled, Collapsed: v.Collapsed}
	}
	return commands.PresentationInput{Theme: theme, Greeting: commands.GreetingInput{First: p.Greeting.First, Second: p.Greeting.Second, SecondHighlight: p.Greeting.SecondHighlight, Third: p.Greeting.Third}, Availability: av, ResumeLink: p.ResumeLink, Sections: sections}
}

func (r ResumeRootRequest) ToCreateCommand(userID uuid.UUID) (commands.CreateResumeRootCommand, error) {
	uid, err := useridentity.IDFromUUID(userID)
	if err != nil {
		return commands.CreateResumeRootCommand{}, err
	}
	return commands.CreateResumeRootCommand{UserID: uid, Title: r.Title, Language: enums.ResumeLanguage(r.Language), Presentation: presentationInput(r.Presentation)}, nil
}
func (r ResumeRootUpdateRequest) ToUpdateCommand(userID uuid.UUID, resumeID valueobjects.ID) (commands.UpdateResumeRootCommand, error) {
	uid, err := useridentity.IDFromUUID(userID)
	if err != nil {
		return commands.UpdateResumeRootCommand{}, err
	}
	return commands.UpdateResumeRootCommand{UserID: uid, ResumeID: resumeID, Title: r.Title, Presentation: presentationInput(r.Presentation)}, nil
}

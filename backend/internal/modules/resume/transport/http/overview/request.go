package overview

import (
	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"

	"github.com/google/uuid"
)

type OverviewRequest struct {
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

func (r OverviewRequest) ToCommand(uid uuid.UUID, id valueobjects.ID) (commands.UpdateOverviewCommand, error) {
	u, e := useridentity.IDFromUUID(uid)
	if e != nil {
		return commands.UpdateOverviewCommand{}, e
	}
	p := commands.OverviewInput{Greeting: commands.GreetingInput{First: r.Presentation.Greeting.First, Second: r.Presentation.Greeting.Second, SecondHighlight: r.Presentation.Greeting.SecondHighlight, Third: r.Presentation.Greeting.Third}, ResumeLink: r.Presentation.ResumeLink}
	if r.Presentation.Theme != nil {
		p.Theme = &commands.ThemeInput{Primary: r.Presentation.Theme.Primary, Secondary: r.Presentation.Theme.Secondary}
	}
	if r.Presentation.Availability != nil {
		p.Availability = &commands.AvailabilityInput{Status: r.Presentation.Availability.Status, From: r.Presentation.Availability.From, Note: r.Presentation.Availability.Note}
	}
	p.Sections = make([]commands.SectionInput, len(r.Presentation.Sections))
	for i, s := range r.Presentation.Sections {
		p.Sections[i] = commands.SectionInput{ID: s.ID, Enabled: s.Enabled, Collapsed: s.Collapsed}
	}
	return commands.UpdateOverviewCommand{ResumeID: id, UserID: u, Presentation: p}, nil
}

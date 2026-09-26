package overview

import (
	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type overviewRequest struct {
	Presentation presentationRequest `json:"presentation"`
}
type presentationRequest struct {
	Theme        *themeRequest        `json:"theme,omitempty"`
	Greeting     greetingRequest      `json:"greeting"`
	Availability *availabilityRequest `json:"availability,omitempty"`
	ResumeLink   *string              `json:"resumeLink,omitempty"`
	Sections     []sectionRequest     `json:"sections"`
}
type themeRequest struct {
	Primary   *string `json:"primary,omitempty"`
	Secondary *string `json:"secondary,omitempty"`
}
type greetingRequest struct {
	First           string `json:"first"`
	Second          string `json:"second"`
	SecondHighlight string `json:"secondHighlight"`
	Third           string `json:"third"`
}
type availabilityRequest struct {
	Status string  `json:"status"`
	From   *string `json:"from,omitempty"`
	Note   *string `json:"note,omitempty"`
}
type sectionRequest struct {
	ID        string `json:"id"`
	Enabled   bool   `json:"enabled"`
	Collapsed bool   `json:"collapsed"`
}

func (r overviewRequest) ToCommand(uid useridentity.ID, id valueobjects.ID) (commands.UpdateOverviewCommand, error) {
	p := commands.OverviewInput{
		Greeting: commands.GreetingInput{
			First:           r.Presentation.Greeting.First,
			Second:          r.Presentation.Greeting.Second,
			SecondHighlight: r.Presentation.Greeting.SecondHighlight,
			Third:           r.Presentation.Greeting.Third,
		},
	}
	if r.Presentation.ResumeLink != nil {
		link, err := valueobjects.NewExternalURL(*r.Presentation.ResumeLink)
		if err != nil {
			return commands.UpdateOverviewCommand{}, err
		}
		p.ResumeLink = &link
	}
	if r.Presentation.Theme != nil {
		theme := &commands.ThemeInput{}
		if r.Presentation.Theme.Primary != nil {
			color, err := valueobjects.NewColor(*r.Presentation.Theme.Primary)
			if err != nil {
				return commands.UpdateOverviewCommand{}, err
			}
			theme.Primary = &color
		}
		if r.Presentation.Theme.Secondary != nil {
			color, err := valueobjects.NewColor(*r.Presentation.Theme.Secondary)
			if err != nil {
				return commands.UpdateOverviewCommand{}, err
			}
			theme.Secondary = &color
		}
		p.Theme = theme
	}
	if r.Presentation.Availability != nil {
		availability := &commands.AvailabilityInput{Status: enums.AvailabilityStatus(r.Presentation.Availability.Status), Note: r.Presentation.Availability.Note}
		if r.Presentation.Availability.From != nil {
			date, err := valueobjects.ParseDate(*r.Presentation.Availability.From)
			if err != nil {
				return commands.UpdateOverviewCommand{}, err
			}
			availability.From = &date
		}
		p.Availability = availability
	}
	p.Sections = make([]commands.SectionInput, len(r.Presentation.Sections))
	for i, section := range r.Presentation.Sections {
		p.Sections[i] = commands.SectionInput{ID: enums.SectionID(section.ID), Enabled: section.Enabled, Collapsed: section.Collapsed}
	}
	return commands.UpdateOverviewCommand{ResumeID: id, UserID: uid, Presentation: p}, nil
}

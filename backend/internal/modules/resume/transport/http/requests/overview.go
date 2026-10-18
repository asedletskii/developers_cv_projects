package requests

import (
	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ResumeOverviewRequest struct {
	Version      int64                        `json:"version"`
	Theme        *OverviewThemeRequest        `json:"theme,omitempty"`
	Greeting     OverviewGreetingRequest      `json:"greeting"`
	Availability *OverviewAvailabilityRequest `json:"availability,omitempty"`
	Sections     []OverviewSectionRequest     `json:"sections"`
}

type OverviewThemeRequest struct {
	Primary   *string `json:"primary,omitempty"`
	Secondary *string `json:"secondary,omitempty"`
}
type OverviewGreetingRequest struct {
	First           string `json:"first"`
	Second          string `json:"second"`
	SecondHighlight string `json:"secondHighlight"`
	Third           string `json:"third"`
}
type OverviewAvailabilityRequest struct {
	Status string  `json:"status"`
	From   *string `json:"from,omitempty"`
	Note   *string `json:"note,omitempty"`
}
type OverviewSectionRequest struct {
	ID        string `json:"id"`
	Enabled   bool   `json:"enabled"`
	Collapsed bool   `json:"collapsed"`
}

func (r ResumeOverviewRequest) ToCommand(uid sharedvalueobjects.ID, id sharedvalueobjects.ID) (commands.ReplaceOverviewCommand, error) {
	input := commands.OverviewInput{Greeting: commands.GreetingInput{First: r.Greeting.First, Second: r.Greeting.Second, SecondHighlight: r.Greeting.SecondHighlight, Third: r.Greeting.Third}}
	if r.Theme != nil {
		theme := &commands.ThemeInput{}
		if r.Theme.Primary != nil {
			color, err := domain.NewColor(*r.Theme.Primary)
			if err != nil {
				return commands.ReplaceOverviewCommand{}, err
			}
			theme.Primary = &color
		}
		if r.Theme.Secondary != nil {
			color, err := domain.NewColor(*r.Theme.Secondary)
			if err != nil {
				return commands.ReplaceOverviewCommand{}, err
			}
			theme.Secondary = &color
		}
		input.Theme = theme
	}
	if r.Availability != nil {
		availability := &commands.AvailabilityInput{Status: domain.AvailabilityStatus(r.Availability.Status), Note: r.Availability.Note}
		if r.Availability.From != nil {
			date, err := domain.ParseDate(*r.Availability.From)
			if err != nil {
				return commands.ReplaceOverviewCommand{}, err
			}
			availability.From = &date
		}
		input.Availability = availability
	}
	input.Sections = make([]commands.SectionInput, len(r.Sections))
	for i, section := range r.Sections {
		input.Sections[i] = commands.SectionInput{ID: domain.SectionID(section.ID), Enabled: section.Enabled, Collapsed: section.Collapsed}
	}
	return commands.ReplaceOverviewCommand{ResumeID: id, UserID: uid, Version: r.Version, Input: input}, nil
}

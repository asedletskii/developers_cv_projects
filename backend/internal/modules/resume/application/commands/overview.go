package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type GreetingInput struct {
	First           string
	Second          string
	SecondHighlight string
	Third           string
}
type AvailabilityInput struct {
	Status domain.AvailabilityStatus
	From   *domain.Date
	Note   *string
}
type ThemeInput struct {
	Primary   *domain.Color
	Secondary *domain.Color
}
type SectionInput struct {
	ID        domain.SectionID
	Enabled   bool
	Collapsed bool
	Style     *string
}
type OverviewInput struct {
	Theme        *ThemeInput
	Greeting     GreetingInput
	Availability *AvailabilityInput
	Sections     []SectionInput
}

func (i OverviewInput) toDomain() (domain.ResumeOverviewParams, error) {
	greeting, err := domain.NewGreeting(domain.GreetingParams{
		First:           i.Greeting.First,
		Second:          i.Greeting.Second,
		SecondHighlight: i.Greeting.SecondHighlight,
		Third:           i.Greeting.Third,
	})
	if err != nil {
		return domain.ResumeOverviewParams{}, err
	}
	var theme *domain.ThemeSettings
	if i.Theme != nil {
		value, err := domain.NewThemeSettings(domain.ThemeSettingsParams{
			Primary: i.Theme.Primary, Secondary: i.Theme.Secondary,
		})
		if err != nil {
			return domain.ResumeOverviewParams{}, err
		}
		theme = &value
	}
	var availability *domain.Availability
	if i.Availability != nil {
		value, err := domain.NewAvailability(domain.AvailabilityParams{
			Status: i.Availability.Status, From: i.Availability.From, Note: i.Availability.Note,
		})
		if err != nil {
			return domain.ResumeOverviewParams{}, err
		}
		availability = &value
	}
	sections := make([]domain.SectionSettings, len(i.Sections))
	for index, section := range i.Sections {
		value, err := domain.NewSectionSettings(domain.SectionSettingsParams{
			ID: section.ID, Enabled: section.Enabled, Collapsed: section.Collapsed, Style: section.Style,
		})
		if err != nil {
			return domain.ResumeOverviewParams{}, err
		}
		sections[index] = value
	}
	return domain.ResumeOverviewParams{
		Theme:        theme,
		Greeting:     greeting,
		Availability: availability,
		Sections:     sections,
	}, nil
}

type ReplaceOverviewCommand struct {
	ResumeID sharedvalueobjects.ID
	UserID   sharedvalueobjects.ID
	Version  int64
	Input    OverviewInput
}

type ReplaceOverviewHandler struct {
	repository ports.ResumeRepository
	uow        ports.UnitOfWork
}

func NewReplaceOverviewHandler(repository ports.ResumeRepository, uow ports.UnitOfWork) *ReplaceOverviewHandler {
	return &ReplaceOverviewHandler{repository: repository, uow: uow}
}

func (h *ReplaceOverviewHandler) Handle(ctx context.Context, c ReplaceOverviewCommand) error {
	if h == nil {
		return fmt.Errorf("replace overview handler dependencies are nil")
	}
	input, err := c.Input.toDomain()
	if err != nil {
		return err
	}
	return executeResumeMutation(ctx, h.repository, h.uow, c.ResumeID, c.UserID, c.Version, "replace overview", ports.ResumeMutationOverview, func(resume *domain.Resume) error {
		return resume.ReplaceOverview(input)
	})
}

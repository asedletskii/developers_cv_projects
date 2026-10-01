package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type GreetingInput struct{ First, Second, SecondHighlight, Third string }
type AvailabilityInput struct {
	Status domain.AvailabilityStatus
	From   *domain.Date
	Note   *string
}
type ThemeInput struct{ Primary, Secondary *domain.Color }
type SectionInput struct {
	ID                 domain.SectionID
	Enabled, Collapsed bool
}
type OverviewInput struct {
	Theme        *ThemeInput
	Greeting     GreetingInput
	Availability *AvailabilityInput
	ResumeLink   *domain.ExternalURL
	Sections     []SectionInput
}

func (i OverviewInput) toDomain() (domain.ResumeOverviewParams, error) {
	g := domain.NewGreeting(domain.GreetingParams{First: i.Greeting.First, Second: i.Greeting.Second, SecondHighlight: i.Greeting.SecondHighlight, Third: i.Greeting.Third})
	var theme *domain.ThemeSettings
	if i.Theme != nil {
		v, err := domain.NewThemeSettings(domain.ThemeSettingsParams{Primary: cloneColor(i.Theme.Primary), Secondary: cloneColor(i.Theme.Secondary)})
		if err != nil {
			return domain.ResumeOverviewParams{}, err
		}
		theme = &v
	}
	var av *domain.Availability
	if i.Availability != nil {
		v, err := domain.NewAvailability(domain.AvailabilityParams{Status: i.Availability.Status, From: i.Availability.From, Note: i.Availability.Note})
		if err != nil {
			return domain.ResumeOverviewParams{}, err
		}
		av = &v
	}
	sections := make([]domain.SectionSettings, len(i.Sections))
	for n, v := range i.Sections {
		x, err := domain.NewSectionSettings(domain.SectionSettingsParams{ID: v.ID, Enabled: v.Enabled, Collapsed: v.Collapsed})
		if err != nil {
			return domain.ResumeOverviewParams{}, err
		}
		sections[n] = x
	}
	return domain.ResumeOverviewParams{Theme: theme, Greeting: g, Availability: av, ResumeLink: cloneExternalURL(i.ResumeLink), Sections: sections}, nil
}

type ReplaceOverviewCommand struct {
	ResumeID identity.ID
	UserID   identity.ID
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
	return executeResumeMutation(ctx, h.repository, h.uow, c.ResumeID, c.UserID, c.Version, "replace overview", func(resume *domain.Resume) error {
		return resume.ReplaceOverview(input)
	})
}

func cloneColor(value *domain.Color) *domain.Color {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func cloneExternalURL(value *domain.ExternalURL) *domain.ExternalURL {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

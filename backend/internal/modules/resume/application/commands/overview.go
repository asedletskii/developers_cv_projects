package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type GreetingInput struct{ First, Second, SecondHighlight, Third string }
type AvailabilityInput struct {
	Status enums.AvailabilityStatus
	From   *valueobjects.Date
	Note   *string
}
type ThemeInput struct{ Primary, Secondary *valueobjects.Color }
type SectionInput struct {
	ID                 enums.SectionID
	Enabled, Collapsed bool
}
type OverviewInput struct {
	Theme        *ThemeInput
	Greeting     GreetingInput
	Availability *AvailabilityInput
	ResumeLink   *valueobjects.ExternalURL
	Sections     []SectionInput
}

func (i OverviewInput) toDomain() (domain.ResumeOverviewParams, error) {
	g := valueobjects.NewGreeting(valueobjects.GreetingParams{First: i.Greeting.First, Second: i.Greeting.Second, SecondHighlight: i.Greeting.SecondHighlight, Third: i.Greeting.Third})
	var theme *valueobjects.ThemeSettings
	if i.Theme != nil {
		v, err := valueobjects.NewThemeSettings(valueobjects.ThemeSettingsParams{Primary: cloneColor(i.Theme.Primary), Secondary: cloneColor(i.Theme.Secondary)})
		if err != nil {
			return domain.ResumeOverviewParams{}, err
		}
		theme = &v
	}
	var av *valueobjects.Availability
	if i.Availability != nil {
		v, err := valueobjects.NewAvailability(valueobjects.AvailabilityParams{Status: i.Availability.Status, From: i.Availability.From, Note: i.Availability.Note})
		if err != nil {
			return domain.ResumeOverviewParams{}, err
		}
		av = &v
	}
	sections := make([]valueobjects.SectionSettings, len(i.Sections))
	for n, v := range i.Sections {
		x, err := valueobjects.NewSectionSettings(valueobjects.SectionSettingsParams{ID: v.ID, Enabled: v.Enabled, Collapsed: v.Collapsed})
		if err != nil {
			return domain.ResumeOverviewParams{}, err
		}
		sections[n] = x
	}
	return domain.ResumeOverviewParams{Theme: theme, Greeting: g, Availability: av, ResumeLink: cloneExternalURL(i.ResumeLink), Sections: sections}, nil
}

type ReplaceOverviewCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
	Version  int64
	Input    OverviewInput
}

type ReplaceOverviewHandler struct {
	repository ports.ResumeRepository
	uow        UnitOfWork
}

func NewReplaceOverviewHandler(repository ports.ResumeRepository, uow UnitOfWork) *ReplaceOverviewHandler {
	return &ReplaceOverviewHandler{repository: repository, uow: uow}
}

func (h *ReplaceOverviewHandler) Handle(ctx context.Context, c ReplaceOverviewCommand) error {
	if h == nil || h.repository == nil || h.uow == nil {
		return fmt.Errorf("replace overview handler dependencies are nil")
	}
	if c.ResumeID.IsZero() || c.UserID.IsZero() || c.Version < 1 {
		return fmt.Errorf("invalid replace overview command")
	}
	input, err := c.Input.toDomain()
	if err != nil {
		return err
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		resume, err := h.repository.Get(tx, ports.GetResumeParams{ResumeID: c.ResumeID, UserID: c.UserID})
		if err != nil {
			return err
		}
		if resume.Version() != c.Version {
			return domain.ErrConcurrencyConflict
		}
		if err := resume.ReplaceOverview(input); err != nil {
			return err
		}
		return h.repository.Save(tx, resume)
	})
}

func cloneColor(value *valueobjects.Color) *valueobjects.Color {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func cloneExternalURL(value *valueobjects.ExternalURL) *valueobjects.ExternalURL {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

package commands

import (
	"context"
	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
	"fmt"
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
		var a, b *valueobjects.Color
		if i.Theme.Primary != nil {
			a = cloneColor(i.Theme.Primary)
		}
		if i.Theme.Secondary != nil {
			b = cloneColor(i.Theme.Secondary)
		}
		v, e := valueobjects.NewThemeSettings(valueobjects.ThemeSettingsParams{Primary: a, Secondary: b})
		if e != nil {
			return domain.ResumeOverviewParams{}, e
		}
		theme = &v
	}
	var av *valueobjects.Availability
	if i.Availability != nil {
		v, e := valueobjects.NewAvailability(valueobjects.AvailabilityParams{Status: i.Availability.Status, From: i.Availability.From, Note: i.Availability.Note})
		if e != nil {
			return domain.ResumeOverviewParams{}, e
		}
		av = &v
	}
	link := cloneExternalURL(i.ResumeLink)
	sections := make([]valueobjects.SectionSettings, len(i.Sections))
	for n, v := range i.Sections {
		x, e := valueobjects.NewSectionSettings(valueobjects.SectionSettingsParams{ID: v.ID, Enabled: v.Enabled, Collapsed: v.Collapsed})
		if e != nil {
			return domain.ResumeOverviewParams{}, e
		}
		sections[n] = x
	}
	return domain.ResumeOverviewParams{Theme: theme, Greeting: g, Availability: av, ResumeLink: link, Sections: sections}, nil
}

type UpdateOverviewCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
	Input    OverviewInput
}

type UpdateOverviewHandler struct {
	updater ports.UpdateOverviewRepository
	uow     UnitOfWork
}

func NewUpdateOverviewHandler(u ports.UpdateOverviewRepository, tx UnitOfWork) *UpdateOverviewHandler {
	return &UpdateOverviewHandler{updater: u, uow: tx}
}
func (h *UpdateOverviewHandler) Handle(ctx context.Context, c UpdateOverviewCommand) error {
	if h == nil || h.updater == nil || h.uow == nil {
		return fmt.Errorf("update overview handler dependencies are nil")
	}
	p, e := c.Input.toDomain()
	if e != nil {
		return e
	}
	o, e := domain.NewResumeOverview(p)
	if e != nil {
		return e
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		return h.updater.UpdateOverview(tx, ports.UpdateOverviewRepositoryCommand{ResumeID: c.ResumeID, UserID: c.UserID, Overview: o})
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

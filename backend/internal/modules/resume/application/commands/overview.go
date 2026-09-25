package commands

import (
	"context"
	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
	"fmt"
)

type GreetingInput struct{ First, Second, SecondHighlight, Third string }
type AvailabilityInput struct {
	Status string
	From   *string
	Note   *string
}
type ThemeInput struct{ Primary, Secondary *string }
type SectionInput struct {
	ID                 string
	Enabled, Collapsed bool
}
type OverviewInput struct {
	Theme        *ThemeInput
	Greeting     GreetingInput
	Availability *AvailabilityInput
	ResumeLink   *string
	Sections     []SectionInput
}

func (i OverviewInput) toDomain() (domain.OverviewPresentation, error) {
	g := valueobjects.NewGreeting(valueobjects.GreetingParams{First: i.Greeting.First, Second: i.Greeting.Second, SecondHighlight: i.Greeting.SecondHighlight, Third: i.Greeting.Third})
	var theme *valueobjects.ThemeSettings
	if i.Theme != nil {
		var a, b *valueobjects.Color
		if i.Theme.Primary != nil {
			v, e := valueobjects.NewColor(*i.Theme.Primary)
			if e != nil {
				return domain.OverviewPresentation{}, e
			}
			a = &v
		}
		if i.Theme.Secondary != nil {
			v, e := valueobjects.NewColor(*i.Theme.Secondary)
			if e != nil {
				return domain.OverviewPresentation{}, e
			}
			b = &v
		}
		v, e := valueobjects.NewThemeSettings(valueobjects.ThemeSettingsParams{Primary: a, Secondary: b})
		if e != nil {
			return domain.OverviewPresentation{}, e
		}
		theme = &v
	}
	var av *valueobjects.Availability
	if i.Availability != nil {
		var from *valueobjects.Date
		if i.Availability.From != nil {
			v, e := valueobjects.ParseDate(*i.Availability.From)
			if e != nil {
				return domain.OverviewPresentation{}, e
			}
			from = &v
		}
		v, e := valueobjects.NewAvailability(valueobjects.AvailabilityParams{Status: enums.AvailabilityStatus(i.Availability.Status), From: from, Note: i.Availability.Note})
		if e != nil {
			return domain.OverviewPresentation{}, e
		}
		av = &v
	}
	var link *valueobjects.AssetPath
	if i.ResumeLink != nil {
		v, e := valueobjects.NewAssetPath(*i.ResumeLink)
		if e != nil {
			return domain.OverviewPresentation{}, e
		}
		link = &v
	}
	sections := make([]valueobjects.SectionSettings, len(i.Sections))
	for n, v := range i.Sections {
		x, e := valueobjects.NewSectionSettings(valueobjects.SectionSettingsParams{ID: enums.SectionID(v.ID), Enabled: v.Enabled, Collapsed: v.Collapsed})
		if e != nil {
			return domain.OverviewPresentation{}, e
		}
		sections[n] = x
	}
	return domain.OverviewPresentation{Theme: theme, Greeting: g, Availability: av, ResumeLink: link, Sections: sections}, nil
}

type UpdateOverviewCommand struct {
	ResumeID     valueobjects.ID
	UserID       useridentity.ID
	Presentation OverviewInput
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
	p, e := c.Presentation.toDomain()
	if e != nil {
		return e
	}
	o, e := domain.NewResumeOverview(domain.ResumeOverviewParams{Presentation: p})
	if e != nil {
		return e
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		return h.updater.UpdateOverview(tx, ports.UpdateOverviewRepositoryCommand{ResumeID: c.ResumeID, UserID: c.UserID, Overview: o})
	})
}

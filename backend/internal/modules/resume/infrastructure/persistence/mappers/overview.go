package mappers

import (
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func OverviewToDBModels(
	id sharedvalueobjects.ID,
	o *domain.ResumeOverview,
) (gormmodel.ResumeOverview, []gormmodel.ResumeOverviewSection, error) {
	if o == nil {
		return gormmodel.ResumeOverview{}, nil, fmt.Errorf("overview is nil")
	}
	theme := o.Theme()
	greeting := o.Greeting()
	availability := o.Availability()
	resumeLink := o.ResumeLink()
	sections := o.Sections()
	m := gormmodel.ResumeOverview{
		ResumeID:          id.UUID(),
		GreetingFirst:     greeting.First(),
		GreetingSecond:    greeting.Second(),
		GreetingHighlight: greeting.SecondHighlight(),
		GreetingThird:     greeting.Third(),
	}
	if availability != nil {
		v := string(availability.Status())
		m.AvailabilityStatus = &v
		if d := availability.From(); d != nil {
			x := d.Time()
			m.AvailabilityFrom = &x
		}
		if n := availability.Note(); n != nil {
			v := *n
			m.AvailabilityNote = &v
		}
	}
	if resumeLink != nil {
		v := resumeLink.String()
		m.ResumeLink = &v
	}
	if theme != nil {
		if theme.Primary != nil {
			v := theme.Primary.String()
			m.ThemePrimary = &v
		}
		if theme.Secondary != nil {
			v := theme.Secondary.String()
			m.ThemeSecondary = &v
		}
	}
	mapped := make([]gormmodel.ResumeOverviewSection, 0, len(sections))
	for i, x := range sections {
		mapped = append(mapped, gormmodel.ResumeOverviewSection{
			ResumeID:  id.UUID(),
			SectionID: string(x.ID()),
			Position:  i,
			Enabled:   x.Enabled(),
			Collapsed: x.Collapsed(),
		})
	}
	return m, mapped, nil
}

func OverviewToDomain(
	m gormmodel.ResumeOverview,
	s []gormmodel.ResumeOverviewSection,
) (*domain.ResumeOverview, error) {
	g := domain.NewGreeting(domain.GreetingParams{
		First:           m.GreetingFirst,
		Second:          m.GreetingSecond,
		SecondHighlight: m.GreetingHighlight,
		Third:           m.GreetingThird,
	})
	var av *domain.Availability
	if m.AvailabilityStatus != nil {
		st := domain.AvailabilityStatus(*m.AvailabilityStatus)
		params := domain.AvailabilityParams{Status: st, Note: m.AvailabilityNote}
		if m.AvailabilityFrom != nil {
			d, err := domain.NewDate(*m.AvailabilityFrom)
			if err != nil {
				return nil, err
			}
			params.From = &d
		}
		v, err := domain.NewAvailability(params)
		if err != nil {
			return nil, err
		}
		av = &v
	}
	var link *domain.ExternalURL
	if m.ResumeLink != nil {
		v, err := domain.NewExternalURL(*m.ResumeLink)
		if err != nil {
			return nil, err
		}
		link = &v
	}
	var theme *domain.ThemeSettings
	if m.ThemePrimary != nil || m.ThemeSecondary != nil {
		var a, b *domain.Color
		if m.ThemePrimary != nil {
			v, err := domain.NewColor(*m.ThemePrimary)
			if err != nil {
				return nil, err
			}
			a = &v
		}
		if m.ThemeSecondary != nil {
			v, err := domain.NewColor(*m.ThemeSecondary)
			if err != nil {
				return nil, err
			}
			b = &v
		}
		v, err := domain.NewThemeSettings(domain.ThemeSettingsParams{Primary: a, Secondary: b})
		if err != nil {
			return nil, err
		}
		theme = &v
	}
	ss := make([]domain.SectionSettings, 0, len(s))
	for _, x := range s {
		v, err := domain.NewSectionSettings(domain.SectionSettingsParams{
			ID:        domain.SectionID(x.SectionID),
			Enabled:   x.Enabled,
			Collapsed: x.Collapsed,
		})
		if err != nil {
			return nil, err
		}
		ss = append(ss, v)
	}
	return domain.NewResumeOverview(domain.ResumeOverviewParams{
		Theme:        theme,
		Greeting:     g,
		Availability: av,
		ResumeLink:   link,
		Sections:     ss,
	})
}

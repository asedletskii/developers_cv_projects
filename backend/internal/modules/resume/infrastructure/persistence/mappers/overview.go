package mappers

import (
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func OverviewToDBModels(
	id sharedvalueobjects.ID,
	o *domain.ResumeOverview,
) (gormmodel.ResumeOverview, []gormmodel.ResumeOverviewSection) {
	theme := o.Theme()
	greeting := o.Greeting()
	availability := o.Availability()
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
	if theme != nil {
		if theme.Primary() != nil {
			v := theme.Primary().String()
			m.ThemePrimary = &v
		}
		if theme.Secondary() != nil {
			v := theme.Secondary().String()
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
			Style:     x.Style(),
		})
	}
	return m, mapped
}

func OverviewToDomain(
	m gormmodel.ResumeOverview,
	s []gormmodel.ResumeOverviewSection,
) *domain.ResumeOverview {

	g := domain.RestoreGreeting(domain.GreetingParams{
		First: m.GreetingFirst, Second: m.GreetingSecond,
		SecondHighlight: m.GreetingHighlight, Third: m.GreetingThird,
	})
	var availability *domain.Availability
	if m.AvailabilityStatus != nil {
		params := domain.AvailabilityParams{Status: domain.AvailabilityStatus(*m.AvailabilityStatus), Note: m.AvailabilityNote}
		if m.AvailabilityFrom != nil {
			date := domain.RestoreDate(*m.AvailabilityFrom)
			params.From = &date
		}
		value := domain.RestoreAvailability(params)
		availability = &value
	}
	var theme *domain.ThemeSettings
	if m.ThemePrimary != nil || m.ThemeSecondary != nil {
		var primary, secondary *domain.Color
		if m.ThemePrimary != nil {
			value := domain.RestoreColor(*m.ThemePrimary)
			primary = &value
		}
		if m.ThemeSecondary != nil {
			value := domain.RestoreColor(*m.ThemeSecondary)
			secondary = &value
		}
		value := domain.RestoreThemeSettings(domain.ThemeSettingsParams{Primary: primary, Secondary: secondary})
		theme = &value
	}
	sections := make([]domain.SectionSettings, 0, len(s))
	for _, item := range s {
		sections = append(sections, domain.RestoreSectionSettings(domain.SectionSettingsParams{
			ID: domain.SectionID(item.SectionID), Enabled: item.Enabled, Collapsed: item.Collapsed, Style: item.Style,
		}))
	}
	return domain.RestoreResumeOverview(domain.ResumeOverviewParams{
		Theme: theme, Greeting: g, Availability: availability, Sections: sections,
	})
}

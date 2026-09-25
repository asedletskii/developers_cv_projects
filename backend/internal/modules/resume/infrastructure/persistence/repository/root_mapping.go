package repository

import (
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

func rootModels(r *domain.ResumeRoot) (gormmodel.ResumeRoot, gormmodel.ResumePresentation, []gormmodel.ResumePresentationSection, error) {
	if r == nil {
		return gormmodel.ResumeRoot{}, gormmodel.ResumePresentation{}, nil, fmt.Errorf("resume root is nil")
	}
	p := r.Presentation()
	m := gormmodel.ResumeRoot{ID: r.ID().UUID(), UserID: r.UserID().UUID(), Title: r.Title()}
	pm := gormmodel.ResumePresentation{ResumeID: r.ID().UUID(), GreetingFirstRU: p.Greeting.First().RU(), GreetingFirstEN: p.Greeting.First().EN(), GreetingSecondRU: p.Greeting.Second().RU(), GreetingSecondEN: p.Greeting.Second().EN(), GreetingHighlightRU: p.Greeting.SecondHighlight().RU(), GreetingHighlightEN: p.Greeting.SecondHighlight().EN(), GreetingThirdRU: p.Greeting.Third().RU(), GreetingThirdEN: p.Greeting.Third().EN()}
	if p.Availability != nil {
		s := string(p.Availability.Status())
		pm.AvailabilityStatus = &s
		if p.Availability.From() != nil {
			v := p.Availability.From().Time()
			pm.AvailabilityFrom = &v
		}
		if p.Availability.Note() != nil {
			ru := p.Availability.Note().RU()
			en := p.Availability.Note().EN()
			pm.AvailabilityNoteRU = &ru
			pm.AvailabilityNoteEN = &en
		}
	}
	if p.ResumeLink != nil {
		v := p.ResumeLink.String()
		pm.ResumeLink = &v
	}
	if p.Theme != nil {
		if p.Theme.Primary != nil {
			v := p.Theme.Primary.String()
			pm.ThemePrimary = &v
		}
		if p.Theme.Secondary != nil {
			v := p.Theme.Secondary.String()
			pm.ThemeSecondary = &v
		}
	}
	ss := make([]gormmodel.ResumePresentationSection, 0, len(p.Sections))
	for i, s := range p.Sections {
		ss = append(ss, gormmodel.ResumePresentationSection{ResumeID: r.ID().UUID(), SectionID: string(s.ID()), Position: i, Enabled: s.Enabled(), Collapsed: s.Collapsed()})
	}
	return m, pm, ss, nil
}
func mapRoot(m gormmodel.ResumeRoot, p gormmodel.ResumePresentation, ss []gormmodel.ResumePresentationSection) (*domain.ResumeRoot, error) {
	id, e := valueobjects.IDFromUUID(m.ID)
	if e != nil {
		return nil, e
	}
	uid, e := useridentity.IDFromUUID(m.UserID)
	if e != nil {
		return nil, e
	}
	f, e := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: p.GreetingFirstRU, EN: p.GreetingFirstEN})
	if e != nil {
		return nil, e
	}
	s, e := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: p.GreetingSecondRU, EN: p.GreetingSecondEN})
	if e != nil {
		return nil, e
	}
	h, e := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: p.GreetingHighlightRU, EN: p.GreetingHighlightEN})
	if e != nil {
		return nil, e
	}
	t, e := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: p.GreetingThirdRU, EN: p.GreetingThirdEN})
	if e != nil {
		return nil, e
	}
	g := valueobjects.NewGreeting(valueobjects.GreetingParams{First: f, Second: s, SecondHighlight: h, Third: t})
	var av *valueobjects.Availability
	if p.AvailabilityStatus != nil {
		var from *valueobjects.Date
		if p.AvailabilityFrom != nil {
			d, e := valueobjects.NewDate(*p.AvailabilityFrom)
			if e != nil {
				return nil, e
			}
			from = &d
		}
		var note *valueobjects.LocalizedText
		if p.AvailabilityNoteRU != nil || p.AvailabilityNoteEN != nil {
			ru, en := "", ""
			if p.AvailabilityNoteRU != nil {
				ru = *p.AvailabilityNoteRU
			}
			if p.AvailabilityNoteEN != nil {
				en = *p.AvailabilityNoteEN
			}
			v, e := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: ru, EN: en})
			if e != nil {
				return nil, e
			}
			note = &v
		}
		v, e := valueobjects.NewAvailability(valueobjects.AvailabilityParams{Status: enums.AvailabilityStatus(*p.AvailabilityStatus), From: from, Note: note})
		if e != nil {
			return nil, e
		}
		av = &v
	}
	var theme *valueobjects.ThemeSettings
	if p.ThemePrimary != nil || p.ThemeSecondary != nil {
		var a, b *valueobjects.Color
		if p.ThemePrimary != nil {
			v, e := valueobjects.NewColor(*p.ThemePrimary)
			if e != nil {
				return nil, e
			}
			a = &v
		}
		if p.ThemeSecondary != nil {
			v, e := valueobjects.NewColor(*p.ThemeSecondary)
			if e != nil {
				return nil, e
			}
			b = &v
		}
		v, e := valueobjects.NewThemeSettings(valueobjects.ThemeSettingsParams{Primary: a, Secondary: b})
		if e != nil {
			return nil, e
		}
		theme = &v
	}
	var link *valueobjects.AssetPath
	if p.ResumeLink != nil {
		v, e := valueobjects.NewAssetPath(*p.ResumeLink)
		if e != nil {
			return nil, e
		}
		link = &v
	}
	sections := make([]valueobjects.SectionSettings, 0, len(ss))
	for _, x := range ss {
		v, e := valueobjects.NewSectionSettings(valueobjects.SectionSettingsParams{ID: enums.SectionID(x.SectionID), Enabled: x.Enabled, Collapsed: x.Collapsed})
		if e != nil {
			return nil, e
		}
		sections = append(sections, v)
	}
	return domain.NewResumeRoot(domain.ResumeRootParams{ID: id, UserID: uid, Title: m.Title, Presentation: domain.ResumePresentation{Theme: theme, Greeting: g, Availability: av, ResumeLink: link, Sections: sections}})
}

package queries

import (
	"time"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

type ResumeRootView struct {
	ID           string                 `json:"id"`
	Title        string                 `json:"title"`
	CreatedAt    time.Time              `json:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt"`
	Presentation ResumePresentationView `json:"presentation"`
}

type ResumePresentationView struct {
	Theme        *ResumeRootThemeSettings    `json:"theme,omitempty"`
	Greeting     ResumeRootGreeting          `json:"greeting"`
	Availability *ResumeRootAvailability     `json:"availability,omitempty"`
	ResumeLink   *string                     `json:"resumeLink,omitempty"`
	Sections     []ResumeRootSectionSettings `json:"sections"`
}

type ResumeRootLocalizedText struct {
	RU string `json:"ru"`
	EN string `json:"en"`
}

type ResumeRootGreeting struct {
	First           ResumeRootLocalizedText `json:"first"`
	Second          ResumeRootLocalizedText `json:"second"`
	SecondHighlight ResumeRootLocalizedText `json:"secondHighlight"`
	Third           ResumeRootLocalizedText `json:"third"`
}

type ResumeRootAvailability struct {
	Status string                   `json:"status"`
	From   *string                  `json:"from,omitempty"`
	Note   *ResumeRootLocalizedText `json:"note,omitempty"`
}

type ResumeRootThemeSettings struct {
	Primary   string `json:"primary,omitempty"`
	Secondary string `json:"secondary,omitempty"`
}

type ResumeRootSectionSettings struct {
	ID        string `json:"id"`
	Enabled   bool   `json:"enabled"`
	Collapsed bool   `json:"collapsed"`
}

func ResumeRootToView(r *domain.ResumeRoot) ResumeRootView {
	p := r.Presentation()
	var theme *ResumeRootThemeSettings
	if p.Theme != nil {
		theme = &ResumeRootThemeSettings{}
		if p.Theme.Primary != nil {
			theme.Primary = p.Theme.Primary.String()
		}
		if p.Theme.Secondary != nil {
			theme.Secondary = p.Theme.Secondary.String()
		}
	}
	var availability *ResumeRootAvailability
	if p.Availability != nil {
		availability = &ResumeRootAvailability{Status: string(p.Availability.Status())}
		if p.Availability.From() != nil {
			v := p.Availability.From().Time().Format("2006-01-02")
			availability.From = &v
		}
		if p.Availability.Note() != nil {
			v := localizedTextView(*p.Availability.Note())
			availability.Note = &v
		}
	}
	var link *string
	if p.ResumeLink != nil {
		v := p.ResumeLink.String()
		link = &v
	}
	return ResumeRootView{
		ID: r.ID().String(), Title: r.Title(), CreatedAt: r.CreatedAt(), UpdatedAt: r.UpdatedAt(),
		Presentation: ResumePresentationView{
			Theme: theme,
			Greeting: ResumeRootGreeting{
				First: localizedTextView(p.Greeting.First()), Second: localizedTextView(p.Greeting.Second()),
				SecondHighlight: localizedTextView(p.Greeting.SecondHighlight()), Third: localizedTextView(p.Greeting.Third()),
			},
			Availability: availability, ResumeLink: link, Sections: mapSections(p.Sections),
		},
	}
}

func mapSections(values []valueobjects.SectionSettings) []ResumeRootSectionSettings {
	result := make([]ResumeRootSectionSettings, 0, len(values))
	for _, v := range values {
		result = append(result, ResumeRootSectionSettings{ID: string(v.ID()), Enabled: v.Enabled(), Collapsed: v.Collapsed()})
	}
	return result
}

func localizedTextView(v valueobjects.LocalizedText) ResumeRootLocalizedText {
	return ResumeRootLocalizedText{RU: v.RU(), EN: v.EN()}
}

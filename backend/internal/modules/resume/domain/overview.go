package domain

import (
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"fmt"
)

type ResumeOverview struct{ presentation OverviewPresentation }
type OverviewPresentation struct {
	Theme        *valueobjects.ThemeSettings
	Greeting     valueobjects.Greeting
	Availability *valueobjects.Availability
	ResumeLink   *valueobjects.AssetPath
	Sections     []valueobjects.SectionSettings
}
type ResumeOverviewParams struct{ Presentation OverviewPresentation }

func NewResumeOverview(params ResumeOverviewParams) (*ResumeOverview, error) {
	if err := validateOverviewSections(params.Presentation.Sections); err != nil {
		return nil, err
	}
	return &ResumeOverview{presentation: cloneOverviewPresentation(params.Presentation)}, nil
}
func (o *ResumeOverview) Presentation() OverviewPresentation {
	return cloneOverviewPresentation(o.presentation)
}
func (o *ResumeOverview) ReplacePresentation(p OverviewPresentation) error {
	if err := validateOverviewSections(p.Sections); err != nil {
		return err
	}
	o.presentation = cloneOverviewPresentation(p)
	return nil
}
func validateOverviewSections(sections []valueobjects.SectionSettings) error {
	seen := make(map[string]struct{}, len(sections))
	for _, s := range sections {
		if !s.ID().IsValid() {
			return fmt.Errorf("%w: invalid section id", valueobjects.ErrInvalidValue)
		}
		k := string(s.ID())
		if _, ok := seen[k]; ok {
			return fmt.Errorf("%w: duplicate resume section", ErrInvalidResume)
		}
		seen[k] = struct{}{}
	}
	return nil
}
func cloneOverviewPresentation(p OverviewPresentation) OverviewPresentation {
	var theme *valueobjects.ThemeSettings
	if p.Theme != nil {
		t := *p.Theme
		if p.Theme.Primary != nil {
			v := *p.Theme.Primary
			t.Primary = &v
		}
		if p.Theme.Secondary != nil {
			v := *p.Theme.Secondary
			t.Secondary = &v
		}
		theme = &t
	}
	var av *valueobjects.Availability
	if p.Availability != nil {
		v := *p.Availability
		av = &v
	}
	var link *valueobjects.AssetPath
	if p.ResumeLink != nil {
		v := *p.ResumeLink
		link = &v
	}
	var sections []valueobjects.SectionSettings
	if p.Sections != nil {
		sections = append([]valueobjects.SectionSettings(nil), p.Sections...)
	}
	return OverviewPresentation{Theme: theme, Greeting: p.Greeting, Availability: av, ResumeLink: link, Sections: sections}
}

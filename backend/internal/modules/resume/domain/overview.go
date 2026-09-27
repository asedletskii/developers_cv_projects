package domain

import (
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

type ResumeOverviewParams struct {
	Theme        *valueobjects.ThemeSettings
	Greeting     valueobjects.Greeting
	Availability *valueobjects.Availability
	ResumeLink   *valueobjects.ExternalURL
	Sections     []valueobjects.SectionSettings
}

type ResumeOverview struct {
	theme        *valueobjects.ThemeSettings
	greeting     valueobjects.Greeting
	availability *valueobjects.Availability
	resumeLink   *valueobjects.ExternalURL
	sections     []valueobjects.SectionSettings
}

func NewResumeOverview(params ResumeOverviewParams) (*ResumeOverview, error) {
	if err := validateOverviewSections(params.Sections); err != nil {
		return nil, err
	}
	return &ResumeOverview{
		theme:        cloneThemeSettings(params.Theme),
		greeting:     params.Greeting,
		availability: cloneAvailability(params.Availability),
		resumeLink:   cloneExternalURL(params.ResumeLink),
		sections:     cloneSectionSettings(params.Sections),
	}, nil
}

func (o *ResumeOverview) Theme() *valueobjects.ThemeSettings { return cloneThemeSettings(o.theme) }
func (o *ResumeOverview) Greeting() valueobjects.Greeting    { return o.greeting }
func (o *ResumeOverview) Availability() *valueobjects.Availability {
	return cloneAvailability(o.availability)
}
func (o *ResumeOverview) ResumeLink() *valueobjects.ExternalURL {
	return cloneExternalURL(o.resumeLink)
}
func (o *ResumeOverview) Sections() []valueobjects.SectionSettings {
	return cloneSectionSettings(o.sections)
}

func (o *ResumeOverview) Replace(params ResumeOverviewParams) error {
	if err := validateOverviewSections(params.Sections); err != nil {
		return err
	}
	o.theme = cloneThemeSettings(params.Theme)
	o.greeting = params.Greeting
	o.availability = cloneAvailability(params.Availability)
	o.resumeLink = cloneExternalURL(params.ResumeLink)
	o.sections = cloneSectionSettings(params.Sections)
	return nil
}

func validateOverviewSections(sections []valueobjects.SectionSettings) error {
	seen := make(map[string]struct{}, len(sections))
	for _, section := range sections {
		if !section.ID().IsValid() {
			return fmt.Errorf("%w: invalid section id", valueobjects.ErrInvalidValue)
		}
		key := string(section.ID())
		if _, ok := seen[key]; ok {
			return fmt.Errorf("%w: duplicate resume section", ErrInvalidResume)
		}
		seen[key] = struct{}{}
	}
	return nil
}

func cloneThemeSettings(value *valueobjects.ThemeSettings) *valueobjects.ThemeSettings {
	if value == nil {
		return nil
	}
	copy := *value
	if value.Primary != nil {
		primary := *value.Primary
		copy.Primary = &primary
	}
	if value.Secondary != nil {
		secondary := *value.Secondary
		copy.Secondary = &secondary
	}
	return &copy
}

func cloneAvailability(value *valueobjects.Availability) *valueobjects.Availability {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

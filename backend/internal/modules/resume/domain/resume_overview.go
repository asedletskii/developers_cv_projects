package domain

import (
	"fmt"
	"regexp"
	"strings"
)

type ResumeOverviewParams struct {
	Theme        *ThemeSettings
	Greeting     Greeting
	Availability *Availability
	ResumeLink   *ExternalURL
	Sections     []SectionSettings
}

type ResumeOverview struct {
	theme        *ThemeSettings
	greeting     Greeting
	availability *Availability
	resumeLink   *ExternalURL
	sections     []SectionSettings
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

func (o *ResumeOverview) Theme() *ThemeSettings { return cloneThemeSettings(o.theme) }
func (o *ResumeOverview) Greeting() Greeting    { return o.greeting }
func (o *ResumeOverview) Availability() *Availability {
	return cloneAvailability(o.availability)
}
func (o *ResumeOverview) ResumeLink() *ExternalURL {
	return cloneExternalURL(o.resumeLink)
}
func (o *ResumeOverview) Sections() []SectionSettings {
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

func validateOverviewSections(sections []SectionSettings) error {
	seen := make(map[string]struct{}, len(sections))
	for _, section := range sections {
		if !section.ID().IsValid() {
			return fmt.Errorf("%w: invalid section id", ErrInvalidValue)
		}
		key := string(section.ID())
		if _, ok := seen[key]; ok {
			return fmt.Errorf("%w: duplicate resume section", ErrInvalidResume)
		}
		seen[key] = struct{}{}
	}
	return nil
}

func cloneThemeSettings(value *ThemeSettings) *ThemeSettings {
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

func cloneAvailability(value *Availability) *Availability {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func (r *Resume) ReplaceOverview(params ResumeOverviewParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	o, err := NewResumeOverview(params)
	if err != nil {
		return err
	}
	r.overview = o
	return nil
}

func (r *Resume) Overview() *ResumeOverview {
	if r == nil || r.overview == nil {
		return nil
	}
	copy := *r.overview
	return &copy
}

var hexColorPattern = regexp.MustCompile(`^#[0-9a-fA-F]{3}(?:[0-9a-fA-F]{3})?$`)

type Color struct {
	value string
}

func NewColor(value string) (Color, error) {
	value = strings.TrimSpace(value)
	if !hexColorPattern.MatchString(value) {
		return Color{}, fmt.Errorf("%w: invalid color", ErrInvalidValue)
	}
	return Color{value: value}, nil
}

func (c Color) String() string { return c.value }

type ThemeSettingsParams struct {
	Primary   *Color
	Secondary *Color
}

type ThemeSettings struct {
	Primary   *Color
	Secondary *Color
}

func NewThemeSettings(params ThemeSettingsParams) (ThemeSettings, error) {
	if params.Primary == nil && params.Secondary == nil {
		return ThemeSettings{}, fmt.Errorf("%w: theme is empty", ErrInvalidValue)
	}
	return ThemeSettings{
		Primary:   cloneColor(params.Primary),
		Secondary: cloneColor(params.Secondary),
	}, nil
}

func cloneColor(value *Color) *Color {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

type AvailabilityStatus string

const (
	AvailabilityStatusOpen        AvailabilityStatus = "open"
	AvailabilityStatusConsidering AvailabilityStatus = "considering"
	AvailabilityStatusBusy        AvailabilityStatus = "busy"
)

func (v AvailabilityStatus) IsValid() bool {
	switch v {
	case AvailabilityStatusOpen, AvailabilityStatusConsidering, AvailabilityStatusBusy:
		return true
	default:
		return false
	}
}

type AvailabilityParams struct {
	Status AvailabilityStatus
	From   *Date
	Note   *string
}

type Availability struct {
	status AvailabilityStatus
	from   *Date
	note   *string
}

func NewAvailability(params AvailabilityParams) (Availability, error) {
	if !params.Status.IsValid() {
		return Availability{}, fmt.Errorf("%w: invalid availability status", ErrInvalidValue)
	}
	return Availability{status: params.Status, from: cloneDate(params.From), note: cloneString(params.Note)}, nil
}

func (a Availability) Status() AvailabilityStatus { return a.status }
func (a Availability) From() *Date                { return cloneDate(a.from) }
func (a Availability) Note() *string              { return cloneString(a.note) }

type GreetingParams struct {
	First           string
	Second          string
	SecondHighlight string
	Third           string
}

type Greeting struct {
	first           string
	second          string
	secondHighlight string
	third           string
}

func NewGreeting(params GreetingParams) Greeting {
	return Greeting{
		first:           strings.TrimSpace(params.First),
		second:          strings.TrimSpace(params.Second),
		secondHighlight: strings.TrimSpace(params.SecondHighlight),
		third:           strings.TrimSpace(params.Third),
	}
}

func (g Greeting) First() string           { return g.first }
func (g Greeting) Second() string          { return g.second }
func (g Greeting) SecondHighlight() string { return g.secondHighlight }
func (g Greeting) Third() string           { return g.third }

type SectionSettingsParams struct {
	ID        SectionID
	Enabled   bool
	Collapsed bool
}

type SectionSettings struct {
	id        SectionID
	enabled   bool
	collapsed bool
}

func NewSectionSettings(params SectionSettingsParams) (SectionSettings, error) {
	if !params.ID.IsValid() {
		return SectionSettings{}, fmt.Errorf("%w: invalid section id", ErrInvalidValue)
	}
	return SectionSettings{id: params.ID, enabled: params.Enabled, collapsed: params.Collapsed}, nil
}

func (s SectionSettings) ID() SectionID   { return s.id }
func (s SectionSettings) Enabled() bool   { return s.enabled }
func (s SectionSettings) Collapsed() bool { return s.collapsed }

type SectionID string

const (
	SectionAbout      SectionID = "about"
	SectionTechStack  SectionID = "techStack"
	SectionExperience SectionID = "experience"
	SectionSkills     SectionID = "skills"
	SectionEducation  SectionID = "education"
	SectionCourses    SectionID = "courses"
	SectionProjects   SectionID = "projects"
	SectionActivities SectionID = "activities"
	SectionContacts   SectionID = "contacts"
)

func (v SectionID) IsValid() bool {
	switch v {
	case SectionAbout, SectionTechStack, SectionExperience, SectionSkills, SectionEducation,
		SectionCourses, SectionProjects, SectionActivities, SectionContacts:
		return true
	default:
		return false
	}
}

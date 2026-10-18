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
	Sections     []SectionSettings
}

type ResumeOverview struct {
	theme        *ThemeSettings
	greeting     Greeting
	availability *Availability
	sections     []SectionSettings
}

func newResumeOverview(params ResumeOverviewParams) (*ResumeOverview, error) {
	overview := RestoreResumeOverview(params)
	if err := overview.validate(); err != nil {
		return nil, err
	}
	return overview, nil
}

func RestoreResumeOverview(params ResumeOverviewParams) *ResumeOverview {
	return &ResumeOverview{
		theme: cloneThemeSettings(params.Theme), greeting: params.Greeting,
		availability: cloneAvailability(params.Availability),
		sections:     cloneSectionSettings(params.Sections),
	}
}

func (o *ResumeOverview) validate() error {
	if o == nil {
		return fmt.Errorf("%w: overview is nil", ErrInvalidValue)
	}
	params := ResumeOverviewParams{
		Theme: o.theme, Greeting: o.greeting, Availability: o.availability,
		Sections: o.sections,
	}
	if err := validateOverviewSections(params.Sections); err != nil {
		return err
	}
	if err := validateOverviewValueObjects(params); err != nil {
		return err
	}
	return params.Greeting.validate()
}

func (o *ResumeOverview) Theme() *ThemeSettings { return cloneThemeSettings(o.theme) }
func (o *ResumeOverview) Greeting() Greeting    { return o.greeting }
func (o *ResumeOverview) Availability() *Availability {
	return cloneAvailability(o.availability)
}
func (o *ResumeOverview) Sections() []SectionSettings {
	return cloneSectionSettings(o.sections)
}

func validateOverviewValueObjects(params ResumeOverviewParams) error {
	if params.Theme != nil {
		if err := params.Theme.validate(); err != nil {
			return err
		}
	}
	if params.Availability != nil {
		if err := params.Availability.validate(); err != nil {
			return err
		}
		if from := params.Availability.From(); from != nil && from.IsZero() {
			return fmt.Errorf("%w: invalid availability date", ErrInvalidValue)
		}
	}
	return nil
}

func validateOverviewSections(sections []SectionSettings) error {
	if len(sections) > 9 {
		return fmt.Errorf("%w: overview cannot contain more than 9 sections", ErrInvalidValue)
	}
	seen := make(map[string]struct{}, len(sections))
	for _, section := range sections {
		if err := section.validate(); err != nil {
			return err
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
	return &ThemeSettings{primary: cloneColor(value.primary), secondary: cloneColor(value.secondary)}
}

func cloneAvailability(value *Availability) *Availability {
	if value == nil {
		return nil
	}
	copy := *value
	copy.note = copyString(value.note)
	copy.from = cloneDate(value.from)
	return &copy
}

func (r *Resume) ReplaceOverview(params ResumeOverviewParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	o, err := newResumeOverview(params)
	if err != nil {
		return err
	}
	r.overview = o
	return nil
}

func (r *Resume) Overview() *ResumeOverview {
	if r == nil {
		return nil
	}
	return cloneResumeOverview(r.overview)
}

var hexColorPattern = regexp.MustCompile(`^#[0-9a-fA-F]{3}(?:[0-9a-fA-F]{3})?$`)

type Color struct {
	value string
}

func NewColor(value string) (Color, error) {
	color := Color{value: strings.TrimSpace(value)}
	if err := color.validate(); err != nil {
		return Color{}, err
	}
	return color, nil
}

func RestoreColor(value string) Color { return Color{value: value} }

func (c Color) validate() error {
	if !hexColorPattern.MatchString(c.value) {
		return fmt.Errorf("%w: invalid color", ErrInvalidValue)
	}
	return nil
}

func (c Color) String() string { return c.value }

type ThemeSettingsParams struct {
	Primary   *Color
	Secondary *Color
}

type ThemeSettings struct {
	primary   *Color
	secondary *Color
}

func NewThemeSettings(params ThemeSettingsParams) (ThemeSettings, error) {
	theme := RestoreThemeSettings(params)
	if err := theme.validate(); err != nil {
		return ThemeSettings{}, err
	}
	return theme, nil
}

func RestoreThemeSettings(params ThemeSettingsParams) ThemeSettings {
	return ThemeSettings{primary: cloneColor(params.Primary), secondary: cloneColor(params.Secondary)}
}

func (t ThemeSettings) validate() error {
	if t.isEmpty() {
		return fmt.Errorf("%w: theme is empty", ErrInvalidValue)
	}
	if t.primary != nil {
		if err := t.primary.validate(); err != nil {
			return fmt.Errorf("primary theme color: %w", err)
		}
	}
	if t.secondary != nil {
		if err := t.secondary.validate(); err != nil {
			return fmt.Errorf("secondary theme color: %w", err)
		}
	}
	return nil
}

func (t ThemeSettings) Primary() *Color { return cloneColor(t.primary) }

func (t ThemeSettings) Secondary() *Color { return cloneColor(t.secondary) }

func (t ThemeSettings) isEmpty() bool { return t.primary == nil && t.secondary == nil }

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

func (v AvailabilityStatus) isValid() bool {
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
	availability := RestoreAvailability(params)
	availability.note = normalizeOptionalString(params.Note)
	if err := availability.validate(); err != nil {
		return Availability{}, err
	}
	return availability, nil
}

func RestoreAvailability(params AvailabilityParams) Availability {
	return Availability{status: params.Status, from: cloneDate(params.From), note: copyString(params.Note)}
}

func (a Availability) validate() error {
	if !a.status.isValid() {
		return fmt.Errorf("%w: invalid availability status", ErrInvalidValue)
	}
	if !optionalTextLengthWithin(a.note, 500) {
		return fmt.Errorf("%w: availability note exceeds 500 characters", ErrInvalidValue)
	}
	return nil
}

func (a Availability) Status() AvailabilityStatus { return a.status }
func (a Availability) From() *Date                { return cloneDate(a.from) }
func (a Availability) Note() *string              { return copyString(a.note) }

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

func NewGreeting(params GreetingParams) (Greeting, error) {
	greeting := RestoreGreeting(GreetingParams{
		First: strings.TrimSpace(params.First), Second: strings.TrimSpace(params.Second),
		SecondHighlight: strings.TrimSpace(params.SecondHighlight), Third: strings.TrimSpace(params.Third),
	})
	if err := greeting.validate(); err != nil {
		return Greeting{}, err
	}
	return greeting, nil
}

func (g Greeting) validate() error {
	if !textLengthWithin(g.first, MaxShortTextLength) || !textLengthWithin(g.second, MaxShortTextLength) ||
		!textLengthWithin(g.secondHighlight, 120) || !textLengthWithin(g.third, MaxShortTextLength) {
		return fmt.Errorf("%w: overview greeting exceeds configured limits", ErrInvalidValue)
	}
	return nil
}

func RestoreGreeting(params GreetingParams) Greeting {
	return Greeting{first: params.First, second: params.Second, secondHighlight: params.SecondHighlight, third: params.Third}
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
	section := RestoreSectionSettings(params)
	if err := section.validate(); err != nil {
		return SectionSettings{}, err
	}
	return section, nil
}

func RestoreSectionSettings(params SectionSettingsParams) SectionSettings {
	return SectionSettings{id: params.ID, enabled: params.Enabled, collapsed: params.Collapsed}
}

func (s SectionSettings) validate() error {
	if !s.id.isValid() {
		return fmt.Errorf("%w: invalid section id", ErrInvalidValue)
	}
	return nil
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

func (v SectionID) isValid() bool {
	switch v {
	case SectionAbout, SectionTechStack, SectionExperience, SectionSkills, SectionEducation,
		SectionCourses, SectionProjects, SectionActivities, SectionContacts:
		return true
	default:
		return false
	}
}

func ParseSectionID(value string) (SectionID, error) {
	section := SectionID(value)
	if !section.isValid() {
		return "", fmt.Errorf("%w: invalid section id", ErrInvalidValue)
	}
	return section, nil
}

func cloneResumeOverview(value *ResumeOverview) *ResumeOverview {
	if value == nil {
		return nil
	}
	copy := *value
	copy.theme = cloneThemeSettings(value.theme)
	copy.availability = cloneAvailability(value.availability)
	copy.sections = cloneSectionSettings(value.sections)
	return &copy
}

package valueobjects

import (
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain/enums"
)

type AvailabilityParams struct {
	Status enums.AvailabilityStatus
	From   *Date
	Note   *string
}

type Availability struct {
	status enums.AvailabilityStatus
	from   *Date
	note   *string
}

func NewAvailability(params AvailabilityParams) (Availability, error) {
	if !params.Status.IsValid() {
		return Availability{}, fmt.Errorf("%w: invalid availability status", ErrInvalidValue)
	}
	return Availability{status: params.Status, from: cloneDate(params.From), note: cloneString(params.Note)}, nil
}

func (a Availability) Status() enums.AvailabilityStatus { return a.status }
func (a Availability) From() *Date                      { return cloneDate(a.from) }
func (a Availability) Note() *string                    { return cloneString(a.note) }

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
	return Greeting{first: params.First, second: params.Second, secondHighlight: params.SecondHighlight, third: params.Third}
}

func (g Greeting) First() string           { return g.first }
func (g Greeting) Second() string          { return g.second }
func (g Greeting) SecondHighlight() string { return g.secondHighlight }
func (g Greeting) Third() string           { return g.third }

type ContactsTextParams struct {
	Title    *string
	Subtitle *string
}

type ContactsText struct {
	title    *string
	subtitle *string
}

func NewContactsText(params ContactsTextParams) (ContactsText, error) {
	if params.Title == nil && params.Subtitle == nil {
		return ContactsText{}, fmt.Errorf("%w: contacts text is empty", ErrInvalidValue)
	}
	return ContactsText{title: cloneString(params.Title), subtitle: cloneString(params.Subtitle)}, nil
}

func (c ContactsText) Title() *string    { return cloneString(c.title) }
func (c ContactsText) Subtitle() *string { return cloneString(c.subtitle) }

type SectionSettingsParams struct {
	ID        enums.SectionID
	Enabled   bool
	Collapsed bool
}

type SectionSettings struct {
	id        enums.SectionID
	enabled   bool
	collapsed bool
}

func NewSectionSettings(params SectionSettingsParams) (SectionSettings, error) {
	if !params.ID.IsValid() {
		return SectionSettings{}, fmt.Errorf("%w: invalid section id", ErrInvalidValue)
	}
	return SectionSettings{id: params.ID, enabled: params.Enabled, collapsed: params.Collapsed}, nil
}

func (s SectionSettings) ID() enums.SectionID { return s.id }
func (s SectionSettings) Enabled() bool       { return s.enabled }
func (s SectionSettings) Collapsed() bool     { return s.collapsed }

func cloneDate(value *Date) *Date {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func cloneString(value *string) *string {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

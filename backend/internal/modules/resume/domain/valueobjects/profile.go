package valueobjects

import (
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain/enums"
)

type AvailabilityParams struct {
	Status enums.AvailabilityStatus
	From   *Date
	Note   *LocalizedText
}

type Availability struct {
	status enums.AvailabilityStatus
	from   *Date
	note   *LocalizedText
}

func NewAvailability(params AvailabilityParams) (Availability, error) {
	if !params.Status.IsValid() {
		return Availability{}, fmt.Errorf("%w: invalid availability status", ErrInvalidValue)
	}
	return Availability{
		status: params.Status,
		from:   cloneDate(params.From),
		note:   cloneLocalizedText(params.Note),
	}, nil
}

func (a Availability) Status() enums.AvailabilityStatus { return a.status }
func (a Availability) From() *Date                      { return cloneDate(a.from) }
func (a Availability) Note() *LocalizedText             { return cloneLocalizedText(a.note) }

type GreetingParams struct {
	First           LocalizedText
	Second          LocalizedText
	SecondHighlight LocalizedText
	Third           LocalizedText
}

type Greeting struct {
	first           LocalizedText
	second          LocalizedText
	secondHighlight LocalizedText
	third           LocalizedText
}

func NewGreeting(params GreetingParams) Greeting {
	return Greeting{first: params.First, second: params.Second, secondHighlight: params.SecondHighlight, third: params.Third}
}

func (g Greeting) First() LocalizedText           { return g.first }
func (g Greeting) Second() LocalizedText          { return g.second }
func (g Greeting) SecondHighlight() LocalizedText { return g.secondHighlight }
func (g Greeting) Third() LocalizedText           { return g.third }

type ContactsTextParams struct {
	Title    *LocalizedText
	Subtitle *LocalizedText
}

type ContactsText struct {
	title    *LocalizedText
	subtitle *LocalizedText
}

func NewContactsText(params ContactsTextParams) (ContactsText, error) {
	if params.Title == nil && params.Subtitle == nil {
		return ContactsText{}, fmt.Errorf("%w: contacts text is empty", ErrInvalidValue)
	}
	return ContactsText{title: cloneLocalizedText(params.Title), subtitle: cloneLocalizedText(params.Subtitle)}, nil
}

func (c ContactsText) Title() *LocalizedText    { return cloneLocalizedText(c.title) }
func (c ContactsText) Subtitle() *LocalizedText { return cloneLocalizedText(c.subtitle) }

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

func cloneLocalizedText(value *LocalizedText) *LocalizedText {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

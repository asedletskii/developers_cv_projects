package root

import (
	"fmt"

	"github.com/google/uuid"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type Request struct {
	Title        string              `json:"title"`
	Presentation PresentationRequest `json:"presentation"`
}
type PresentationRequest struct {
	Theme        *ThemeRequest        `json:"theme,omitempty"`
	Greeting     GreetingRequest      `json:"greeting"`
	Availability *AvailabilityRequest `json:"availability,omitempty"`
	ResumeLink   *string              `json:"resumeLink,omitempty"`
	Sections     []SectionRequest     `json:"sections"`
}
type ThemeRequest struct {
	Primary   *string `json:"primary,omitempty"`
	Secondary *string `json:"secondary,omitempty"`
}
type LocalizedRequest struct {
	RU string `json:"ru"`
	EN string `json:"en"`
}
type GreetingRequest struct {
	First           LocalizedRequest `json:"first"`
	Second          LocalizedRequest `json:"second"`
	SecondHighlight LocalizedRequest `json:"secondHighlight"`
	Third           LocalizedRequest `json:"third"`
}
type AvailabilityRequest struct {
	Status string            `json:"status"`
	From   *string           `json:"from,omitempty"`
	Note   *LocalizedRequest `json:"note,omitempty"`
}
type SectionRequest struct {
	ID        string `json:"id"`
	Enabled   bool   `json:"enabled"`
	Collapsed bool   `json:"collapsed"`
}

func (r Request) toPresentation() (domain.ResumePresentation, error) {
	greeting, err := newGreeting(r.Presentation.Greeting)
	if err != nil {
		return domain.ResumePresentation{}, err
	}
	availability, err := newAvailability(r.Presentation.Availability)
	if err != nil {
		return domain.ResumePresentation{}, err
	}
	theme, err := newTheme(r.Presentation.Theme)
	if err != nil {
		return domain.ResumePresentation{}, err
	}
	link, err := optionalAssetPath(r.Presentation.ResumeLink)
	if err != nil {
		return domain.ResumePresentation{}, err
	}
	sections := make([]valueobjects.SectionSettings, 0, len(r.Presentation.Sections))
	for _, s := range r.Presentation.Sections {
		v, err := valueobjects.NewSectionSettings(valueobjects.SectionSettingsParams{ID: enums.SectionID(s.ID), Enabled: s.Enabled, Collapsed: s.Collapsed})
		if err != nil {
			return domain.ResumePresentation{}, err
		}
		sections = append(sections, v)
	}
	return domain.ResumePresentation{Theme: theme, Greeting: greeting, Availability: availability, ResumeLink: link, Sections: sections}, nil
}
func (r Request) ToCreateCommand(userID uuid.UUID) (commands.CreateResumeRootCommand, error) {
	uid, err := useridentity.IDFromUUID(userID)
	if err != nil {
		return commands.CreateResumeRootCommand{}, err
	}
	p, err := r.toPresentation()
	if err != nil {
		return commands.CreateResumeRootCommand{}, err
	}
	return commands.CreateResumeRootCommand{UserID: uid, Title: r.Title, Presentation: p}, nil
}
func (r Request) ToUpdateCommand(userID uuid.UUID, resumeID valueobjects.ID) (commands.UpdateResumeRootCommand, error) {
	uid, err := useridentity.IDFromUUID(userID)
	if err != nil {
		return commands.UpdateResumeRootCommand{}, err
	}
	p, err := r.toPresentation()
	if err != nil {
		return commands.UpdateResumeRootCommand{}, err
	}
	return commands.UpdateResumeRootCommand{UserID: uid, ResumeID: resumeID, Title: r.Title, Presentation: p}, nil
}
func newGreeting(r GreetingRequest) (valueobjects.Greeting, error) {
	f, e := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: r.First.RU, EN: r.First.EN})
	if e != nil {
		return valueobjects.Greeting{}, e
	}
	s, e := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: r.Second.RU, EN: r.Second.EN})
	if e != nil {
		return valueobjects.Greeting{}, e
	}
	h, e := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: r.SecondHighlight.RU, EN: r.SecondHighlight.EN})
	if e != nil {
		return valueobjects.Greeting{}, e
	}
	t, e := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: r.Third.RU, EN: r.Third.EN})
	if e != nil {
		return valueobjects.Greeting{}, e
	}
	return valueobjects.NewGreeting(valueobjects.GreetingParams{First: f, Second: s, SecondHighlight: h, Third: t}), nil
}
func newAvailability(r *AvailabilityRequest) (*valueobjects.Availability, error) {
	if r == nil {
		return nil, nil
	}
	var from *valueobjects.Date
	if r.From != nil {
		v, e := valueobjects.ParseDate(*r.From)
		if e != nil {
			return nil, e
		}
		from = &v
	}
	var note *valueobjects.LocalizedText
	if r.Note != nil {
		v, e := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: r.Note.RU, EN: r.Note.EN})
		if e != nil {
			return nil, e
		}
		note = &v
	}
	v, e := valueobjects.NewAvailability(valueobjects.AvailabilityParams{Status: enums.AvailabilityStatus(r.Status), From: from, Note: note})
	if e != nil {
		return nil, e
	}
	return &v, nil
}
func newTheme(r *ThemeRequest) (*valueobjects.ThemeSettings, error) {
	if r == nil {
		return nil, nil
	}
	var primary, secondary *valueobjects.Color
	if r.Primary != nil {
		v, e := valueobjects.NewColor(*r.Primary)
		if e != nil {
			return nil, e
		}
		primary = &v
	}
	if r.Secondary != nil {
		v, e := valueobjects.NewColor(*r.Secondary)
		if e != nil {
			return nil, e
		}
		secondary = &v
	}
	v, e := valueobjects.NewThemeSettings(valueobjects.ThemeSettingsParams{Primary: primary, Secondary: secondary})
	if e != nil {
		return nil, e
	}
	return &v, nil
}
func optionalAssetPath(v *string) (*valueobjects.AssetPath, error) {
	if v == nil {
		return nil, nil
	}
	x, e := valueobjects.NewAssetPath(*v)
	if e != nil {
		return nil, fmt.Errorf("resume link: %w", e)
	}
	return &x, nil
}

package domain

import (
	"fmt"
	"strings"

	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type ResumeRoot struct {
	id           valueobjects.ID
	userID       useridentity.ID
	title        string
	presentation ResumePresentation
}

type ResumePresentation struct {
	Theme        *valueobjects.ThemeSettings
	Greeting     valueobjects.Greeting
	Availability *valueobjects.Availability
	ResumeLink   *valueobjects.AssetPath
	Sections     []valueobjects.SectionSettings
}

type ResumeRootParams struct {
	ID           valueobjects.ID
	UserID       useridentity.ID
	Title        string
	Presentation ResumePresentation
}

func NewResumeRoot(params ResumeRootParams) (*ResumeRoot, error) {
	title := strings.TrimSpace(params.Title)
	if params.ID.IsZero() || params.UserID.IsZero() || title == "" {
		return nil, fmt.Errorf("%w: invalid resume root identity or title", ErrInvalidResume)
	}
	if err := ValidateResumeTitle(title); err != nil {
		return nil, err
	}
	if err := validateSections(params.Presentation.Sections); err != nil {
		return nil, err
	}
	return &ResumeRoot{
		id: params.ID, userID: params.UserID, title: title,
		presentation: cloneResumePresentation(params.Presentation),
	}, nil
}

func (r *ResumeRoot) ID() valueobjects.ID     { return r.id }
func (r *ResumeRoot) UserID() useridentity.ID { return r.userID }
func (r *ResumeRoot) Title() string           { return r.title }
func (r *ResumeRoot) Presentation() ResumePresentation {
	return cloneResumePresentation(r.presentation)
}

func (r *ResumeRoot) Replace(params ResumeRootParams) error {
	params.ID, params.UserID = r.id, r.userID
	updated, err := NewResumeRoot(params)
	if err != nil {
		return err
	}
	*r = *updated
	return nil
}

func (r *ResumeRoot) Rename(title string) error {
	title = strings.TrimSpace(title)
	if err := ValidateResumeTitle(title); err != nil {
		return err
	}
	r.title = title
	return nil
}

func (r *ResumeRoot) ReplacePresentation(p ResumePresentation) error {
	if err := validateSections(p.Sections); err != nil {
		return err
	}
	r.presentation = cloneResumePresentation(p)
	return nil
}

func cloneResumePresentation(p ResumePresentation) ResumePresentation {
	result := ResumePresentation{
		Theme: cloneTheme(p.Theme), Greeting: p.Greeting,
		Availability: cloneAvailability(p.Availability), ResumeLink: cloneAssetPath(p.ResumeLink),
		Sections: cloneSections(p.Sections),
	}
	return result
}

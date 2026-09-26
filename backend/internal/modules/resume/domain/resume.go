package domain

import (
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
	"fmt"
	"strings"
)

type Resume struct {
	id       valueobjects.ID
	userID   useridentity.ID
	title    string
	language enums.ResumeLanguage
}
type ResumeParams struct {
	ID       valueobjects.ID
	UserID   useridentity.ID
	Title    string
	Language enums.ResumeLanguage
}

func NewResume(params ResumeParams) (*Resume, error) {
	title := strings.TrimSpace(params.Title)
	if params.ID.IsZero() || params.UserID.IsZero() || title == "" {
		return nil, fmt.Errorf("%w: invalid resume identity or title", ErrInvalidResume)
	}
	if !params.Language.IsValid() {
		return nil, fmt.Errorf("%w: invalid resume language", ErrInvalidResume)
	}
	if err := ValidateResumeTitle(title); err != nil {
		return nil, err
	}
	return &Resume{id: params.ID, userID: params.UserID, title: title, language: params.Language}, nil
}
func (r *Resume) ID() valueobjects.ID            { return r.id }
func (r *Resume) UserID() useridentity.ID        { return r.userID }
func (r *Resume) Title() string                  { return r.title }
func (r *Resume) Language() enums.ResumeLanguage { return r.language }
func (r *Resume) Rename(title string) error {
	title = strings.TrimSpace(title)
	if err := ValidateResumeTitle(title); err != nil {
		return err
	}
	r.title = title
	return nil
}
func ValidateResumeTitle(title string) error {
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("%w: resume title is empty", ErrInvalidResume)
	}
	return nil
}

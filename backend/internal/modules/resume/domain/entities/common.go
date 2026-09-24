package entities

import (
	"fmt"
	"strings"

	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

type TechnologyParams struct {
	Name  string
	Grade *enums.TechnologyGrade
}

type Technology struct {
	name  string
	grade *enums.TechnologyGrade
}

func NewTechnology(params TechnologyParams) (Technology, error) {
	name := strings.TrimSpace(params.Name)
	if name == "" {
		return Technology{}, fmt.Errorf("%w: technology name is empty", valueobjects.ErrInvalidValue)
	}
	if params.Grade != nil && !params.Grade.IsValid() {
		return Technology{}, fmt.Errorf("%w: invalid technology grade", valueobjects.ErrInvalidValue)
	}
	return Technology{name: name, grade: cloneTechnologyGrade(params.Grade)}, nil
}

func (t Technology) Name() string                  { return t.name }
func (t Technology) Grade() *enums.TechnologyGrade { return cloneTechnologyGrade(t.grade) }

func cloneTechnologyGrade(value *enums.TechnologyGrade) *enums.TechnologyGrade {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

type ResumeLanguageParams struct {
	Language enums.LanguageType
	Grade    *enums.LanguageGrade
}

type ResumeLanguage struct {
	language enums.LanguageType
	grade    *enums.LanguageGrade
}

func NewResumeLanguage(params ResumeLanguageParams) (ResumeLanguage, error) {
	if !params.Language.IsValid() {
		return ResumeLanguage{}, fmt.Errorf("%w: invalid language", valueobjects.ErrInvalidValue)
	}
	if params.Grade != nil && !params.Grade.IsValid() {
		return ResumeLanguage{}, fmt.Errorf("%w: invalid language grade", valueobjects.ErrInvalidValue)
	}
	return ResumeLanguage{language: params.Language, grade: cloneLanguageGrade(params.Grade)}, nil
}

func (l ResumeLanguage) Language() enums.LanguageType { return l.language }
func (l ResumeLanguage) Grade() *enums.LanguageGrade  { return cloneLanguageGrade(l.grade) }

func cloneLanguageGrade(value *enums.LanguageGrade) *enums.LanguageGrade {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

type ContributorParams struct {
	Name     string
	Surname  string
	Contacts Contacts
	Role     string
}

type Contributor struct {
	name     string
	surname  string
	contacts Contacts
	role     string
}

func NewContributor(params ContributorParams) (Contributor, error) {
	name := strings.TrimSpace(params.Name)
	surname := strings.TrimSpace(params.Surname)
	role := strings.TrimSpace(params.Role)
	if name == "" || surname == "" || role == "" {
		return Contributor{}, fmt.Errorf("%w: contributor fields are required", valueobjects.ErrInvalidValue)
	}
	return Contributor{name: name, surname: surname, contacts: params.Contacts, role: role}, nil
}

func (c Contributor) Name() string       { return c.name }
func (c Contributor) Surname() string    { return c.surname }
func (c Contributor) Role() string       { return c.role }
func (c Contributor) Contacts() Contacts { return c.contacts }

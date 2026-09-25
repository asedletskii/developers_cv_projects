package domain

import (
	"fmt"
	"strings"

	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

type ExperienceTechnology struct {
	name  string
	grade *enums.TechnologyGrade
}

type ExperienceTechnologyParams struct {
	Name  string
	Grade *enums.TechnologyGrade
}

func NewExperienceTechnology(params ExperienceTechnologyParams) (ExperienceTechnology, error) {
	name := strings.TrimSpace(params.Name)
	if name == "" {
		return ExperienceTechnology{}, fmt.Errorf("%w: experience technology name is empty", ErrInvalidEntity)
	}
	if params.Grade != nil && !params.Grade.IsValid() {
		return ExperienceTechnology{}, fmt.Errorf("%w: invalid experience technology grade", ErrInvalidEntity)
	}
	var grade *enums.TechnologyGrade
	if params.Grade != nil {
		value := *params.Grade
		grade = &value
	}
	return ExperienceTechnology{name: name, grade: grade}, nil
}

func (e *Experience) Replace(params ExperienceParams) error {
	params.ID, params.ResumeID = e.id, e.resumeID
	updated, err := NewExperience(params)
	if err != nil {
		return err
	}
	*e = *updated
	return nil
}

func (t ExperienceTechnology) Name() string { return t.name }

func (t ExperienceTechnology) Grade() *enums.TechnologyGrade {
	if t.grade == nil {
		return nil
	}
	value := *t.grade
	return &value
}

type ExperienceParams struct {
	ID           valueobjects.ID
	ResumeID     valueobjects.ID
	Company      string
	Position     string
	Location     *string
	StartDate    valueobjects.YearMonth
	EndDate      *valueobjects.YearMonth
	Description  *string
	Achievements []string
	Technologies []ExperienceTechnologyParams
}

type Experience struct {
	id           valueobjects.ID
	resumeID     valueobjects.ID
	company      string
	position     string
	location     *string
	startDate    valueobjects.YearMonth
	endDate      *valueobjects.YearMonth
	description  *string
	achievements []string
	technologies []ExperienceTechnology
}

func NewExperience(params ExperienceParams) (*Experience, error) {
	company := strings.TrimSpace(params.Company)
	position := strings.TrimSpace(params.Position)
	if params.ID.IsZero() || params.ResumeID.IsZero() || company == "" || position == "" || params.StartDate.IsZero() {
		return nil, fmt.Errorf("%w: invalid experience identity or required fields", ErrInvalidEntity)
	}
	if params.EndDate != nil && params.EndDate.Before(params.StartDate) {
		return nil, fmt.Errorf("%w: experience end date precedes start date", ErrInvalidEntity)
	}

	achievements := make([]string, len(params.Achievements))
	for i, achievement := range params.Achievements {
		achievements[i] = strings.TrimSpace(achievement)
		if achievements[i] == "" {
			return nil, fmt.Errorf("%w: experience achievement is empty", ErrInvalidEntity)
		}
	}

	technologies := make([]ExperienceTechnology, 0, len(params.Technologies))
	for _, technology := range params.Technologies {
		item, err := NewExperienceTechnology(technology)
		if err != nil {
			return nil, err
		}
		technologies = append(technologies, item)
	}

	return &Experience{
		id:           params.ID,
		resumeID:     params.ResumeID,
		company:      company,
		position:     position,
		location:     cloneString(params.Location),
		startDate:    params.StartDate,
		endDate:      cloneYearMonth(params.EndDate),
		description:  cloneString(params.Description),
		achievements: cloneStrings(achievements),
		technologies: technologies,
	}, nil
}

func cloneYearMonth(value *valueobjects.YearMonth) *valueobjects.YearMonth {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func cloneStrings(values []string) []string {
	if values == nil {
		return nil
	}
	return append([]string{}, values...)
}

func (e *Experience) ID() valueobjects.ID               { return e.id }
func (e *Experience) ResumeID() valueobjects.ID         { return e.resumeID }
func (e *Experience) Company() string                   { return e.company }
func (e *Experience) Position() string                  { return e.position }
func (e *Experience) Location() *string                 { return cloneString(e.location) }
func (e *Experience) StartDate() valueobjects.YearMonth { return e.startDate }
func (e *Experience) EndDate() *valueobjects.YearMonth  { return cloneYearMonth(e.endDate) }
func (e *Experience) Description() *string              { return cloneString(e.description) }
func (e *Experience) Achievements() []string            { return cloneStrings(e.achievements) }
func (e *Experience) Technologies() []ExperienceTechnology {
	return append([]ExperienceTechnology(nil), e.technologies...)
}

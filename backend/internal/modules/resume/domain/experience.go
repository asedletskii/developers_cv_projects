package domain

import (
	"fmt"
	"strings"

	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

type ExperienceStackParams struct {
	Name  string
	Grade *enums.TechnologyGrade
}

type ExperienceParams struct {
	ID, ResumeID      valueobjects.ID
	Company, Position string
	Location          *string
	StartDate         valueobjects.YearMonth
	EndDate           *valueobjects.YearMonth
	Description       *string
	Achievements      []string
	Stack             []ExperienceStackParams
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
	stack        []entities.Technology
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
	stack := make([]entities.Technology, 0, len(params.Stack))
	for _, item := range params.Stack {
		technology, err := entities.NewTechnology(entities.TechnologyParams{Name: item.Name, Grade: item.Grade})
		if err != nil {
			return nil, err
		}
		stack = append(stack, technology)
	}
	achievements := make([]string, len(params.Achievements))
	for i, achievement := range params.Achievements {
		value := strings.TrimSpace(achievement)
		if value == "" {
			return nil, fmt.Errorf("%w: experience achievement is empty", ErrInvalidEntity)
		}
		achievements[i] = value
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
		achievements: achievements,
		stack:        stack,
	}, nil
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

func (e *Experience) ID() valueobjects.ID               { return e.id }
func (e *Experience) ResumeID() valueobjects.ID         { return e.resumeID }
func (e *Experience) Company() string                   { return e.company }
func (e *Experience) Position() string                  { return e.position }
func (e *Experience) Location() *string                 { return cloneString(e.location) }
func (e *Experience) StartDate() valueobjects.YearMonth { return e.startDate }
func (e *Experience) EndDate() *valueobjects.YearMonth  { return cloneYearMonth(e.endDate) }
func (e *Experience) Description() *string              { return cloneString(e.description) }
func (e *Experience) Achievements() []string            { return append([]string(nil), e.achievements...) }
func (e *Experience) Stack() []entities.Technology {
	return append([]entities.Technology(nil), e.stack...)
}

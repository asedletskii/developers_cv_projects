package entities

import (
	"fmt"
	"strings"

	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

type ExperienceParams struct {
	ID           valueobjects.ID
	Company      string
	Position     string
	Location     *string
	StartDate    valueobjects.YearMonth
	EndDate      *valueobjects.YearMonth
	Description  *string
	Achievements []string
	Stack        []Technology
}

type Experience struct {
	id           valueobjects.ID
	company      string
	position     string
	location     *string
	startDate    valueobjects.YearMonth
	endDate      *valueobjects.YearMonth
	description  *string
	achievements []string
	stack        []Technology
}

func NewExperience(params ExperienceParams) (Experience, error) {
	company := strings.TrimSpace(params.Company)
	position := strings.TrimSpace(params.Position)
	if params.ID.IsZero() || company == "" || position == "" || params.StartDate.IsZero() {
		return Experience{}, fmt.Errorf(
			"%w: invalid experience identity or required fields", valueobjects.ErrInvalidValue,
		)
	}
	if params.EndDate != nil && params.EndDate.Before(params.StartDate) {
		return Experience{}, fmt.Errorf(
			"%w: experience end date precedes start date", valueobjects.ErrInvalidValue,
		)
	}
	return Experience{
		id:           params.ID,
		company:      company,
		position:     position,
		location:     cloneString(params.Location),
		startDate:    params.StartDate,
		endDate:      cloneYearMonth(params.EndDate),
		description:  cloneString(params.Description),
		achievements: cloneStrings(params.Achievements),
		stack:        cloneTechnologies(params.Stack),
	}, nil
}

func (e Experience) ID() valueobjects.ID               { return e.id }
func (e Experience) Company() string                   { return e.company }
func (e Experience) Position() string                  { return e.position }
func (e Experience) Location() *string                 { return cloneString(e.location) }
func (e Experience) StartDate() valueobjects.YearMonth { return e.startDate }
func (e Experience) EndDate() *valueobjects.YearMonth  { return cloneYearMonth(e.endDate) }
func (e Experience) Description() *string              { return cloneString(e.description) }
func (e Experience) Achievements() []string            { return cloneStrings(e.achievements) }
func (e Experience) Stack() []Technology               { return cloneTechnologies(e.stack) }

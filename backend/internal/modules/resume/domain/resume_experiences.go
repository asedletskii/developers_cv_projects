package domain

import (
	"fmt"
	"strings"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ExperienceStackParams struct {
	Name  string
	Grade *TechnologyGrade
}

type ExperienceParams struct {
	ID, ResumeID      sharedvalueobjects.ID
	Company, Position string
	Location          *string
	StartDate         YearMonth
	EndDate           *YearMonth
	Description       *string
	Achievements      []string
	Tasks             []string
	Stack             []ExperienceStackParams
}

type Experience struct {
	id           sharedvalueobjects.ID
	resumeID     sharedvalueobjects.ID
	company      string
	position     string
	location     *string
	startDate    YearMonth
	endDate      *YearMonth
	description  *string
	achievements []string
	tasks        []string
	stack        []Technology
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
	stack := make([]Technology, 0, len(params.Stack))
	for _, item := range params.Stack {
		technology, err := NewTechnology(TechnologyParams(item))
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
	tasks := make([]string, len(params.Tasks))
	for i, task := range params.Tasks {
		value := strings.TrimSpace(task)
		if value == "" {
			return nil, fmt.Errorf("%w: experience task is empty", ErrInvalidEntity)
		}
		tasks[i] = value
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
		tasks:        tasks,
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

func (e *Experience) ID() sharedvalueobjects.ID       { return e.id }
func (e *Experience) ResumeID() sharedvalueobjects.ID { return e.resumeID }
func (e *Experience) Company() string                 { return e.company }
func (e *Experience) Position() string                { return e.position }
func (e *Experience) Location() *string               { return cloneString(e.location) }
func (e *Experience) StartDate() YearMonth            { return e.startDate }
func (e *Experience) EndDate() *YearMonth             { return cloneYearMonth(e.endDate) }
func (e *Experience) Description() *string            { return cloneString(e.description) }
func (e *Experience) Achievements() []string          { return append([]string(nil), e.achievements...) }
func (e *Experience) Tasks() []string                 { return append([]string(nil), e.tasks...) }
func (e *Experience) Stack() []Technology {
	return append([]Technology(nil), e.stack...)
}

func (r *Resume) ReplaceExperiences(params []ExperienceParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	existing := make(map[sharedvalueobjects.ID]*Experience, len(r.experiences))
	for _, item := range r.experiences {
		if item == nil {
			return fmt.Errorf("%w: nil experience", ErrResumeInvariant)
		}
		existing[item.ID()] = item
	}
	seen := make(map[sharedvalueobjects.ID]struct{}, len(params))
	replaced := make([]*Experience, 0, len(params))
	for _, item := range params {
		if item.ID.IsZero() {
			item.ID = sharedvalueobjects.NewID()
		} else if _, ok := existing[item.ID]; !ok {
			return fmt.Errorf("%w: experience does not belong to resume", ErrResumeInvariant)
		}
		if _, ok := seen[item.ID]; ok {
			return fmt.Errorf("%w: duplicate experience id", ErrInvalidEntity)
		}
		seen[item.ID] = struct{}{}
		item.ResumeID = r.id
		experience, err := NewExperience(item)
		if err != nil {
			return err
		}
		replaced = append(replaced, experience)
	}
	r.experiences = replaced
	return nil
}

func (r *Resume) Experiences() []*Experience {
	if r == nil {
		return nil
	}
	return append([]*Experience(nil), r.experiences...)
}

package domain

import (
	"fmt"
	"strings"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ExperienceStackParams struct {
	Name string
}

type ExperienceParams struct {
	ID                sharedvalueobjects.ID
	Company, Position string
	Location          *string
	StartDate         YearMonth
	EndDate           *YearMonth
	Description       *string
	Achievements      []string
	Tasks             []string
	Stack             []ExperienceStackParams
}

type ExperienceSnapshot struct {
	ID                sharedvalueobjects.ID
	ResumeID          sharedvalueobjects.ID
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

func newExperience(id, resumeID sharedvalueobjects.ID, params ExperienceParams) (*Experience, error) {
	experience := &Experience{
		id: id, resumeID: resumeID, company: strings.TrimSpace(params.Company),
		position: strings.TrimSpace(params.Position), location: normalizeOptionalString(params.Location), startDate: params.StartDate,
		endDate: cloneYearMonth(params.EndDate), description: normalizeOptionalString(params.Description),
		achievements: make([]string, len(params.Achievements)), tasks: make([]string, len(params.Tasks)),
		stack: make([]Technology, 0, len(params.Stack)),
	}
	for i, value := range params.Achievements {
		experience.achievements[i] = strings.TrimSpace(value)
	}
	for i, value := range params.Tasks {
		experience.tasks[i] = strings.TrimSpace(value)
	}
	for _, item := range params.Stack {
		experience.stack = append(experience.stack, RestoreTechnology(TechnologyParams{Name: strings.TrimSpace(item.Name)}))
	}
	if err := experience.validate(); err != nil {
		return nil, err
	}
	return experience, nil
}

func RestoreExperience(snapshot ExperienceSnapshot) *Experience {
	achievements := append([]string(nil), snapshot.Achievements...)
	tasks := append([]string(nil), snapshot.Tasks...)
	stack := make([]Technology, 0, len(snapshot.Stack))
	for _, item := range snapshot.Stack {
		stack = append(stack, Technology{name: item.Name})
	}
	return &Experience{
		id: snapshot.ID, resumeID: snapshot.ResumeID, company: snapshot.Company, position: snapshot.Position,
		location: copyString(snapshot.Location), startDate: snapshot.StartDate, endDate: cloneYearMonth(snapshot.EndDate),
		description: copyString(snapshot.Description), achievements: achievements, tasks: tasks,
		stack: stack,
	}
}

func (e *Experience) validate() error {
	if e == nil || e.id.IsZero() || e.resumeID.IsZero() || e.company == "" || e.position == "" || e.startDate.IsZero() {
		return fmt.Errorf("%w: invalid experience identity or required fields", ErrInvalidEntity)
	}
	if !textLengthWithin(e.company, MaxCompanyLength) || !textLengthWithin(e.position, MaxProfilePositionLength) ||
		!optionalTextLengthWithin(e.location, MaxLocationLength) || !optionalTextLengthWithin(e.description, MaxDescriptionLength) ||
		len(e.tasks) > MaxExperienceTasks || len(e.achievements) > MaxExperienceAchievements || len(e.stack) > MaxExperienceStack {
		return fmt.Errorf("%w: experience exceeds configured limits", ErrInvalidEntity)
	}
	if e.endDate != nil {
		if err := e.endDate.validate(); err != nil || e.endDate.before(e.startDate) {
			return fmt.Errorf("%w: experience end date precedes start date", ErrInvalidEntity)
		}
	}
	seenTechnologies := make(map[string]struct{}, len(e.stack))
	for _, technology := range e.stack {
		if err := technology.validate(); err != nil {
			return fmt.Errorf("%w: invalid experience technology", ErrInvalidEntity)
		}
		if _, exists := seenTechnologies[technology.Name()]; exists {
			return fmt.Errorf("%w: duplicate experience technology", ErrDuplicateEntity)
		}
		seenTechnologies[technology.Name()] = struct{}{}
	}
	for _, achievement := range e.achievements {
		if achievement == "" || !textLengthWithin(achievement, MaxListItemTextLength) {
			return fmt.Errorf("%w: experience achievement is empty or too long", ErrInvalidEntity)
		}
	}
	for _, task := range e.tasks {
		if task == "" || !textLengthWithin(task, MaxListItemTextLength) {
			return fmt.Errorf("%w: experience task is empty or too long", ErrInvalidEntity)
		}
	}
	return nil
}

func (e *Experience) ID() sharedvalueobjects.ID       { return e.id }
func (e *Experience) ResumeID() sharedvalueobjects.ID { return e.resumeID }
func (e *Experience) Company() string                 { return e.company }
func (e *Experience) Position() string                { return e.position }
func (e *Experience) Location() *string               { return copyString(e.location) }
func (e *Experience) StartDate() YearMonth            { return e.startDate }
func (e *Experience) EndDate() *YearMonth             { return cloneYearMonth(e.endDate) }
func (e *Experience) Description() *string            { return copyString(e.description) }
func (e *Experience) Achievements() []string          { return append([]string(nil), e.achievements...) }
func (e *Experience) Tasks() []string                 { return append([]string(nil), e.tasks...) }
func (e *Experience) Stack() []Technology {
	return append([]Technology(nil), e.stack...)
}

func (r *Resume) ReplaceExperiences(params []ExperienceParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	if len(params) > MaxExperiencesPerResume {
		return fmt.Errorf("%w: resume cannot contain more than %d experiences", ErrInvalidEntity, MaxExperiencesPerResume)
	}
	existing := make(map[sharedvalueobjects.ID]struct{}, len(r.experiences))
	for _, current := range r.experiences {
		if current == nil {
			return fmt.Errorf("%w: nil experience", ErrResumeInvariant)
		}
		existing[current.ID()] = struct{}{}
	}
	seen := make(map[sharedvalueobjects.ID]struct{}, len(params))
	replaced := make([]*Experience, 0, len(params))
	for _, item := range params {
		resolvedID, err := resolveCollectionItemID(item.ID, existing, "experience")
		if err != nil {
			return err
		}
		if _, ok := seen[resolvedID]; ok {
			return fmt.Errorf("%w: duplicate experience id", ErrInvalidEntity)
		}
		seen[resolvedID] = struct{}{}
		experience, err := newExperience(resolvedID, r.id, item)
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
	return cloneExperiences(r.experiences)
}

func cloneExperiences(values []*Experience) []*Experience {
	if values == nil {
		return nil
	}
	result := make([]*Experience, len(values))
	for i, value := range values {
		if value == nil {
			continue
		}
		copy := *value
		copy.location = copyString(value.location)
		copy.endDate = cloneYearMonth(value.endDate)
		copy.description = copyString(value.description)
		copy.achievements = append([]string(nil), value.achievements...)
		copy.tasks = append([]string(nil), value.tasks...)
		copy.stack = append([]Technology(nil), value.stack...)
		result[i] = &copy
	}
	return result
}

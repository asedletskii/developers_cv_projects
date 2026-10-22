package domain

import (
	"fmt"
	"strings"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type EducationLevel int

const (
	EducationLevelSecondary    EducationLevel = 0
	EducationLevelBachelor     EducationLevel = 1
	EducationLevelSpecialist   EducationLevel = 2
	EducationLevelMaster       EducationLevel = 3
	EducationLevelPostgraduate EducationLevel = 4
	EducationLevelCandidate    EducationLevel = 5
	EducationLevelDoctor       EducationLevel = 6
)

func (v EducationLevel) isValid() bool {
	return v >= EducationLevelSecondary && v <= EducationLevelDoctor
}

type EducationParams struct {
	ID          sharedvalueobjects.ID
	Location    string
	University  string
	Grade       EducationLevel
	ProgramName string
	ProfileName *string
	Faculty     string
	Course      *int
	EndDate     *Year
}

type EducationSnapshot struct {
	ID          sharedvalueobjects.ID
	ResumeID    sharedvalueobjects.ID
	Location    string
	University  string
	Grade       EducationLevel
	ProgramName string
	ProfileName *string
	Faculty     string
	Course      *int
	EndDate     *Year
}

type Education struct {
	id          sharedvalueobjects.ID
	resumeID    sharedvalueobjects.ID
	location    string
	university  string
	grade       EducationLevel
	programName string
	profileName *string
	faculty     string
	course      *int
	endDate     *Year
}

func newEducation(id, resumeID sharedvalueobjects.ID, params EducationParams) (*Education, error) {
	education := &Education{
		id: id, resumeID: resumeID, location: strings.TrimSpace(params.Location),
		university: strings.TrimSpace(params.University), grade: params.Grade,
		programName: strings.TrimSpace(params.ProgramName), profileName: normalizeOptionalString(params.ProfileName),
		faculty: strings.TrimSpace(params.Faculty), course: cloneInt(params.Course), endDate: cloneYear(params.EndDate),
	}
	if err := education.validate(); err != nil {
		return nil, err
	}
	return education, nil
}

func RestoreEducation(snapshot EducationSnapshot) *Education {
	return &Education{
		id: snapshot.ID, resumeID: snapshot.ResumeID, location: snapshot.Location, university: snapshot.University,
		grade: snapshot.Grade, programName: snapshot.ProgramName, profileName: copyString(snapshot.ProfileName),
		faculty: snapshot.Faculty, course: cloneInt(snapshot.Course), endDate: cloneYear(snapshot.EndDate),
	}
}

func (e *Education) validate() error {
	if e == nil || e.id.IsZero() || e.resumeID.IsZero() || e.location == "" || e.university == "" ||
		e.programName == "" || e.faculty == "" || !e.grade.isValid() {
		return fmt.Errorf("%w: invalid education fields", ErrInvalidEntity)
	}
	if !textLengthWithin(e.location, MaxLocationLength) || !textLengthWithin(e.university, MaxShortTextLength) ||
		!textLengthWithin(e.programName, MaxShortTextLength) || !textLengthWithin(e.faculty, MaxShortTextLength) ||
		!optionalTextLengthWithin(e.profileName, MaxShortTextLength) {
		return fmt.Errorf("%w: education text exceeds configured limits", ErrInvalidEntity)
	}
	if e.course != nil && *e.course < 1 {
		return fmt.Errorf("%w: invalid education course", ErrInvalidEntity)
	}
	if e.endDate != nil {
		if err := e.endDate.validate(); err != nil {
			return fmt.Errorf("%w: invalid education end date", ErrInvalidEntity)
		}
	}
	return nil
}

func (e *Education) ID() sharedvalueobjects.ID       { return e.id }
func (e *Education) ResumeID() sharedvalueobjects.ID { return e.resumeID }
func (e *Education) Location() string                { return e.location }
func (e *Education) University() string              { return e.university }
func (e *Education) Grade() EducationLevel           { return e.grade }
func (e *Education) ProgramName() string             { return e.programName }
func (e *Education) ProfileName() *string            { return copyString(e.profileName) }
func (e *Education) Faculty() string                 { return e.faculty }
func (e *Education) Course() *int                    { return cloneInt(e.course) }
func (e *Education) EndDate() *Year                  { return cloneYear(e.endDate) }

func (r *Resume) ReplaceEducations(params []EducationParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	if len(params) > MaxEducationsPerResume {
		return fmt.Errorf("%w: resume cannot contain more than %d educations", ErrInvalidEntity, MaxEducationsPerResume)
	}
	existing := make(map[sharedvalueobjects.ID]struct{}, len(r.educations))
	for _, current := range r.educations {
		if current == nil {
			return fmt.Errorf("%w: nil education", ErrResumeInvariant)
		}
		existing[current.ID()] = struct{}{}
	}
	seen := make(map[sharedvalueobjects.ID]struct{}, len(params))
	replaced := make([]*Education, 0, len(params))
	for _, item := range params {
		resolvedID, err := resolveCollectionItemID(item.ID, existing, "education")
		if err != nil {
			return err
		}
		if _, ok := seen[resolvedID]; ok {
			return fmt.Errorf("%w: duplicate education id", ErrInvalidEntity)
		}
		seen[resolvedID] = struct{}{}
		education, err := newEducation(resolvedID, r.id, item)
		if err != nil {
			return err
		}
		replaced = append(replaced, education)
	}
	r.educations = replaced
	return nil
}

func (r *Resume) Educations() []*Education {
	if r == nil {
		return nil
	}
	return cloneEducations(r.educations)
}

func cloneEducations(values []*Education) []*Education {
	if values == nil {
		return nil
	}
	result := make([]*Education, len(values))
	for i, value := range values {
		if value == nil {
			continue
		}
		copy := *value
		copy.profileName = copyString(value.profileName)
		copy.course = cloneInt(value.course)
		copy.endDate = cloneYear(value.endDate)
		result[i] = &copy
	}
	return result
}

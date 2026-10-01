package domain

import (
	"fmt"
	"strings"

	identity "developers-cv/backend/internal/shared/domain/valueobjects"
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

func (v EducationLevel) IsValid() bool {
	return v >= EducationLevelSecondary && v <= EducationLevelDoctor
}

type EducationParams struct {
	ID          identity.ID
	ResumeID    identity.ID
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
	id          identity.ID
	resumeID    identity.ID
	location    string
	university  string
	grade       EducationLevel
	programName string
	profileName *string
	faculty     string
	course      *int
	endDate     *Year
}

func NewEducation(params EducationParams) (*Education, error) {
	location := strings.TrimSpace(params.Location)
	university := strings.TrimSpace(params.University)
	programName := strings.TrimSpace(params.ProgramName)
	faculty := strings.TrimSpace(params.Faculty)
	if params.ID.IsZero() || params.ResumeID.IsZero() || location == "" || university == "" ||
		programName == "" || faculty == "" || !params.Grade.IsValid() {
		return nil, fmt.Errorf("%w: invalid education fields", ErrInvalidEntity)
	}
	if params.Course != nil && *params.Course < 1 {
		return nil, fmt.Errorf("%w: invalid education course", ErrInvalidEntity)
	}
	return &Education{
		id: params.ID, resumeID: params.ResumeID, location: location, university: university,
		grade: params.Grade, programName: programName, profileName: cloneString(params.ProfileName),
		faculty: faculty, course: cloneInt(params.Course), endDate: cloneYear(params.EndDate),
	}, nil
}

func (e *Education) Replace(params EducationParams) error {
	params.ID, params.ResumeID = e.id, e.resumeID
	updated, err := NewEducation(params)
	if err != nil {
		return err
	}
	*e = *updated
	return nil
}

func (e *Education) ID() identity.ID       { return e.id }
func (e *Education) ResumeID() identity.ID { return e.resumeID }
func (e *Education) Location() string      { return e.location }
func (e *Education) University() string    { return e.university }
func (e *Education) Grade() EducationLevel { return e.grade }
func (e *Education) ProgramName() string   { return e.programName }
func (e *Education) ProfileName() *string  { return cloneString(e.profileName) }
func (e *Education) Faculty() string       { return e.faculty }
func (e *Education) Course() *int          { return cloneInt(e.course) }
func (e *Education) EndDate() *Year        { return cloneYear(e.endDate) }

func (r *Resume) ReplaceEducations(params []EducationParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	existing := make(map[identity.ID]struct{}, len(r.educations))
	for _, item := range r.educations {
		if item == nil {
			return fmt.Errorf("%w: nil education", ErrResumeInvariant)
		}
		existing[item.ID()] = struct{}{}
	}
	seen := make(map[identity.ID]struct{}, len(params))
	replaced := make([]*Education, 0, len(params))
	for _, item := range params {
		if item.ID.IsZero() {
			item.ID = identity.NewID()
		} else if _, ok := existing[item.ID]; !ok {
			return fmt.Errorf("%w: education does not belong to resume", ErrResumeInvariant)
		}
		if _, ok := seen[item.ID]; ok {
			return fmt.Errorf("%w: duplicate education id", ErrInvalidEntity)
		}
		seen[item.ID] = struct{}{}
		item.ResumeID = r.id
		education, err := NewEducation(item)
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
	return append([]*Education(nil), r.educations...)
}

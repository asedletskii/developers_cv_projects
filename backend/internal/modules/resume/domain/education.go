package domain

import (
	"fmt"
	"strings"

	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

type EducationParams struct {
	ID          valueobjects.ID
	ResumeID    valueobjects.ID
	Location    string
	University  string
	Grade       enums.EducationLevel
	ProgramName string
	ProfileName *string
	Faculty     string
	Course      *int
	EndDate     *valueobjects.Year
}

type Education struct {
	id          valueobjects.ID
	resumeID    valueobjects.ID
	location    string
	university  string
	grade       enums.EducationLevel
	programName string
	profileName *string
	faculty     string
	course      *int
	endDate     *valueobjects.Year
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

func (e *Education) ID() valueobjects.ID         { return e.id }
func (e *Education) ResumeID() valueobjects.ID   { return e.resumeID }
func (e *Education) Location() string            { return e.location }
func (e *Education) University() string          { return e.university }
func (e *Education) Grade() enums.EducationLevel { return e.grade }
func (e *Education) ProgramName() string         { return e.programName }
func (e *Education) ProfileName() *string        { return cloneString(e.profileName) }
func (e *Education) Faculty() string             { return e.faculty }
func (e *Education) Course() *int                { return cloneInt(e.course) }
func (e *Education) EndDate() *valueobjects.Year { return cloneYear(e.endDate) }

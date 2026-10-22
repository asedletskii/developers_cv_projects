package domain

import (
	"fmt"
	"strings"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type CourseParams struct {
	ID              sharedvalueobjects.ID
	Name            string
	Organization    string
	EndDate         YearMonth
	Hours           *int
	CertificateLink *ExternalURL
}

type CourseSnapshot struct {
	ID              sharedvalueobjects.ID
	ResumeID        sharedvalueobjects.ID
	Name            string
	Organization    string
	EndDate         YearMonth
	Hours           *int
	CertificateLink *ExternalURL
}

type Course struct {
	id              sharedvalueobjects.ID
	resumeID        sharedvalueobjects.ID
	name            string
	organization    string
	endDate         YearMonth
	hours           *int
	certificateLink *ExternalURL
}

func newCourse(id, resumeID sharedvalueobjects.ID, params CourseParams) (*Course, error) {
	course := &Course{
		id: id, resumeID: resumeID, name: strings.TrimSpace(params.Name),
		organization: strings.TrimSpace(params.Organization), endDate: params.EndDate,
		hours: cloneInt(params.Hours), certificateLink: cloneExternalURL(params.CertificateLink),
	}
	if err := course.validate(); err != nil {
		return nil, err
	}
	return course, nil
}

func RestoreCourse(snapshot CourseSnapshot) *Course {
	return &Course{
		id: snapshot.ID, resumeID: snapshot.ResumeID, name: snapshot.Name, organization: snapshot.Organization,
		endDate: snapshot.EndDate, hours: cloneInt(snapshot.Hours), certificateLink: cloneExternalURL(snapshot.CertificateLink),
	}
}

func (c *Course) validate() error {
	if c == nil || c.id.IsZero() || c.resumeID.IsZero() || c.name == "" || c.organization == "" || c.endDate.IsZero() {
		return fmt.Errorf("%w: invalid course fields", ErrInvalidEntity)
	}
	if !textLengthWithin(c.name, MaxShortTextLength) || !textLengthWithin(c.organization, MaxShortTextLength) {
		return fmt.Errorf("%w: course text exceeds configured limits", ErrInvalidEntity)
	}
	if c.certificateLink != nil {
		if err := c.certificateLink.validate(); err != nil {
			return fmt.Errorf("%w: invalid course certificate URL", ErrInvalidEntity)
		}
	}
	if c.hours != nil && *c.hours <= 0 {
		return fmt.Errorf("%w: course hours must be positive", ErrInvalidEntity)
	}
	return nil
}

func (c *Course) ID() sharedvalueobjects.ID       { return c.id }
func (c *Course) ResumeID() sharedvalueobjects.ID { return c.resumeID }
func (c *Course) Name() string                    { return c.name }
func (c *Course) Organization() string            { return c.organization }
func (c *Course) EndDate() YearMonth              { return c.endDate }
func (c *Course) Hours() *int                     { return cloneInt(c.hours) }
func (c *Course) CertificateLink() *ExternalURL {
	return cloneExternalURL(c.certificateLink)
}

func (r *Resume) ReplaceCourses(params []CourseParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	if len(params) > MaxCoursesPerResume {
		return fmt.Errorf("%w: resume cannot contain more than %d courses", ErrInvalidEntity, MaxCoursesPerResume)
	}
	existing := make(map[sharedvalueobjects.ID]struct{}, len(r.courses))
	for _, current := range r.courses {
		if current == nil {
			return fmt.Errorf("%w: nil course", ErrResumeInvariant)
		}
		existing[current.ID()] = struct{}{}
	}
	seen := make(map[sharedvalueobjects.ID]struct{}, len(params))
	replaced := make([]*Course, 0, len(params))
	for _, item := range params {
		resolvedID, err := resolveCollectionItemID(item.ID, existing, "course")
		if err != nil {
			return err
		}
		if _, ok := seen[resolvedID]; ok {
			return fmt.Errorf("%w: duplicate course id", ErrInvalidEntity)
		}
		seen[resolvedID] = struct{}{}
		course, err := newCourse(resolvedID, r.id, item)
		if err != nil {
			return err
		}
		replaced = append(replaced, course)
	}
	r.courses = replaced
	return nil
}

func (r *Resume) Courses() []*Course {
	if r == nil {
		return nil
	}
	return cloneCourses(r.courses)
}

func cloneCourses(values []*Course) []*Course {
	if values == nil {
		return nil
	}
	result := make([]*Course, len(values))
	for i, value := range values {
		if value == nil {
			continue
		}
		copy := *value
		copy.hours = cloneInt(value.hours)
		copy.certificateLink = cloneExternalURL(value.certificateLink)
		result[i] = &copy
	}
	return result
}

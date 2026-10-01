package domain

import (
	"fmt"
	"strings"

	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type CourseParams struct {
	ID              identity.ID
	ResumeID        identity.ID
	Name            string
	Organization    string
	EndDate         YearMonth
	Hours           *int
	CertificateLink *ExternalURL
}

type Course struct {
	id              identity.ID
	resumeID        identity.ID
	name            string
	organization    string
	endDate         YearMonth
	hours           *int
	certificateLink *ExternalURL
}

func NewCourse(params CourseParams) (*Course, error) {
	name := strings.TrimSpace(params.Name)
	organization := strings.TrimSpace(params.Organization)
	if params.ID.IsZero() || params.ResumeID.IsZero() || name == "" || organization == "" || params.EndDate.IsZero() {
		return nil, fmt.Errorf("%w: invalid course fields", ErrInvalidEntity)
	}
	if params.Hours != nil && *params.Hours <= 0 {
		return nil, fmt.Errorf("%w: course hours must be positive", ErrInvalidEntity)
	}
	return &Course{
		id: params.ID, resumeID: params.ResumeID, name: name, organization: organization,
		endDate: params.EndDate, hours: cloneInt(params.Hours), certificateLink: cloneExternalURL(params.CertificateLink),
	}, nil
}

func (c *Course) Replace(params CourseParams) error {
	params.ID, params.ResumeID = c.id, c.resumeID
	updated, err := NewCourse(params)
	if err != nil {
		return err
	}
	*c = *updated
	return nil
}

func (c *Course) ID() identity.ID       { return c.id }
func (c *Course) ResumeID() identity.ID { return c.resumeID }
func (c *Course) Name() string          { return c.name }
func (c *Course) Organization() string  { return c.organization }
func (c *Course) EndDate() YearMonth    { return c.endDate }
func (c *Course) Hours() *int           { return cloneInt(c.hours) }
func (c *Course) CertificateLink() *ExternalURL {
	return cloneExternalURL(c.certificateLink)
}

func (r *Resume) ReplaceCourses(params []CourseParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	existing := make(map[identity.ID]struct{}, len(r.courses))
	for _, item := range r.courses {
		if item == nil {
			return fmt.Errorf("%w: nil course", ErrResumeInvariant)
		}
		existing[item.ID()] = struct{}{}
	}
	seen := make(map[identity.ID]struct{}, len(params))
	replaced := make([]*Course, 0, len(params))
	for _, item := range params {
		if item.ID.IsZero() {
			item.ID = identity.NewID()
		} else if _, ok := existing[item.ID]; !ok {
			return fmt.Errorf("%w: course does not belong to resume", ErrResumeInvariant)
		}
		if _, ok := seen[item.ID]; ok {
			return fmt.Errorf("%w: duplicate course id", ErrInvalidEntity)
		}
		seen[item.ID] = struct{}{}
		item.ResumeID = r.id
		course, err := NewCourse(item)
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
	return append([]*Course(nil), r.courses...)
}

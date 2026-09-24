package domain

import (
	"fmt"
	"strings"

	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

type CourseParams struct {
	ID              valueobjects.ID
	ResumeID        valueobjects.ID
	Name            string
	Organization    string
	EndDate         valueobjects.YearMonth
	Hours           *int
	CertificateLink *valueobjects.ExternalURL
}

type Course struct {
	id              valueobjects.ID
	resumeID        valueobjects.ID
	name            string
	organization    string
	endDate         valueobjects.YearMonth
	hours           *int
	certificateLink *valueobjects.ExternalURL
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

func (c *Course) ID() valueobjects.ID             { return c.id }
func (c *Course) ResumeID() valueobjects.ID       { return c.resumeID }
func (c *Course) Name() string                    { return c.name }
func (c *Course) Organization() string            { return c.organization }
func (c *Course) EndDate() valueobjects.YearMonth { return c.endDate }
func (c *Course) Hours() *int                     { return cloneInt(c.hours) }
func (c *Course) CertificateLink() *valueobjects.ExternalURL {
	return cloneExternalURL(c.certificateLink)
}

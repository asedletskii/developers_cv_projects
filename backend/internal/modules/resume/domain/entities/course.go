package entities

import (
	"fmt"
	"strings"

	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

type CourseParams struct {
	ID              valueobjects.ID
	Name            string
	Organization    string
	EndDate         valueobjects.YearMonth
	Hours           *int
	CertificateLink *valueobjects.ExternalURL
}

type Course struct {
	id              valueobjects.ID
	name            string
	organization    string
	endDate         valueobjects.YearMonth
	hours           *int
	certificateLink *valueobjects.ExternalURL
}

func NewCourse(params CourseParams) (Course, error) {
	name := strings.TrimSpace(params.Name)
	organization := strings.TrimSpace(params.Organization)
	if params.ID.IsZero() || name == "" || organization == "" || params.EndDate.IsZero() {
		return Course{}, fmt.Errorf("%w: invalid course fields", valueobjects.ErrInvalidValue)
	}
	if params.Hours != nil && *params.Hours <= 0 {
		return Course{}, fmt.Errorf("%w: course hours must be positive", valueobjects.ErrInvalidValue)
	}
	return Course{
		id:              params.ID,
		name:            name,
		organization:    organization,
		endDate:         params.EndDate,
		hours:           cloneInt(params.Hours),
		certificateLink: cloneExternalURL(params.CertificateLink),
	}, nil
}

func (c Course) ID() valueobjects.ID             { return c.id }
func (c Course) Name() string                    { return c.name }
func (c Course) Organization() string            { return c.organization }
func (c Course) EndDate() valueobjects.YearMonth { return c.endDate }
func (c Course) Hours() *int                     { return cloneInt(c.hours) }
func (c Course) CertificateLink() *valueobjects.ExternalURL {
	return cloneExternalURL(c.certificateLink)
}

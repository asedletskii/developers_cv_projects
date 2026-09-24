package domain_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

func TestCourseValidation(t *testing.T) {
	endDate, err := valueobjects.ParseYearMonth("2024-06")
	require.NoError(t, err)
	certificate, err := valueobjects.NewExternalURL("https://example.com/certificate")
	require.NoError(t, err)
	hours := 40
	item, err := domain.NewCourse(domain.CourseParams{ID: valueobjects.NewID(), ResumeID: valueobjects.NewID(), Name: "Go", Organization: "Academy", EndDate: endDate, Hours: &hours, CertificateLink: &certificate})
	require.NoError(t, err)
	require.Equal(t, "2024-06", item.EndDate().String())
	require.Equal(t, certificate.String(), item.CertificateLink().String())
}

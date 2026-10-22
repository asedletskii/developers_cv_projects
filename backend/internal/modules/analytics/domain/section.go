package domain

import "fmt"

type AnalyticsSectionID string

const (
	AnalyticsSectionAbout      AnalyticsSectionID = "about"
	AnalyticsSectionTechStack  AnalyticsSectionID = "techStack"
	AnalyticsSectionExperience AnalyticsSectionID = "experience"
	AnalyticsSectionSkills     AnalyticsSectionID = "skills"
	AnalyticsSectionEducation  AnalyticsSectionID = "education"
	AnalyticsSectionCourses    AnalyticsSectionID = "courses"
	AnalyticsSectionProjects   AnalyticsSectionID = "projects"
	AnalyticsSectionActivities AnalyticsSectionID = "activities"
	AnalyticsSectionContacts   AnalyticsSectionID = "contacts"
)

func ParseAnalyticsSectionID(value string) (AnalyticsSectionID, error) {
	section := AnalyticsSectionID(value)
	switch section {
	case AnalyticsSectionAbout, AnalyticsSectionTechStack, AnalyticsSectionExperience,
		AnalyticsSectionSkills, AnalyticsSectionEducation, AnalyticsSectionCourses,
		AnalyticsSectionProjects, AnalyticsSectionActivities, AnalyticsSectionContacts:
		return section, nil
	default:
		return "", fmt.Errorf("%w: invalid section %q", ErrInvalidAnalyticsSection, value)
	}
}

func IsValidAnalyticsContactKey(value string) bool {
	switch value {
	case "phone", "telegram", "email":
		return true
	default:
		return false
	}
}

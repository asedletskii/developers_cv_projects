package queries

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
)

func TestMapPublicResumeCopiesFullProviderShape(t *testing.T) {
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
	patronymic := "Ivanovna"
	from := "2026-03-01"
	note := "Remote only"
	phone := "+33123456789"
	telegram := "ada"
	email := "ada@example.com"
	contactsTitle := "Contacts"
	contactsSubtitle := "Public channels"
	grade := "C1"
	location := "Paris"
	endDate := "2025-12-31"
	description := "Built public APIs"
	profileName := "Applied Mathematics"
	course := 4
	hours := 40
	certificateLink := "https://example.com/cert"
	activityDate := "2025-05-01"
	activityLink := "https://example.com/activity"
	projectDescription := "Developer CV"
	projectLink := "https://example.com/project"
	contributorPhone := "+33111222333"
	contributorEmail := "grace@example.com"

	view := resumecontracts.PublicResumeView{
		Title:     "Backend Engineer",
		Language:  "en",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Overview: resumecontracts.PublicOverviewView{
			Theme: &resumecontracts.PublicThemeSettings{Primary: "#112233", Secondary: "#445566"},
			Greeting: resumecontracts.PublicGreeting{
				First: "Hello", Second: "I am", SecondHighlight: "Ada", Third: "Lovelace",
			},
			Availability: &resumecontracts.PublicAvailability{Status: "open", From: &from, Note: &note},
			Sections: []resumecontracts.PublicSectionSettings{
				{ID: "about", Enabled: true, Collapsed: false},
				{ID: "skills", Enabled: false, Collapsed: true},
			},
		},
		Profile: &resumecontracts.PublicProfileView{
			ID: "profile-1", Name: "Ada", Surname: "Lovelace", Patronymic: &patronymic,
			Position: "Backend Engineer", PhotoURL: "/photo.png", DateOfBirth: "1990-01-01",
			Location: "Paris", Relocation: true, WorkFormat: []string{"remote", "hybrid"},
		},
		Contacts: &resumecontracts.PublicContactsView{
			ID: "contacts-1", Phone: &phone, Telegram: &telegram, Email: &email,
			HiddenContacts: []string{"email"},
			ContactsText:   &resumecontracts.PublicContactsTextView{Title: &contactsTitle, Subtitle: &contactsSubtitle},
		},
		Skills: &resumecontracts.PublicSkillsView{
			ID: "skills-1", HardSkills: []string{"Go", "PostgreSQL"}, SoftSkills: []string{"Communication"},
		},
		Languages: []resumecontracts.PublicLanguageView{
			{ID: "language-1", Language: "English", Grade: &grade, Position: 1},
			{ID: "language-2", Language: "French", Position: 2},
		},
		Technologies: []resumecontracts.PublicTechnologyView{
			{ID: "technology-1", Name: "Go", Position: 1},
			{ID: "technology-2", Name: "gRPC", Position: 2},
		},
		Experiences: []resumecontracts.PublicExperienceView{{
			ID: "experience-1", Company: "Example", Position: "Backend Engineer", Location: &location,
			StartDate: "2022-01-01", EndDate: &endDate, Description: &description,
			Achievements: []string{"Shipped API"}, Tasks: []string{"Design services"},
			Stack: []resumecontracts.PublicExperienceStackItemView{{Name: "Go"}, {Name: "gRPC"}},
		}},
		Educations: []resumecontracts.PublicEducationView{{
			ID: "education-1", Location: "Paris", University: "Sorbonne", Grade: 5, ProgramName: "CS",
			ProfileName: &profileName, Faculty: "Science", Course: &course, EndDate: &endDate,
		}},
		Courses: []resumecontracts.PublicCourseView{{
			ID: "course-1", Name: "Go Deep Dive", Organization: "Example School", EndDate: "2025-06-01",
			Hours: &hours, CertificateLink: &certificateLink,
		}},
		Activities: []resumecontracts.PublicActivityView{{
			ID: "activity-1", Type: "open_source", Title: "Contribution", Description: &description,
			Date: &activityDate, Link: &activityLink,
		}},
		Projects: []resumecontracts.PublicProjectView{{
			ID: "project-1", Name: "CV", Description: &projectDescription,
			Screenshots: []string{"https://example.com/shot.png"}, Link: &projectLink,
			Stack: []resumecontracts.PublicProjectTechnologyView{{Name: "Go"}},
			Contributors: []resumecontracts.PublicProjectContributorView{{
				Name: "Grace", Surname: "Hopper", Role: "Advisor",
				Contacts: resumecontracts.PublicProjectContactsView{Phone: &contributorPhone, Email: &contributorEmail},
			}},
		}},
	}

	expected := PublicResumeResult{
		Title:     "Backend Engineer",
		Language:  "en",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Overview: PublicOverviewResult{
			Theme: &PublicThemeSettingsResult{Primary: "#112233", Secondary: "#445566"},
			Greeting: PublicGreetingResult{
				First: "Hello", Second: "I am", SecondHighlight: "Ada", Third: "Lovelace",
			},
			Availability: &PublicAvailabilityResult{Status: "open", From: &from, Note: &note},
			Sections: []PublicSectionSettingsResult{
				{ID: "about", Enabled: true, Collapsed: false},
				{ID: "skills", Enabled: false, Collapsed: true},
			},
		},
		Profile: &PublicProfileResult{
			ID: "profile-1", Name: "Ada", Surname: "Lovelace", Patronymic: &patronymic,
			Position: "Backend Engineer", PhotoURL: "/photo.png", DateOfBirth: "1990-01-01",
			Location: "Paris", Relocation: true, WorkFormat: []string{"remote", "hybrid"},
		},
		Contacts: &PublicContactsResult{
			ID: "contacts-1", Phone: &phone, Telegram: &telegram, Email: &email,
			HiddenContacts: []string{"email"},
			ContactsText:   &PublicContactsTextResult{Title: &contactsTitle, Subtitle: &contactsSubtitle},
		},
		Skills: &PublicSkillsResult{
			ID: "skills-1", HardSkills: []string{"Go", "PostgreSQL"}, SoftSkills: []string{"Communication"},
		},
		Languages: []PublicLanguageResult{
			{ID: "language-1", Language: "English", Grade: &grade, Position: 1},
			{ID: "language-2", Language: "French", Position: 2},
		},
		Technologies: []PublicTechnologyResult{
			{ID: "technology-1", Name: "Go", Position: 1},
			{ID: "technology-2", Name: "gRPC", Position: 2},
		},
		Experiences: []PublicExperienceResult{{
			ID: "experience-1", Company: "Example", Position: "Backend Engineer", Location: &location,
			StartDate: "2022-01-01", EndDate: &endDate, Description: &description,
			Achievements: []string{"Shipped API"}, Tasks: []string{"Design services"},
			Stack: []PublicExperienceStackItemResult{{Name: "Go"}, {Name: "gRPC"}},
		}},
		Educations: []PublicEducationResult{{
			ID: "education-1", Location: "Paris", University: "Sorbonne", Grade: 5, ProgramName: "CS",
			ProfileName: &profileName, Faculty: "Science", Course: &course, EndDate: &endDate,
		}},
		Courses: []PublicCourseResult{{
			ID: "course-1", Name: "Go Deep Dive", Organization: "Example School", EndDate: "2025-06-01",
			Hours: &hours, CertificateLink: &certificateLink,
		}},
		Activities: []PublicActivityResult{{
			ID: "activity-1", Type: "open_source", Title: "Contribution", Description: &description,
			Date: &activityDate, Link: &activityLink,
		}},
		Projects: []PublicProjectResult{{
			ID: "project-1", Name: "CV", Description: &projectDescription,
			Screenshots: []string{"https://example.com/shot.png"}, Link: &projectLink,
			Stack: []PublicProjectTechnologyResult{{Name: "Go"}},
			Contributors: []PublicProjectContributorResult{{
				Name: "Grace", Surname: "Hopper", Role: "Advisor",
				Contacts: PublicProjectContactsResult{Phone: &contributorPhone, Email: &contributorEmail},
			}},
		}},
	}

	actual := mapPublicResume(view)

	require.Equal(t, expected, actual)
	require.NotSame(t, view.Profile, actual.Profile)
	require.NotSame(t, &view.Languages[0], &actual.Languages[0])
}

func TestMapPublicResumeNormalizesRequiredSlices(t *testing.T) {
	result := mapPublicResume(resumecontracts.PublicResumeView{})

	require.Nil(t, result.Profile)
	require.Nil(t, result.Contacts)
	require.Nil(t, result.Skills)
	require.Nil(t, result.Overview.Theme)
	require.Nil(t, result.Overview.Availability)
	require.NotNil(t, result.Overview.Sections)
	require.NotNil(t, result.Languages)
	require.NotNil(t, result.Technologies)
	require.NotNil(t, result.Experiences)
	require.NotNil(t, result.Educations)
	require.NotNil(t, result.Courses)
	require.NotNil(t, result.Activities)
	require.NotNil(t, result.Projects)
}

func TestPublicResumeResultDoesNotExposeResumeID(t *testing.T) {
	_, exists := reflect.TypeOf(PublicResumeResult{}).FieldByName("ResumeID")
	require.False(t, exists, "public consumer result must not expose resumeId")
}

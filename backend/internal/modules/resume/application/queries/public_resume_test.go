package queries

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPublicResumeViewJSONGolden(t *testing.T) {
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	view := PublicResumeView{
		ID: "resume-1", Name: "Ada", Surname: "Lovelace", Position: "Engineer", Title: "CV", Language: "en",
		Availability:   &Availability{Status: "open"},
		UpdatedAt:      &now,
		ResumeLink:     stringPtr("https://example.com/cv/resume.pdf"),
		Theme:          &ThemeSettings{Primary: "#fff", Secondary: "#000"},
		Sections:       []SectionSettings{{ID: "about", Enabled: true}},
		HiddenContacts: []string{"email"},
		ContactsText:   &ContactsText{Title: stringPtr("Contact"), Subtitle: stringPtr("Say hello")},
		Greeting:       Greeting{First: "Hello", Second: "I am", SecondHighlight: "Ada", Third: "Engineer"},
		PhotoPath:      "https://example.com/assets/photo.svg", DateOfBirth: "1998-04-12", Location: "Paris", Relocation: true,
		WorkFormat: []string{"remote"}, Contacts: Contacts{Email: stringPtr("ada@example.com")},
		Languages:    []ResumeLanguage{{Language: "en", Grade: stringPtr("C1")}},
		Experience:   []Experience{{ID: "exp-1", Company: "Acme", Position: "Engineer", StartDate: "2024-01", Achievements: []string{"shipped"}, Stack: []Technology{{Name: "Go"}}}},
		Education:    []Education{{ID: "edu-1", Location: "Paris", University: "University", Grade: 2, ProgramName: "CS", Faculty: "Engineering"}},
		Courses:      []Course{{ID: "course-1", Name: "Go", Organization: "Academy", EndDate: "2025-01"}},
		Activities:   []Activity{{ID: "activity-1", Type: "openSource", Title: "Project"}},
		HardSkills:   []string{"Go"},
		SoftSkills:   []string{"Communication"},
		Technologies: []Technology{{Name: "PostgreSQL"}},
		Projects:     []ProjectSummary{{ID: "project-1", Name: "CV", Stack: []Technology{{Name: "Go"}}, Link: stringPtr("https://example.com")}},
	}

	actual, err := json.MarshalIndent(view, "", "  ")
	require.NoError(t, err)

	golden, err := os.ReadFile(filepath.Join("testdata", "public_resume.json"))
	require.NoError(t, err)
	require.JSONEq(t, string(golden), string(actual))
}

func stringPtr(value string) *string { return &value }

func TestProjectPublicResumeMapsCanonicalResumeView(t *testing.T) {
	phone := "+33123456789"
	view := ResumeView{
		ID: "resume-1", Title: "Backend Developer", Language: "en",
		Overview: OverviewPresentationView{
			Greeting: OverviewGreeting{First: "Hello", Second: "I am", SecondHighlight: "Ada", Third: "Engineer"},
			Sections: []OverviewSectionSettings{{ID: "experience", Enabled: true, Collapsed: false}},
		},
		Profile:      &ResumeProfileView{Name: "Ada", Surname: "Lovelace", Position: "Backend Developer", PhotoPath: "https://example.com/avatar.jpg", WorkFormat: []string{"remote"}},
		Contacts:     &ResumeContactsView{Phone: &phone, HiddenContacts: []string{"email"}},
		Skills:       &SkillsView{HardSkills: []string{"Go"}, SoftSkills: []string{"Communication"}},
		Languages:    []LanguageView{{Language: "en", Grade: stringPtr("C1")}},
		Technologies: []TechnologyView{{Name: "PostgreSQL"}},
		Projects:     []ProjectView{{ID: "project-1", Name: "CV", Stack: []ProjectTechnologyView{{Name: "Go"}}}},
	}

	result := projectPublicResume(view)
	require.Equal(t, "resume-1", result.ID)
	require.Equal(t, "Ada", result.Name)
	require.Equal(t, "Backend Developer", result.Position)
	require.Equal(t, []string{"remote"}, result.WorkFormat)
	require.Equal(t, &phone, result.Contacts.Phone)
	require.Equal(t, []string{"email"}, result.HiddenContacts)
	require.Equal(t, []string{"Go"}, result.HardSkills)
	require.Len(t, result.Languages, 1)
	require.Len(t, result.Projects, 1)
	require.Equal(t, []string{"Go"}, []string{result.Projects[0].Stack[0].Name})
}

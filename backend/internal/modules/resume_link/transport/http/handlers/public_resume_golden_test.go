package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type goldenLinkResolver struct{ resumeID sharedvalueobjects.ID }

func (s goldenLinkResolver) Resolve(context.Context, sharedvalueobjects.ID) (sharedvalueobjects.ID, error) {
	return s.resumeID, nil
}

type goldenPublicResumeProvider struct {
	view resumecontracts.PublicResumeView
}

func (s goldenPublicResumeProvider) Handle(context.Context, resumecontracts.PublicResumeQuery) (resumecontracts.PublicResumeView, error) {
	return s.view, nil
}

func TestHandlerPublicResumeGolden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	linkID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	view := richPublicResumeView()

	golden, err := os.ReadFile("testdata/public_resume.golden.json")
	require.NoError(t, err)

	handler := queries.NewPublicResumeLinkQueryHandler(
		goldenLinkResolver{resumeID: resumeID},
		goldenPublicResumeProvider{view: view},
	)
	engine := gin.New()
	engine.GET("/api/v1/cv/:linkId", NewPublicResumeLinkHandler(handler).Handle)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cv/"+linkID.String(), nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "no-store", w.Header().Get("Cache-Control"))
	require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	require.Empty(t, w.Header().Get("X-View-Id"))
	require.Empty(t, w.Header().Get("Set-Cookie"))
	require.NotContains(t, w.Body.String(), `"resumeId"`)
	require.NotContains(t, w.Body.String(), "null")
	require.JSONEq(t, string(golden), w.Body.String())
}

func richPublicResumeView() resumecontracts.PublicResumeView {
	from := "2026-03-01"
	note := "Remote only"
	patronymic := "Ivanovna"
	phone := "+33123456789"
	telegram := "ada"
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
	style := "second"

	return resumecontracts.PublicResumeView{
		Title:     "Backend Engineer",
		Language:  "en",
		CreatedAt: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
		UpdatedAt: time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC),
		Overview: resumecontracts.PublicOverviewView{
			Theme: &resumecontracts.PublicThemeSettings{Primary: "#112233", Secondary: "#445566"},
			Greeting: resumecontracts.PublicGreeting{
				First: "Hello", Second: "I am", SecondHighlight: "Ada", Third: "Lovelace",
			},
			Availability: &resumecontracts.PublicAvailability{Status: "open", From: &from, Note: &note},
			Sections: []resumecontracts.PublicSectionSettings{
				{ID: "about", Enabled: true, Collapsed: false, Style: &style},
				{ID: "skills", Enabled: false, Collapsed: true},
			},
		},
		Profile: &resumecontracts.PublicProfileView{
			ID: "profile-1", Name: "Ada", Surname: "Lovelace", Patronymic: &patronymic,
			Position: "Backend Engineer", PhotoURL: "/photo.png", DateOfBirth: "1990-01-01",
			Location: "Paris", Relocation: true, WorkFormat: []string{"remote", "hybrid"},
		},
		Contacts: &resumecontracts.PublicContactsView{
			ID: "contacts-1", Phone: &phone, Telegram: &telegram,
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
		Projects: []resumecontracts.PublicProjectView{
			{
				ID: "project-1", Name: "CV", Description: &projectDescription,
				Screenshots: []string{"https://example.com/shot.png"}, Link: &projectLink,
				Stack: []resumecontracts.PublicProjectTechnologyView{{Name: "Go"}},
				Contributors: []resumecontracts.PublicProjectContributorView{{
					Name: "Grace", Surname: "Hopper", Role: "Advisor",
					Contacts: resumecontracts.PublicProjectContactsView{Phone: &contributorPhone, Email: &contributorEmail},
				}},
			},
			{
				ID: "project-2", Name: "Empty arrays",
				Screenshots: []string{}, Stack: []resumecontracts.PublicProjectTechnologyView{}, Contributors: []resumecontracts.PublicProjectContributorView{},
			},
		},
	}
}

package queries

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func findAnalyticsLayoutSectionHash(t *testing.T, view PublicResumeView, sectionID string) string {
	t.Helper()
	layout, err := ToResumeAnalyticsLayout(view)
	require.NoError(t, err)
	for _, section := range layout.Sections {
		if section.ID == sectionID {
			return section.Hash
		}
	}
	t.Fatalf("section %q not found in analytics layout", sectionID)
	return ""
}

func TestToResumeAnalyticsLayoutContainsOnlySectionHashesAndTechnologies(t *testing.T) {
	profileName := "Private Profile Name"
	phone := "+33123456789"
	view := PublicResumeView{
		Language: "en",
		Overview: PublicOverviewView{
			Sections: []PublicSectionSettings{
				{ID: "projects", Enabled: true},
				{ID: "contacts", Enabled: true},
			},
		},
		Profile:      &PublicProfileView{Name: profileName},
		Contacts:     &PublicContactsView{Phone: &phone},
		Projects:     []PublicProjectView{{Name: "Private project title"}},
		Skills:       &PublicSkillsView{HardSkills: []string{"Go"}},
		Technologies: []PublicTechnologyView{{Name: "PostgreSQL"}},
	}

	layout, err := ToResumeAnalyticsLayout(view)

	require.NoError(t, err)
	require.Equal(t, []string{"PostgreSQL", "Go"}, layout.Technologies)
	require.Len(t, layout.Sections, 3)
	require.Equal(t, "about", layout.Sections[0].ID)
	for _, section := range layout.Sections {
		require.Len(t, section.Hash, 16)
	}
	encoded, err := json.Marshal(layout)
	require.NoError(t, err)
	require.NotContains(t, string(encoded), profileName)
	require.NotContains(t, string(encoded), phone)
	require.NotContains(t, string(encoded), "Private project title")
}

func TestToResumeAnalyticsLayoutPreservesSectionVisibilityRules(t *testing.T) {
	view := PublicResumeView{
		Overview: PublicOverviewView{
			Sections: []PublicSectionSettings{
				{ID: "projects", Enabled: true},
				{ID: "about", Enabled: false},
				{ID: "techStack", Enabled: true},
				{ID: "experience", Enabled: false},
			},
		},
		Projects:     []PublicProjectView{{Name: "Project"}},
		Experiences:  []PublicExperienceView{{Company: "Company"}},
		Technologies: []PublicTechnologyView{{Name: "Go"}},
	}

	layout, err := ToResumeAnalyticsLayout(view)

	require.NoError(t, err)
	require.Len(t, layout.Sections, 2)
	require.Equal(t, []string{"about", "projects"}, []string{layout.Sections[0].ID, layout.Sections[1].ID})
}

func TestToResumeAnalyticsLayoutIgnoresProjectIDsAndTracksContent(t *testing.T) {
	base := PublicResumeView{
		Language: "en",
		Overview: PublicOverviewView{Sections: []PublicSectionSettings{{ID: "projects", Enabled: true}}},
		Projects: []PublicProjectView{{ID: "project-a", Name: "Project"}},
	}
	changedID := base
	changedID.Projects = []PublicProjectView{{ID: "project-b", Name: "Project"}}
	changedContent := base
	changedContent.Projects = []PublicProjectView{{ID: "project-a", Name: "Updated"}}

	require.Equal(t, findAnalyticsLayoutSectionHash(t, base, "projects"), findAnalyticsLayoutSectionHash(t, changedID, "projects"))
	require.NotEqual(t, findAnalyticsLayoutSectionHash(t, base, "projects"), findAnalyticsLayoutSectionHash(t, changedContent, "projects"))
}

func TestToResumeAnalyticsLayoutIgnoresProfileIDButTracksProfileContent(t *testing.T) {
	base := PublicResumeView{
		Language: "en",
		Profile:  &PublicProfileView{ID: "profile-a", Name: "Ada", Surname: "Lovelace"},
	}
	changedID := base
	changedID.Profile = &PublicProfileView{ID: "profile-b", Name: "Ada", Surname: "Lovelace"}
	changedContent := base
	changedContent.Profile = &PublicProfileView{ID: "profile-a", Name: "Augusta", Surname: "Lovelace"}

	require.Equal(t, findAnalyticsLayoutSectionHash(t, base, "about"), findAnalyticsLayoutSectionHash(t, changedID, "about"))
	require.NotEqual(t, findAnalyticsLayoutSectionHash(t, base, "about"), findAnalyticsLayoutSectionHash(t, changedContent, "about"))
}

func TestToResumeAnalyticsLayoutTracksExperienceContent(t *testing.T) {
	base := PublicResumeView{
		Language:    "en",
		Overview:    PublicOverviewView{Sections: []PublicSectionSettings{{ID: "experience", Enabled: true}}},
		Experiences: []PublicExperienceView{{ID: "experience-a", Company: "Example Inc", Position: "Engineer", StartDate: "2024-01"}},
	}
	changed := base
	changed.Experiences = []PublicExperienceView{{ID: "experience-a", Company: "Example Inc", Position: "Senior Engineer", StartDate: "2024-01"}}

	require.NotEqual(t, findAnalyticsLayoutSectionHash(t, base, "experience"), findAnalyticsLayoutSectionHash(t, changed, "experience"))
}

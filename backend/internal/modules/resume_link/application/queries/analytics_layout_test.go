package queries

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
)

func findAnalyticsLayoutSectionHash(t *testing.T, view resumecontracts.PublicResumeView, sectionID string) string {
	t.Helper()
	layout, err := BuildAnalyticsLayout(view)
	require.NoError(t, err)
	for _, section := range layout.Sections {
		if section.ID == sectionID {
			return section.Hash
		}
	}
	t.Fatalf("section %q not found in analytics layout", sectionID)
	return ""
}

func TestBuildAnalyticsLayoutContainsOnlySectionHashesAndTechnologies(t *testing.T) {
	profileName := "Private Profile Name"
	phone := "+33123456789"
	view := resumecontracts.PublicResumeView{
		Language: "en",
		Overview: resumecontracts.PublicOverviewView{
			Sections: []resumecontracts.PublicSectionSettings{
				{ID: "projects", Enabled: true},
				{ID: "contacts", Enabled: true},
			},
		},
		Profile:      &resumecontracts.PublicProfileView{Name: profileName},
		Contacts:     &resumecontracts.PublicContactsView{Phone: &phone},
		Projects:     []resumecontracts.PublicProjectView{{Name: "Private project title"}},
		Skills:       &resumecontracts.PublicSkillsView{HardSkills: []string{"Go"}},
		Technologies: []resumecontracts.PublicTechnologyView{{Name: "PostgreSQL"}},
	}

	layout, err := BuildAnalyticsLayout(view)

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

func TestBuildAnalyticsLayoutPreservesSectionVisibilityRules(t *testing.T) {
	view := resumecontracts.PublicResumeView{
		Overview: resumecontracts.PublicOverviewView{
			Sections: []resumecontracts.PublicSectionSettings{
				{ID: "projects", Enabled: true},
				{ID: "about", Enabled: false},
				{ID: "techStack", Enabled: true},
				{ID: "experience", Enabled: false},
			},
		},
		Projects:     []resumecontracts.PublicProjectView{{Name: "Project"}},
		Experiences:  []resumecontracts.PublicExperienceView{{Company: "Company"}},
		Technologies: []resumecontracts.PublicTechnologyView{{Name: "Go"}},
	}

	layout, err := BuildAnalyticsLayout(view)

	require.NoError(t, err)
	require.Len(t, layout.Sections, 2)
	require.Equal(t, []string{"about", "projects"}, []string{layout.Sections[0].ID, layout.Sections[1].ID})
}

func TestBuildAnalyticsLayoutKeepsStableHashes(t *testing.T) {
	phone := "+33123456789"
	view := resumecontracts.PublicResumeView{
		Language: "en",
		Overview: resumecontracts.PublicOverviewView{
			Sections: []resumecontracts.PublicSectionSettings{
				{ID: "projects", Enabled: true},
				{ID: "contacts", Enabled: true},
				{ID: "experience", Enabled: true},
			},
		},
		Profile:      &resumecontracts.PublicProfileView{Name: "Ada", Surname: "Lovelace"},
		Contacts:     &resumecontracts.PublicContactsView{Phone: &phone},
		Projects:     []resumecontracts.PublicProjectView{{Name: "Private project title"}},
		Experiences:  []resumecontracts.PublicExperienceView{{Company: "Example", Position: "Engineer", StartDate: "2024-01"}},
		Skills:       &resumecontracts.PublicSkillsView{HardSkills: []string{"Go"}},
		Technologies: []resumecontracts.PublicTechnologyView{{Name: "PostgreSQL"}},
	}

	layout, err := BuildAnalyticsLayout(view)

	require.NoError(t, err)
	require.Equal(t, []string{"PostgreSQL", "Go"}, layout.Technologies)
	ids := make([]string, 0, len(layout.Sections))
	hashes := make([]string, 0, len(layout.Sections))
	for _, section := range layout.Sections {
		ids = append(ids, section.ID)
		hashes = append(hashes, section.Hash)
	}
	require.Equal(t, []string{"about", "projects", "contacts", "experience"}, ids)
	require.Equal(t, []string{"4569dcc4e09676ca", "d54f9a6e0bebe49b", "32f412160315ade6", "d45de8fb9051cb16"}, hashes)
}

func TestBuildAnalyticsLayoutIgnoresProjectIDsAndTracksContent(t *testing.T) {
	base := resumecontracts.PublicResumeView{
		Language: "en",
		Overview: resumecontracts.PublicOverviewView{Sections: []resumecontracts.PublicSectionSettings{{ID: "projects", Enabled: true}}},
		Projects: []resumecontracts.PublicProjectView{{ID: "project-a", Name: "Project"}},
	}
	changedID := base
	changedID.Projects = []resumecontracts.PublicProjectView{{ID: "project-b", Name: "Project"}}
	changedContent := base
	changedContent.Projects = []resumecontracts.PublicProjectView{{ID: "project-a", Name: "Updated"}}

	require.Equal(t, findAnalyticsLayoutSectionHash(t, base, "projects"), findAnalyticsLayoutSectionHash(t, changedID, "projects"))
	require.NotEqual(t, findAnalyticsLayoutSectionHash(t, base, "projects"), findAnalyticsLayoutSectionHash(t, changedContent, "projects"))
}

func TestBuildAnalyticsLayoutIgnoresProfileIDButTracksProfileContent(t *testing.T) {
	base := resumecontracts.PublicResumeView{
		Language: "en",
		Profile:  &resumecontracts.PublicProfileView{ID: "profile-a", Name: "Ada", Surname: "Lovelace"},
	}
	changedID := base
	changedID.Profile = &resumecontracts.PublicProfileView{ID: "profile-b", Name: "Ada", Surname: "Lovelace"}
	changedContent := base
	changedContent.Profile = &resumecontracts.PublicProfileView{ID: "profile-a", Name: "Augusta", Surname: "Lovelace"}

	require.Equal(t, findAnalyticsLayoutSectionHash(t, base, "about"), findAnalyticsLayoutSectionHash(t, changedID, "about"))
	require.NotEqual(t, findAnalyticsLayoutSectionHash(t, base, "about"), findAnalyticsLayoutSectionHash(t, changedContent, "about"))
}

func TestBuildAnalyticsLayoutTracksExperienceContent(t *testing.T) {
	base := resumecontracts.PublicResumeView{
		Language:    "en",
		Overview:    resumecontracts.PublicOverviewView{Sections: []resumecontracts.PublicSectionSettings{{ID: "experience", Enabled: true}}},
		Experiences: []resumecontracts.PublicExperienceView{{ID: "experience-a", Company: "Example Inc", Position: "Engineer", StartDate: "2024-01"}},
	}
	changed := base
	changed.Experiences = []resumecontracts.PublicExperienceView{{ID: "experience-a", Company: "Example Inc", Position: "Senior Engineer", StartDate: "2024-01"}}

	require.NotEqual(t, findAnalyticsLayoutSectionHash(t, base, "experience"), findAnalyticsLayoutSectionHash(t, changed, "experience"))
}

func TestBuildAnalyticsLayoutContainsPublicSectionInputsOnly(t *testing.T) {
	phone := "+33123456789"
	view := resumecontracts.PublicResumeView{
		Title:    "Private editor metadata",
		Language: "en",
		Overview: resumecontracts.PublicOverviewView{
			Sections: []resumecontracts.PublicSectionSettings{{ID: "projects", Enabled: true}},
		},
		Profile:  &resumecontracts.PublicProfileView{Name: "Ada"},
		Contacts: &resumecontracts.PublicContactsView{Phone: &phone},
	}

	layout, err := BuildAnalyticsLayout(view)

	require.NoError(t, err)
	require.Len(t, layout.Sections, 1)
	require.Equal(t, "about", layout.Sections[0].ID)
	require.Len(t, layout.Sections[0].Hash, 16)
	require.Empty(t, layout.Technologies)
	encoded, err := json.Marshal(layout)
	require.NoError(t, err)
	require.NotContains(t, string(encoded), "Ada")
	require.NotContains(t, string(encoded), phone)
	require.NotContains(t, string(encoded), "Profile")
	require.NotContains(t, string(encoded), "Contacts")
}

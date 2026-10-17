package analytics

import (
	"testing"

	"github.com/stretchr/testify/require"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumedomain "developers-cv/backend/internal/modules/resume/domain"
)

func TestLayoutSnapshotMatchesFrontendVisibilityRules(t *testing.T) {
	builder := NewLayoutSnapshotBuilder()
	resume := resumequeries.PublicResumeView{
		Language: "en",
		Overview: resumequeries.PublicOverviewView{Sections: []resumequeries.PublicSectionSettings{
			{ID: string(resumedomain.SectionProjects), Enabled: true},
			{ID: string(resumedomain.SectionAbout), Enabled: false},
			{ID: string(resumedomain.SectionTechStack), Enabled: true},
		}},
		Projects:     []resumequeries.PublicProjectView{{Name: "Project"}},
		Technologies: []resumequeries.PublicTechnologyView{{Name: "Go"}},
	}
	layout, err := builder.Build(resume)
	require.NoError(t, err)
	require.Equal(t, []resumedomain.SectionID{resumedomain.SectionAbout, resumedomain.SectionProjects}, layout.Order())
	require.Equal(t, []string{"Go"}, layout.Technologies())
	_, ok := layout.HashFor(resumedomain.SectionTechStack)
	require.False(t, ok)
}

func TestLayoutHashIgnoresProjectIDs(t *testing.T) {
	base := resumequeries.PublicResumeView{
		Language: "en",
		Overview: resumequeries.PublicOverviewView{Sections: []resumequeries.PublicSectionSettings{{ID: "projects", Enabled: true}}},
		Projects: []resumequeries.PublicProjectView{{ID: "project-a", Name: "Project"}},
	}
	changed := base
	changed.Projects = []resumequeries.PublicProjectView{{ID: "project-b", Name: "Project"}}

	require.Equal(t, sectionHashForTest(t, base, resumedomain.SectionProjects), sectionHashForTest(t, changed, resumedomain.SectionProjects))
}

func TestLayoutHashChangesWhenPublicProjectContentChanges(t *testing.T) {
	base := resumequeries.PublicResumeView{
		Language: "en",
		Overview: resumequeries.PublicOverviewView{Sections: []resumequeries.PublicSectionSettings{{ID: "projects", Enabled: true}}},
		Projects: []resumequeries.PublicProjectView{{ID: "project-a", Name: "Project", Description: stringPointer("Original description")}},
	}
	changed := base
	changed.Projects = []resumequeries.PublicProjectView{{ID: "project-a", Name: "Project", Description: stringPointer("Updated description")}}

	require.NotEqual(t, sectionHashForTest(t, base, resumedomain.SectionProjects), sectionHashForTest(t, changed, resumedomain.SectionProjects))
}

func TestLayoutHashIgnoresProfileIDButTracksProfileContent(t *testing.T) {
	base := resumequeries.PublicResumeView{
		Language: "en",
		Profile:  &resumequeries.PublicProfileView{ID: "profile-a", Name: "Ada", Surname: "Lovelace"},
	}
	changedID := base
	changedID.Profile = &resumequeries.PublicProfileView{ID: "profile-b", Name: "Ada", Surname: "Lovelace"}
	changedContent := base
	changedContent.Profile = &resumequeries.PublicProfileView{ID: "profile-a", Name: "Augusta", Surname: "Lovelace"}

	baseHash := sectionHashForTest(t, base, resumedomain.SectionAbout)
	require.Equal(t, baseHash, sectionHashForTest(t, changedID, resumedomain.SectionAbout))
	require.NotEqual(t, baseHash, sectionHashForTest(t, changedContent, resumedomain.SectionAbout))
}

func TestLayoutHashChangesWhenPublicExperienceContentChanges(t *testing.T) {
	base := resumequeries.PublicResumeView{
		Language:    "en",
		Overview:    resumequeries.PublicOverviewView{Sections: []resumequeries.PublicSectionSettings{{ID: "experience", Enabled: true}}},
		Experiences: []resumequeries.PublicExperienceView{{ID: "experience-a", Company: "Example Inc", Position: "Engineer", StartDate: "2024-01"}},
	}
	changed := base
	changed.Experiences = []resumequeries.PublicExperienceView{{ID: "experience-a", Company: "Example Inc", Position: "Senior Engineer", StartDate: "2024-01"}}

	require.NotEqual(t, sectionHashForTest(t, base, resumedomain.SectionExperience), sectionHashForTest(t, changed, resumedomain.SectionExperience))
}

func sectionHashForTest(t *testing.T, resume resumequeries.PublicResumeView, section resumedomain.SectionID) string {
	t.Helper()
	hash, err := sectionHash(section, resume)
	require.NoError(t, err)
	return hash
}

func stringPointer(value string) *string { return &value }

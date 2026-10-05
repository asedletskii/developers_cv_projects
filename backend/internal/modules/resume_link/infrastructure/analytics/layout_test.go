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
		Overview: resumequeries.ResumeOverviewView{Sections: []resumequeries.OverviewSectionSettings{
			{ID: string(resumedomain.SectionProjects), Enabled: true},
			{ID: string(resumedomain.SectionAbout), Enabled: false},
			{ID: string(resumedomain.SectionTechStack), Enabled: true},
		}},
		Projects:     []resumequeries.ProjectView{{Name: "Project"}},
		Technologies: []resumequeries.TechnologyView{{Name: "Go"}},
	}
	layout, err := builder.Build(resume)
	require.NoError(t, err)
	require.Equal(t, []resumedomain.SectionID{resumedomain.SectionAbout, resumedomain.SectionProjects}, layout.Order())
	_, ok := layout.HashFor(resumedomain.SectionTechStack)
	require.False(t, ok)
}

func TestLayoutHashIgnoresProjectIDs(t *testing.T) {
	builder := NewLayoutSnapshotBuilder()
	base := resumequeries.PublicResumeView{Language: "en", Overview: resumequeries.ResumeOverviewView{Sections: []resumequeries.OverviewSectionSettings{{ID: "projects", Enabled: true}}}, Projects: []resumequeries.ProjectView{{ID: "a", Name: "Project"}}}
	changed := base
	changed.Projects = []resumequeries.ProjectView{{ID: "b", Name: "Project"}}
	a, err := builder.Build(base)
	require.NoError(t, err)
	b, err := builder.Build(changed)
	require.NoError(t, err)
	ha, _ := a.HashFor(resumedomain.SectionProjects)
	hb, _ := b.HashFor(resumedomain.SectionProjects)
	require.Equal(t, ha, hb)
}

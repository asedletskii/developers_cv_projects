package analytics

import (
	"testing"

	"github.com/stretchr/testify/require"

	analyticsports "developers-cv/backend/internal/modules/analytics/application/ports"
	analyticsdomain "developers-cv/backend/internal/modules/analytics/domain"
)

func TestLayoutSnapshotBuilderConsumesSectionHashesAndTechnologyCatalog(t *testing.T) {
	builder := NewLayoutSnapshotBuilder()
	resume := analyticsports.ResumeAnalyticsLayout{
		Sections: []analyticsports.ResumeAnalyticsSection{
			{ID: "about", Hash: "0123456789abcdef"},
			{ID: "projects", Hash: "fedcba9876543210"},
		},
		Technologies: []string{"Go", "PostgreSQL"},
	}

	layout, err := builder.Build(resume)

	require.NoError(t, err)
	require.Equal(t, []analyticsdomain.AnalyticsSectionID{analyticsdomain.AnalyticsSectionAbout, analyticsdomain.AnalyticsSectionProjects}, layout.Order())
	require.Equal(t, "0123456789abcdef", layout.Hashes()[analyticsdomain.AnalyticsSectionAbout])
	require.Equal(t, []string{"Go", "PostgreSQL"}, layout.Technologies())
}

func TestLayoutSnapshotBuilderRejectsInvalidSectionOrHash(t *testing.T) {
	builder := NewLayoutSnapshotBuilder()

	_, err := builder.Build(analyticsports.ResumeAnalyticsLayout{
		Sections: []analyticsports.ResumeAnalyticsSection{{ID: "invalid", Hash: "0123456789abcdef"}},
	})
	require.ErrorIs(t, err, analyticsdomain.ErrInvalidAnalyticsSection)

	_, err = builder.Build(analyticsports.ResumeAnalyticsLayout{
		Sections: []analyticsports.ResumeAnalyticsSection{{ID: "about", Hash: "invalid"}},
	})
	require.ErrorIs(t, err, analyticsdomain.ErrInvalidLayoutSnapshot)
}

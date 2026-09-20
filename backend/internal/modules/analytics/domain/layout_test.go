package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewLayoutSnapshotRejectsTechStackAndInvalidHash(t *testing.T) {
	_, err := NewLayoutSnapshot(
		[]AnalyticsSectionID{AnalyticsSectionTechStack},
		map[AnalyticsSectionID]string{AnalyticsSectionTechStack: "0123456789abcdef"},
	)
	require.ErrorIs(t, err, ErrInvalidLayoutSnapshot)

	_, err = NewLayoutSnapshot(
		[]AnalyticsSectionID{AnalyticsSectionAbout},
		map[AnalyticsSectionID]string{AnalyticsSectionAbout: "invalid"},
	)
	require.ErrorIs(t, err, ErrInvalidLayoutSnapshot)
}

func TestLayoutSnapshotReturnsCopies(t *testing.T) {
	layout, err := NewLayoutSnapshot(
		[]AnalyticsSectionID{AnalyticsSectionAbout},
		map[AnalyticsSectionID]string{AnalyticsSectionAbout: "0123456789abcdef"},
	)
	require.NoError(t, err)

	order := layout.Order()
	order[0] = AnalyticsSectionProjects
	hashes := layout.Hashes()
	hashes[AnalyticsSectionAbout] = "fedcba9876543210"

	require.Equal(t, AnalyticsSectionAbout, layout.Order()[0])
	require.Equal(t, "0123456789abcdef", layout.Hashes()[AnalyticsSectionAbout])
}

func TestLayoutSnapshotNormalizesAndCopiesTechnologyCatalog(t *testing.T) {
	layout, err := NewLayoutSnapshot(
		[]AnalyticsSectionID{AnalyticsSectionAbout},
		map[AnalyticsSectionID]string{AnalyticsSectionAbout: "0123456789abcdef"},
		[]string{" Go ", "go", "PostgreSQL"},
	)
	require.NoError(t, err)
	require.Equal(t, []string{"Go", "PostgreSQL"}, layout.Technologies())
	canonical, ok := layout.NormalizeTechnology("  gO ")
	require.True(t, ok)
	require.Equal(t, "Go", canonical)
	_, ok = layout.NormalizeTechnology("MadeUp")
	require.False(t, ok)

	technologies := layout.Technologies()
	technologies[0] = "mutated"
	require.Equal(t, "Go", layout.Technologies()[0])
}

func TestRestoreLayoutSnapshotRejectsUnnormalizedTechnologyCatalog(t *testing.T) {
	_, err := RestoreLayoutSnapshot(
		[]AnalyticsSectionID{AnalyticsSectionAbout},
		map[AnalyticsSectionID]string{AnalyticsSectionAbout: "0123456789abcdef"},
		[]string{" Go "},
	)
	require.ErrorIs(t, err, ErrInvalidLayoutSnapshot)
}

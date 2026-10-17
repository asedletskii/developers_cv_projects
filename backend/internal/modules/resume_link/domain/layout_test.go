package domain

import (
	"testing"

	resumedomain "developers-cv/backend/internal/modules/resume/domain"

	"github.com/stretchr/testify/require"
)

func TestNewLayoutSnapshotRejectsTechStackAndInvalidHash(t *testing.T) {
	_, err := NewLayoutSnapshot(
		[]resumedomain.SectionID{resumedomain.SectionTechStack},
		map[resumedomain.SectionID]string{resumedomain.SectionTechStack: "0123456789abcdef"},
	)
	require.ErrorIs(t, err, ErrInvalidLayoutSnapshot)

	_, err = NewLayoutSnapshot(
		[]resumedomain.SectionID{resumedomain.SectionAbout},
		map[resumedomain.SectionID]string{resumedomain.SectionAbout: "invalid"},
	)
	require.ErrorIs(t, err, ErrInvalidLayoutSnapshot)
}

func TestLayoutSnapshotReturnsCopies(t *testing.T) {
	layout, err := NewLayoutSnapshot(
		[]resumedomain.SectionID{resumedomain.SectionAbout},
		map[resumedomain.SectionID]string{resumedomain.SectionAbout: "0123456789abcdef"},
	)
	require.NoError(t, err)

	order := layout.Order()
	order[0] = resumedomain.SectionProjects
	hashes := layout.Hashes()
	hashes[resumedomain.SectionAbout] = "fedcba9876543210"

	require.Equal(t, resumedomain.SectionAbout, layout.Order()[0])
	require.Equal(t, "0123456789abcdef", layout.Hashes()[resumedomain.SectionAbout])
}

func TestLayoutSnapshotNormalizesAndCopiesTechnologyCatalog(t *testing.T) {
	layout, err := NewLayoutSnapshot(
		[]resumedomain.SectionID{resumedomain.SectionAbout},
		map[resumedomain.SectionID]string{resumedomain.SectionAbout: "0123456789abcdef"},
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
		[]resumedomain.SectionID{resumedomain.SectionAbout},
		map[resumedomain.SectionID]string{resumedomain.SectionAbout: "0123456789abcdef"},
		[]string{" Go "},
	)
	require.ErrorIs(t, err, ErrInvalidLayoutSnapshot)
}

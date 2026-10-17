package domain

import (
	"testing"
	"time"

	resumedomain "developers-cv/backend/internal/modules/resume/domain"

	"github.com/stretchr/testify/require"
)

func testViewStateLayout(t *testing.T) LayoutSnapshot {
	t.Helper()
	layout, err := NewLayoutSnapshot(
		[]resumedomain.SectionID{
			resumedomain.SectionAbout,
			resumedomain.SectionExperience,
			resumedomain.SectionProjects,
		},
		map[resumedomain.SectionID]string{
			resumedomain.SectionAbout:      "0123456789abcdef",
			resumedomain.SectionExperience: "fedcba9876543210",
			resumedomain.SectionProjects:   "0011223344556677",
		},
	)
	require.NoError(t, err)
	return layout
}

func TestViewStateApplyIsMonotonic(t *testing.T) {
	layout := testViewStateLayout(t)
	var state ViewState

	first := time.Date(2026, 10, 5, 10, 0, 5, 0, time.UTC)
	require.NoError(t, state.Apply(ViewStateUpdate{
		Seq: 1, ActiveMS: 7000, MaxSection: resumedomain.SectionExperience,
		ExitSection: resumedomain.SectionExperience, SeenAt: first,
	}, layout))
	require.Equal(t, int64(7000), state.ActiveMS())
	require.Equal(t, resumedomain.SectionExperience, state.MaxSection())
	require.Equal(t, first, state.ConfirmedAt())

	second := first.Add(time.Second)
	require.NoError(t, state.Apply(ViewStateUpdate{
		Seq: 2, ActiveMS: 3000, MaxSection: resumedomain.SectionAbout,
		ExitSection: resumedomain.SectionProjects, SeenAt: second, Final: true,
	}, layout))
	require.Equal(t, int64(7000), state.ActiveMS())
	require.Equal(t, resumedomain.SectionExperience, state.MaxSection())
	require.Equal(t, resumedomain.SectionProjects, state.ExitSection())
	require.Equal(t, second, state.EndedAt())

	require.NoError(t, state.Apply(ViewStateUpdate{
		Seq: 2, ActiveMS: 7000, SeenAt: second.Add(-time.Second), Final: true,
	}, layout))
	require.Equal(t, second, state.EndedAt())

	require.NoError(t, state.Apply(ViewStateUpdate{
		Seq: 3, ActiveMS: 8000, MaxSection: resumedomain.SectionProjects,
		SeenAt: second.Add(time.Second),
	}, layout))
	require.True(t, state.EndedAt().IsZero())
	require.Equal(t, int64(8000), state.ActiveMS())

	stateBeforeStale := state
	require.NoError(t, state.Apply(ViewStateUpdate{
		Seq: 2, ActiveMS: 1, MaxSection: resumedomain.SectionAbout,
		ExitSection: resumedomain.SectionAbout, SeenAt: second.Add(2 * time.Second), Final: true,
	}, layout))
	require.Equal(t, stateBeforeStale.ActiveMS(), state.ActiveMS())
	require.Equal(t, stateBeforeStale.ExitSection(), state.ExitSection())
	require.True(t, state.EndedAt().IsZero())
}

func TestViewStateFromPersistenceRestoresValues(t *testing.T) {
	confirmed := time.Date(2026, 10, 5, 10, 0, 5, 0, time.UTC)
	lastSeen := confirmed.Add(time.Second)
	ended := lastSeen.Add(time.Second)

	state := ViewStateFromPersistence(
		7, &confirmed, &lastSeen, 8000,
		string(resumedomain.SectionExperience), string(resumedomain.SectionProjects), &ended,
	)

	require.Equal(t, int64(7), state.LastSeq())
	require.Equal(t, confirmed, state.ConfirmedAt())
	require.Equal(t, lastSeen, state.LastSeenAt())
	require.Equal(t, int64(8000), state.ActiveMS())
	require.Equal(t, resumedomain.SectionExperience, state.MaxSection())
	require.Equal(t, resumedomain.SectionProjects, state.ExitSection())
	require.Equal(t, ended, state.EndedAt())
}

func TestViewStateApplyRejectsActiveTimeAboveStorageLimit(t *testing.T) {
	state := ViewState{}
	err := state.Apply(ViewStateUpdate{
		Seq: 1, ActiveMS: MaxActiveMS + 1,
		SeenAt: time.Date(2026, 10, 5, 10, 0, 1, 0, time.UTC),
	}, testViewLayout(t))
	require.ErrorIs(t, err, ErrInvalidViewState)
}

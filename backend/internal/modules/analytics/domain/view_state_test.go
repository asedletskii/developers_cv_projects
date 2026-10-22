package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestViewStateApplyIsMonotonic(t *testing.T) {
	layout := testViewLayout(t)
	var state ViewState

	first := time.Date(2026, 10, 5, 10, 0, 5, 0, time.UTC)
	require.NoError(t, state.Apply(ViewStateUpdate{
		Seq: 1, ActiveMS: 7000, MaxSection: AnalyticsSectionExperience,
		ExitSection: AnalyticsSectionExperience, SeenAt: first,
	}, layout))
	require.Equal(t, int64(7000), state.ActiveMS())
	require.Equal(t, AnalyticsSectionExperience, state.MaxSection())
	require.Equal(t, first, state.ConfirmedAt())

	second := first.Add(time.Second)
	require.NoError(t, state.Apply(ViewStateUpdate{
		Seq: 2, ActiveMS: 3000, MaxSection: AnalyticsSectionAbout,
		ExitSection: AnalyticsSectionProjects, SeenAt: second, Final: true,
	}, layout))
	require.Equal(t, int64(7000), state.ActiveMS())
	require.Equal(t, AnalyticsSectionExperience, state.MaxSection())
	require.Equal(t, AnalyticsSectionProjects, state.ExitSection())
	require.Equal(t, second, state.EndedAt())

	require.NoError(t, state.Apply(ViewStateUpdate{
		Seq: 2, ActiveMS: 7000, SeenAt: second.Add(-time.Second), Final: true,
	}, layout))
	require.Equal(t, second, state.EndedAt())

	require.NoError(t, state.Apply(ViewStateUpdate{
		Seq: 3, ActiveMS: 8000, MaxSection: AnalyticsSectionProjects,
		SeenAt: second.Add(time.Second),
	}, layout))
	require.True(t, state.EndedAt().IsZero())
	require.Equal(t, int64(8000), state.ActiveMS())

	stateBeforeStale := state
	require.NoError(t, state.Apply(ViewStateUpdate{
		Seq: 2, ActiveMS: 1, MaxSection: AnalyticsSectionAbout,
		ExitSection: AnalyticsSectionAbout, SeenAt: second.Add(2 * time.Second), Final: true,
	}, layout))
	require.Equal(t, stateBeforeStale.ActiveMS(), state.ActiveMS())
	require.Equal(t, stateBeforeStale.ExitSection(), state.ExitSection())
	require.True(t, state.EndedAt().IsZero())
}

func TestViewStateSnapshotPreservesValues(t *testing.T) {
	confirmed := time.Date(2026, 10, 5, 10, 0, 5, 0, time.UTC)
	lastSeen := confirmed.Add(time.Second)
	ended := lastSeen.Add(time.Second)

	state := ViewState{
		lastSeq: 7, confirmedAt: confirmed, lastSeenAt: lastSeen, activeMS: 8000,
		maxSection: AnalyticsSectionExperience, exitSection: AnalyticsSectionProjects, endedAt: ended,
	}.Snapshot()

	require.Equal(t, int64(7), state.LastSeq)
	require.Equal(t, confirmed, state.ConfirmedAt)
	require.Equal(t, lastSeen, state.LastSeenAt)
	require.Equal(t, int64(8000), state.ActiveMS)
	require.Equal(t, AnalyticsSectionExperience, state.MaxSection)
	require.Equal(t, AnalyticsSectionProjects, state.ExitSection)
	require.Equal(t, ended, state.EndedAt)
}

func TestViewStateApplyRejectsActiveTimeAboveStorageLimit(t *testing.T) {
	state := ViewState{}
	err := state.Apply(ViewStateUpdate{
		Seq: 1, ActiveMS: MaxActiveMS + 1,
		SeenAt: time.Date(2026, 10, 5, 10, 0, 1, 0, time.UTC),
	}, testViewLayout(t))
	require.ErrorIs(t, err, ErrInvalidViewState)
}

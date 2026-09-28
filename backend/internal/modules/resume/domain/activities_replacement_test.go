package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func activityParams(id valueobjects.ID, title string) ActivityParams {
	return ActivityParams{ID: id, Type: enums.ActivityTypeOpenSource, Title: title}
}

func newActivityResume(t *testing.T) *Resume {
	t.Helper()
	r, err := NewResume(ResumeParams{ID: valueobjects.NewID(), UserID: useridentity.NewID(), Title: "Resume", Language: enums.ResumeLanguageEN})
	require.NoError(t, err)
	return r
}

func TestResumeReplaceActivitiesCreatesAndPreservesIDsAndOrder(t *testing.T) {
	r := newActivityResume(t)
	require.NoError(t, r.ReplaceActivities([]ActivityParams{activityParams(valueobjects.ID{}, "First"), activityParams(valueobjects.ID{}, "Second")}))
	items := r.Activities()
	require.Len(t, items, 2)
	first, second := items[0].ID(), items[1].ID()
	require.False(t, first.IsZero())
	require.False(t, second.IsZero())

	require.NoError(t, r.ReplaceActivities([]ActivityParams{activityParams(second, "Updated second"), activityParams(first, "Updated first")}))
	items = r.Activities()
	require.Equal(t, second, items[0].ID())
	require.Equal(t, first, items[1].ID())
	require.Equal(t, "Updated second", items[0].Title())
}

func TestResumeReplaceActivitiesRejectsUnknownAndDuplicateIDs(t *testing.T) {
	r := newActivityResume(t)
	require.NoError(t, r.ReplaceActivities([]ActivityParams{activityParams(valueobjects.ID{}, "Existing")}))
	existing := r.Activities()[0].ID()

	err := r.ReplaceActivities([]ActivityParams{activityParams(valueobjects.NewID(), "Unknown")})
	require.ErrorIs(t, err, ErrResumeInvariant)

	err = r.ReplaceActivities([]ActivityParams{activityParams(existing, "A"), activityParams(existing, "B")})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeReplaceActivitiesRemovesMissingItems(t *testing.T) {
	r := newActivityResume(t)
	require.NoError(t, r.ReplaceActivities([]ActivityParams{activityParams(valueobjects.ID{}, "Keep"), activityParams(valueobjects.ID{}, "Remove")}))
	keep := r.Activities()[0].ID()
	require.NoError(t, r.ReplaceActivities([]ActivityParams{activityParams(keep, "Keep updated")}))
	require.Len(t, r.Activities(), 1)
	require.Equal(t, keep, r.Activities()[0].ID())
}

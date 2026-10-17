package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestActivityValidation(t *testing.T) {
	date, err := ParseYearMonth("2025-01")
	require.NoError(t, err)
	id, resumeID := sharedvalueobjects.NewID(), sharedvalueobjects.NewID()
	item, err := newActivity(id, resumeID, ActivityParams{ID: id, Type: ActivityTypeOpenSource, Title: "Open Source", Date: &date})
	require.NoError(t, err)
	require.Equal(t, ActivityTypeOpenSource, item.Type())
	require.Equal(t, "2025-01", item.Date().String())

	_, err = newActivity(id, resumeID, ActivityParams{ID: id, Type: ActivityType("invalid"), Title: "Open Source"})
	require.Error(t, err)
}

func activityParams(id sharedvalueobjects.ID, title string) ActivityParams {
	return ActivityParams{ID: id, Type: ActivityTypeOpenSource, Title: title}
}

func TestResumeReplaceActivitiesCreatesAndPreservesIDsAndOrder(t *testing.T) {
	r := newTestResume(t)
	require.NoError(t, r.ReplaceActivities([]ActivityParams{activityParams(sharedvalueobjects.ID{}, "First"), activityParams(sharedvalueobjects.ID{}, "Second")}))
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
	r := newTestResume(t)
	require.NoError(t, r.ReplaceActivities([]ActivityParams{activityParams(sharedvalueobjects.ID{}, "Existing")}))
	existing := r.Activities()[0].ID()

	err := r.ReplaceActivities([]ActivityParams{activityParams(sharedvalueobjects.NewID(), "Unknown")})
	require.ErrorIs(t, err, ErrResumeInvariant)

	err = r.ReplaceActivities([]ActivityParams{activityParams(existing, "A"), activityParams(existing, "B")})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeReplaceActivitiesRemovesMissingItems(t *testing.T) {
	r := newTestResume(t)
	require.NoError(t, r.ReplaceActivities([]ActivityParams{activityParams(sharedvalueobjects.ID{}, "Keep"), activityParams(sharedvalueobjects.ID{}, "Remove")}))
	keep := r.Activities()[0].ID()
	require.NoError(t, r.ReplaceActivities([]ActivityParams{activityParams(keep, "Keep updated")}))
	require.Len(t, r.Activities(), 1)
	require.Equal(t, keep, r.Activities()[0].ID())
}

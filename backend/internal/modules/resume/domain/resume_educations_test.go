package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestEducationValidation(t *testing.T) {
	id := identity.NewID()
	resumeID := identity.NewID()
	endDate, err := ParseYear("2024")
	require.NoError(t, err)

	item, err := NewEducation(EducationParams{
		ID: id, ResumeID: resumeID, Location: "Paris", University: "University",
		Grade: EducationLevelMaster, ProgramName: "Computer Science", Faculty: "Engineering", EndDate: &endDate,
	})
	require.NoError(t, err)
	require.Equal(t, id, item.ID())
	require.Equal(t, resumeID, item.ResumeID())
	require.Equal(t, "2024", item.EndDate().String())

	_, err = NewEducation(EducationParams{ID: id, ResumeID: resumeID, Location: "", University: "University", Grade: EducationLevelMaster, ProgramName: "CS", Faculty: "Engineering"})
	require.Error(t, err)
}

func educationParams(r *Resume, id identity.ID, name string) *EducationParams {
	return &EducationParams{ID: id, ResumeID: r.ID(), Location: "Berlin", University: name, Grade: EducationLevelMaster, ProgramName: "Computer Science", Faculty: "Engineering"}
}

func TestResumeReplaceEducations(t *testing.T) {
	r := newTestResume(t)
	p := educationParams(r, identity.ID{}, "First")
	require.NoError(t, r.ReplaceEducations([]EducationParams{*p}))
	id := r.Educations()[0].ID()
	require.NotEqual(t, identity.ID{}, id)
	p = educationParams(r, id, "Updated")
	require.NoError(t, r.ReplaceEducations([]EducationParams{*p}))
	require.Equal(t, id, r.Educations()[0].ID())
	require.Equal(t, "Updated", r.Educations()[0].University())
	require.NoError(t, r.ReplaceEducations(nil))
	require.Empty(t, r.Educations())
}

func TestResumeReplaceEducationsRejectsUnknownID(t *testing.T) {
	r := newTestResume(t)
	err := r.ReplaceEducations([]EducationParams{*educationParams(r, identity.NewID(), "Unknown")})
	require.ErrorIs(t, err, ErrResumeInvariant)
}

func TestResumeReplaceEducationsRejectsDuplicateID(t *testing.T) {
	r := newTestResume(t)
	p := educationParams(r, identity.ID{}, "First")
	require.NoError(t, r.ReplaceEducations([]EducationParams{*p}))
	id := r.Educations()[0].ID()
	a := educationParams(r, id, "A")
	b := educationParams(r, id, "B")
	require.ErrorIs(t, r.ReplaceEducations([]EducationParams{*a, *b}), ErrInvalidEntity)
}

func TestResumeReplaceEducationsPreservesInputOrder(t *testing.T) {
	r := newTestResume(t)
	a := educationParams(r, identity.ID{}, "A")
	b := educationParams(r, identity.ID{}, "B")
	require.NoError(t, r.ReplaceEducations([]EducationParams{*a, *b}))
	items := r.Educations()
	require.Len(t, items, 2)
	require.Equal(t, "A", items[0].University())
	require.Equal(t, "B", items[1].University())
}

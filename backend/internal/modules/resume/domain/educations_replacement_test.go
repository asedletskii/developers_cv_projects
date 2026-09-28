package domain

import (
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
	"testing"

	"github.com/stretchr/testify/require"
)

func newEducationTestResume(t *testing.T) *Resume {
	t.Helper()
	r, err := NewResume(ResumeParams{ID: valueobjects.NewID(), UserID: useridentity.NewID(), Title: "Resume", Language: enums.ResumeLanguageEN})
	require.NoError(t, err)
	return r
}
func educationParams(r *Resume, id valueobjects.ID, name string) *EducationParams {
	return &EducationParams{ID: id, ResumeID: r.ID(), Location: "Berlin", University: name, Grade: enums.EducationLevelMaster, ProgramName: "Computer Science", Faculty: "Engineering"}
}
func TestResumeReplaceEducations(t *testing.T) {
	r := newEducationTestResume(t)
	p := educationParams(r, valueobjects.ID{}, "First")
	require.NoError(t, r.ReplaceEducations([]EducationParams{*p}))
	id := r.Educations()[0].ID()
	require.NotEqual(t, valueobjects.ID{}, id)
	p = educationParams(r, id, "Updated")
	require.NoError(t, r.ReplaceEducations([]EducationParams{*p}))
	require.Equal(t, id, r.Educations()[0].ID())
	require.Equal(t, "Updated", r.Educations()[0].University())
	require.NoError(t, r.ReplaceEducations(nil))
	require.Empty(t, r.Educations())
}
func TestResumeReplaceEducationsRejectsUnknownID(t *testing.T) {
	r := newEducationTestResume(t)
	err := r.ReplaceEducations([]EducationParams{*educationParams(r, valueobjects.NewID(), "Unknown")})
	require.ErrorIs(t, err, ErrResumeInvariant)
}
func TestResumeReplaceEducationsRejectsDuplicateID(t *testing.T) {
	r := newEducationTestResume(t)
	p := educationParams(r, valueobjects.ID{}, "First")
	require.NoError(t, r.ReplaceEducations([]EducationParams{*p}))
	id := r.Educations()[0].ID()
	a := educationParams(r, id, "A")
	b := educationParams(r, id, "B")
	require.ErrorIs(t, r.ReplaceEducations([]EducationParams{*a, *b}), ErrInvalidEntity)
}
func TestResumeReplaceEducationsPreservesInputOrder(t *testing.T) {
	r := newEducationTestResume(t)
	a := educationParams(r, valueobjects.ID{}, "A")
	b := educationParams(r, valueobjects.ID{}, "B")
	require.NoError(t, r.ReplaceEducations([]EducationParams{*a, *b}))
	items := r.Educations()
	require.Len(t, items, 2)
	require.Equal(t, "A", items[0].University())
	require.Equal(t, "B", items[1].University())
}

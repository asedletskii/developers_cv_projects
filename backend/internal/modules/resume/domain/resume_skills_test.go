package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestSkillsAcceptsArbitraryHardAndSoftSkillValues(t *testing.T) {
	hard := []string{"architecture", "testing"}
	soft := []string{"communication", "teamwork"}

	skills, err := NewSkills(SkillsParams{
		ID:       sharedvalueobjects.NewID(),
		ResumeID: sharedvalueobjects.NewID(),
		Hard:     hard,
		Soft:     soft,
	})
	require.NoError(t, err)
	require.Equal(t, hard, skills.HardSkills())
	require.Equal(t, soft, skills.SoftSkills())
}

func TestSkillsRejectsDuplicatesWithinEachCategory(t *testing.T) {
	_, err := NewSkills(SkillsParams{
		ID:       sharedvalueobjects.NewID(),
		ResumeID: sharedvalueobjects.NewID(),
		Hard:     []string{"testing", "testing"},
	})
	require.ErrorIs(t, err, ErrInvalidEntity)

	_, err = NewSkills(SkillsParams{
		ID:       sharedvalueobjects.NewID(),
		ResumeID: sharedvalueobjects.NewID(),
		Soft:     []string{"teamwork", "teamwork"},
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeReplaceSkillsCreatesAndReplaces(t *testing.T) {
	resume := newTestResume(t)
	require.NoError(t, resume.ReplaceSkills(SkillsParams{
		Hard: []string{"testing"},
		Soft: []string{"teamwork"},
	}))
	firstID := resume.Skills().ID()
	require.False(t, firstID.IsZero())
	require.Equal(t, resume.ID(), resume.Skills().ResumeID())

	require.NoError(t, resume.ReplaceSkills(SkillsParams{
		Hard: []string{"architecture"},
		Soft: []string{"communication"},
	}))
	require.Equal(t, firstID, resume.Skills().ID())
	require.Equal(t, []string{"architecture"}, resume.Skills().HardSkills())
	require.Equal(t, []string{"communication"}, resume.Skills().SoftSkills())
}

func TestResumeReplaceSkillsRejectsUnknownID(t *testing.T) {
	resume := newTestResume(t)
	require.NoError(t, resume.ReplaceSkills(SkillsParams{Hard: []string{"testing"}}))
	err := resume.ReplaceSkills(SkillsParams{ID: sharedvalueobjects.NewID(), Hard: []string{"architecture"}})
	require.ErrorIs(t, err, ErrResumeInvariant)
}

func TestResumeReplaceSkillsAcceptsArbitraryReplacement(t *testing.T) {
	resume := newTestResume(t)
	err := resume.ReplaceSkills(SkillsParams{Hard: []string{"Rust 2026", "Kubernetes", "whatever frontend sends"}})
	require.NoError(t, err)
	require.Equal(t, []string{"Rust 2026", "Kubernetes", "whatever frontend sends"}, resume.Skills().HardSkills())
}

func TestResumeReplaceSkillsPreservesExplicitExistingID(t *testing.T) {
	resume := newTestResume(t)
	require.NoError(t, resume.ReplaceSkills(SkillsParams{Hard: []string{"testing"}}))
	id := resume.Skills().ID()
	require.NoError(t, resume.ReplaceSkills(SkillsParams{ID: id, Hard: []string{"architecture"}}))
	require.Equal(t, id, resume.Skills().ID())
}

package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestSkillsKeepsHardAndSoftSkillTypesSeparate(t *testing.T) {
	hard := []HardSkill{HardSkillArchitecture, HardSkillTesting}
	soft := []SoftSkill{SoftSkillCommunication, SoftSkillTeamwork}

	skills, err := NewSkills(SkillsParams{
		ID:       identity.NewID(),
		ResumeID: identity.NewID(),
		Hard:     hard,
		Soft:     soft,
	})
	require.NoError(t, err)
	require.Equal(t, hard, skills.HardSkills())
	require.Equal(t, soft, skills.SoftSkills())
}

func TestSkillsRejectsDuplicatesWithinEachCategory(t *testing.T) {
	_, err := NewSkills(SkillsParams{
		ID:       identity.NewID(),
		ResumeID: identity.NewID(),
		Hard:     []HardSkill{HardSkillTesting, HardSkillTesting},
	})
	require.ErrorIs(t, err, ErrInvalidEntity)

	_, err = NewSkills(SkillsParams{
		ID:       identity.NewID(),
		ResumeID: identity.NewID(),
		Soft:     []SoftSkill{SoftSkillTeamwork, SoftSkillTeamwork},
	})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeReplaceSkillsCreatesAndReplaces(t *testing.T) {
	resume := newTestResume(t)
	require.NoError(t, resume.ReplaceSkills(SkillsParams{
		Hard: []HardSkill{HardSkillTesting},
		Soft: []SoftSkill{SoftSkillTeamwork},
	}))
	firstID := resume.Skills().ID()
	require.False(t, firstID.IsZero())
	require.Equal(t, resume.ID(), resume.Skills().ResumeID())

	require.NoError(t, resume.ReplaceSkills(SkillsParams{
		Hard: []HardSkill{HardSkillArchitecture},
		Soft: []SoftSkill{SoftSkillCommunication},
	}))
	require.Equal(t, firstID, resume.Skills().ID())
	require.Equal(t, []HardSkill{HardSkillArchitecture}, resume.Skills().HardSkills())
	require.Equal(t, []SoftSkill{SoftSkillCommunication}, resume.Skills().SoftSkills())
}

func TestResumeReplaceSkillsRejectsUnknownID(t *testing.T) {
	resume := newTestResume(t)
	require.NoError(t, resume.ReplaceSkills(SkillsParams{Hard: []HardSkill{HardSkillTesting}}))
	err := resume.ReplaceSkills(SkillsParams{ID: identity.NewID(), Hard: []HardSkill{HardSkillArchitecture}})
	require.ErrorIs(t, err, ErrResumeInvariant)
}

func TestResumeReplaceSkillsRejectsInvalidReplacement(t *testing.T) {
	resume := newTestResume(t)
	err := resume.ReplaceSkills(SkillsParams{Hard: []HardSkill{HardSkill("invalid")}})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeReplaceSkillsPreservesExplicitExistingID(t *testing.T) {
	resume := newTestResume(t)
	require.NoError(t, resume.ReplaceSkills(SkillsParams{Hard: []HardSkill{HardSkillTesting}}))
	id := resume.Skills().ID()
	require.NoError(t, resume.ReplaceSkills(SkillsParams{ID: id, Hard: []HardSkill{HardSkillArchitecture}}))
	require.Equal(t, id, resume.Skills().ID())
}

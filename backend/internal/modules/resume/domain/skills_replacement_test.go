package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

func TestResumeReplaceSkillsCreatesAndReplaces(t *testing.T) {
	resume := newTestResumeForAggregate(t)
	require.NoError(t, resume.ReplaceSkills(SkillsParams{
		Hard: []enums.HardSkill{enums.HardSkillTesting},
		Soft: []enums.SoftSkill{enums.SoftSkillTeamwork},
	}))
	firstID := resume.Skills().ID()
	require.False(t, firstID.IsZero())
	require.Equal(t, resume.ID(), resume.Skills().ResumeID())

	require.NoError(t, resume.ReplaceSkills(SkillsParams{
		Hard: []enums.HardSkill{enums.HardSkillArchitecture},
		Soft: []enums.SoftSkill{enums.SoftSkillCommunication},
	}))
	require.Equal(t, firstID, resume.Skills().ID())
	require.Equal(t, []enums.HardSkill{enums.HardSkillArchitecture}, resume.Skills().HardSkills())
	require.Equal(t, []enums.SoftSkill{enums.SoftSkillCommunication}, resume.Skills().SoftSkills())
}

func TestResumeReplaceSkillsRejectsUnknownID(t *testing.T) {
	resume := newTestResumeForAggregate(t)
	require.NoError(t, resume.ReplaceSkills(SkillsParams{Hard: []enums.HardSkill{enums.HardSkillTesting}}))
	err := resume.ReplaceSkills(SkillsParams{ID: valueobjects.NewID(), Hard: []enums.HardSkill{enums.HardSkillArchitecture}})
	require.ErrorIs(t, err, ErrResumeInvariant)
}

func TestResumeReplaceSkillsRejectsInvalidReplacement(t *testing.T) {
	resume := newTestResumeForAggregate(t)
	err := resume.ReplaceSkills(SkillsParams{Hard: []enums.HardSkill{enums.HardSkill("invalid")}})
	require.ErrorIs(t, err, ErrInvalidEntity)
}

func TestResumeReplaceSkillsPreservesExplicitExistingID(t *testing.T) {
	resume := newTestResumeForAggregate(t)
	require.NoError(t, resume.ReplaceSkills(SkillsParams{Hard: []enums.HardSkill{enums.HardSkillTesting}}))
	id := resume.Skills().ID()
	require.NoError(t, resume.ReplaceSkills(SkillsParams{ID: id, Hard: []enums.HardSkill{enums.HardSkillArchitecture}}))
	require.Equal(t, id, resume.Skills().ID())
}

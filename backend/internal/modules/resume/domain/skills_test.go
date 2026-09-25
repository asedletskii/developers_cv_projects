package domain_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

func TestSkillsKeepsHardAndSoftSkillTypesSeparate(t *testing.T) {
	hard := []enums.HardSkill{enums.HardSkillArchitecture, enums.HardSkillTesting}
	soft := []enums.SoftSkill{enums.SoftSkillCommunication, enums.SoftSkillTeamwork}

	skills, err := domain.NewSkills(domain.SkillsParams{
		ID:       valueobjects.NewID(),
		ResumeID: valueobjects.NewID(),
		Hard:     hard,
		Soft:     soft,
	})
	require.NoError(t, err)
	require.Equal(t, hard, skills.HardSkills())
	require.Equal(t, soft, skills.SoftSkills())
}

func TestSkillsRejectsDuplicatesWithinEachCategory(t *testing.T) {
	_, err := domain.NewSkills(domain.SkillsParams{
		ID:       valueobjects.NewID(),
		ResumeID: valueobjects.NewID(),
		Hard:     []enums.HardSkill{enums.HardSkillTesting, enums.HardSkillTesting},
	})
	require.ErrorIs(t, err, domain.ErrInvalidEntity)

	_, err = domain.NewSkills(domain.SkillsParams{
		ID:       valueobjects.NewID(),
		ResumeID: valueobjects.NewID(),
		Soft:     []enums.SoftSkill{enums.SoftSkillTeamwork, enums.SoftSkillTeamwork},
	})
	require.ErrorIs(t, err, domain.ErrInvalidEntity)
}

package mappers

import (
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func SkillsToDomain(
	model gormmodel.ResumeSkills,
	hardRows []gormmodel.ResumeSkillHard,
	softRows []gormmodel.ResumeSkillSoft,
) (*domain.Skills, error) {
	id, err := sharedvalueobjects.IDFromUUID(model.ID)
	if err != nil {
		return nil, err
	}
	hard := make([]string, len(hardRows))
	for i, row := range hardRows {
		hard[i] = row.Value
	}
	soft := make([]string, len(softRows))
	for i, row := range softRows {
		soft[i] = row.Value
	}
	resumeID, err := sharedvalueobjects.IDFromUUID(model.ResumeID)
	if err != nil {
		return nil, err
	}
	return domain.RestoreSkills(domain.SkillsSnapshot{ID: id, ResumeID: resumeID, Hard: hard, Soft: soft}), nil
}

func SkillsToDBModels(
	skills *domain.Skills,
) (gormmodel.ResumeSkills, []gormmodel.ResumeSkillHard, []gormmodel.ResumeSkillSoft) {
	m := gormmodel.ResumeSkills{ID: skills.ID().UUID(), ResumeID: skills.ResumeID().UUID()}
	hard := make([]gormmodel.ResumeSkillHard, 0, len(skills.HardSkills()))
	for i, value := range skills.HardSkills() {
		hard = append(hard, gormmodel.ResumeSkillHard{SkillsID: m.ID, Value: value, Position: i})
	}
	soft := make([]gormmodel.ResumeSkillSoft, 0, len(skills.SoftSkills()))
	for i, value := range skills.SoftSkills() {
		soft = append(soft, gormmodel.ResumeSkillSoft{SkillsID: m.ID, Value: value, Position: i})
	}
	return m, hard, soft
}

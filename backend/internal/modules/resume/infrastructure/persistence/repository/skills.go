package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

var _ commands.SkillsRepository = (*ResumeRepository)(nil)

func (r *ResumeRepository) CreateSkills(ctx context.Context, a *domain.Skills, uid useridentity.ID) error {
	tx, e := r.owned(ctx, a.ResumeID(), uid)
	if e != nil {
		return e
	}
	var n int64
	tx.Model(&gormmodel.ResumeSkillsV2{}).Where("resume_id=?", a.ResumeID().UUID()).Count(&n)
	if n > 0 {
		return domain.ErrDuplicateEntity
	}
	m := gormmodel.ResumeSkillsV2{ID: a.ID().UUID(), ResumeID: a.ResumeID().UUID()}
	if e = tx.Create(&m).Error; e != nil {
		return e
	}
	for i, v := range a.HardSkills() {
		if e = tx.Create(&gormmodel.ResumeSkillHardV2{SkillsID: m.ID, Value: string(v), Position: i}).Error; e != nil {
			return e
		}
	}
	for i, v := range a.SoftSkills() {
		if e = tx.Create(&gormmodel.ResumeSkillSoftV2{SkillsID: m.ID, Value: string(v), Position: i}).Error; e != nil {
			return e
		}
	}
	return nil
}
func (r *ResumeRepository) GetSkills(ctx context.Context, q commands.GetSkillsQuery) (*domain.Skills, error) {
	tx, e := r.owned(ctx, q.ResumeID, q.UserID)
	if e != nil {
		return nil, e
	}
	var m gormmodel.ResumeSkillsV2
	if e = tx.Where("resume_id=?", q.ResumeID.UUID()).First(&m).Error; e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, e
	}
	var h []gormmodel.ResumeSkillHardV2
	var s []gormmodel.ResumeSkillSoftV2
	tx.Where("skills_id=?", m.ID).Order("position").Find(&h)
	tx.Where("skills_id=?", m.ID).Order("position").Find(&s)
	hp := make([]enums.HardSkill, len(h))
	sp := make([]enums.SoftSkill, len(s))
	for i, v := range h {
		hp[i] = enums.HardSkill(v.Value)
	}
	for i, v := range s {
		sp[i] = enums.SoftSkill(v.Value)
	}
	id, err := valueobjects.IDFromUUID(m.ID)
	if err != nil {
		return nil, err
	}
	rid, err := valueobjects.IDFromUUID(m.ResumeID)
	if err != nil {
		return nil, err
	}
	return domain.NewSkills(domain.SkillsParams{ID: id, ResumeID: rid, Hard: hp, Soft: sp})
}
func (r *ResumeRepository) UpdateSkills(ctx context.Context, a *domain.Skills, uid useridentity.ID) error {
	tx, e := r.owned(ctx, a.ResumeID(), uid)
	if e != nil {
		return e
	}
	var m gormmodel.ResumeSkillsV2
	if e = tx.Where("id=? AND resume_id=?", a.ID().UUID(), a.ResumeID().UUID()).First(&m).Error; e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return domain.ErrEntityNotFound
		}
		return e
	}
	if e = tx.Where("skills_id=?", m.ID).Delete(&gormmodel.ResumeSkillHardV2{}).Error; e != nil {
		return e
	}
	if e = tx.Where("skills_id=?", m.ID).Delete(&gormmodel.ResumeSkillSoftV2{}).Error; e != nil {
		return e
	}
	for i, v := range a.HardSkills() {
		if e = tx.Create(&gormmodel.ResumeSkillHardV2{SkillsID: m.ID, Value: string(v), Position: i}).Error; e != nil {
			return e
		}
	}
	for i, v := range a.SoftSkills() {
		if e = tx.Create(&gormmodel.ResumeSkillSoftV2{SkillsID: m.ID, Value: string(v), Position: i}).Error; e != nil {
			return e
		}
	}
	return nil
}

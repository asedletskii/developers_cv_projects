package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

func (r *SkillsRepository) CreateSkills(ctx context.Context, c ports.CreateSkillsRepositoryCommand) error {
	if c.Skills == nil {
		return fmt.Errorf("skills is nil")
	}
	a, uid := c.Skills, c.UserID
	tx, e := r.ownedInTransaction(ctx, a.ResumeID(), uid)
	if e != nil {
		return e
	}
	m := gormmodel.ResumeSkills{ID: a.ID().UUID(), ResumeID: a.ResumeID().UUID()}
	if e = tx.Create(&m).Error; e != nil {
		return mapDBError(e)
	}
	for i, v := range a.HardSkills() {
		if e = tx.Create(&gormmodel.ResumeSkillHard{SkillsID: m.ID, Value: string(v), Position: i}).Error; e != nil {
			return e
		}
	}
	for i, v := range a.SoftSkills() {
		if e = tx.Create(&gormmodel.ResumeSkillSoft{SkillsID: m.ID, Value: string(v), Position: i}).Error; e != nil {
			return e
		}
	}
	return nil
}
func (r *SkillsRepository) GetSkills(ctx context.Context, q ports.GetSkillsRepositoryCommand) (*domain.Skills, error) {
	tx, e := r.owned(ctx, q.ResumeID, q.UserID)
	if e != nil {
		return nil, e
	}
	var m gormmodel.ResumeSkills
	if e = tx.Where("resume_id=?", q.ResumeID.UUID()).First(&m).Error; e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, e
	}
	var h []gormmodel.ResumeSkillHard
	var s []gormmodel.ResumeSkillSoft
	if e = tx.Where("skills_id=?", m.ID).Order("position").Find(&h).Error; e != nil {
		return nil, e
	}
	if e = tx.Where("skills_id=?", m.ID).Order("position").Find(&s).Error; e != nil {
		return nil, e
	}
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
func (r *SkillsRepository) UpdateSkills(ctx context.Context, c ports.UpdateSkillsRepositoryCommand) error {
	if c.Hard == nil || c.Soft == nil {
		return fmt.Errorf("skills is nil")
	}
	tx, e := r.owned(ctx, c.ResumeID, c.UserID)
	if e != nil {
		return e
	}
	var m gormmodel.ResumeSkills
	if e = tx.Where("resume_id=?", c.ResumeID.UUID()).First(&m).Error; e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return domain.ErrEntityNotFound
		}
		return e
	}
	if e = tx.Where("skills_id=?", m.ID).Delete(&gormmodel.ResumeSkillHard{}).Error; e != nil {
		return e
	}
	if e = tx.Where("skills_id=?", m.ID).Delete(&gormmodel.ResumeSkillSoft{}).Error; e != nil {
		return e
	}
	for i, v := range c.Hard {
		if e = tx.Create(&gormmodel.ResumeSkillHard{SkillsID: m.ID, Value: string(v), Position: i}).Error; e != nil {
			return e
		}
	}
	for i, v := range c.Soft {
		if e = tx.Create(&gormmodel.ResumeSkillSoft{SkillsID: m.ID, Value: string(v), Position: i}).Error; e != nil {
			return e
		}
	}
	return nil
}

package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

var _ commands.CreateTechnologyRepository = (*TechnologiesRepository)(nil)
var _ commands.GetTechnologyRepository = (*TechnologiesRepository)(nil)
var _ commands.UpdateTechnologyRepository = (*TechnologiesRepository)(nil)
var _ commands.DeleteTechnologyRepository = (*TechnologiesRepository)(nil)

func (r *TechnologiesRepository) CreateTechnology(ctx context.Context, c commands.CreateTechnologyRepositoryCommand) error {
	tx, err := r.owned(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	m, err := technologyModel(c.Technology, c.ResumeID)
	if err != nil {
		return err
	}
	if err = tx.Create(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrDuplicateEntity
		}
		return err
	}
	return nil
}

func (r *TechnologiesRepository) GetTechnology(ctx context.Context, q commands.GetTechnologyCommand) (domain.ResumeTechnologyEntry, error) {
	tx, err := r.owned(ctx, q.ResumeID, q.UserID)
	if err != nil {
		return domain.ResumeTechnologyEntry{}, err
	}
	var m gormmodel.ResumeTechnology
	if err = tx.Where("id=? AND resume_id=?", q.TechnologyID.UUID(), q.ResumeID.UUID()).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ResumeTechnologyEntry{}, domain.ErrEntityNotFound
		}
		return domain.ResumeTechnologyEntry{}, err
	}
	return mapTechnology(m)
}

func (r *TechnologiesRepository) UpdateTechnology(ctx context.Context, c commands.UpdateTechnologyRepositoryCommand) error {
	tx, err := r.owned(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	m, err := technologyModel(c.Technology, c.ResumeID)
	if err != nil {
		return err
	}
	res := tx.Model(&gormmodel.ResumeTechnology{}).Where("id=? AND resume_id=?", c.Technology.ID().UUID(), c.ResumeID.UUID()).Updates(map[string]any{"name": m.Name, "position": m.Position, "grade": m.Grade})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
}

func (r *TechnologiesRepository) DeleteTechnology(ctx context.Context, c commands.DeleteTechnologyCommand) error {
	tx, err := r.owned(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	res := tx.Where("id=? AND resume_id=?", c.TechnologyID.UUID(), c.ResumeID.UUID()).Delete(&gormmodel.ResumeTechnology{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
}

func technologyModel(x domain.ResumeTechnologyEntry, rid valueobjects.ID) (gormmodel.ResumeTechnology, error) {
	if x.ID().IsZero() {
		return gormmodel.ResumeTechnology{}, fmt.Errorf("technology id is required")
	}
	if x.Name() == "" {
		return gormmodel.ResumeTechnology{}, fmt.Errorf("technology name is required")
	}
	m := gormmodel.ResumeTechnology{ID: x.ID().UUID(), ResumeID: rid.UUID(), Name: x.Name(), Position: x.Position()}
	if g := x.Grade(); g != nil {
		v := string(*g)
		m.Grade = &v
	}
	return m, nil
}

func mapTechnology(m gormmodel.ResumeTechnology) (domain.ResumeTechnologyEntry, error) {
	id, err := valueobjects.IDFromUUID(m.ID)
	if err != nil {
		return domain.ResumeTechnologyEntry{}, err
	}
	var g *enums.TechnologyGrade
	if m.Grade != nil {
		x := enums.TechnologyGrade(*m.Grade)
		if !x.IsValid() {
			return domain.ResumeTechnologyEntry{}, fmt.Errorf("invalid technology grade")
		}
		g = &x
	}
	return domain.NewResumeTechnologyEntry(domain.ResumeTechnologyEntryParams{ID: id, Name: m.Name, Grade: g, Position: m.Position})
}

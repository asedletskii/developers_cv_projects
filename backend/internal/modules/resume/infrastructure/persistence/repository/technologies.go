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
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

var _ commands.TechnologiesRepository = (*TechnologiesRepository)(nil)

func (r *TechnologiesRepository) CreateTechnology(ctx context.Context, in commands.TechnologyInput, rid valueobjects.ID, uid useridentity.ID) error {
	tx, e := r.owned(ctx, rid, uid)
	if e != nil {
		return e
	}
	m, e := technologyModel(in, rid)
	if e != nil {
		return e
	}
	if e = tx.Create(&m).Error; e != nil {
		if errors.Is(e, gorm.ErrDuplicatedKey) {
			return domain.ErrDuplicateEntity
		}
		return e
	}
	return nil
}
func (r *TechnologiesRepository) ListTechnologies(ctx context.Context, q commands.ListTechnologiesQuery) ([]domain.ResumeTechnologyEntry, error) {
	tx, e := r.owned(ctx, q.ResumeID, q.UserID)
	if e != nil {
		return nil, e
	}
	var ms []gormmodel.ResumeTechnology
	if e = tx.Where("resume_id=?", q.ResumeID.UUID()).Order("position").Find(&ms).Error; e != nil {
		return nil, e
	}
	out := make([]domain.ResumeTechnologyEntry, 0, len(ms))
	for _, m := range ms {
		x, err := mapTechnology(m)
		if err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, nil
}
func (r *TechnologiesRepository) GetTechnology(ctx context.Context, q commands.GetTechnologyQuery) (domain.ResumeTechnologyEntry, error) {
	tx, e := r.owned(ctx, q.ResumeID, q.UserID)
	if e != nil {
		return domain.ResumeTechnologyEntry{}, e
	}
	var m gormmodel.ResumeTechnology
	if e = tx.Where("id=? AND resume_id=?", q.TechnologyID.UUID(), q.ResumeID.UUID()).First(&m).Error; e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return domain.ResumeTechnologyEntry{}, domain.ErrEntityNotFound
		}
		return domain.ResumeTechnologyEntry{}, e
	}
	return mapTechnology(m)
}
func (r *TechnologiesRepository) UpdateTechnology(ctx context.Context, c commands.UpdateTechnologyCommand) error {
	tx, e := r.owned(ctx, c.ResumeID, c.UserID)
	if e != nil {
		return e
	}
	m, e := technologyModel(c.Input, c.ResumeID)
	if e != nil {
		return e
	}
	res := tx.Model(&gormmodel.ResumeTechnology{}).Where("id=? AND resume_id=?", c.TechnologyID.UUID(), c.ResumeID.UUID()).Updates(map[string]any{"name": m.Name, "position": m.Position, "grade": m.Grade})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
}
func (r *TechnologiesRepository) DeleteTechnology(ctx context.Context, c commands.DeleteTechnologyCommand) error {
	tx, e := r.owned(ctx, c.ResumeID, c.UserID)
	if e != nil {
		return e
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
func technologyModel(in commands.TechnologyInput, rid valueobjects.ID) (gormmodel.ResumeTechnology, error) {
	if in.ID.IsZero() {
		return gormmodel.ResumeTechnology{}, fmt.Errorf("technology id is required")
	}
	entry, err := domain.NewResumeTechnologyEntry(domain.ResumeTechnologyEntryParams{
		ID: in.ID, Name: in.Name, Grade: in.Grade, Position: in.Position,
	})
	if err != nil {
		return gormmodel.ResumeTechnology{}, err
	}
	m := gormmodel.ResumeTechnology{ID: in.ID.UUID(), ResumeID: rid.UUID(), Name: entry.Name(), Position: entry.Position()}
	if in.Grade != nil {
		x := string(*in.Grade)
		m.Grade = &x
	}
	return m, nil
}
func mapTechnology(m gormmodel.ResumeTechnology) (domain.ResumeTechnologyEntry, error) {
	id, e := valueobjects.IDFromUUID(m.ID)
	if e != nil {
		return domain.ResumeTechnologyEntry{}, e
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

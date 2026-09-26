package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

func (r *TechnologiesRepository) CreateTechnology(ctx context.Context, c ports.CreateTechnologyRepositoryCommand) error {
	tx, err := r.ownedForWrite(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	m, err := technologyModel(c.Technology, c.ResumeID)
	if err != nil {
		return err
	}
	m.Position, err = nextSortOrder(tx, m.TableName(), c.ResumeID.UUID())
	if err != nil {
		return err
	}
	if err = tx.Create(&m).Error; err != nil {
		return mapDBError(err)
	}
	return nil
}

func (r *TechnologiesRepository) GetTechnology(ctx context.Context, q ports.GetTechnologyRepositoryCommand) (domain.ResumeTechnologyEntry, error) {
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

func (r *TechnologiesRepository) UpdateTechnology(ctx context.Context, c ports.UpdateTechnologyRepositoryCommand) error {
	tx, err := r.owned(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	m, err := technologyModel(c.Technology, c.ResumeID)
	if err != nil {
		return err
	}
	result := tx.Model(&gormmodel.ResumeTechnology{}).
		Where("id=? AND resume_id=?", c.Technology.ID().UUID(), c.ResumeID.UUID()).
		Updates(map[string]any{"name": m.Name, "grade": m.Grade})
	if result.Error != nil {
		return mapDBError(result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
}

func (r *TechnologiesRepository) ReorderTechnologies(ctx context.Context, c ports.ReorderTechnologiesRepositoryCommand) error {
	tx, err := r.ownedForWrite(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	ids := make([]uuid.UUID, len(c.TechnologyIDs))
	for i, id := range c.TechnologyIDs {
		if id.IsZero() {
			return fmt.Errorf("technology id is required: %w", domain.ErrInvalidEntity)
		}
		ids[i] = id.UUID()
	}
	return reorderPositions(tx, "resume_technologies", c.ResumeID.UUID(), ids)
}

func (r *TechnologiesRepository) DeleteTechnology(ctx context.Context, c ports.DeleteTechnologyRepositoryCommand) error {
	tx, err := r.ownedForWrite(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	result := tx.Where("id=? AND resume_id=?", c.TechnologyID.UUID(), c.ResumeID.UUID()).Delete(&gormmodel.ResumeTechnology{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
}

func technologyModel(x domain.ResumeTechnologyEntry, rid valueobjects.ID) (gormmodel.ResumeTechnology, error) {
	if x.ID().IsZero() {
		return gormmodel.ResumeTechnology{}, fmt.Errorf("technology id is required")
	}
	m := gormmodel.ResumeTechnology{ID: x.ID().UUID(), ResumeID: rid.UUID(), Name: x.Name(), Position: x.Position()}
	if g := x.Grade(); g != nil {
		if !g.IsValid() {
			return gormmodel.ResumeTechnology{}, fmt.Errorf("invalid technology grade")
		}
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
	var grade *enums.TechnologyGrade
	if m.Grade != nil {
		value := enums.TechnologyGrade(*m.Grade)
		if !value.IsValid() {
			return domain.ResumeTechnologyEntry{}, domain.ErrInvalidEntity
		}
		grade = &value
	}
	return domain.NewResumeTechnologyEntry(domain.ResumeTechnologyEntryParams{ID: id, Name: m.Name, Grade: grade, Position: m.Position})
}

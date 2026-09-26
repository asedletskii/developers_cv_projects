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

func (r *LanguagesRepository) CreateLanguage(ctx context.Context, c ports.CreateLanguageRepositoryCommand) error {
	tx, err := r.ownedForWrite(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	m, err := languageModel(c.Language, c.ResumeID)
	if err != nil {
		return err
	}
	m.Position, err = nextPosition(tx, m.TableName(), c.ResumeID.UUID())
	if err != nil {
		return err
	}
	if err = tx.Create(&m).Error; err != nil {
		return mapDBError(err)
	}
	return nil
}

func (r *LanguagesRepository) GetLanguage(ctx context.Context, q ports.GetLanguageRepositoryCommand) (domain.ResumeLanguageEntry, error) {
	tx, err := r.owned(ctx, q.ResumeID, q.UserID)
	if err != nil {
		return domain.ResumeLanguageEntry{}, err
	}
	var m gormmodel.ResumeLanguage
	if err = tx.Where("id=? AND resume_id=?", q.LanguageID.UUID(), q.ResumeID.UUID()).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ResumeLanguageEntry{}, domain.ErrEntityNotFound
		}
		return domain.ResumeLanguageEntry{}, err
	}
	return mapLanguage(m)
}

func (r *LanguagesRepository) UpdateLanguage(ctx context.Context, c ports.UpdateLanguageRepositoryCommand) error {
	tx, err := r.owned(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	m, err := languageModel(c.Language, c.ResumeID)
	if err != nil {
		return err
	}
	result := tx.Model(&gormmodel.ResumeLanguage{}).
		Where("id=? AND resume_id=?", c.Language.ID().UUID(), c.ResumeID.UUID()).
		Updates(map[string]any{"language": m.Language, "grade": m.Grade})
	if result.Error != nil {
		return mapDBError(result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
}

func (r *LanguagesRepository) ReorderLanguages(ctx context.Context, c ports.ReorderLanguagesRepositoryCommand) error {
	tx, err := r.ownedForWrite(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	ids := make([]uuid.UUID, len(c.LanguageIDs))
	for i, id := range c.LanguageIDs {
		if id.IsZero() {
			return fmt.Errorf("language id is required: %w", domain.ErrInvalidEntity)
		}
		ids[i] = id.UUID()
	}
	return reorderPositions(tx, "resume_languages", c.ResumeID.UUID(), ids)
}

func (r *LanguagesRepository) DeleteLanguage(ctx context.Context, c ports.DeleteLanguageRepositoryCommand) error {
	tx, err := r.ownedForWrite(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	result := tx.Where("id=? AND resume_id=?", c.LanguageID.UUID(), c.ResumeID.UUID()).Delete(&gormmodel.ResumeLanguage{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
}

func languageModel(x domain.ResumeLanguageEntry, rid valueobjects.ID) (gormmodel.ResumeLanguage, error) {
	if x.ID().IsZero() {
		return gormmodel.ResumeLanguage{}, fmt.Errorf("language id is required")
	}
	m := gormmodel.ResumeLanguage{ID: x.ID().UUID(), ResumeID: rid.UUID(), Language: string(x.Language()), Position: x.Position()}
	if g := x.Grade(); g != nil {
		v := string(*g)
		m.Grade = &v
	}
	if !enums.SpokenLanguage(m.Language).IsValid() {
		return gormmodel.ResumeLanguage{}, fmt.Errorf("invalid language")
	}
	return m, nil
}

func mapLanguage(m gormmodel.ResumeLanguage) (domain.ResumeLanguageEntry, error) {
	id, err := valueobjects.IDFromUUID(m.ID)
	if err != nil {
		return domain.ResumeLanguageEntry{}, err
	}
	language := enums.SpokenLanguage(m.Language)
	if !language.IsValid() {
		return domain.ResumeLanguageEntry{}, domain.ErrInvalidEntity
	}
	var grade *enums.LanguageGrade
	if m.Grade != nil {
		value := enums.LanguageGrade(*m.Grade)
		if !value.IsValid() {
			return domain.ResumeLanguageEntry{}, domain.ErrInvalidEntity
		}
		grade = &value
	}
	return domain.NewResumeLanguageEntry(domain.ResumeLanguageEntryParams{ID: id, Language: language, Grade: grade, Position: m.Position})
}

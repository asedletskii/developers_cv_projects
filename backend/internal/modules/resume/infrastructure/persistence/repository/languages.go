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

var _ commands.CreateLanguageRepository = (*LanguagesRepository)(nil)
var _ commands.GetLanguageRepository = (*LanguagesRepository)(nil)
var _ commands.UpdateLanguageRepository = (*LanguagesRepository)(nil)
var _ commands.DeleteLanguageRepository = (*LanguagesRepository)(nil)

func (r *LanguagesRepository) CreateLanguage(ctx context.Context, c commands.CreateLanguageRepositoryCommand) error {
	tx, err := r.owned(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	m, err := languageModel(c.Language, c.ResumeID)
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

func (r *LanguagesRepository) GetLanguage(ctx context.Context, q commands.GetLanguageCommand) (domain.ResumeLanguageEntry, error) {
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

func (r *LanguagesRepository) UpdateLanguage(ctx context.Context, c commands.UpdateLanguageRepositoryCommand) error {
	tx, err := r.owned(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	m, err := languageModel(c.Language, c.ResumeID)
	if err != nil {
		return err
	}
	res := tx.Model(&gormmodel.ResumeLanguage{}).Where("id=? AND resume_id=?", c.Language.ID().UUID(), c.ResumeID.UUID()).Updates(map[string]any{"language": m.Language, "position": m.Position, "grade": m.Grade})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
}

func (r *LanguagesRepository) DeleteLanguage(ctx context.Context, c commands.DeleteLanguageCommand) error {
	tx, err := r.owned(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	res := tx.Where("id=? AND resume_id=?", c.LanguageID.UUID(), c.ResumeID.UUID()).Delete(&gormmodel.ResumeLanguage{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
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
	if !enums.LanguageType(m.Language).IsValid() {
		return gormmodel.ResumeLanguage{}, fmt.Errorf("invalid language")
	}
	return m, nil
}

func mapLanguage(m gormmodel.ResumeLanguage) (domain.ResumeLanguageEntry, error) {
	id, err := valueobjects.IDFromUUID(m.ID)
	if err != nil {
		return domain.ResumeLanguageEntry{}, err
	}
	var g *enums.LanguageGrade
	if m.Grade != nil {
		x := enums.LanguageGrade(*m.Grade)
		if !x.IsValid() {
			return domain.ResumeLanguageEntry{}, fmt.Errorf("invalid language grade")
		}
		g = &x
	}
	return domain.NewResumeLanguageEntry(domain.ResumeLanguageEntryParams{ID: id, Language: enums.LanguageType(m.Language), Grade: g, Position: m.Position})
}

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

var _ commands.LanguagesRepository = (*LanguagesRepository)(nil)

func (r *LanguagesRepository) CreateLanguage(ctx context.Context, in commands.LanguageInput, rid valueobjects.ID, uid useridentity.ID) error {
	tx, e := r.owned(ctx, rid, uid)
	if e != nil {
		return e
	}
	m, e := languageModel(in, rid)
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
func (r *LanguagesRepository) ListLanguages(ctx context.Context, q commands.ListLanguagesQuery) ([]domain.ResumeLanguageEntry, error) {
	tx, e := r.owned(ctx, q.ResumeID, q.UserID)
	if e != nil {
		return nil, e
	}
	var ms []gormmodel.ResumeLanguageV2
	if e = tx.Where("resume_id=?", q.ResumeID.UUID()).Order("position").Find(&ms).Error; e != nil {
		return nil, e
	}
	out := make([]domain.ResumeLanguageEntry, 0, len(ms))
	for _, m := range ms {
		x, err := mapLanguage(m)
		if err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, nil
}
func (r *LanguagesRepository) GetLanguage(ctx context.Context, q commands.GetLanguageQuery) (domain.ResumeLanguageEntry, error) {
	tx, e := r.owned(ctx, q.ResumeID, q.UserID)
	if e != nil {
		return domain.ResumeLanguageEntry{}, e
	}
	var m gormmodel.ResumeLanguageV2
	if e = tx.Where("id=? AND resume_id=?", q.LanguageID.UUID(), q.ResumeID.UUID()).First(&m).Error; e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return domain.ResumeLanguageEntry{}, domain.ErrEntityNotFound
		}
		return domain.ResumeLanguageEntry{}, e
	}
	return mapLanguage(m)
}
func (r *LanguagesRepository) UpdateLanguage(ctx context.Context, c commands.UpdateLanguageCommand) error {
	tx, e := r.owned(ctx, c.ResumeID, c.UserID)
	if e != nil {
		return e
	}
	m, e := languageModel(c.Input, c.ResumeID)
	if e != nil {
		return e
	}
	res := tx.Model(&gormmodel.ResumeLanguageV2{}).Where("id=? AND resume_id=?", c.LanguageID.UUID(), c.ResumeID.UUID()).Updates(map[string]any{"language": m.Language, "position": m.Position, "grade": m.Grade})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
}
func (r *LanguagesRepository) DeleteLanguage(ctx context.Context, c commands.DeleteLanguageCommand) error {
	tx, e := r.owned(ctx, c.ResumeID, c.UserID)
	if e != nil {
		return e
	}
	res := tx.Where("id=? AND resume_id=?", c.LanguageID.UUID(), c.ResumeID.UUID()).Delete(&gormmodel.ResumeLanguageV2{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
}
func languageModel(in commands.LanguageInput, rid valueobjects.ID) (gormmodel.ResumeLanguageV2, error) {
	if in.ID.IsZero() {
		return gormmodel.ResumeLanguageV2{}, fmt.Errorf("language id is required")
	}
	if _, err := domain.NewResumeLanguageEntry(domain.ResumeLanguageEntryParams{
		ID: in.ID, Language: in.Language, Grade: in.Grade, Position: in.Position,
	}); err != nil {
		return gormmodel.ResumeLanguageV2{}, err
	}
	m := gormmodel.ResumeLanguageV2{ID: in.ID.UUID(), ResumeID: rid.UUID(), Language: string(in.Language), Position: in.Position}
	if in.Grade != nil {
		x := string(*in.Grade)
		m.Grade = &x
	}
	return m, nil
}
func mapLanguage(m gormmodel.ResumeLanguageV2) (domain.ResumeLanguageEntry, error) {
	id, e := valueobjects.IDFromUUID(m.ID)
	if e != nil {
		return domain.ResumeLanguageEntry{}, e
	}
	g, err := gradeLang(m.Grade)
	if err != nil {
		return domain.ResumeLanguageEntry{}, err
	}
	return domain.NewResumeLanguageEntry(domain.ResumeLanguageEntryParams{ID: id, Language: enums.LanguageType(m.Language), Grade: g, Position: m.Position})
}
func gradeLang(v *string) (*enums.LanguageGrade, error) {
	if v == nil {
		return nil, nil
	}
	g := enums.LanguageGrade(*v)
	if !g.IsValid() {
		return nil, fmt.Errorf("invalid language grade")
	}
	return &g, nil
}

package readers

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

type SkillsReader struct{ db *gorm.DB }

func NewSkillsReader(db *gorm.DB) *SkillsReader { return &SkillsReader{db: db} }

func (r *SkillsReader) GetSkills(ctx context.Context, q queries.GetSkillsQuery) (queries.SkillsView, error) {
	db, err := resolveDB(ctx, r.db, "get skills")
	if err != nil {
		return queries.SkillsView{}, err
	}
	if q.ResumeID.IsZero() || q.UserID.IsZero() {
		return queries.SkillsView{}, fmt.Errorf("get skills: invalid query")
	}

	var model gormmodel.ResumeSkills
	err = db.Joins("JOIN resumes ON resumes.id = resume_skills.resume_id").Where("resume_skills.resume_id = ? AND resumes.user_id = ?", q.ResumeID.UUID(), q.UserID.UUID()).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return queries.SkillsView{}, queries.ErrResumeNotFound
	}
	if err != nil {
		return queries.SkillsView{}, fmt.Errorf("get skills: %w", err)
	}

	var hard []gormmodel.ResumeSkillHard
	if err := db.Where("skills_id = ?", model.ID).Order("position").Find(&hard).Error; err != nil {
		return queries.SkillsView{}, fmt.Errorf("get hard skills: %w", err)
	}
	var soft []gormmodel.ResumeSkillSoft
	if err := db.Where("skills_id = ?", model.ID).Order("position").Find(&soft).Error; err != nil {
		return queries.SkillsView{}, fmt.Errorf("get soft skills: %w", err)
	}

	view := queries.SkillsView{ID: model.ID.String(), ResumeID: model.ResumeID.String(), HardSkills: make([]string, 0, len(hard)), SoftSkills: make([]string, 0, len(soft))}
	for _, item := range hard {
		view.HardSkills = append(view.HardSkills, item.Value)
	}
	for _, item := range soft {
		view.SoftSkills = append(view.SoftSkills, item.Value)
	}
	return view, nil
}

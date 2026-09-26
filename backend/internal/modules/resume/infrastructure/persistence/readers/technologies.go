package readers

import (
	"context"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

type TechnologiesReader struct{ db *gorm.DB }

func NewTechnologiesReader(db *gorm.DB) *TechnologiesReader { return &TechnologiesReader{db: db} }

func (r *TechnologiesReader) ListTechnologies(ctx context.Context, q queries.ListTechnologiesQuery) ([]queries.TechnologyView, error) {
	db, err := resolveDB(ctx, r.db, "list technologies")
	if err != nil {
		return nil, err
	}
	if q.ResumeID.IsZero() || q.UserID.IsZero() {
		return nil, fmt.Errorf("list technologies: invalid query")
	}
	var models []gormmodel.ResumeTechnology
	readDB := db.Joins("JOIN resumes ON resumes.id = resume_technologies.resume_id").Where("resume_technologies.resume_id = ?", q.ResumeID.UUID())
	readDB = readDB.Where("resumes.user_id = ?", q.UserID.UUID())
	if err := readDB.Order("position").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("list technologies: %w", err)
	}
	out := make([]queries.TechnologyView, 0, len(models))
	for _, model := range models {
		out = append(out, queries.TechnologyView{ID: model.ID.String(), ResumeID: model.ResumeID.String(), Name: model.Name, Grade: model.Grade, Position: model.Position})
	}
	return out, nil
}

func (r *TechnologiesReader) ReadTechnologies(ctx context.Context, resumeID valueobjects.ID) ([]queries.TechnologyView, error) {
	db, err := resolveDB(ctx, r.db, "list technologies")
	if err != nil {
		return nil, err
	}
	if resumeID.IsZero() {
		return nil, fmt.Errorf("list technologies: invalid query")
	}
	var models []gormmodel.ResumeTechnology
	readDB := db.Joins("JOIN resumes ON resumes.id = resume_technologies.resume_id").Where("resume_technologies.resume_id = ?", resumeID.UUID())
	if err := readDB.Order("position").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("list technologies: %w", err)
	}
	out := make([]queries.TechnologyView, 0, len(models))
	for _, model := range models {
		out = append(out, queries.TechnologyView{ID: model.ID.String(), ResumeID: model.ResumeID.String(), Name: model.Name, Grade: model.Grade, Position: model.Position})
	}
	return out, nil
}

func (r *TechnologiesReader) GetTechnology(ctx context.Context, q queries.GetTechnologyQuery) (queries.TechnologyView, error) {
	db, err := resolveDB(ctx, r.db, "get technology")
	if err != nil {
		return queries.TechnologyView{}, err
	}
	if q.ResumeID.IsZero() || q.TechnologyID.IsZero() || q.UserID.IsZero() {
		return queries.TechnologyView{}, fmt.Errorf("get technology: invalid query")
	}
	var model gormmodel.ResumeTechnology
	readDB := db.Joins("JOIN resumes ON resumes.id = resume_technologies.resume_id").Where("resume_technologies.id = ? AND resume_technologies.resume_id = ?", q.TechnologyID.UUID(), q.ResumeID.UUID())
	readDB = readDB.Where("resumes.user_id = ?", q.UserID.UUID())
	err = readDB.First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return queries.TechnologyView{}, queries.ErrResumeNotFound
	}
	if err != nil {
		return queries.TechnologyView{}, fmt.Errorf("get technology: %w", err)
	}
	return queries.TechnologyView{ID: model.ID.String(), ResumeID: model.ResumeID.String(), Name: model.Name, Grade: model.Grade, Position: model.Position}, nil
}

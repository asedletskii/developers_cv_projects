package readers

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

type EducationReader struct{ db *gorm.DB }

func NewEducationReader(db *gorm.DB) *EducationReader { return &EducationReader{db: db} }

func (r *EducationReader) ListEducations(ctx context.Context, q queries.ListEducationsQuery) ([]queries.EducationView, error) {
	db, err := resolveDB(ctx, r.db, "list educations")
	if err != nil {
		return nil, err
	}
	if q.ResumeID.IsZero() || q.UserID.IsZero() {
		return nil, fmt.Errorf("list educations: invalid query")
	}
	var models []gormmodel.ResumeEducation
	if err := db.Joins("JOIN resumes ON resumes.id = resume_educations.resume_id").Where("resume_educations.resume_id = ? AND resumes.user_id = ?", q.ResumeID.UUID(), q.UserID.UUID()).Order("sort_order").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("list educations: %w", err)
	}
	out := make([]queries.EducationView, 0, len(models))
	for _, model := range models {
		out = append(out, queries.EducationView{ID: model.ID.String(), ResumeID: model.ResumeID.String(), Location: model.Location, University: model.University, Grade: model.Grade, ProgramName: model.ProgramName, ProfileName: model.ProfileName, Faculty: model.Faculty, Course: model.Course, EndDate: optionalYear(model.EndDate)})
	}
	return out, nil
}

func (r *EducationReader) GetEducation(ctx context.Context, q queries.GetEducationQuery) (queries.EducationView, error) {
	db, err := resolveDB(ctx, r.db, "get education")
	if err != nil {
		return queries.EducationView{}, err
	}
	if q.ResumeID.IsZero() || q.EducationID.IsZero() || q.UserID.IsZero() {
		return queries.EducationView{}, fmt.Errorf("get education: invalid query")
	}
	var model gormmodel.ResumeEducation
	err = db.Joins("JOIN resumes ON resumes.id = resume_educations.resume_id").Where("resume_educations.id = ? AND resume_educations.resume_id = ? AND resumes.user_id = ?", q.EducationID.UUID(), q.ResumeID.UUID(), q.UserID.UUID()).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return queries.EducationView{}, queries.ErrResumeNotFound
	}
	if err != nil {
		return queries.EducationView{}, fmt.Errorf("get education: %w", err)
	}
	return queries.EducationView{ID: model.ID.String(), ResumeID: model.ResumeID.String(), Location: model.Location, University: model.University, Grade: model.Grade, ProgramName: model.ProgramName, ProfileName: model.ProfileName, Faculty: model.Faculty, Course: model.Course, EndDate: optionalYear(model.EndDate)}, nil
}

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

type ActivityReader struct{ db *gorm.DB }

func NewActivityReader(db *gorm.DB) *ActivityReader { return &ActivityReader{db: db} }

func (r *ActivityReader) ListActivities(ctx context.Context, q queries.ListActivitiesQuery) ([]queries.ActivityView, error) {
	db, err := resolveDB(ctx, r.db, "list activities")
	if err != nil {
		return nil, err
	}
	if q.ResumeID.IsZero() || q.UserID.IsZero() {
		return nil, fmt.Errorf("list activities: invalid query")
	}
	var models []gormmodel.ResumeActivity
	readDB := db.Joins("JOIN resumes ON resumes.id = resume_activities.resume_id").Where("resume_activities.resume_id = ?", q.ResumeID.UUID())
	readDB = readDB.Where("resumes.user_id = ?", q.UserID.UUID())
	if err := readDB.Order("sort_order").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("list activities: %w", err)
	}
	out := make([]queries.ActivityView, 0, len(models))
	for _, model := range models {
		out = append(out, queries.ActivityView{ID: model.ID.String(), ResumeID: model.ResumeID.String(), Type: model.Type, Title: model.Title, Description: model.Description, Date: optionalYearMonth(model.Date), Link: model.Link})
	}
	return out, nil
}

func (r *ActivityReader) ReadActivities(ctx context.Context, resumeID valueobjects.ID) ([]queries.ActivityView, error) {
	db, err := resolveDB(ctx, r.db, "list activities")
	if err != nil {
		return nil, err
	}
	if resumeID.IsZero() {
		return nil, fmt.Errorf("list activities: invalid query")
	}
	var models []gormmodel.ResumeActivity
	readDB := db.Joins("JOIN resumes ON resumes.id = resume_activities.resume_id").Where("resume_activities.resume_id = ?", resumeID.UUID())
	if err := readDB.Order("sort_order").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("list activities: %w", err)
	}
	out := make([]queries.ActivityView, 0, len(models))
	for _, model := range models {
		out = append(out, queries.ActivityView{ID: model.ID.String(), ResumeID: model.ResumeID.String(), Type: model.Type, Title: model.Title, Description: model.Description, Date: optionalYearMonth(model.Date), Link: model.Link})
	}
	return out, nil
}

func (r *ActivityReader) GetActivity(ctx context.Context, q queries.GetActivityQuery) (queries.ActivityView, error) {
	db, err := resolveDB(ctx, r.db, "get activity")
	if err != nil {
		return queries.ActivityView{}, err
	}
	if q.ResumeID.IsZero() || q.ActivityID.IsZero() || q.UserID.IsZero() {
		return queries.ActivityView{}, fmt.Errorf("get activity: invalid query")
	}
	var model gormmodel.ResumeActivity
	readDB := db.Joins("JOIN resumes ON resumes.id = resume_activities.resume_id").Where("resume_activities.id = ? AND resume_activities.resume_id = ?", q.ActivityID.UUID(), q.ResumeID.UUID())
	readDB = readDB.Where("resumes.user_id = ?", q.UserID.UUID())
	err = readDB.First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return queries.ActivityView{}, queries.ErrResumeNotFound
	}
	if err != nil {
		return queries.ActivityView{}, fmt.Errorf("get activity: %w", err)
	}
	return queries.ActivityView{ID: model.ID.String(), ResumeID: model.ResumeID.String(), Type: model.Type, Title: model.Title, Description: model.Description, Date: optionalYearMonth(model.Date), Link: model.Link}, nil
}

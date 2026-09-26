package readers

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

type CourseReader struct{ db *gorm.DB }

func NewCourseReader(db *gorm.DB) *CourseReader { return &CourseReader{db: db} }

func (r *CourseReader) ListCourses(ctx context.Context, q queries.ListCoursesQuery) ([]queries.CourseView, error) {
	db, err := resolveDB(ctx, r.db, "list courses")
	if err != nil {
		return nil, err
	}
	if q.ResumeID.IsZero() || q.UserID.IsZero() {
		return nil, fmt.Errorf("list courses: invalid query")
	}
	var models []gormmodel.ResumeCourse
	if err := db.Joins("JOIN resumes ON resumes.id = resume_courses.resume_id").Where("resume_courses.resume_id = ? AND resumes.user_id = ?", q.ResumeID.UUID(), q.UserID.UUID()).Order("sort_order").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("list courses: %w", err)
	}
	out := make([]queries.CourseView, 0, len(models))
	for _, model := range models {
		out = append(out, queries.CourseView{ID: model.ID.String(), ResumeID: model.ResumeID.String(), Name: model.Name, Organization: model.Organization, EndDate: yearMonth(model.EndDate), Hours: model.Hours, CertificateLink: model.CertificateLink})
	}
	return out, nil
}

func (r *CourseReader) GetCourse(ctx context.Context, q queries.GetCourseQuery) (queries.CourseView, error) {
	db, err := resolveDB(ctx, r.db, "get course")
	if err != nil {
		return queries.CourseView{}, err
	}
	if q.ResumeID.IsZero() || q.CourseID.IsZero() || q.UserID.IsZero() {
		return queries.CourseView{}, fmt.Errorf("get course: invalid query")
	}
	var model gormmodel.ResumeCourse
	err = db.Joins("JOIN resumes ON resumes.id = resume_courses.resume_id").Where("resume_courses.id = ? AND resume_courses.resume_id = ? AND resumes.user_id = ?", q.CourseID.UUID(), q.ResumeID.UUID(), q.UserID.UUID()).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return queries.CourseView{}, queries.ErrResumeNotFound
	}
	if err != nil {
		return queries.CourseView{}, fmt.Errorf("get course: %w", err)
	}
	return queries.CourseView{ID: model.ID.String(), ResumeID: model.ResumeID.String(), Name: model.Name, Organization: model.Organization, EndDate: yearMonth(model.EndDate), Hours: model.Hours, CertificateLink: model.CertificateLink}, nil
}

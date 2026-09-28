package repository

import (
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"

	"gorm.io/gorm"
)

func persistCourses(tx *gorm.DB, resume *domain.Resume) error {
	models := make([]gormmodel.ResumeCourse, 0, len(resume.Courses()))
	for position, item := range resume.Courses() {
		if item == nil {
			return fmt.Errorf("%w: nil course", domain.ErrResumeInvariant)
		}
		if item.ResumeID() != resume.ID() {
			return fmt.Errorf("%w: course does not belong to resume", domain.ErrResumeInvariant)
		}
		model, err := courseModel(item)
		if err != nil {
			return err
		}
		model.SortOrder = position
		models = append(models, model)
	}

	if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeCourse{}).Error; err != nil {
		return fmt.Errorf("replace resume courses: %w", mapDBError(err))
	}
	if len(models) > 0 {
		if err := tx.Create(&models).Error; err != nil {
			return fmt.Errorf("create resume courses: %w", mapDBError(err))
		}
	}
	return nil
}

func loadCourses(tx *gorm.DB, resumeID valueobjects.ID) ([]*domain.Course, error) {
	var rows []gormmodel.ResumeCourse
	if err := tx.Where("resume_id = ?", resumeID.UUID()).Order("sort_order").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("get resume courses: %w", mapDBError(err))
	}
	out := make([]*domain.Course, 0, len(rows))
	for _, row := range rows {
		x, err := mapCourse(row)
		if err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	if len(out) == 0 {
		return nil, nil
	}
	return out, nil
}

func courseModel(item *domain.Course) (gormmodel.ResumeCourse, error) {
	if item == nil {
		return gormmodel.ResumeCourse{}, fmt.Errorf("course is nil")
	}
	return gormmodel.ResumeCourse{ID: item.ID().UUID(), ResumeID: item.ResumeID().UUID(), Name: item.Name(), Organization: item.Organization(), EndDate: yearMonthDate(item.EndDate()), Hours: item.Hours(), CertificateLink: externalURLString(item.CertificateLink())}, nil
}
func mapCourse(model gormmodel.ResumeCourse) (*domain.Course, error) {
	id, err := valueobjects.IDFromUUID(model.ID)
	if err != nil {
		return nil, err
	}
	rid, err := valueobjects.IDFromUUID(model.ResumeID)
	if err != nil {
		return nil, err
	}
	end, err := valueobjects.ParseYearMonth(model.EndDate.Format("2006-01"))
	if err != nil {
		return nil, err
	}
	var link *valueobjects.ExternalURL
	if model.CertificateLink != nil {
		v, e := valueobjects.NewExternalURL(*model.CertificateLink)
		if e != nil {
			return nil, e
		}
		link = &v
	}
	return domain.NewCourse(domain.CourseParams{ID: id, ResumeID: rid, Name: model.Name, Organization: model.Organization, EndDate: end, Hours: model.Hours, CertificateLink: link})
}
func externalURLString(value *valueobjects.ExternalURL) *string {
	if value == nil {
		return nil
	}
	result := value.String()
	return &result
}

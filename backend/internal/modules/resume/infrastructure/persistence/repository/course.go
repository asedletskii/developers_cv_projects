package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

var _ commands.CoursesRepository = (*ResumeRepository)(nil)

func (r *ResumeRepository) CreateCourse(ctx context.Context, c commands.CreateCourseCommand) error {
	tx, err := r.owned(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	item, err := courseModel(c.Input)
	if err != nil {
		return err
	}
	item.SortOrder, err = nextSortOrder(tx, item.TableName(), c.ResumeID.UUID())
	if err != nil {
		return err
	}
	if err := tx.Create(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrDuplicateEntity
		}
		return err
	}
	return nil
}
func (r *ResumeRepository) ListCourses(ctx context.Context, q commands.ListCoursesQuery) ([]*domain.Course, error) {
	tx, err := r.owned(ctx, q.ResumeID, q.UserID)
	if err != nil {
		return nil, err
	}
	var models []gormmodel.ResumeCourse
	if err := tx.Where("resume_id = ?", q.ResumeID.UUID()).Order("sort_order").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]*domain.Course, 0, len(models))
	for _, model := range models {
		item, err := mapCourse(model)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, nil
}
func (r *ResumeRepository) GetCourse(ctx context.Context, q commands.GetCourseQuery) (*domain.Course, error) {
	tx, err := r.owned(ctx, q.ResumeID, q.UserID)
	if err != nil {
		return nil, err
	}
	var model gormmodel.ResumeCourse
	if err := tx.Where("id = ? AND resume_id = ?", q.CourseID.UUID(), q.ResumeID.UUID()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, err
	}
	return mapCourse(model)
}
func (r *ResumeRepository) UpdateCourse(ctx context.Context, c commands.UpdateCourseCommand) error {
	tx, err := r.owned(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	model, err := courseModel(c.Input)
	if err != nil {
		return err
	}
	result := tx.Model(&gormmodel.ResumeCourse{}).Where("id = ? AND resume_id = ?", c.CourseID.UUID(), c.ResumeID.UUID()).Updates(map[string]any{"name": model.Name, "organization": model.Organization, "end_date": model.EndDate, "hours": model.Hours, "certificate_link": model.CertificateLink})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
}
func (r *ResumeRepository) DeleteCourse(ctx context.Context, c commands.DeleteCourseCommand) error {
	tx, err := r.owned(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	result := tx.Where("id = ? AND resume_id = ?", c.CourseID.UUID(), c.ResumeID.UUID()).Delete(&gormmodel.ResumeCourse{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
}
func courseModel(input commands.CourseInput) (gormmodel.ResumeCourse, error) {
	item, err := domain.NewCourse(domain.CourseParams{ID: input.ID, ResumeID: input.ResumeID, Name: input.Name, Organization: input.Organization, EndDate: input.EndDate, Hours: input.Hours, CertificateLink: input.CertificateLink})
	if err != nil {
		return gormmodel.ResumeCourse{}, err
	}
	return gormmodel.ResumeCourse{ID: item.ID().UUID(), ResumeID: item.ResumeID().UUID(), Name: item.Name(), Organization: item.Organization(), EndDate: item.EndDate().String() + "-01", Hours: item.Hours(), CertificateLink: externalURLString(item.CertificateLink())}, nil
}
func mapCourse(model gormmodel.ResumeCourse) (*domain.Course, error) {
	id, err := valueobjects.IDFromUUID(model.ID)
	if err != nil {
		return nil, err
	}
	resumeID, err := valueobjects.IDFromUUID(model.ResumeID)
	if err != nil {
		return nil, err
	}
	endDate, err := valueobjects.ParseYearMonth(model.EndDate[:7])
	if err != nil {
		return nil, err
	}
	var link *valueobjects.ExternalURL
	if model.CertificateLink != nil {
		value, err := valueobjects.NewExternalURL(*model.CertificateLink)
		if err != nil {
			return nil, err
		}
		link = &value
	}
	return domain.NewCourse(domain.CourseParams{ID: id, ResumeID: resumeID, Name: model.Name, Organization: model.Organization, EndDate: endDate, Hours: model.Hours, CertificateLink: link})
}
func externalURLString(value *valueobjects.ExternalURL) *string {
	if value == nil {
		return nil
	}
	result := value.String()
	return &result
}

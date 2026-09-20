package mappers

import (
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func CoursesToDomain(rows []gormmodel.ResumeCourse) ([]*domain.Course, error) {
	courses := make([]*domain.Course, 0, len(rows))
	for i, row := range rows {
		course, err := CourseToDomain(row)
		if err != nil {
			return nil, fmt.Errorf("map course at index %d: %w", i, err)
		}
		courses = append(courses, course)
	}
	return courses, nil
}

func CourseToDB(item *domain.Course) gormmodel.ResumeCourse {
	return gormmodel.ResumeCourse{
		ID:              item.ID().UUID(),
		ResumeID:        item.ResumeID().UUID(),
		Name:            item.Name(),
		Organization:    item.Organization(),
		EndDate:         YearMonthToTime(item.EndDate()),
		Hours:           item.Hours(),
		CertificateLink: ExternalURLToDB(item.CertificateLink()),
	}
}

func CourseToDomain(model gormmodel.ResumeCourse) (*domain.Course, error) {
	id, err := sharedvalueobjects.IDFromUUID(model.ID)
	if err != nil {
		return nil, err
	}
	rid, err := sharedvalueobjects.IDFromUUID(model.ResumeID)
	if err != nil {
		return nil, err
	}
	end := domain.RestoreYearMonth(model.EndDate)
	var link *domain.ExternalURL
	if model.CertificateLink != nil {
		v := domain.RestoreExternalURL(*model.CertificateLink)
		link = &v
	}
	return domain.RestoreCourse(domain.CourseSnapshot{
		ID:              id,
		ResumeID:        rid,
		Name:            model.Name,
		Organization:    model.Organization,
		EndDate:         end,
		Hours:           model.Hours,
		CertificateLink: link,
	}), nil
}

func ExternalURLToDB(value *domain.ExternalURL) *string {
	if value == nil {
		return nil
	}
	result := value.String()
	return &result
}

func CoursesToDBModels(items []*domain.Course) []gormmodel.ResumeCourse {
	models := make([]gormmodel.ResumeCourse, 0, len(items))
	for position, item := range items {
		model := CourseToDB(item)
		model.SortOrder = position
		models = append(models, model)
	}
	return models
}

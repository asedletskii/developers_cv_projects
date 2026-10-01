package mappers

import (
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func CourseToDB(item *domain.Course) (gormmodel.ResumeCourse, error) {
	if item == nil {
		return gormmodel.ResumeCourse{}, fmt.Errorf("course is nil")
	}
	return gormmodel.ResumeCourse{
		ID:              item.ID().UUID(),
		ResumeID:        item.ResumeID().UUID(),
		Name:            item.Name(),
		Organization:    item.Organization(),
		EndDate:         YearMonthToTime(item.EndDate()),
		Hours:           item.Hours(),
		CertificateLink: ExternalURLToDB(item.CertificateLink()),
	}, nil
}

func CourseToDomain(model gormmodel.ResumeCourse) (*domain.Course, error) {
	id, err := identity.IDFromUUID(model.ID)
	if err != nil {
		return nil, err
	}
	rid, err := identity.IDFromUUID(model.ResumeID)
	if err != nil {
		return nil, err
	}
	end, err := domain.ParseYearMonth(model.EndDate.Format("2006-01"))
	if err != nil {
		return nil, err
	}
	var link *domain.ExternalURL
	if model.CertificateLink != nil {
		v, e := domain.NewExternalURL(*model.CertificateLink)
		if e != nil {
			return nil, e
		}
		link = &v
	}
	return domain.NewCourse(domain.CourseParams{
		ID:              id,
		ResumeID:        rid,
		Name:            model.Name,
		Organization:    model.Organization,
		EndDate:         end,
		Hours:           model.Hours,
		CertificateLink: link,
	})
}

func ExternalURLToDB(value *domain.ExternalURL) *string {
	if value == nil {
		return nil
	}
	result := value.String()
	return &result
}

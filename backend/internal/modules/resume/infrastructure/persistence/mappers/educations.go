package mappers

import (
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func EducationToDB(item *domain.Education) (gormmodel.ResumeEducation, error) {
	if item == nil {
		return gormmodel.ResumeEducation{}, fmt.Errorf("education is nil")
	}
	m := gormmodel.ResumeEducation{
		ID:          item.ID().UUID(),
		ResumeID:    item.ResumeID().UUID(),
		Location:    item.Location(),
		University:  item.University(),
		Grade:       int(item.Grade()),
		ProgramName: item.ProgramName(),
		ProfileName: item.ProfileName(),
		Faculty:     item.Faculty(),
		Course:      item.Course(),
	}
	if item.EndDate() != nil {
		v := YearToTime(*item.EndDate())
		m.EndDate = &v
	}
	return m, nil
}

func EducationToDomain(model gormmodel.ResumeEducation) (*domain.Education, error) {
	id, err := identity.IDFromUUID(model.ID)
	if err != nil {
		return nil, err
	}
	rid, err := identity.IDFromUUID(model.ResumeID)
	if err != nil {
		return nil, err
	}
	grade := domain.EducationLevel(model.Grade)
	if !grade.IsValid() {
		return nil, domain.ErrInvalidEntity
	}
	var end *domain.Year
	if model.EndDate != nil {
		v, e := domain.ParseYear(model.EndDate.Format("2006"))
		if e != nil {
			return nil, e
		}
		end = &v
	}
	return domain.NewEducation(domain.EducationParams{
		ID:          id,
		ResumeID:    rid,
		Location:    model.Location,
		University:  model.University,
		Grade:       grade,
		ProgramName: model.ProgramName,
		ProfileName: model.ProfileName,
		Faculty:     model.Faculty,
		Course:      model.Course,
		EndDate:     end,
	})
}

package mappers

import (
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func EducationsToDomain(rows []gormmodel.ResumeEducation) ([]*domain.Education, error) {
	educations := make([]*domain.Education, 0, len(rows))
	for i, row := range rows {
		education, err := EducationToDomain(row)
		if err != nil {
			return nil, fmt.Errorf("map education at index %d: %w", i, err)
		}
		educations = append(educations, education)
	}
	return educations, nil
}

func EducationToDB(item *domain.Education) gormmodel.ResumeEducation {
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
	return m
}

func EducationToDomain(model gormmodel.ResumeEducation) (*domain.Education, error) {
	id, err := sharedvalueobjects.IDFromUUID(model.ID)
	if err != nil {
		return nil, err
	}
	rid, err := sharedvalueobjects.IDFromUUID(model.ResumeID)
	if err != nil {
		return nil, err
	}
	grade := domain.EducationLevel(model.Grade)
	var end *domain.Year
	if model.EndDate != nil {
		v := domain.RestoreYear(model.EndDate.Year())
		end = &v
	}
	return domain.RestoreEducation(domain.EducationSnapshot{
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
	}), nil
}

func EducationsToDBModels(items []*domain.Education) []gormmodel.ResumeEducation {
	models := make([]gormmodel.ResumeEducation, 0, len(items))
	for position, item := range items {
		model := EducationToDB(item)
		model.SortOrder = position
		models = append(models, model)
	}
	return models
}

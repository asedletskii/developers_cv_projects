package repository

import (
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"fmt"

	"gorm.io/gorm"
)

func persistEducations(tx *gorm.DB, resume *domain.Resume) error {
	models := make([]gormmodel.ResumeEducation, 0, len(resume.Educations()))
	for position, item := range resume.Educations() {
		if item == nil {
			return fmt.Errorf("%w: nil education", domain.ErrResumeInvariant)
		}
		if item.ResumeID() != resume.ID() {
			return fmt.Errorf("%w: education does not belong to resume", domain.ErrResumeInvariant)
		}
		model, err := educationModel(item)
		if err != nil {
			return err
		}
		model.SortOrder = position
		models = append(models, model)
	}

	if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeEducation{}).Error; err != nil {
		return fmt.Errorf("replace resume educations: %w", mapDBError(err))
	}
	if len(models) > 0 {
		if err := tx.Create(&models).Error; err != nil {
			return fmt.Errorf("create resume educations: %w", mapDBError(err))
		}
	}
	return nil
}
func loadEducations(tx *gorm.DB, resumeID valueobjects.ID) ([]*domain.Education, error) {
	var rows []gormmodel.ResumeEducation
	if err := tx.Where("resume_id = ?", resumeID.UUID()).Order("sort_order").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("get resume educations: %w", mapDBError(err))
	}
	out := make([]*domain.Education, 0, len(rows))
	for _, row := range rows {
		x, err := mapEducation(row)
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
func educationModel(item *domain.Education) (gormmodel.ResumeEducation, error) {
	if item == nil {
		return gormmodel.ResumeEducation{}, fmt.Errorf("education is nil")
	}
	m := gormmodel.ResumeEducation{ID: item.ID().UUID(), ResumeID: item.ResumeID().UUID(), Location: item.Location(), University: item.University(), Grade: int(item.Grade()), ProgramName: item.ProgramName(), ProfileName: item.ProfileName(), Faculty: item.Faculty(), Course: item.Course()}
	if item.EndDate() != nil {
		v := yearDate(*item.EndDate())
		m.EndDate = &v
	}
	return m, nil
}
func mapEducation(model gormmodel.ResumeEducation) (*domain.Education, error) {
	id, err := valueobjects.IDFromUUID(model.ID)
	if err != nil {
		return nil, err
	}
	rid, err := valueobjects.IDFromUUID(model.ResumeID)
	if err != nil {
		return nil, err
	}
	grade := enums.EducationLevel(model.Grade)
	if !grade.IsValid() {
		return nil, domain.ErrInvalidEntity
	}
	var end *valueobjects.Year
	if model.EndDate != nil {
		v, e := valueobjects.ParseYear(model.EndDate.Format("2006"))
		if e != nil {
			return nil, e
		}
		end = &v
	}
	return domain.NewEducation(domain.EducationParams{ID: id, ResumeID: rid, Location: model.Location, University: model.University, Grade: grade, ProgramName: model.ProgramName, ProfileName: model.ProfileName, Faculty: model.Faculty, Course: model.Course, EndDate: end})
}

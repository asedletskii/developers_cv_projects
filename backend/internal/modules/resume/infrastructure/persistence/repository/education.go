package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

func (r *EducationRepository) CreateEducation(ctx context.Context, c ports.CreateEducationRepositoryCommand) error {
	tx, err := r.ownedForWrite(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	item, err := educationModel(c.Education)
	if err != nil {
		return err
	}
	item.SortOrder, err = nextSortOrder(tx, item.TableName(), c.ResumeID.UUID())
	if err != nil {
		return err
	}
	if err := tx.Create(&item).Error; err != nil {
		return mapDBError(err)
	}
	return nil
}

func (r *EducationRepository) GetEducation(ctx context.Context, q ports.GetEducationRepositoryCommand) (*domain.Education, error) {
	tx, err := r.owned(ctx, q.ResumeID, q.UserID)
	if err != nil {
		return nil, err
	}
	var model gormmodel.ResumeEducation
	if err := tx.Where("id = ? AND resume_id = ?", q.EducationID.UUID(), q.ResumeID.UUID()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, err
	}
	return mapEducation(model)
}

func (r *EducationRepository) UpdateEducation(ctx context.Context, c ports.UpdateEducationRepositoryCommand) error {
	tx, err := r.owned(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	model, err := educationModel(c.Education)
	if err != nil {
		return err
	}
	result := tx.Model(&gormmodel.ResumeEducation{}).Where("id = ? AND resume_id = ?", c.Education.ID().UUID(), c.ResumeID.UUID()).Updates(map[string]any{
		"location": model.Location, "university": model.University, "grade": model.Grade,
		"program_name": model.ProgramName, "profile_name": model.ProfileName, "faculty": model.Faculty,
		"course": model.Course, "end_date": model.EndDate,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
}

func (r *EducationRepository) DeleteEducation(ctx context.Context, c ports.DeleteEducationRepositoryCommand) error {
	tx, err := r.owned(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	result := tx.Where("id = ? AND resume_id = ?", c.EducationID.UUID(), c.ResumeID.UUID()).Delete(&gormmodel.ResumeEducation{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
}

func educationModel(item *domain.Education) (gormmodel.ResumeEducation, error) {
	model := gormmodel.ResumeEducation{ID: item.ID().UUID(), ResumeID: item.ResumeID().UUID(), Location: item.Location(), University: item.University(), Grade: int(item.Grade()), ProgramName: item.ProgramName(), ProfileName: item.ProfileName(), Faculty: item.Faculty(), Course: item.Course()}
	if item.EndDate() != nil {
		value := yearDate(*item.EndDate())
		model.EndDate = &value
	}
	return model, nil
}

func mapEducation(model gormmodel.ResumeEducation) (*domain.Education, error) {
	id, err := valueobjects.IDFromUUID(model.ID)
	if err != nil {
		return nil, err
	}
	resumeID, err := valueobjects.IDFromUUID(model.ResumeID)
	if err != nil {
		return nil, err
	}
	grade := enums.EducationLevel(model.Grade)
	if !grade.IsValid() {
		return nil, domain.ErrInvalidEntity
	}
	var endDate *valueobjects.Year
	if model.EndDate != nil {
		value, err := valueobjects.ParseYear(model.EndDate.Format("2006"))
		if err != nil {
			return nil, err
		}
		endDate = &value
	}
	return domain.NewEducation(domain.EducationParams{ID: id, ResumeID: resumeID, Location: model.Location, University: model.University, Grade: grade, ProgramName: model.ProgramName, ProfileName: model.ProfileName, Faculty: model.Faculty, Course: model.Course, EndDate: endDate})
}

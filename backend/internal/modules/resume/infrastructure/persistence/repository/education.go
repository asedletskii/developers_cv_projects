package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

var _ commands.EducationsRepository = (*EducationRepository)(nil)

func (r *EducationRepository) CreateEducation(ctx context.Context, c commands.CreateEducationCommand) error {
	tx, err := r.owned(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	item, err := educationModel(c.Input)
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

func (r *EducationRepository) ListEducations(ctx context.Context, q commands.ListEducationsQuery) ([]*domain.Education, error) {
	tx, err := r.owned(ctx, q.ResumeID, q.UserID)
	if err != nil {
		return nil, err
	}
	var models []gormmodel.ResumeEducation
	if err := tx.Where("resume_id = ?", q.ResumeID.UUID()).Order("sort_order").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]*domain.Education, 0, len(models))
	for _, model := range models {
		item, err := mapEducation(model)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, nil
}

func (r *EducationRepository) GetEducation(ctx context.Context, q commands.GetEducationQuery) (*domain.Education, error) {
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

func (r *EducationRepository) UpdateEducation(ctx context.Context, c commands.UpdateEducationCommand) error {
	tx, err := r.owned(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	model, err := educationModel(c.Input)
	if err != nil {
		return err
	}
	result := tx.Model(&gormmodel.ResumeEducation{}).Where("id = ? AND resume_id = ?", c.EducationID.UUID(), c.ResumeID.UUID()).Updates(map[string]any{
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

func (r *EducationRepository) DeleteEducation(ctx context.Context, c commands.DeleteEducationCommand) error {
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

func educationModel(input commands.EducationInput) (gormmodel.ResumeEducation, error) {
	item, err := domain.NewEducation(domain.EducationParams{ID: input.ID, ResumeID: input.ResumeID, Location: input.Location, University: input.University, Grade: input.Grade, ProgramName: input.ProgramName, ProfileName: input.ProfileName, Faculty: input.Faculty, Course: input.Course, EndDate: input.EndDate})
	if err != nil {
		return gormmodel.ResumeEducation{}, err
	}
	model := gormmodel.ResumeEducation{ID: item.ID().UUID(), ResumeID: item.ResumeID().UUID(), Location: item.Location(), University: item.University(), Grade: int(item.Grade()), ProgramName: item.ProgramName(), ProfileName: item.ProfileName(), Faculty: item.Faculty(), Course: item.Course()}
	if item.EndDate() != nil {
		value := item.EndDate().String() + "-01-01"
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
		value, err := valueobjects.ParseYear((*model.EndDate)[:4])
		if err != nil {
			return nil, err
		}
		endDate = &value
	}
	return domain.NewEducation(domain.EducationParams{ID: id, ResumeID: resumeID, Location: model.Location, University: model.University, Grade: grade, ProgramName: model.ProgramName, ProfileName: model.ProfileName, Faculty: model.Faculty, Course: model.Course, EndDate: endDate})
}

func nextSortOrder(tx *gorm.DB, table string, resumeID uuid.UUID) (int, error) {
	var maxOrder *int
	if err := tx.Table(table).Where("resume_id = ?", resumeID).Select("MAX(sort_order)").Scan(&maxOrder).Error; err != nil {
		return 0, err
	}
	if maxOrder == nil {
		return 0, nil
	}
	return *maxOrder + 1, nil
}

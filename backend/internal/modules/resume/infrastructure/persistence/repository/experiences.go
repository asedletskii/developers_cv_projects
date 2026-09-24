package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

var _ commands.ExperiencesRepository = (*ExperiencesRepository)(nil)

func (r *ExperiencesRepository) CreateExperience(ctx context.Context, c commands.CreateExperienceCommand) error {
	tx, err := r.owned(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	item, err := experienceModel(c.Input)
	if err != nil {
		return err
	}
	var maxOrder *int
	if err := tx.Model(&gormmodel.ResumeExperience{}).
		Where("resume_id = ?", c.ResumeID.UUID()).
		Select("MAX(sort_order)").Scan(&maxOrder).Error; err != nil {
		return err
	}
	if maxOrder == nil {
		item.SortOrder = 0
	} else {
		item.SortOrder = *maxOrder + 1
	}
	if err := tx.Create(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrDuplicateEntity
		}
		return err
	}
	return replaceExperienceChildren(tx, item.ID, c.Input)
}

func (r *ExperiencesRepository) ListExperiences(ctx context.Context, q commands.ListExperiencesQuery) ([]*domain.Experience, error) {
	tx, err := r.owned(ctx, q.ResumeID, q.UserID)
	if err != nil {
		return nil, err
	}
	var models []gormmodel.ResumeExperience
	if err := tx.Where("resume_id = ?", q.ResumeID.UUID()).Order("sort_order").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]*domain.Experience, 0, len(models))
	for _, model := range models {
		item, err := mapExperience(tx, model)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, nil
}

func (r *ExperiencesRepository) GetExperience(ctx context.Context, q commands.GetExperienceQuery) (*domain.Experience, error) {
	tx, err := r.owned(ctx, q.ResumeID, q.UserID)
	if err != nil {
		return nil, err
	}
	var model gormmodel.ResumeExperience
	if err := tx.Where("id = ? AND resume_id = ?", q.ExperienceID.UUID(), q.ResumeID.UUID()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, err
	}
	return mapExperience(tx, model)
}

func (r *ExperiencesRepository) UpdateExperience(ctx context.Context, c commands.UpdateExperienceCommand) error {
	tx, err := r.owned(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	model, err := experienceModel(c.Input)
	if err != nil {
		return err
	}
	result := tx.Model(&gormmodel.ResumeExperience{}).
		Where("id = ? AND resume_id = ?", c.ExperienceID.UUID(), c.ResumeID.UUID()).
		Updates(map[string]any{
			"company":     model.Company,
			"position":    model.Position,
			"location":    model.Location,
			"start_date":  model.StartDate,
			"end_date":    model.EndDate,
			"description": model.Description,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	if err := replaceExperienceChildren(tx, c.ExperienceID.UUID(), c.Input); err != nil {
		return err
	}
	return nil
}

func (r *ExperiencesRepository) DeleteExperience(ctx context.Context, c commands.DeleteExperienceCommand) error {
	tx, err := r.owned(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	result := tx.Where("id = ? AND resume_id = ?", c.ExperienceID.UUID(), c.ResumeID.UUID()).Delete(&gormmodel.ResumeExperience{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
}

func experienceModel(input commands.ExperienceInput) (gormmodel.ResumeExperience, error) {
	if input.ID.IsZero() || input.ResumeID.IsZero() {
		return gormmodel.ResumeExperience{}, fmt.Errorf("experience identity is required")
	}
	item, err := newExperienceFromInput(input)
	if err != nil {
		return gormmodel.ResumeExperience{}, err
	}
	model := gormmodel.ResumeExperience{
		ID:          item.ID().UUID(),
		ResumeID:    item.ResumeID().UUID(),
		Company:     item.Company(),
		Position:    item.Position(),
		Location:    item.Location(),
		StartDate:   item.StartDate().String() + "-01",
		Description: item.Description(),
	}
	if item.EndDate() != nil {
		value := item.EndDate().String() + "-01"
		model.EndDate = &value
	}
	return model, nil
}

func replaceExperienceChildren(tx *gorm.DB, experienceID uuid.UUID, input commands.ExperienceInput) error {
	item, err := newExperienceFromInput(input)
	if err != nil {
		return err
	}
	if err := tx.Where("experience_id = ?", experienceID).Delete(&gormmodel.ResumeExperienceAchievement{}).Error; err != nil {
		return err
	}
	if err := tx.Where("experience_id = ?", experienceID).Delete(&gormmodel.ResumeExperienceTechnology{}).Error; err != nil {
		return err
	}
	for i, achievement := range item.Achievements() {
		if err := tx.Create(&gormmodel.ResumeExperienceAchievement{ExperienceID: experienceID, Position: i, Value: achievement}).Error; err != nil {
			return err
		}
	}
	for i, technology := range item.Technologies() {
		var grade *string
		if value := technology.Grade(); value != nil {
			text := string(*value)
			grade = &text
		}
		if err := tx.Create(&gormmodel.ResumeExperienceTechnology{ExperienceID: experienceID, Position: i, Name: technology.Name(), Grade: grade}).Error; err != nil {
			return err
		}
	}
	return nil
}

func mapExperience(tx *gorm.DB, model gormmodel.ResumeExperience) (*domain.Experience, error) {
	id, err := valueobjects.IDFromUUID(model.ID)
	if err != nil {
		return nil, err
	}
	resumeID, err := valueobjects.IDFromUUID(model.ResumeID)
	if err != nil {
		return nil, err
	}
	start, err := valueobjects.ParseYearMonth(model.StartDate[:7])
	if err != nil {
		return nil, err
	}
	var end *valueobjects.YearMonth
	if model.EndDate != nil {
		value, err := valueobjects.ParseYearMonth((*model.EndDate)[:7])
		if err != nil {
			return nil, err
		}
		end = &value
	}
	var achievements []gormmodel.ResumeExperienceAchievement
	if err := tx.Where("experience_id = ?", model.ID).Order("position").Find(&achievements).Error; err != nil {
		return nil, err
	}
	var technologies []gormmodel.ResumeExperienceTechnology
	if err := tx.Where("experience_id = ?", model.ID).Order("position").Find(&technologies).Error; err != nil {
		return nil, err
	}
	technologyParams := make([]domain.ExperienceTechnologyParams, 0, len(technologies))
	for _, technology := range technologies {
		var grade *enums.TechnologyGrade
		if technology.Grade != nil {
			value := enums.TechnologyGrade(*technology.Grade)
			if !value.IsValid() {
				return nil, fmt.Errorf("invalid experience technology grade")
			}
			grade = &value
		}
		technologyParams = append(technologyParams, domain.ExperienceTechnologyParams{Name: technology.Name, Grade: grade})
	}
	achievementValues := make([]string, len(achievements))
	for i, achievement := range achievements {
		achievementValues[i] = achievement.Value
	}
	return domain.NewExperience(domain.ExperienceParams{
		ID:           id,
		ResumeID:     resumeID,
		Company:      model.Company,
		Position:     model.Position,
		Location:     model.Location,
		StartDate:    start,
		EndDate:      end,
		Description:  model.Description,
		Achievements: achievementValues,
		Technologies: technologyParams,
	})
}

func newExperienceFromInput(input commands.ExperienceInput) (*domain.Experience, error) {
	technologies := make([]domain.ExperienceTechnologyParams, len(input.Technologies))
	for i, technology := range input.Technologies {
		technologies[i] = domain.ExperienceTechnologyParams{Name: technology.Name, Grade: technology.Grade}
	}
	return domain.NewExperience(domain.ExperienceParams{
		ID: input.ID, ResumeID: input.ResumeID, Company: input.Company, Position: input.Position,
		Location: input.Location, StartDate: input.StartDate, EndDate: input.EndDate,
		Description: input.Description, Achievements: input.Achievements, Technologies: technologies,
	})
}

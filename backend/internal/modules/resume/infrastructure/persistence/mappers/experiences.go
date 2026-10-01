package mappers

import (
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func ExperienceToDB(item *domain.Experience) (gormmodel.ResumeExperience, error) {
	if item == nil {
		return gormmodel.ResumeExperience{}, fmt.Errorf("experience is nil")
	}
	if item.ID().IsZero() || item.ResumeID().IsZero() {
		return gormmodel.ResumeExperience{}, fmt.Errorf("experience identity is required")
	}
	model := gormmodel.ResumeExperience{
		ID:          item.ID().UUID(),
		ResumeID:    item.ResumeID().UUID(),
		Company:     item.Company(),
		Position:    item.Position(),
		Location:    item.Location(),
		StartDate:   YearMonthToTime(item.StartDate()),
		Description: item.Description(),
	}
	if item.EndDate() != nil {
		value := YearMonthToTime(*item.EndDate())
		model.EndDate = &value
	}
	return model, nil
}

func ExperienceToDomain(
	model gormmodel.ResumeExperience,
	achievements []gormmodel.ResumeExperienceAchievement,
	tasks []gormmodel.ResumeExperienceTask,
	stack []gormmodel.ResumeExperienceStack,
) (*domain.Experience, error) {
	id, err := identity.IDFromUUID(model.ID)
	if err != nil {
		return nil, err
	}
	resumeID, err := identity.IDFromUUID(model.ResumeID)
	if err != nil {
		return nil, err
	}
	start, err := domain.ParseYearMonth(model.StartDate.Format("2006-01"))
	if err != nil {
		return nil, err
	}
	var end *domain.YearMonth
	if model.EndDate != nil {
		value, err := domain.ParseYearMonth(model.EndDate.Format("2006-01"))
		if err != nil {
			return nil, err
		}
		end = &value
	}
	stackValues := make([]domain.ExperienceStackParams, 0, len(stack))
	for _, technology := range stack {
		var grade *domain.TechnologyGrade
		if technology.Grade != nil {
			value := domain.TechnologyGrade(*technology.Grade)
			if !value.IsValid() {
				return nil, domain.ErrInvalidEntity
			}
			grade = &value
		}
		stackValues = append(stackValues, domain.ExperienceStackParams{
			Name:  technology.Name,
			Grade: grade,
		})
	}
	achievementValues := make([]string, len(achievements))
	for i, achievement := range achievements {
		achievementValues[i] = achievement.Value
	}
	taskValues := make([]string, len(tasks))
	for i, task := range tasks {
		taskValues[i] = task.Value
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
		Tasks:        taskValues,
		Stack:        stackValues,
	})
}

func ExperienceAchievementsToDB(experienceID identity.ID, values []string) []gormmodel.ResumeExperienceAchievement {
	models := make([]gormmodel.ResumeExperienceAchievement, 0, len(values))
	for position, value := range values {
		models = append(models, gormmodel.ResumeExperienceAchievement{
			ExperienceID: experienceID.UUID(),
			Position:     position,
			Value:        value,
		})
	}
	return models
}

func ExperienceTasksToDB(experienceID identity.ID, values []string) []gormmodel.ResumeExperienceTask {
	models := make([]gormmodel.ResumeExperienceTask, 0, len(values))
	for position, value := range values {
		models = append(models, gormmodel.ResumeExperienceTask{
			ExperienceID: experienceID.UUID(),
			Position:     position,
			Value:        value,
		})
	}
	return models
}

func ExperienceStackToDB(experienceID identity.ID, values []domain.Technology) []gormmodel.ResumeExperienceStack {
	models := make([]gormmodel.ResumeExperienceStack, 0, len(values))
	for position, value := range values {
		models = append(models, gormmodel.ResumeExperienceStack{
			ExperienceID: experienceID.UUID(),
			Position:     position,
			Name:         value.Name(),
			Grade:        TechnologyGradeToDB(value.Grade()),
		})
	}
	return models
}

package mappers

import (
	"fmt"

	"github.com/google/uuid"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func ExperiencesToDomain(
	rows []gormmodel.ResumeExperience,
	achievements map[uuid.UUID][]gormmodel.ResumeExperienceAchievement,
	tasks map[uuid.UUID][]gormmodel.ResumeExperienceTask,
	stack map[uuid.UUID][]gormmodel.ResumeExperienceStack,
) ([]*domain.Experience, error) {
	experiences := make([]*domain.Experience, 0, len(rows))
	for i, row := range rows {
		experience, err := ExperienceToDomain(row, achievements[row.ID], tasks[row.ID], stack[row.ID])
		if err != nil {
			return nil, fmt.Errorf("map experience at index %d: %w", i, err)
		}
		experiences = append(experiences, experience)
	}
	return experiences, nil
}

func ExperienceToDB(item *domain.Experience) gormmodel.ResumeExperience {
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
	return model
}

func ExperienceToDomain(
	model gormmodel.ResumeExperience,
	achievements []gormmodel.ResumeExperienceAchievement,
	tasks []gormmodel.ResumeExperienceTask,
	stack []gormmodel.ResumeExperienceStack,
) (*domain.Experience, error) {
	id, err := sharedvalueobjects.IDFromUUID(model.ID)
	if err != nil {
		return nil, err
	}
	resumeID, err := sharedvalueobjects.IDFromUUID(model.ResumeID)
	if err != nil {
		return nil, err
	}
	start := domain.RestoreYearMonth(model.StartDate)
	var end *domain.YearMonth
	if model.EndDate != nil {
		value := domain.RestoreYearMonth(*model.EndDate)
		end = &value
	}
	stackValues := make([]domain.ExperienceStackParams, 0, len(stack))
	for _, technology := range stack {
		stackValues = append(stackValues, domain.ExperienceStackParams{Name: technology.Name})
	}
	achievementValues := make([]string, len(achievements))
	for i, achievement := range achievements {
		achievementValues[i] = achievement.Value
	}
	taskValues := make([]string, len(tasks))
	for i, task := range tasks {
		taskValues[i] = task.Value
	}
	return domain.RestoreExperience(domain.ExperienceSnapshot{
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
	}), nil
}

func ExperienceAchievementsToDB(experienceID sharedvalueobjects.ID, values []string) []gormmodel.ResumeExperienceAchievement {
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

func ExperienceTasksToDB(experienceID sharedvalueobjects.ID, values []string) []gormmodel.ResumeExperienceTask {
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

func ExperienceStackToDB(experienceID sharedvalueobjects.ID, values []domain.Technology) []gormmodel.ResumeExperienceStack {
	models := make([]gormmodel.ResumeExperienceStack, 0, len(values))
	for position, value := range values {
		models = append(models, gormmodel.ResumeExperienceStack{
			ExperienceID: experienceID.UUID(),
			Position:     position,
			Name:         value.Name(),
		})
	}
	return models
}

func ExperiencesToDBModels(items []*domain.Experience) (
	[]gormmodel.ResumeExperience,
	[]gormmodel.ResumeExperienceAchievement,
	[]gormmodel.ResumeExperienceTask,
	[]gormmodel.ResumeExperienceStack,
) {
	models := make([]gormmodel.ResumeExperience, 0, len(items))
	achievements := make([]gormmodel.ResumeExperienceAchievement, 0)
	tasks := make([]gormmodel.ResumeExperienceTask, 0)
	stack := make([]gormmodel.ResumeExperienceStack, 0)
	for position, item := range items {
		model := ExperienceToDB(item)
		model.SortOrder = position
		models = append(models, model)
		achievements = append(achievements, ExperienceAchievementsToDB(item.ID(), item.Achievements())...)
		tasks = append(tasks, ExperienceTasksToDB(item.ID(), item.Tasks())...)
		stack = append(stack, ExperienceStackToDB(item.ID(), item.Stack())...)
	}
	return models, achievements, tasks, stack
}

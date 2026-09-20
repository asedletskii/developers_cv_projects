package mappers

import (
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func ActivitiesToDomain(rows []gormmodel.ResumeActivity) ([]*domain.Activity, error) {
	activities := make([]*domain.Activity, 0, len(rows))
	for i, row := range rows {
		activity, err := ActivityToDomain(row)
		if err != nil {
			return nil, fmt.Errorf("map activity at index %d: %w", i, err)
		}
		activities = append(activities, activity)
	}
	return activities, nil
}

func ActivityToDB(item *domain.Activity) gormmodel.ResumeActivity {
	model := gormmodel.ResumeActivity{
		ID:          item.ID().UUID(),
		ResumeID:    item.ResumeID().UUID(),
		Type:        string(item.Type()),
		Title:       item.Title(),
		Description: item.Description(),
		Link:        ExternalURLToDB(item.Link()),
	}
	if item.Date() != nil {
		value := YearMonthToTime(*item.Date())
		model.Date = &value
	}
	return model
}

func ActivityToDomain(model gormmodel.ResumeActivity) (*domain.Activity, error) {
	id, err := sharedvalueobjects.IDFromUUID(model.ID)
	if err != nil {
		return nil, err
	}
	resumeID, err := sharedvalueobjects.IDFromUUID(model.ResumeID)
	if err != nil {
		return nil, err
	}
	activityType := domain.ActivityType(model.Type)
	var date *domain.YearMonth
	if model.Date != nil {
		value := domain.RestoreYearMonth(*model.Date)
		date = &value
	}
	var link *domain.ExternalURL
	if model.Link != nil {
		value := domain.RestoreExternalURL(*model.Link)
		link = &value
	}
	return domain.RestoreActivity(domain.ActivitySnapshot{
		ID:          id,
		ResumeID:    resumeID,
		Type:        activityType,
		Title:       model.Title,
		Description: model.Description,
		Date:        date,
		Link:        link,
	}), nil
}

func ActivitiesToDBModels(items []*domain.Activity) []gormmodel.ResumeActivity {
	models := make([]gormmodel.ResumeActivity, 0, len(items))
	for position, item := range items {
		model := ActivityToDB(item)
		model.SortOrder = position
		models = append(models, model)
	}
	return models
}

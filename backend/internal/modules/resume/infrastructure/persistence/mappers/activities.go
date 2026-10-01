package mappers

import (
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func ActivityToDB(item *domain.Activity) (gormmodel.ResumeActivity, error) {
	if item == nil {
		return gormmodel.ResumeActivity{}, fmt.Errorf("activity is nil")
	}
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
	return model, nil
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
	if !activityType.IsValid() {
		return nil, domain.ErrInvalidEntity
	}
	var date *domain.YearMonth
	if model.Date != nil {
		value, err := domain.ParseYearMonth(model.Date.Format("2006-01"))
		if err != nil {
			return nil, err
		}
		date = &value
	}
	var link *domain.ExternalURL
	if model.Link != nil {
		value, err := domain.NewExternalURL(*model.Link)
		if err != nil {
			return nil, err
		}
		link = &value
	}
	return domain.NewActivity(domain.ActivityParams{
		ID:          id,
		ResumeID:    resumeID,
		Type:        activityType,
		Title:       model.Title,
		Description: model.Description,
		Date:        date,
		Link:        link,
	})
}

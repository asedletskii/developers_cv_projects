package repository

import (
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

func persistActivities(tx *gorm.DB, resume *domain.Resume) error {
	models := make([]gormmodel.ResumeActivity, 0, len(resume.Activities()))
	for position, item := range resume.Activities() {
		if item == nil {
			return fmt.Errorf("%w: nil activity", domain.ErrResumeInvariant)
		}
		if item.ResumeID() != resume.ID() {
			return fmt.Errorf("%w: activity does not belong to resume", domain.ErrResumeInvariant)
		}
		model, err := activityModel(item)
		if err != nil {
			return err
		}
		model.SortOrder = position
		models = append(models, model)
	}

	if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeActivity{}).Error; err != nil {
		return fmt.Errorf("replace resume activities: %w", mapDBError(err))
	}
	if len(models) > 0 {
		if err := tx.Create(&models).Error; err != nil {
			return fmt.Errorf("create resume activities: %w", mapDBError(err))
		}
	}
	return nil
}

func loadActivities(tx *gorm.DB, resumeID valueobjects.ID) ([]*domain.Activity, error) {
	var rows []gormmodel.ResumeActivity
	if err := tx.Where("resume_id = ?", resumeID.UUID()).Order("sort_order").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("get resume activities: %w", mapDBError(err))
	}
	out := make([]*domain.Activity, 0, len(rows))
	for _, row := range rows {
		x, err := mapActivity(row)
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

func activityModel(item *domain.Activity) (gormmodel.ResumeActivity, error) {
	if item == nil {
		return gormmodel.ResumeActivity{}, fmt.Errorf("activity is nil")
	}
	model := gormmodel.ResumeActivity{
		ID: item.ID().UUID(), ResumeID: item.ResumeID().UUID(), Type: string(item.Type()), Title: item.Title(),
		Description: item.Description(), Link: externalURLString(item.Link()),
	}
	if item.Date() != nil {
		value := yearMonthDate(*item.Date())
		model.Date = &value
	}
	return model, nil
}

func mapActivity(model gormmodel.ResumeActivity) (*domain.Activity, error) {
	id, err := valueobjects.IDFromUUID(model.ID)
	if err != nil {
		return nil, err
	}
	resumeID, err := valueobjects.IDFromUUID(model.ResumeID)
	if err != nil {
		return nil, err
	}
	activityType := enums.ActivityType(model.Type)
	if !activityType.IsValid() {
		return nil, domain.ErrInvalidEntity
	}
	var date *valueobjects.YearMonth
	if model.Date != nil {
		value, err := valueobjects.ParseYearMonth(model.Date.Format("2006-01"))
		if err != nil {
			return nil, err
		}
		date = &value
	}
	var link *valueobjects.ExternalURL
	if model.Link != nil {
		value, err := valueobjects.NewExternalURL(*model.Link)
		if err != nil {
			return nil, err
		}
		link = &value
	}
	return domain.NewActivity(domain.ActivityParams{ID: id, ResumeID: resumeID, Type: activityType, Title: model.Title, Description: model.Description, Date: date, Link: link})
}

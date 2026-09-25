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

var _ ports.CreateActivityRepository = (*ActivityRepository)(nil)
var _ ports.GetActivityRepository = (*ActivityRepository)(nil)
var _ ports.UpdateActivityRepository = (*ActivityRepository)(nil)
var _ ports.DeleteActivityRepository = (*ActivityRepository)(nil)

func (r *ActivityRepository) CreateActivity(ctx context.Context, c ports.CreateActivityRepositoryCommand) error {
	tx, err := r.owned(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	item, err := activityModel(c.Activity)
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
func (r *ActivityRepository) GetActivity(ctx context.Context, q ports.GetActivityRepositoryCommand) (*domain.Activity, error) {
	tx, err := r.owned(ctx, q.ResumeID, q.UserID)
	if err != nil {
		return nil, err
	}
	var model gormmodel.ResumeActivity
	if err := tx.Where("id = ? AND resume_id = ?", q.ActivityID.UUID(), q.ResumeID.UUID()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, err
	}
	return mapActivity(model)
}
func (r *ActivityRepository) UpdateActivity(ctx context.Context, c ports.UpdateActivityRepositoryCommand) error {
	tx, err := r.owned(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	model, err := activityModel(c.Activity)
	if err != nil {
		return err
	}
	result := tx.Model(&gormmodel.ResumeActivity{}).Where("id = ? AND resume_id = ?", c.Activity.ID().UUID(), c.ResumeID.UUID()).Updates(map[string]any{"type": model.Type, "title": model.Title, "description": model.Description, "date": model.Date, "link": model.Link})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
}
func (r *ActivityRepository) DeleteActivity(ctx context.Context, c ports.DeleteActivityRepositoryCommand) error {
	tx, err := r.owned(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	result := tx.Where("id = ? AND resume_id = ?", c.ActivityID.UUID(), c.ResumeID.UUID()).Delete(&gormmodel.ResumeActivity{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
}
func activityModel(item *domain.Activity) (gormmodel.ResumeActivity, error) {
	model := gormmodel.ResumeActivity{ID: item.ID().UUID(), ResumeID: item.ResumeID().UUID(), Type: string(item.Type()), Title: item.Title(), Description: item.Description(), Link: externalURLString(item.Link())}
	if item.Date() != nil {
		value := item.Date().String() + "-01"
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
		value, err := valueobjects.ParseYearMonth((*model.Date)[:7])
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

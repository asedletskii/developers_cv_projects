package readers

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

func (r *ResumeReader) loadProfile(ctx context.Context, db *gorm.DB, resumeID uuid.UUID, view *queries.ResumeView) error {
	var model gormmodel.ResumeProfile
	if err := db.WithContext(ctx).Where("resume_id = ?", resumeID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return fmt.Errorf("get resume profile: %w", err)
	}
	var formats []gormmodel.ResumeProfileWorkFormat
	if err := db.WithContext(ctx).Where("profile_id = ?", model.ID).Order("position").Find(&formats).Error; err != nil {
		return fmt.Errorf("get profile work formats: %w", err)
	}
	workFormat := make([]string, 0, len(formats))
	for _, format := range formats {
		workFormat = append(workFormat, format.Value)
	}
	view.Profile = &queries.ResumeProfileView{
		ID: model.ID.String(), Name: model.Name, Surname: model.Surname,
		Patronymic: model.Patronymic, Position: model.Position, PhotoURL: model.PhotoURL,
		DateOfBirth: model.DateOfBirth.Format("2006-01-02"), Location: model.Location, Relocation: model.Relocation,
		WorkFormat: workFormat,
	}
	return nil
}

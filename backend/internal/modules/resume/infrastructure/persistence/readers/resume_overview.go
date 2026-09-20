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

func (r *ResumeReader) loadOverview(ctx context.Context, db *gorm.DB, resumeID uuid.UUID, view *queries.ResumeView) error {
	var model gormmodel.ResumeOverview
	if err := db.WithContext(ctx).Where("resume_id = ?", resumeID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			view.Overview = queries.ResumeOverviewView{
				Sections: make([]queries.OverviewSectionSettings, 0),
			}
			return nil
		}
		return fmt.Errorf("get overview: %w", err)
	}
	result := queries.ResumeOverviewView{
		Greeting: queries.OverviewGreeting{First: model.GreetingFirst, Second: model.GreetingSecond, SecondHighlight: model.GreetingHighlight, Third: model.GreetingThird},
		Sections: make([]queries.OverviewSectionSettings, 0),
	}
	if model.AvailabilityStatus != nil {
		result.Availability = &queries.OverviewAvailability{Status: *model.AvailabilityStatus, Note: model.AvailabilityNote}
		if model.AvailabilityFrom != nil {
			value := model.AvailabilityFrom.Format("2006-01-02")
			result.Availability.From = &value
		}
	}
	if model.ThemePrimary != nil || model.ThemeSecondary != nil {
		result.Theme = &queries.OverviewThemeSettings{Primary: valueOrEmpty(model.ThemePrimary), Secondary: valueOrEmpty(model.ThemeSecondary)}
	}
	var sections []gormmodel.ResumeOverviewSection
	if err := db.WithContext(ctx).Where("resume_id = ?", resumeID).Order("position").Find(&sections).Error; err != nil {
		return fmt.Errorf("get overview sections: %w", err)
	}
	for _, section := range sections {
		result.Sections = append(result.Sections, queries.OverviewSectionSettings{ID: section.SectionID, Enabled: section.Enabled, Collapsed: section.Collapsed, Style: section.Style})
	}
	view.Overview = result
	return nil
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

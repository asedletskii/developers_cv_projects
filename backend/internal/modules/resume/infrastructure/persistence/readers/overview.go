package readers

import (
	"context"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

type OverviewReader struct{ db *gorm.DB }

func NewOverviewReader(db *gorm.DB) *OverviewReader { return &OverviewReader{db: db} }

func (r *OverviewReader) GetOverview(ctx context.Context, q queries.GetOverviewQuery) (queries.ResumeOverviewView, error) {
	db, err := resolveDB(ctx, r.db, "get overview")
	if err != nil {
		return queries.ResumeOverviewView{}, err
	}
	if q.ResumeID.IsZero() || q.UserID.IsZero() {
		return queries.ResumeOverviewView{}, fmt.Errorf("get overview: invalid query")
	}

	var model gormmodel.ResumeOverview
	readDB := db.Joins("JOIN resumes ON resumes.id = resume_overviews.resume_id").Where("resume_overviews.resume_id = ?", q.ResumeID.UUID())
	readDB = readDB.Where("resumes.user_id = ?", q.UserID.UUID())
	err = readDB.First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return queries.ResumeOverviewView{}, queries.ErrResumeNotFound
	}
	if err != nil {
		return queries.ResumeOverviewView{}, fmt.Errorf("get overview: %w", err)
	}

	view := queries.ResumeOverviewView{Presentation: queries.OverviewPresentationView{
		Greeting:   queries.OverviewGreeting{First: model.GreetingFirst, Second: model.GreetingSecond, SecondHighlight: model.GreetingHighlight, Third: model.GreetingThird},
		ResumeLink: model.ResumeLink,
		Sections:   make([]queries.OverviewSectionSettings, 0),
	}}
	if model.AvailabilityStatus != nil {
		availability := &queries.OverviewAvailability{Status: *model.AvailabilityStatus, From: nil, Note: model.AvailabilityNote}
		if model.AvailabilityFrom != nil {
			value := model.AvailabilityFrom.Format("2006-01-02")
			availability.From = &value
		}
		view.Presentation.Availability = availability
	}
	if model.ThemePrimary != nil || model.ThemeSecondary != nil {
		view.Presentation.Theme = &queries.OverviewThemeSettings{Primary: valueOrEmpty(model.ThemePrimary), Secondary: valueOrEmpty(model.ThemeSecondary)}
	}

	var sections []gormmodel.ResumeOverviewSection
	if err := db.Where("resume_id = ?", model.ResumeID).Order("position").Find(&sections).Error; err != nil {
		return queries.ResumeOverviewView{}, fmt.Errorf("get overview sections: %w", err)
	}
	for _, section := range sections {
		view.Presentation.Sections = append(view.Presentation.Sections, queries.OverviewSectionSettings{ID: section.SectionID, Enabled: section.Enabled, Collapsed: section.Collapsed})
	}
	return view, nil
}

func (r *OverviewReader) ReadOverview(ctx context.Context, resumeID valueobjects.ID) (queries.ResumeOverviewView, error) {
	db, err := resolveDB(ctx, r.db, "get overview")
	if err != nil {
		return queries.ResumeOverviewView{}, err
	}
	if resumeID.IsZero() {
		return queries.ResumeOverviewView{}, fmt.Errorf("get overview: invalid query")
	}

	var model gormmodel.ResumeOverview
	readDB := db.Joins("JOIN resumes ON resumes.id = resume_overviews.resume_id").Where("resume_overviews.resume_id = ?", resumeID.UUID())
	err = readDB.First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return queries.ResumeOverviewView{}, queries.ErrResumeNotFound
	}
	if err != nil {
		return queries.ResumeOverviewView{}, fmt.Errorf("get overview: %w", err)
	}

	view := queries.ResumeOverviewView{Presentation: queries.OverviewPresentationView{
		Greeting:   queries.OverviewGreeting{First: model.GreetingFirst, Second: model.GreetingSecond, SecondHighlight: model.GreetingHighlight, Third: model.GreetingThird},
		ResumeLink: model.ResumeLink,
		Sections:   make([]queries.OverviewSectionSettings, 0),
	}}
	if model.AvailabilityStatus != nil {
		availability := &queries.OverviewAvailability{Status: *model.AvailabilityStatus, From: nil, Note: model.AvailabilityNote}
		if model.AvailabilityFrom != nil {
			value := model.AvailabilityFrom.Format("2006-01-02")
			availability.From = &value
		}
		view.Presentation.Availability = availability
	}
	if model.ThemePrimary != nil || model.ThemeSecondary != nil {
		view.Presentation.Theme = &queries.OverviewThemeSettings{Primary: valueOrEmpty(model.ThemePrimary), Secondary: valueOrEmpty(model.ThemeSecondary)}
	}

	var sections []gormmodel.ResumeOverviewSection
	if err := db.Where("resume_id = ?", model.ResumeID).Order("position").Find(&sections).Error; err != nil {
		return queries.ResumeOverviewView{}, fmt.Errorf("get overview sections: %w", err)
	}
	for _, section := range sections {
		view.Presentation.Sections = append(view.Presentation.Sections, queries.OverviewSectionSettings{ID: section.SectionID, Enabled: section.Enabled, Collapsed: section.Collapsed})
	}
	return view, nil
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

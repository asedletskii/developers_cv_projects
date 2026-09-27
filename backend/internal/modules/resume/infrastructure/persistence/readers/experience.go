package readers

import (
	"context"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"fmt"

	"github.com/google/uuid"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

type ExperiencesReader struct{ db *gorm.DB }

func NewExperiencesReader(db *gorm.DB) *ExperiencesReader { return &ExperiencesReader{db: db} }

type experienceRow struct {
	ID          string
	ResumeID    string
	Company     string
	Position    string
	Location    *string
	StartDate   string
	EndDate     *string
	Description *string
}

func (r *ExperiencesReader) ListExperiences(ctx context.Context, q queries.ListExperiencesQuery) ([]queries.ExperienceView, error) {
	db, err := resolveDB(ctx, r.db, "list experiences")
	if err != nil {
		return nil, err
	}
	if q.ResumeID.IsZero() || q.UserID.IsZero() {
		return nil, fmt.Errorf("list experiences: invalid query")
	}
	var rows []experienceRow
	readDB := db.Table("resume_experiences AS e").Select("e.id::text AS id, e.resume_id::text AS resume_id, e.company, e.position, e.location, to_char(e.start_date, 'YYYY-MM') AS start_date, CASE WHEN e.end_date IS NULL THEN NULL ELSE to_char(e.end_date, 'YYYY-MM') END AS end_date, e.description").Joins("JOIN resumes AS r ON r.id = e.resume_id").Where("e.resume_id = ?", q.ResumeID.UUID())
	readDB = readDB.Where("r.user_id = ?", q.UserID.UUID())
	err = readDB.Order("e.sort_order").Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("list experiences: %w", err)
	}
	out := make([]queries.ExperienceView, 0, len(rows))
	for _, row := range rows {
		view := queries.ExperienceView{ID: row.ID, ResumeID: row.ResumeID, Company: row.Company, Position: row.Position, Location: row.Location, StartDate: row.StartDate, EndDate: row.EndDate, Description: row.Description}
		if err := r.loadExperienceChildren(ctx, db, row.ID, &view); err != nil {
			return nil, err
		}
		out = append(out, view)
	}
	return out, nil
}

func (r *ExperiencesReader) ReadExperiences(ctx context.Context, resumeID valueobjects.ID) ([]queries.ExperienceView, error) {
	db, err := resolveDB(ctx, r.db, "list experiences")
	if err != nil {
		return nil, err
	}
	if resumeID.IsZero() {
		return nil, fmt.Errorf("list experiences: invalid query")
	}
	var rows []experienceRow
	readDB := db.Table("resume_experiences AS e").Select("e.id::text AS id, e.resume_id::text AS resume_id, e.company, e.position, e.location, to_char(e.start_date, 'YYYY-MM') AS start_date, CASE WHEN e.end_date IS NULL THEN NULL ELSE to_char(e.end_date, 'YYYY-MM') END AS end_date, e.description").Joins("JOIN resumes AS r ON r.id = e.resume_id").Where("e.resume_id = ?", resumeID.UUID())
	err = readDB.Order("e.sort_order").Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("list experiences: %w", err)
	}
	out := make([]queries.ExperienceView, 0, len(rows))
	for _, row := range rows {
		view := queries.ExperienceView{ID: row.ID, ResumeID: row.ResumeID, Company: row.Company, Position: row.Position, Location: row.Location, StartDate: row.StartDate, EndDate: row.EndDate, Description: row.Description}
		if err := r.loadExperienceChildren(ctx, db, row.ID, &view); err != nil {
			return nil, err
		}
		out = append(out, view)
	}
	return out, nil
}

func (r *ExperiencesReader) GetExperience(ctx context.Context, q queries.GetExperienceQuery) (queries.ExperienceView, error) {
	db, err := resolveDB(ctx, r.db, "get experience")
	if err != nil {
		return queries.ExperienceView{}, err
	}
	if q.ResumeID.IsZero() || q.ExperienceID.IsZero() || q.UserID.IsZero() {
		return queries.ExperienceView{}, fmt.Errorf("get experience: invalid query")
	}
	var row experienceRow
	readDB := db.Table("resume_experiences AS e").Select("e.id::text AS id, e.resume_id::text AS resume_id, e.company, e.position, e.location, to_char(e.start_date, 'YYYY-MM') AS start_date, CASE WHEN e.end_date IS NULL THEN NULL ELSE to_char(e.end_date, 'YYYY-MM') END AS end_date, e.description").Joins("JOIN resumes AS r ON r.id = e.resume_id").Where("e.id = ? AND e.resume_id = ?", q.ExperienceID.UUID(), q.ResumeID.UUID())
	readDB = readDB.Where("r.user_id = ?", q.UserID.UUID())
	err = readDB.Scan(&row).Error
	if err != nil {
		return queries.ExperienceView{}, fmt.Errorf("get experience: %w", err)
	}
	if row.ID == "" {
		return queries.ExperienceView{}, queries.ErrResumeNotFound
	}
	view := queries.ExperienceView{ID: row.ID, ResumeID: row.ResumeID, Company: row.Company, Position: row.Position, Location: row.Location, StartDate: row.StartDate, EndDate: row.EndDate, Description: row.Description}
	if err := r.loadExperienceChildren(ctx, db, row.ID, &view); err != nil {
		return queries.ExperienceView{}, err
	}
	return view, nil
}

func (r *ExperiencesReader) loadExperienceChildren(ctx context.Context, db *gorm.DB, experienceID string, view *queries.ExperienceView) error {
	id, err := uuidParse(experienceID)
	if err != nil {
		return fmt.Errorf("load experience children: %w", err)
	}
	var achievements []gormmodel.ResumeExperienceAchievement
	if err := db.WithContext(ctx).Where("experience_id = ?", id).Order("position").Find(&achievements).Error; err != nil {
		return fmt.Errorf("get experience achievements: %w", err)
	}
	view.Achievements = make([]string, 0, len(achievements))
	for _, item := range achievements {
		view.Achievements = append(view.Achievements, item.Value)
	}
	var stack []gormmodel.ResumeExperienceStack
	if err := db.WithContext(ctx).Where("experience_id = ?", id).Order("position").Find(&stack).Error; err != nil {
		return fmt.Errorf("get experience stack: %w", err)
	}
	view.Stack = make([]queries.ExperienceStackItemView, 0, len(stack))
	for _, item := range stack {
		view.Stack = append(view.Stack, queries.ExperienceStackItemView{Name: item.Name, Grade: item.Grade})
	}
	return nil
}

func uuidParse(value string) (uuid.UUID, error) {
	return uuid.Parse(value)
}

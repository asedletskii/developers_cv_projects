package readers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

type experienceRow struct {
	ID, ResumeID, Company, Position string
	Location                        *string
	StartDate                       string
	EndDate                         *string
	Description                     *string
}

func (r *ResumeReader) loadExperiences(ctx context.Context, db *gorm.DB, resumeID uuid.UUID, view *queries.ResumeView) error {
	var rows []experienceRow
	if err := db.WithContext(ctx).Table("resume_experiences AS e").Select("e.id::text AS id, e.resume_id::text AS resume_id, e.company, e.position, e.location, to_char(e.start_date, 'YYYY-MM') AS start_date, CASE WHEN e.end_date IS NULL THEN NULL ELSE to_char(e.end_date, 'YYYY-MM') END AS end_date, e.description").Where("e.resume_id = ?", resumeID).Order("e.sort_order").Scan(&rows).Error; err != nil {
		return fmt.Errorf("list experiences: %w", err)
	}
	result := make([]queries.ExperienceView, 0, len(rows))
	if len(rows) == 0 {
		view.Experiences = result
		return nil
	}
	ids := make([]uuid.UUID, 0, len(rows))
	for _, row := range rows {
		id, err := uuid.Parse(row.ID)
		if err != nil {
			return fmt.Errorf("list experiences: invalid experience id: %w", err)
		}
		ids = append(ids, id)
	}
	achievements := make(map[uuid.UUID][]string, len(ids))
	var achievementRows []gormmodel.ResumeExperienceAchievement
	if err := db.WithContext(ctx).Where("experience_id IN ?", ids).Order("experience_id, position").Find(&achievementRows).Error; err != nil {
		return fmt.Errorf("get experience achievements: %w", err)
	}
	for _, item := range achievementRows {
		achievements[item.ExperienceID] = append(achievements[item.ExperienceID], item.Value)
	}
	tasks := make(map[uuid.UUID][]string, len(ids))
	var taskRows []gormmodel.ResumeExperienceTask
	if err := db.WithContext(ctx).Where("experience_id IN ?", ids).Order("experience_id, position").Find(&taskRows).Error; err != nil {
		return fmt.Errorf("get experience tasks: %w", err)
	}
	for _, item := range taskRows {
		tasks[item.ExperienceID] = append(tasks[item.ExperienceID], item.Value)
	}
	stacks := make(map[uuid.UUID][]queries.ExperienceStackItemView, len(ids))
	var stackRows []gormmodel.ResumeExperienceStack
	if err := db.WithContext(ctx).Where("experience_id IN ?", ids).Order("experience_id, position").Find(&stackRows).Error; err != nil {
		return fmt.Errorf("get experience stack: %w", err)
	}
	for _, item := range stackRows {
		stacks[item.ExperienceID] = append(stacks[item.ExperienceID], queries.ExperienceStackItemView{Name: item.Name, Grade: item.Grade})
	}
	for _, row := range rows {
		id, err := uuid.Parse(row.ID)
		if err != nil {
			return err
		}
		result = append(result, queries.ExperienceView{ID: row.ID, ResumeID: row.ResumeID, Company: row.Company, Position: row.Position, Location: row.Location, StartDate: row.StartDate, EndDate: row.EndDate, Description: row.Description, Achievements: append([]string(nil), achievements[id]...), Tasks: append([]string(nil), tasks[id]...), Stack: append([]queries.ExperienceStackItemView(nil), stacks[id]...)})
	}
	for i := range result {
		if result[i].Achievements == nil {
			result[i].Achievements = make([]string, 0)
		}
		if result[i].Tasks == nil {
			result[i].Tasks = make([]string, 0)
		}
		if result[i].Stack == nil {
			result[i].Stack = make([]queries.ExperienceStackItemView, 0)
		}
	}
	view.Experiences = result
	return nil
}

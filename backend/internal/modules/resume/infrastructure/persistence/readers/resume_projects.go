package readers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

type projectRow struct {
	ID, ResumeID, Name               string
	Description, Role, Link, GitLink *string
}

func int64Ptr(value int64) *int64 { return &value }

func (r *ResumeReader) loadProjects(ctx context.Context, db *gorm.DB, resumeID uuid.UUID, view *queries.ResumeView) error {
	var rows []projectRow
	if err := db.WithContext(ctx).Table("resume_projects AS p").Select("p.id::text AS id, p.resume_id::text AS resume_id, p.name, p.description, p.role, p.link, p.git_link").Where("p.resume_id = ?", resumeID).Order("p.sort_order").Scan(&rows).Error; err != nil {
		return fmt.Errorf("list projects: %w", err)
	}
	result := make([]queries.ProjectView, 0, len(rows))
	if len(rows) == 0 {
		view.Projects = result
		return nil
	}
	ids := make([]uuid.UUID, 0, len(rows))
	for _, row := range rows {
		id, err := uuid.Parse(row.ID)
		if err != nil {
			return fmt.Errorf("list projects: invalid project id: %w", err)
		}
		ids = append(ids, id)
	}
	screenshots := make(map[uuid.UUID][]string, len(ids))
	var screenshotRows []gormmodel.ResumeProjectScreenshot
	if err := db.WithContext(ctx).Where("project_id IN ?", ids).Order("project_id, position").Find(&screenshotRows).Error; err != nil {
		return fmt.Errorf("list project screenshots: %w", err)
	}
	for _, item := range screenshotRows {
		screenshots[item.ProjectID] = append(screenshots[item.ProjectID], item.Path)
	}
	stacks := make(map[uuid.UUID][]queries.ProjectTechnologyView, len(ids))
	var stackRows []gormmodel.ResumeProjectTechnology
	if err := db.WithContext(ctx).Where("project_id IN ?", ids).Order("project_id, position").Find(&stackRows).Error; err != nil {
		return fmt.Errorf("list project technologies: %w", err)
	}
	for _, item := range stackRows {
		stacks[item.ProjectID] = append(stacks[item.ProjectID], queries.ProjectTechnologyView{Name: item.Name, Grade: item.Grade})
	}
	contributors := make(map[uuid.UUID][]queries.ProjectContributorView, len(ids))
	var contributorRows []gormmodel.ResumeProjectContributor
	if err := db.WithContext(ctx).Where("project_id IN ?", ids).Order("project_id, position").Find(&contributorRows).Error; err != nil {
		return fmt.Errorf("list project contributors: %w", err)
	}
	for _, item := range contributorRows {
		contributors[item.ProjectID] = append(contributors[item.ProjectID], queries.ProjectContributorView{Name: item.Name, Surname: item.Surname, Role: item.Role, Contacts: queries.ProjectContactsView{Phone: item.Phone, Telegram: item.Telegram, Email: item.Email}})
	}
	for _, row := range rows {
		id, err := uuid.Parse(row.ID)
		if err != nil {
			return err
		}
		result = append(result, queries.ProjectView{ID: row.ID, ResumeID: row.ResumeID, Name: row.Name, Description: row.Description, Role: row.Role, Link: row.Link, GitLink: row.GitLink, Screenshots: append([]string(nil), screenshots[id]...), Stack: append([]queries.ProjectTechnologyView(nil), stacks[id]...), Contributors: append([]queries.ProjectContributorView(nil), contributors[id]...)})
	}
	for i := range result {
		if result[i].Screenshots == nil {
			result[i].Screenshots = make([]string, 0)
		}
		if result[i].Stack == nil {
			result[i].Stack = make([]queries.ProjectTechnologyView, 0)
		}
		if result[i].Contributors == nil {
			result[i].Contributors = make([]queries.ProjectContributorView, 0)
		}
	}
	view.Projects = result
	return nil
}

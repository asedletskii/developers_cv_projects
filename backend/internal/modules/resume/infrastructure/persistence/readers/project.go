package readers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

type ProjectsReader struct{ db *gorm.DB }

func NewProjectsReader(db *gorm.DB) *ProjectsReader { return &ProjectsReader{db: db} }

type projectRow struct {
	ID          string
	ResumeID    string
	Name        string
	Description *string
	Role        *string
	Link        *string
	GitLink     *string
}

func (r *ProjectsReader) ListProjects(ctx context.Context, q queries.ListProjectsQuery) ([]queries.ProjectView, error) {
	db, err := resolveDB(ctx, r.db, "list projects")
	if err != nil {
		return nil, err
	}
	if q.ResumeID.IsZero() || q.UserID.IsZero() {
		return nil, fmt.Errorf("list projects: invalid query")
	}
	var rows []projectRow
	if err := db.Table("resume_projects AS p").Select("p.id::text AS id, p.resume_id::text AS resume_id, p.name, p.description, p.role, p.link, p.git_link").Joins("JOIN resumes AS r ON r.id = p.resume_id").Where("p.resume_id = ? AND r.user_id = ?", q.ResumeID.UUID(), q.UserID.UUID()).Order("p.sort_order").Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	out := make([]queries.ProjectView, 0, len(rows))
	for _, row := range rows {
		view := queries.ProjectView{ID: row.ID, ResumeID: row.ResumeID, Name: row.Name, Description: row.Description, Role: row.Role, Link: row.Link, GitLink: row.GitLink}
		if err := r.loadProjectChildren(ctx, db, row.ID, &view); err != nil {
			return nil, err
		}
		out = append(out, view)
	}
	return out, nil
}

func (r *ProjectsReader) GetProject(ctx context.Context, q queries.GetProjectQuery) (queries.ProjectView, error) {
	db, err := resolveDB(ctx, r.db, "get project")
	if err != nil {
		return queries.ProjectView{}, err
	}
	if q.ResumeID.IsZero() || q.ProjectID.IsZero() || q.UserID.IsZero() {
		return queries.ProjectView{}, fmt.Errorf("get project: invalid query")
	}
	var row projectRow
	if err := db.Table("resume_projects AS p").Select("p.id::text AS id, p.resume_id::text AS resume_id, p.name, p.description, p.role, p.link, p.git_link").Joins("JOIN resumes AS r ON r.id = p.resume_id").Where("p.id = ? AND p.resume_id = ? AND r.user_id = ?", q.ProjectID.UUID(), q.ResumeID.UUID(), q.UserID.UUID()).Scan(&row).Error; err != nil {
		return queries.ProjectView{}, fmt.Errorf("get project: %w", err)
	}
	if row.ID == "" {
		return queries.ProjectView{}, queries.ErrResumeNotFound
	}
	view := queries.ProjectView{ID: row.ID, ResumeID: row.ResumeID, Name: row.Name, Description: row.Description, Role: row.Role, Link: row.Link, GitLink: row.GitLink}
	if err := r.loadProjectChildren(ctx, db, row.ID, &view); err != nil {
		return queries.ProjectView{}, err
	}
	return view, nil
}

func (r *ProjectsReader) loadProjectChildren(ctx context.Context, db *gorm.DB, projectID string, view *queries.ProjectView) error {
	id, err := uuid.Parse(projectID)
	if err != nil {
		return fmt.Errorf("load project children: %w", err)
	}
	var screenshots []gormmodel.ResumeProjectScreenshot
	if err := db.WithContext(ctx).Where("project_id = ?", id).Order("position").Find(&screenshots).Error; err != nil {
		return fmt.Errorf("get project screenshots: %w", err)
	}
	view.Screenshots = make([]string, 0, len(screenshots))
	for _, item := range screenshots {
		view.Screenshots = append(view.Screenshots, item.Path)
	}
	var technologies []gormmodel.ResumeProjectTechnology
	if err := db.WithContext(ctx).Where("project_id = ?", id).Order("position").Find(&technologies).Error; err != nil {
		return fmt.Errorf("get project technologies: %w", err)
	}
	view.Stack = make([]queries.ProjectTechnologyView, 0, len(technologies))
	for _, item := range technologies {
		view.Stack = append(view.Stack, queries.ProjectTechnologyView{Name: item.Name, Grade: item.Grade})
	}
	var contributors []gormmodel.ResumeProjectContributor
	if err := db.WithContext(ctx).Where("project_id = ?", id).Order("position").Find(&contributors).Error; err != nil {
		return fmt.Errorf("get project contributors: %w", err)
	}
	view.Contributors = make([]queries.ProjectContributorView, 0, len(contributors))
	for _, item := range contributors {
		view.Contributors = append(view.Contributors, queries.ProjectContributorView{Name: item.Name, Surname: item.Surname, Role: item.Role, Contacts: queries.ProjectContactsView{Phone: item.Phone, Telegram: item.Telegram, Email: item.Email}})
	}
	return nil
}

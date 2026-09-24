package readers

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"developers-cv/backend/internal/modules/resume/application/queries"
)

type ProjectsReader struct{ pool *pgxpool.Pool }

func NewProjectsReader(pool *pgxpool.Pool) *ProjectsReader {
	return &ProjectsReader{pool: pool}
}

func (r *ProjectsReader) ListProjects(ctx context.Context, q queries.ListProjectsQuery) ([]queries.ProjectView, error) {
	if r == nil || r.pool == nil {
		return nil, fmt.Errorf("list projects: reader dependencies are nil")
	}
	rows, err := r.pool.Query(ctx, `
		SELECT p.id::text,p.resume_id::text,p.name,p.description,p.role,p.link,p.git_link
		FROM resume_projects p
		JOIN resume_roots rr ON rr.id=p.resume_id
		WHERE p.resume_id=$1 AND rr.user_id=$2
		ORDER BY p.sort_order`, q.ResumeID.UUID(), q.UserID.UUID())
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	defer rows.Close()

	out := make([]queries.ProjectView, 0)
	for rows.Next() {
		var v queries.ProjectView
		if err := rows.Scan(&v.ID, &v.ResumeID, &v.Name, &v.Description, &v.Role, &v.Link, &v.GitLink); err != nil {
			return nil, fmt.Errorf("scan project: %w", err)
		}
		if err := r.loadChildren(ctx, v.ID, &v); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate projects: %w", err)
	}
	return out, nil
}

func (r *ProjectsReader) GetProject(ctx context.Context, q queries.GetProjectQuery) (queries.ProjectView, error) {
	if r == nil || r.pool == nil {
		return queries.ProjectView{}, fmt.Errorf("get project: reader dependencies are nil")
	}
	var v queries.ProjectView
	err := r.pool.QueryRow(ctx, `
		SELECT p.id::text,p.resume_id::text,p.name,p.description,p.role,p.link,p.git_link
		FROM resume_projects p
		JOIN resume_roots rr ON rr.id=p.resume_id
		WHERE p.id=$1 AND p.resume_id=$2 AND rr.user_id=$3`, q.ProjectID.UUID(), q.ResumeID.UUID(), q.UserID.UUID()).
		Scan(&v.ID, &v.ResumeID, &v.Name, &v.Description, &v.Role, &v.Link, &v.GitLink)
	if errors.Is(err, pgx.ErrNoRows) {
		return v, queries.ErrResumeNotFound
	}
	if err != nil {
		return v, fmt.Errorf("get project: %w", err)
	}
	if err := r.loadChildren(ctx, v.ID, &v); err != nil {
		return queries.ProjectView{}, err
	}
	return v, nil
}

func (r *ProjectsReader) loadChildren(ctx context.Context, projectID string, view *queries.ProjectView) error {
	rows, err := r.pool.Query(ctx, `SELECT path FROM resume_project_screenshots WHERE project_id=$1 ORDER BY position`, projectID)
	if err != nil {
		return fmt.Errorf("load project screenshots: %w", err)
	}
	view.Screenshots = make([]string, 0)
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			rows.Close()
			return fmt.Errorf("scan project screenshot: %w", err)
		}
		view.Screenshots = append(view.Screenshots, path)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("iterate project screenshots: %w", err)
	}
	rows.Close()

	rows, err = r.pool.Query(ctx, `SELECT name,grade FROM resume_project_technologies WHERE project_id=$1 ORDER BY position`, projectID)
	if err != nil {
		return fmt.Errorf("load project technologies: %w", err)
	}
	view.Stack = make([]queries.ProjectTechnologyView, 0)
	for rows.Next() {
		var item queries.ProjectTechnologyView
		if err := rows.Scan(&item.Name, &item.Grade); err != nil {
			rows.Close()
			return fmt.Errorf("scan project technology: %w", err)
		}
		view.Stack = append(view.Stack, item)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("iterate project technologies: %w", err)
	}
	rows.Close()

	rows, err = r.pool.Query(ctx, `SELECT name,surname,role,phone,telegram,email FROM resume_project_contributors WHERE project_id=$1 ORDER BY position`, projectID)
	if err != nil {
		return fmt.Errorf("load project contributors: %w", err)
	}
	view.Contributors = make([]queries.ProjectContributorView, 0)
	for rows.Next() {
		var item queries.ProjectContributorView
		var phone, telegram, email *string
		if err := rows.Scan(&item.Name, &item.Surname, &item.Role, &phone, &telegram, &email); err != nil {
			rows.Close()
			return fmt.Errorf("scan project contributor: %w", err)
		}
		item.Contacts.Phone = phone
		item.Contacts.Telegram = telegram
		item.Contacts.Email = email
		view.Contributors = append(view.Contributors, item)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("iterate project contributors: %w", err)
	}
	rows.Close()
	return nil
}

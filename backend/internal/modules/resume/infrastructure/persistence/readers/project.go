package readers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"developers-cv/backend/internal/modules/resume/application/queries"
)

type ProjectsReader struct{ pool *pgxpool.Pool }

func NewProjectsReader(pool *pgxpool.Pool) *ProjectsReader { return &ProjectsReader{pool: pool} }

const projectSelect = `
SELECT
 p.id::text,
 p.resume_id::text,
 p.name,
 p.description,
 p.role,
 p.link,
 p.git_link,
 COALESCE((SELECT jsonb_agg(s.path ORDER BY s.position) FROM resume_project_screenshots s WHERE s.project_id=p.id), '[]'::jsonb),
 COALESCE((SELECT jsonb_agg(jsonb_build_object('name', t.name, 'grade', t.grade) ORDER BY t.position) FROM resume_project_technologies t WHERE t.project_id=p.id), '[]'::jsonb),
 COALESCE((SELECT jsonb_agg(jsonb_build_object('name', c.name, 'surname', c.surname, 'role', c.role, 'contacts', jsonb_build_object('phone', c.phone, 'telegram', c.telegram, 'email', c.email)) ORDER BY c.position) FROM resume_project_contributors c WHERE c.project_id=p.id), '[]'::jsonb)
FROM resume_projects p
JOIN resume_roots rr ON rr.id=p.resume_id
`

func (r *ProjectsReader) ListProjects(ctx context.Context, q queries.ListProjectsQuery) ([]queries.ProjectView, error) {
	if r == nil || r.pool == nil {
		return nil, fmt.Errorf("list projects: reader dependencies are nil")
	}
	rows, err := r.pool.Query(ctx, projectSelect+`WHERE p.resume_id=$1 AND rr.user_id=$2 ORDER BY p.sort_order`, q.ResumeID.UUID(), q.UserID.UUID())
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	defer rows.Close()
	out := make([]queries.ProjectView, 0)
	for rows.Next() {
		var v queries.ProjectView
		var screenshots, stack, contributors []byte
		if err := rows.Scan(&v.ID, &v.ResumeID, &v.Name, &v.Description, &v.Role, &v.Link, &v.GitLink, &screenshots, &stack, &contributors); err != nil {
			return nil, fmt.Errorf("scan project: %w", err)
		}
		if err := decodeProjectChildren(screenshots, stack, contributors, &v); err != nil {
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
	var screenshots, stack, contributors []byte
	err := r.pool.QueryRow(ctx, projectSelect+`WHERE p.id=$1 AND p.resume_id=$2 AND rr.user_id=$3`, q.ProjectID.UUID(), q.ResumeID.UUID(), q.UserID.UUID()).
		Scan(&v.ID, &v.ResumeID, &v.Name, &v.Description, &v.Role, &v.Link, &v.GitLink, &screenshots, &stack, &contributors)
	if errors.Is(err, pgx.ErrNoRows) {
		return v, queries.ErrResumeNotFound
	}
	if err != nil {
		return v, fmt.Errorf("get project: %w", err)
	}
	if err := decodeProjectChildren(screenshots, stack, contributors, &v); err != nil {
		return queries.ProjectView{}, err
	}
	return v, nil
}

func decodeProjectChildren(screenshots, stack, contributors []byte, view *queries.ProjectView) error {
	if err := json.Unmarshal(screenshots, &view.Screenshots); err != nil {
		return fmt.Errorf("decode project screenshots: %w", err)
	}
	if err := json.Unmarshal(stack, &view.Stack); err != nil {
		return fmt.Errorf("decode project technologies: %w", err)
	}
	if err := json.Unmarshal(contributors, &view.Contributors); err != nil {
		return fmt.Errorf("decode project contributors: %w", err)
	}
	if view.Screenshots == nil {
		view.Screenshots = []string{}
	}
	if view.Stack == nil {
		view.Stack = []queries.ProjectTechnologyView{}
	}
	if view.Contributors == nil {
		view.Contributors = []queries.ProjectContributorView{}
	}
	return nil
}

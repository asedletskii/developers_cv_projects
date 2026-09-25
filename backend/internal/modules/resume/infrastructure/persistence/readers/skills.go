package readers

import (
	"context"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SkillsReader struct{ pool *pgxpool.Pool }

func NewSkillsReader(pool *pgxpool.Pool) *SkillsReader { return &SkillsReader{pool: pool} }
func (r *SkillsReader) GetSkills(ctx context.Context, q queries.GetSkillsQuery) (queries.SkillsView, error) {
	if r == nil || r.pool == nil {
		return queries.SkillsView{}, fmt.Errorf("get skills: reader dependencies are nil")
	}
	var v queries.SkillsView
	if err := r.pool.QueryRow(ctx, `SELECT s.id::text,s.resume_id::text FROM resume_skills s JOIN resume_roots r ON r.id=s.resume_id WHERE s.resume_id=$1 AND r.user_id=$2`, q.ResumeID.UUID(), q.UserID.UUID()).Scan(&v.ID, &v.ResumeID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return v, queries.ErrResumeNotFound
		}
		return v, fmt.Errorf("get skills: %w", err)
	}
	rows, err := r.pool.Query(ctx, `SELECT value FROM resume_skill_hard WHERE skills_id=(SELECT id FROM resume_skills WHERE resume_id=$1) ORDER BY position`, q.ResumeID.UUID())
	if err != nil {
		return v, fmt.Errorf("get hard skills: %w", err)
	}
	defer rows.Close()
	v.HardSkills = []string{}
	for rows.Next() {
		var x string
		if err := rows.Scan(&x); err != nil {
			return v, fmt.Errorf("scan hard skill: %w", err)
		}
		v.HardSkills = append(v.HardSkills, x)
	}
	if err := rows.Err(); err != nil {
		return v, fmt.Errorf("iterate hard skills: %w", err)
	}
	rows2, err := r.pool.Query(ctx, `SELECT value FROM resume_skill_soft WHERE skills_id=(SELECT id FROM resume_skills WHERE resume_id=$1) ORDER BY position`, q.ResumeID.UUID())
	if err != nil {
		return v, fmt.Errorf("get soft skills: %w", err)
	}
	defer rows2.Close()
	v.SoftSkills = []string{}
	for rows2.Next() {
		var x string
		if err := rows2.Scan(&x); err != nil {
			return v, fmt.Errorf("scan soft skill: %w", err)
		}
		v.SoftSkills = append(v.SoftSkills, x)
	}
	if err := rows2.Err(); err != nil {
		return v, fmt.Errorf("iterate soft skills: %w", err)
	}
	return v, nil
}

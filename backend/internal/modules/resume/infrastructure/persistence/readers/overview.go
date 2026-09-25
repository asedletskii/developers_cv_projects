package readers

import (
	"context"
	"errors"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/queries"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OverviewReader struct{ pool *pgxpool.Pool }

func NewOverviewReader(pool *pgxpool.Pool) *OverviewReader { return &OverviewReader{pool: pool} }

func (r *OverviewReader) GetOverview(ctx context.Context, q queries.GetOverviewQuery) (queries.ResumeOverviewView, error) {
	if r == nil || r.pool == nil {
		return queries.ResumeOverviewView{}, fmt.Errorf("get overview: reader dependencies are nil")
	}
	var v queries.ResumeOverviewView
	var from, note, link, tp, ts, status *string
	var first, second, highlight, third string
	err := r.pool.QueryRow(ctx, `SELECT p.greeting_first,p.greeting_second,p.greeting_highlight,p.greeting_third,p.availability_status,p.availability_from::text,p.availability_note,p.resume_link,p.theme_primary,p.theme_secondary FROM resume_overviews p JOIN resumes r ON r.id=p.resume_id WHERE p.resume_id=$1 AND r.user_id=$2`, q.ResumeID.UUID(), q.UserID.UUID()).Scan(&first, &second, &highlight, &third, &status, &from, &note, &link, &tp, &ts)
	if errors.Is(err, pgx.ErrNoRows) {
		return queries.ResumeOverviewView{}, queries.ErrResumeNotFound
	}
	if err != nil {
		return queries.ResumeOverviewView{}, fmt.Errorf("get overview: %w", err)
	}
	v.Presentation.Greeting = queries.OverviewGreeting{First: first, Second: second, SecondHighlight: highlight, Third: third}
	if status != nil {
		v.Presentation.Availability = &queries.OverviewAvailability{Status: *status, From: from, Note: note}
	}
	v.Presentation.ResumeLink = link
	if tp != nil || ts != nil {
		v.Presentation.Theme = &queries.OverviewThemeSettings{Primary: valueOrEmpty(tp), Secondary: valueOrEmpty(ts)}
	}
	rows, err := r.pool.Query(ctx, `SELECT s.section_id,s.enabled,s.collapsed FROM resume_overview_sections s JOIN resumes r ON r.id=s.resume_id WHERE s.resume_id=$1 AND r.user_id=$2 ORDER BY s.position`, q.ResumeID.UUID(), q.UserID.UUID())
	if err != nil {
		return queries.ResumeOverviewView{}, fmt.Errorf("get overview sections: %w", err)
	}
	defer rows.Close()
	v.Presentation.Sections = make([]queries.OverviewSectionSettings, 0)
	for rows.Next() {
		var s queries.OverviewSectionSettings
		if err = rows.Scan(&s.ID, &s.Enabled, &s.Collapsed); err != nil {
			return queries.ResumeOverviewView{}, fmt.Errorf("scan overview section: %w", err)
		}
		v.Presentation.Sections = append(v.Presentation.Sections, s)
	}
	if err = rows.Err(); err != nil {
		return queries.ResumeOverviewView{}, fmt.Errorf("iterate overview sections: %w", err)
	}
	return v, nil
}

func valueOrEmpty(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

package readers

import (
	"context"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfileReader struct{ pool *pgxpool.Pool }

func NewProfileReader(pool *pgxpool.Pool) *ProfileReader { return &ProfileReader{pool: pool} }
func (r *ProfileReader) GetResumeProfile(ctx context.Context, q queries.GetResumeProfileQuery) (queries.ResumeProfileView, error) {
	if r == nil || r.pool == nil {
		return queries.ResumeProfileView{}, fmt.Errorf("get resume profile: reader dependencies are nil")
	}
	var v queries.ResumeProfileView
	var dob time.Time
	var patronymic *string
	if err := r.pool.QueryRow(ctx, `SELECT p.id::text,p.resume_id::text,p.name,p.surname,p.patronymic,p.position,p.photo_path,p.date_of_birth,p.location,p.relocation FROM resume_profiles p JOIN resume_roots r ON r.id=p.resume_id WHERE p.resume_id=$1 AND r.user_id=$2`, q.ResumeID.UUID(), q.UserID.UUID()).Scan(&v.ID, &v.ResumeID, &v.Name, &v.Surname, &patronymic, &v.Position, &v.PhotoPath, &dob, &v.Location, &v.Relocation); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return queries.ResumeProfileView{}, queries.ErrResumeNotFound
		}
		return queries.ResumeProfileView{}, fmt.Errorf("get resume profile: %w", err)
	}
	v.Patronymic = patronymic
	v.DateOfBirth = dob.Format("2006-01-02")
	rows, err := r.pool.Query(ctx, `SELECT value FROM resume_profile_work_formats WHERE profile_id=(SELECT id FROM resume_profiles WHERE resume_id=$1) ORDER BY position`, q.ResumeID.UUID())
	if err != nil {
		return queries.ResumeProfileView{}, fmt.Errorf("get profile work formats: %w", err)
	}
	defer rows.Close()
	v.WorkFormat = []string{}
	for rows.Next() {
		var x string
		if err := rows.Scan(&x); err != nil {
			return queries.ResumeProfileView{}, fmt.Errorf("scan profile work format: %w", err)
		}
		v.WorkFormat = append(v.WorkFormat, x)
	}
	if err := rows.Err(); err != nil {
		return queries.ResumeProfileView{}, fmt.Errorf("iterate profile work formats: %w", err)
	}
	return v, nil
}

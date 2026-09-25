package readers

import (
	"context"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ContactsReader struct{ pool *pgxpool.Pool }

func NewContactsReader(pool *pgxpool.Pool) *ContactsReader { return &ContactsReader{pool: pool} }
func (r *ContactsReader) GetResumeContacts(ctx context.Context, q queries.GetResumeContactsQuery) (queries.ResumeContactsView, error) {
	if r == nil || r.pool == nil {
		return queries.ResumeContactsView{}, fmt.Errorf("get resume contacts: reader dependencies are nil")
	}
	var v queries.ResumeContactsView
	var phone, tg, email *string
	var cid uuid.UUID
	if err := r.pool.QueryRow(ctx, `SELECT c.id::text,c.resume_id::text,c.phone,c.telegram,c.email FROM resume_contacts c JOIN resume_roots r ON r.id=c.resume_id WHERE c.resume_id=$1 AND r.user_id=$2`, q.ResumeID.UUID(), q.UserID.UUID()).Scan(&cid, &v.ResumeID, &phone, &tg, &email); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return v, queries.ErrResumeNotFound
		}
		return v, fmt.Errorf("get resume contacts: %w", err)
	}
	v.ID = cid.String()
	v.Phone = phone
	v.Telegram = tg
	v.Email = email
	rows, err := r.pool.Query(ctx, `SELECT value FROM resume_contact_hidden WHERE contacts_id=$1 ORDER BY position`, cid)
	if err != nil {
		return v, fmt.Errorf("get hidden contacts: %w", err)
	}
	defer rows.Close()
	v.HiddenContacts = []string{}
	for rows.Next() {
		var x string
		if err := rows.Scan(&x); err != nil {
			return v, fmt.Errorf("scan hidden contact: %w", err)
		}
		v.HiddenContacts = append(v.HiddenContacts, x)
	}
	if err := rows.Err(); err != nil {
		return v, fmt.Errorf("iterate hidden contacts: %w", err)
	}
	var tr, te, sr, se *string
	err = r.pool.QueryRow(ctx, `SELECT title_ru,title_en,subtitle_ru,subtitle_en FROM resume_contacts_text WHERE contacts_id=$1`, cid).Scan(&tr, &te, &sr, &se)
	if err == nil {
		v.ContactsText = &queries.ResumeContactsTextView{Title: localizedPtr(tr, te), Subtitle: localizedPtr(sr, se)}
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return v, fmt.Errorf("get contacts text: %w", err)
	}
	return v, nil
}
func localizedPtr(ru, en *string) *queries.ResumeLocalizedTextView {
	if ru == nil && en == nil {
		return nil
	}
	return &queries.ResumeLocalizedTextView{RU: valueOrEmpty(ru), EN: valueOrEmpty(en)}
}
